package main

import (
	"context"
	"log"

	protos "github.com/learn_grpc/protos/currency"
	"google.golang.org/grpc"
)

func main() {

	//creating client conn for currency server
	conn,err := grpc.Dial("localhost:9092",grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	defer conn.Close()

	currClient := protos.NewCurrencyClient(conn)

	rr := &protos.RateRequest{
		Base: protos.Currencies(protos.Currencies_value["INR"]),
		Destination: protos.Currencies(protos.Currencies_value["USD"]),
	}

	resp,err := currClient.GetRate(context.Background(),rr)
	if err != nil {
		log.Fatalf("couldn't getRate : %v",err)
	}

	log.Printf("rate: %v",resp.Rate)

}