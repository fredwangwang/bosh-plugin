package main

import (
	"code.cloudfoundry.org/go-loggregator"
	"fmt"
	"github.com/jasonlvhit/gocron"
	"log"
	"net/http"
	"os"
)

func main() {
	port, ok := os.LookupEnv("PORT")
	if !ok {
		panic("PORT is not provided")
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello world!")
	})

	go func() { http.ListenAndServe(":"+port, nil) }()

	tlsConfig, err := loggregator.NewIngressTLSConfig(
		os.Getenv("METRON_CA_CERT_PATH"),
		os.Getenv("METRON_CERT_PATH"),
		os.Getenv("METRON_KEY_PATH"),
	)
	if err != nil {
		log.Fatal("Could not create TLS config ", err)
	}

	client, err := loggregator.NewIngressClient(
		tlsConfig,
		loggregator.WithAddr("localhost:3458"),
	)
	if err != nil {
		log.Fatal("Could not connect to metron ", err)
	}

	gocron.Every(2).Seconds().Do(emitCounter, client)
	<-gocron.Start()
}

func emitCounter(client *loggregator.IngressClient) {
	client.EmitCounter("sample-plugin")
	fmt.Println("hello")
}
