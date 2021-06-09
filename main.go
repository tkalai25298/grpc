package main

import (
	"fmt"
	"net"
	"os"

	"github.com/hashicorp/go-hclog"
	protos "github.com/learn_grpc/protos/currency"
	data "github.com/learn_grpc/data"
	"github.com/learn_grpc/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	log := hclog.Default()

	//getting all the rates from ecb based on euro
	rates,err := data.NewRates(log)

	if err != nil {
		log.Error("Unable to generate rates", "error", err)
		os.Exit(1)
	}

	// create a new gRPC server, use WithInsecure to allow http connections
	gs := grpc.NewServer()

	// create an instance of the Currency server
	c := server.NewCurrency(log,rates)

	// register the currency server
	protos.RegisterCurrencyServer(gs, c)

	// register the reflection service which allows clients to determine the methods
	// for this gRPC service
	reflection.Register(gs)

	// create a TCP socket for inbound server connections
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", 9092))
	if err != nil {
		log.Error("Unable to create listener", "error", err)
		os.Exit(1)
	}
	
	// listen for requests
	gs.Serve(l)
}