package main

import (
	"flag"
	"fmt"
	"github.com/eclipse/che-go-jsonrpc"
	"github.com/eclipse/che-go-jsonrpc/jsonrpcws"
	"github.com/eclipse/che-plugin-broker/model"
	"log"
	"os"
	"sync"
	"time"
)



func ConnectOrFail(endpoint string, token string) *jsonrpc.Tunnel {
	tunnel, err := Connect(endpoint, token)
	if err != nil {
		log.Fatalf("Couldn't connect to endpoint '%s', due to error '%s'", endpoint, err)
	}
	return tunnel
}
func Connect(endpoint string, token string) (*jsonrpc.Tunnel, error) {
	conn, err := jsonrpcws.Dial(endpoint, token)
	if err != nil {
		return nil, err
	}
	return jsonrpc.NewManagedTunnel(conn), nil
}

func main() {
	log.SetOutput(os.Stdout)
	var cheUrl string
	flag.StringVar(&cheUrl, "cheurl", "ws://che-eclipse-che.192.168.64.12.nip.io/api/websocket", "Che endpoint url")
	var cheToken string
	flag.StringVar(&cheToken, "token", "", "Che token")

	numOfThreads := flag.Int("tnum", 100, "number of threads sending messages")
	numOfMessages := flag.Int("mnum", 10, "number of messages to send")
	flag.PrintDefaults()
	flag.Parse()
	wg := &sync.WaitGroup{}
	wg.Add(*numOfThreads)

	for i := 0; i < *numOfThreads; i++ {
		go func(endpoint string, token string, th int) {
			tunnel:= ConnectOrFail(endpoint, token)
			defer tunnel.Close()

			tunnel.Conn()

			for j := 0; j < *numOfMessages; j++ {
				message := fmt.Sprintf("Messaget from thread %d number %d", th, j)
				event := &model.PluginBrokerLogEvent{
					RuntimeID: model.RuntimeID{Workspace:"ws1", Environment:"e1", OwnerId:"own1"},
					Text:      message ,
					Time:      time.Now(),
				}
				log.Print(message)
				if err := tunnel.Notify(event.Type(), event); err != nil {
					log.Fatalf("Trying to send event of type '%s' to closed tunnel '%s'", event.Type(), tunnel.ID())
				}
			}

			wg.Done()

		}(cheUrl,cheToken, i)
	}
	wg.Wait()
}