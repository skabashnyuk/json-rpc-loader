package main

import (
	"fmt"
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
var RateCounter = ratecounter.NewRateCounter(1 * time.Second)
var WG = &sync.WaitGroup{}

type Configuration struct {
	CheHost   string        `required:"true" split_words:"true" desc:"Che Server host"`
	CheToken  string        `split_words:"true" desc:"User token for multi-user che"`
	Threads   int           `default:"10" split_words:"true" desc:"Number of clients used to send messages"`
	WsTimeout time.Duration `default:"10s" split_words:"true" desc:"Websocket connection timeout "`
	Secure    bool          `default:"false" split_words:"true" desc:"Whatever secure websocket aka wss connection should be used"`
}

func PrintRate() {
	for !Suspending {
		fmt.Printf("Iteration at %d/5s \n", RateCounter.Rate())
		time.Sleep(5 * time.Second)
	}
	WG.Done()
}

func SendMessagesInLoop(wsUrlMajor, wsUrlMinor, token, senderId string) {
	loader := &Loader{}
	defer loader.Close()
	loader.Init(wsUrlMajor, wsUrlMinor, token)
	for !Suspending {
		loader.Start()
		RateCounter.Incr(1)
	}
	fmt.Printf("Sending complete %s\n", senderId)
	WG.Done()
}

func Init(configuration *Configuration) {

	WG.Add(configuration.Threads + 1)
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

	format := "Configuration is set to:\nCheHost: %s\nToken: %s\nThreads: %d\nTimeout: %s\nSecure: %v\n"
	_, err = fmt.Printf(format, configuration.CheHost, configuration.CheToken, configuration.Threads, configuration.WsTimeout, configuration.Secure)
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

	for i := 0; i < configuration.Threads; i++ {
		go SendMessagesInLoop(major.String(), minor.String(), configuration.CheToken, "Loader "+strconv.Itoa(i))
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
