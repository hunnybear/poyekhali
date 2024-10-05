package main

import (
	"context"

	krpcgo "github.com/atburke/krpc-go"

	//"github.com/atburke/krpc-go/krpc"
	"github.com/atburke/krpc-go/spacecenter"
)

func main() {
	// Connect to the kRPC server with all default parameters.
	client := krpcgo.DefaultKRPCClient()
	err := client.Connect(context.Background())
	if err != nil {
		panic(err)
	}
	defer client.Close()

	sc := spacecenter.New(client)

	vessel, _ := sc.ActiveVessel()
	control, _ := vessel.Control()

	control.SetSAS(true)
	control.SetRCS(false)
	control.SetThrottle(1.0)
	control.ActivateNextStage()
}
