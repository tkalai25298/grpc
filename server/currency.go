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
	subscriptions map[protos.Currency_GetStreamRateServer][]*protos.RateRequest //maintaining different clients & sending resp to particular client request
}

// NewCurrency creates a new Currency 
func NewCurrency(l hclog.Logger) *Currency {

	c := &Currency{l,make(map[protos.Currency_GetStreamRateServer][]*protos.RateRequest)}
	go c.handleSubscriptions()

	return c
}


//sending subscriptions for every 5seconds
func (c *Currency) handleSubscriptions() {
	
		for { //continously sending server rates to clients
			c.log.Info("server sending rates")


			// loop over subscribed clients
			for k, v := range c.subscriptions {

				// loop over subscribed rates
				for _, rr := range v {
					r, err := c.GetRate(context.Background(),rr)
					if err != nil {
						c.log.Error("Unable to get update rate", "base", rr.Base, "destination", rr.Destination)
					}
	
					err = k.Send(&protos.RateResponse{Base: rr.Base, Destination: rr.Destination, Rate: r.Rate})
					if err != nil {
						c.log.Error("Unable to send updated rate", "base", rr.Base, "destination", rr.Destination)
					}
					
				}
			}
			time.Sleep(5*time.Second)
		}
	
}

// GetRate implements the CurrencyServer GetRate method and returns the currency exchange rate
// for the two given currencies.
func (c *Currency) GetRate(ctx context.Context, rr *protos.RateRequest) (*protos.RateResponse, error) {
	c.log.Info("Handle request for GetRate", "base", rr.GetBase(), "dest", rr.GetDestination())
	return &protos.RateResponse{Rate: 0.5}, nil
}

func (c *Currency) GetStreamRate(str protos.Currency_GetStreamRateServer) error {

	//using go func to simultaneously handle client and server
	// go func ()  {

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

			//caching the client requests to subscriptions map 
			//creating subscriptions and adding the list of each rateRequests from client to the server streaming
			rrs,ok := c.subscriptions[str]
			//if the key doesnt exists create new 
			if !ok {
				rrs = []*protos.RateRequest{}
			}

			//appending all the client requests made
			rrs = append(rrs, rr)
			c.subscriptions[str] = rrs
		}
	// }()

	return nil
} 