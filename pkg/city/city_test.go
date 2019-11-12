package city

import (
	"context"
	"encoding/json"
	"log"
	"testing"
	"time"

	"google.golang.org/grpc"
)

const (
	address = "localhost:9090"
)

func TestCityClient(t *testing.T) {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := NewCityManagerClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.CityList(ctx, &CityQuery{Citycode: ""})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	cbts, _ := json.Marshal(r.Citys)
	log.Printf("Citys: %s", cbts)

	ad, err := c.AreaList(ctx, &AreaQuery{Citycode: "0571"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	abts, _ := json.Marshal(ad.Areas)
	log.Printf("Areas: %s", abts)
}
