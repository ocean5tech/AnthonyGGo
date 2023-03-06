package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/ocean5tech/pricefetcher/client"
)

func main() {
	client := client.New("http://localhost:9091")
	price, err := client.FetchPrice(context.Background(), "GG1")

	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", price)

	listenAddr := flag.String("listenAddr", ":9091", "listen address the service is running")
	flag.Parse()
	svc := NewLoggingService(NewMetricService(&priceFetcher{}))

	server := NewJSONAPIServer(*listenAddr, svc)

	server.Run()

}
