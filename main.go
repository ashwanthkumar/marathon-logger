package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gambol99/go-marathon"
	flag "github.com/spf13/pflag"
)

var marathonURI string
var mesosSlavePort int
var appCheckInterval time.Duration

var appMonitor AppMonitor

func main() {
	os.Args[0] = "marathon-logger"
	flag.Parse()

	client, err := marathonClient(marathonURI)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	appMonitor = AppMonitor{
		Client:        client,
		CheckInterval: appCheckInterval,
	}
	appMonitor.Start()
	appMonitor.RunWaitGroup.Wait()
}

func init() {
	flag.StringVar(&marathonURI, "uri", "", "Marathon URI to connect")
	flag.IntVar(&mesosSlavePort, "slave-port", 5051, "Mesos slave port")
	flag.DurationVar(&appCheckInterval, "check-interval", 30*time.Second, "Frequency at which we check for new tasks")
}

func marathonClient(uri string) (marathon.Marathon, error) {
	config := marathon.NewDefaultConfig()
	config.URL = uri
	config.HTTPClient = &http.Client{
		Timeout: (30 * time.Second),
	}

	return marathon.NewClient(config)
}
