package detail

import (
	"context"
	"encoding/json"
	"log"
	"testing"
	"time"

	"google.golang.org/grpc"
)

const (
	address = "localhost:9091"
)

func TestDetailClient(t *testing.T) {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := NewDetailClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	ds, err := c.GetDetail(ctx, &DetailQuery{Citycode: "010", Date: "2019-11-12"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	cbts, _ := json.Marshal(ds)
	log.Printf("Greeting: %s", cbts)
}
