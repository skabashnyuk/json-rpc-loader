package main

import (
	"encoding/json"
	"fmt"
	. "github.com/Nerzal/gocloak"
	"github.com/gorilla/websocket"
	"github.com/kelseyhightower/envconfig"
	"github.com/paulbellamy/ratecounter"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"time"
)

var Suspending = false
var RateCounter = ratecounter.NewRateCounter(1 * time.Second)
var WG = &sync.WaitGroup{}

type Configuration struct {
	CheHost      string        `required:"true" split_words:"true" desc:"Che Server host"`
	Client       int           `default:"10" split_words:"true" desc:"Number of clients used to send messages"`
	WsTimeout    time.Duration `default:"10s" split_words:"true" desc:"Websocket connection timeout"`
	Secure       bool          `default:"false" split_words:"true" desc:"Whether or not to use secure websocket aka wss connection"`
	Multiplexing bool          `default:"false" split_words:"true" desc:"Whether or not to use single websocket connection by each client to send request"`
	MultiUser    bool          `default:"false" split_words:"false" desc:"Use che in multi-user mode"`
	UserName     string        `default:"admin" split_words:"false" desc:"Che user name"`
	UserPassword string        `default:"admin" split_words:"false" desc:"Che user password"`
	СheRealm     string        `default:"che" split_words:"false" desc:"Multi user  Che realm"`
	СheClientId  string        `default:"che-public" split_words:"false" desc:"Keycloak client id of Che"`
	WorkspaceId  string        `required:"true" split_words:"false" desc:"Workspace ide used to generate load"`
}

type GetToken func() string

type KeycloakTokenProvider struct {
	CheHost         string
	UserName        string
	UserPassword    string
	Realm           string
	ClientId        string
	token           *JWT
	goCloak         GoCloak
	lastRefreshTime time.Time
}

func dummyToken() string {
	return ""
}

func newKeycloakTokenProvider(cheHost, userName, userPassword, realm, clientId string) *KeycloakTokenProvider {

	resp, err := http.Get("http://" + cheHost + "/api/keycloak/settings")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	var result map[string]interface{}

	// Unmarshal or Decode the JSON to the interface.
	json.Unmarshal(body, &result)
	authServerUrl := result["che.keycloak.auth_server_url"].(string)
	gcloak := NewClient(authServerUrl[:len(authServerUrl)-5])

	jwtToken, err := gcloak.Login(clientId, "null", realm, userName, userPassword)
	if err != nil {
		panic("Something wrong with the credentials or url")
	}
	return &KeycloakTokenProvider{
		CheHost:         cheHost,
		UserName:        userName,
		UserPassword:    userPassword,
		Realm:           realm,
		ClientId:        clientId,
		goCloak:         gcloak,
		token:           jwtToken,
		lastRefreshTime: time.Now(),
	}
}

func (provider *KeycloakTokenProvider) getToken() string {

	if provider.lastRefreshTime.Add(time.Duration(provider.token.ExpiresIn-5) * time.Second).Before(time.Now()) {
		newToken, err := provider.goCloak.RefreshToken(provider.token.RefreshToken, provider.ClientId, "null", provider.Realm)
		if err != nil {
			panic("Something wrong with the credentials or url")
		}

		provider.lastRefreshTime = time.Now()
		fmt.Printf("Token refreshed \n")
		provider.token = newToken
	}
	return provider.token.AccessToken
}

func PrintRate() {
	for !Suspending {
		fmt.Printf("Iteration at %d/s \n", RateCounter.Rate())
		time.Sleep(5 * time.Second)
	}
	WG.Done()
}

func SendMessagesInLoop(wsUrlMajor, wsUrlMinor, senderId, workspaceId string, multiplexing bool, tokenProvider GetToken) {

	if multiplexing {
		loader := &Loader{}
		defer loader.Close()
		loader.Init(wsUrlMajor, wsUrlMinor, workspaceId, tokenProvider)
		for !Suspending {
			loader.Start()
			RateCounter.Incr(1)
		}
	} else {
		for !Suspending {
			loader := &Loader{}
			loader.Init(wsUrlMajor, wsUrlMinor, workspaceId, tokenProvider)
			loader.Start()
			loader.Close()
			RateCounter.Incr(1)
		}
	}

	fmt.Printf("Sending complete %s\n", senderId)
	WG.Done()
}

func Init(configuration *Configuration) {

	WG.Add(configuration.Client + 1)
	log.SetOutput(os.Stdout)
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func main() {

	var configuration Configuration
	envconfig.Usage("JsonRpcLoader", &configuration)
	err := envconfig.Process("JsonRpcLoader", &configuration)
	if err != nil {
		log.Fatal(err.Error())
	}

	format := "Configuration is set to:\nCheHost: %s\nToken: %s\nThreads: %d\nTimeout: %s\nSecure: %v\nMultiplexing: %v\n"
	_, err = fmt.Printf(format, configuration.CheHost, "", configuration.Client, configuration.WsTimeout, configuration.Secure, configuration.Multiplexing)
	if err != nil {
		log.Fatal(err.Error())
	}
	Init(&configuration)
	websocket.DefaultDialer.HandshakeTimeout = configuration.WsTimeout

	var major, minor strings.Builder
	if configuration.Secure {
		major.WriteString("wss://")
		minor.WriteString("wss://")
	} else {
		major.WriteString("ws://")
		minor.WriteString("ws://")
	}
	major.WriteString(configuration.CheHost)
	major.WriteString("/api/websocket")
	minor.WriteString(configuration.CheHost)
	minor.WriteString("/api/websocket-minor")

	tokeProvider := dummyToken

	if configuration.MultiUser {
		tp := newKeycloakTokenProvider(configuration.CheHost, configuration.UserName, configuration.UserPassword, configuration.СheRealm, configuration.СheClientId)
		tokeProvider = tp.getToken
	}

	for i := 0; i < configuration.Client; i++ {
		go SendMessagesInLoop(major.String(), minor.String(), "Loader "+strconv.Itoa(i), configuration.WorkspaceId, configuration.Multiplexing, tokeProvider)
	}

	go PrintRate()
	// After setting everything up!
	// Wait for a SIGINT (perhaps triggered by user with CTRL-C)
	// Run cleanup when signal is received
	signalChan := make(chan os.Signal, 1)
	cleanupDone := make(chan struct{})
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		<-signalChan
		fmt.Println("\nReceived an interrupt, stopping services...")
		Suspending = true
		WG.Wait()
		close(cleanupDone)
	}()
	<-cleanupDone
}
