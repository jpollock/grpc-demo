package main

import (
	"context"
	"io"
	"log"
	"math"
	"math/rand"
	"time"
	"os"

	pb "github.com/jpollock/grpc-demo/go/pb.go"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const (
	//address = "0.0.0.0:50051"
	address = "pubnub-arke.prd-eks-bom-1.prd-eks.ps.pn:80"
	// address = "13.56.150.134:50051"
	channel = "demo"
)


func main() {
	pubkey  := os.Getenv("PUBLISH_KEY")
	subkey  := os.Getenv("SUBSCRIBE_KEY")
	

	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewDriverTrackingClient(conn)

	// Create a context with account credentials
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1800)
	defer cancel()
	ctx = metadata.AppendToOutgoingContext(
		ctx,
		"publish-key", pubkey,
		"subscribe-key", subkey,
	)

	// Create a publish stream.
	pubStream, err := client.Publish(ctx)
	if err != nil {
		log.Fatalf("%v.Publish(_) = _, %v", client, err)
	}

	pubc := make(chan struct{})

	go func() {
		data := pb.DriverTrackingMessage{
			DriverId:     "JoeBob",
			OrderId: 	  "Jeremy",
			DriverStatus: pb.DriverTrackingMessage_WAITING_FOR_ASSIGNMENT,
			Location: &pb.Location{
				Latitude:  37.7749,
				Longitude: 122.4194,
			},
			Heading:  0.0,
			Velocity: 115.0,
		}

		for {
			simulate(&data)

			message := pb.DriverTrackingEnvelope{
				Channel: channel,
				Data:    &data,
			}

			log.Printf("Sending message: `%v`", message)
			if err := pubStream.Send(&message); err != nil {
				log.Fatalf("Failed to send a message: %v", err)
				close(pubc)
				return
			}

			in, err := pubStream.Recv()
			if err == io.EOF {
				// read done.
				close(pubc)
				return
			}
			if err != nil {
				log.Fatalf("Failed to receive a status: %v", err)
			}

			log.Printf("Got status: `%s`", in.Message)

			// Sleep for 5 seconds
			duration, _ := time.ParseDuration("5s")
			time.Sleep(duration)
		}
	}()

	// Wait for goroutine loop to exit
	<-pubc
}

func simulate(data *pb.DriverTrackingMessage) {
	// Randomize heading (in radians)
	heading := (rand.Float64() - 0.5) * math.Pi / 4.0
	data.Heading += float32(heading)
	if data.Heading < 0.0 {
		data.Heading += float32(math.Pi * 2.0)
	}

	// Randomize velocity (in meters per second)
	velocity := math.Max(math.Min(float64(data.Velocity)+(rand.Float64()-0.5)/3600.0, 1.0/60.0), 0.0)
	data.Velocity = float32(velocity)

	// Update location
	latitude := math.Mod(float64(data.Location.Latitude)+math.Sin(heading)*velocity, 180.0)
	longitude := math.Mod(float64(data.Location.Longitude)+math.Cos(heading)*velocity, 360.0)
	data.Location.Latitude = float32(latitude)
	data.Location.Longitude = float32(longitude)
}
