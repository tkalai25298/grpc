package server

import (
	"context"
	"io"
	"time"

	"github.com/hashicorp/go-hclog"
	protos "github.com/learn_grpc/protos/currency"
)

// Currency is a gRPC server it implements the methods defined by the CurrencyServer interface
type Currency struct {
	log hclog.Logger
}

// NewCurrency creates a new Currency 
func NewCurrency(l hclog.Logger) *Currency {
	return &Currency{l}
}

// GetRate implements the CurrencyServer GetRate method and returns the currency exchange rate
// for the two given currencies.
func (c *Currency) GetRate(ctx context.Context, rr *protos.RateRequest) (*protos.RateResponse, error) {
	c.log.Info("Handle request for GetRate", "base", rr.GetBase(), "dest", rr.GetDestination())
	return &protos.RateResponse{Rate: 0.5}, nil
}

func (c *Currency) GetStreamRate(str protos.Currency_GetStreamRateServer) error {

	//using go func to simultaneously handle client and server
	go func ()  {

		//receiving client stream 
		//looping over to receive continous requests
		for {
			rr,err := str.Recv()

			//if client closes outbound connections
			if err == io.EOF {
				c.log.Info("Client closed the connection")
				break
			}

			//for generic errors
			if err != nil {
				c.log.Error("Unable to read from client","error",err)
				break
			}
			c.log.Info("Handling client requests","request",rr)
		}
	}()

	//sending server stream msg
	for {
		err := str.Send(&protos.RateResponse{Rate: 25})
		if err != nil {
			return err
		}
		time.Sleep(5 *time.Second)
	}
} 