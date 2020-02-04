//go:generate protoc -I ../proto/ --go_out=plugins=grpc:./pb.go/ ../proto/pubnub.proto
//go:generate protoc -I ../proto/ --go_out=plugins=grpc:./pb.go/ ../proto/pubnub.tracking.proto
//go:generate protoc -I ../proto/ --go_out=plugins=grpc:./pb.go/ ../proto/pubnub.types.proto

package main

import (
	"context"
	"io"
	"log"
	"time"

	"github.com/golang/protobuf/jsonpb"
	pb "github.com/jpollock/grpc-demo/go/pb.go"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const (
	address = "pubnub-arke.prd-eks-bom-1.prd-eks.ps.pn:80"
	pubkey  = "<insert publish key>"
	subkey  = "<insert subscribe key>"
	channel = "demo"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewPubSubClient(conn)

	// Create a context with account credentials
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	ctx = metadata.AppendToOutgoingContext(
		ctx,
		"publish-key", pubkey,
		"subscribe-key", subkey,
	)

	// Subscribe to the channel
	subscription := pb.Subscription{Channel: channel}
	subStream, err := client.Subscribe(ctx, &subscription)
	if err != nil {
		log.Fatalf("%v.Subscribe(_) = _, %v", client, err)
	}

	// Wait for stream to be established
	subStream.Header()
	/*
	messages := []*pb.Message{{}, {}}

	jsonpb.UnmarshalString(`{
		"channel": "demo",
		"data": {
			"foo": "123"
		}
	}`, messages[0])

	jsonpb.UnmarshalString(`{
		"channel": "demo",
		"data": {
			"foo": "456"
		}
	}`, messages[1])

		
	pubStream, err := client.StreamingPublish(ctx)
	if err != nil {
		log.Fatalf("%v.StreamingPublish(_) = _, %v", client, err)
	}

	pubc := make(chan struct{})
	go func() {
		for {
			in, err := pubStream.Recv()
			if err == io.EOF {
				// read done.
				close(pubc)
				return
			}
			if err != nil {
				log.Fatalf("Failed to receive a status: %v", err)
			}
			log.Printf("Got status: `%v`", in)
		}
	}()
	for _, message := range messages {
		if err := pubStream.Send(message); err != nil {
			log.Fatalf("Failed to send a message: %v", err)
		}
	}
	pubStream.CloseSend()
	<-pubc

	// Send a message synchronously
	message := pb.Message{}
	jsonpb.UnmarshalString(`{
		"channel": "demo",
		"data": {
			"foo": "789"
		}
	}`, &message)

	status, err := client.Publish(ctx, &message)
	if err != nil {
		log.Fatalf("%v.Publish(_) = _, %v", client, err)
	}

	log.Printf("Got status: `%v`", status)
	*/
	// Recieve messages
	subc := make(chan struct{})
	go func() {
		counter := 0
		marshaler := jsonpb.Marshaler{Indent: "  "}
		for {
			message, err := subStream.Recv()
			if err == io.EOF {
				// read done.
				close(subc)
				return
			}
			if err != nil {
				log.Fatalf("Failed to receive a message: %v", err)
			}

			pretty, err := marshaler.MarshalToString(message)
			if err != nil {
				log.Fatalf("Failed to pretty-print message: %v", err)
			} else {
				log.Printf("Got message: %s", pretty)
			}

			counter++
			if counter >= 3 {
				close(subc)
				return
			}
		}
	}()

	subStream.CloseSend()
	<-subc
}
