package main

import (
	"fmt"
	"github.com/cenkalti/backoff"
	"github.com/eclipse/che-go-jsonrpc"
	"github.com/eclipse/che-go-jsonrpc/jsonrpcws"
	"github.com/eclipse/che-plugin-broker/model"
	"github.com/gorilla/websocket"
	"github.com/kelseyhightower/envconfig"
	"github.com/paulbellamy/ratecounter"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"time"
)

var Suspending = false
var MajorCounter = ratecounter.NewRateCounter(1 * time.Second)
var MinorCounter = ratecounter.NewRateCounter(1 * time.Second)
var WG = &sync.WaitGroup{}

type Configuration struct {
	CheHost      string        `required:"true" split_words:"true" desc:"Che Server host"`
	CheToken     string        `split_words:"true" desc:"User token for multi-user che"`
	MajorThreads int           `default:"10" split_words:"true" desc:"Number of clients used to send message to major websocket endpoint"`
	MinorThreads int           `default:"10" split_words:"true" desc:"Number of clients used to send message to minor websocket endpoint"`
	WsTimeout    time.Duration `default:"10s" split_words:"true" desc:"Websocket connection timeout "`
	Secure       bool          `default:"false" split_words:"true" desc:"Whatever secure websocket aka wss connection should be used"`
}

func ConnectRetryOrFail(endpoint string, token string) *jsonrpc.Tunnel {
	var result *jsonrpc.Tunnel
	err := backoff.Retry(func() error {
		tun, err2 := Connect(endpoint, token)
		if err2 != nil {
			return err2
		} else {
			result = tun
			return nil
		}

	}, backoff.NewExponentialBackOff())
	if err != nil {
		log.Panic("Couldn't connect to endpoint '%s', due to error '%s'", endpoint, err)
	}

	return result
}

func Connect(endpoint string, token string) (*jsonrpc.Tunnel, error) {
	conn, err := jsonrpcws.Dial(endpoint, token)
	if err != nil {
		return nil, err
	}
	return jsonrpc.NewManagedTunnel(conn), nil
}

func PrintRate() {
	for !Suspending {
		fmt.Printf("Major rate %d/s  Minor rate %d/s\n", MajorCounter.Rate(), MinorCounter.Rate())
		time.Sleep(5 * time.Second)
	}
	WG.Done()
}

func SendMessagesInLoop(wsUrl, token, senderId string, ratecounter *ratecounter.RateCounter) {

	tunnel := ConnectRetryOrFail(wsUrl, token)
	defer tunnel.Close()

	tunnel.Conn()

	for !Suspending {
		message := fmt.Sprintf("Message %s sent  from %s", RandStringRunes(rand.Intn(200)), senderId)
		event := &model.PluginBrokerLogEvent{
			RuntimeID: model.RuntimeID{Workspace: "ws1", Environment: "e1", OwnerId: "own1"},
			Text:      message,
			Time:      time.Now(),
		}
		//log.Print(message)
		if err := tunnel.Notify(event.Type(), event); err != nil {
			log.Printf("Trying to send event of type '%s' to closed tunnel '%s'", event.Type(), tunnel.ID())
			return
		}
		ratecounter.Incr(1)
	}
	fmt.Printf("Sending complete %s\n", senderId)
	WG.Done()
}

func Init(configuration *Configuration) {

	WG.Add(configuration.MinorThreads + configuration.MajorThreads + 1)
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

	format := "Configuration is set to:\nCheHost: %s\nToken: %s\nMajorThreads: %d\nMinorThreads: %d\nTimeout: %s\nSecure: %v\n"
	_, err = fmt.Printf(format, configuration.CheHost, configuration.CheToken, configuration.MajorThreads, configuration.MinorThreads, configuration.WsTimeout, configuration.Secure)
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

	for i := 0; i < configuration.MajorThreads; i++ {
		go SendMessagesInLoop(major.String(), configuration.CheToken, "major"+strconv.Itoa(i), MajorCounter)
	}

	for i := 0; i < configuration.MinorThreads; i++ {
		go SendMessagesInLoop(minor.String(), configuration.CheToken, "minor"+strconv.Itoa(i), MinorCounter)
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
