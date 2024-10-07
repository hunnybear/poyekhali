package main

import (
	"context"

	"flag"

	krpcgo "github.com/atburke/krpc-go"

	//"github.com/atburke/krpc-go/krpc"
	"github.com/atburke/krpc-go/spacecenter"
	"github.com/hunnybear/poyekhali/ui"
)

func setFlag(flagPtr *bool, value bool) {
	*flagPtr = value
}

func main() {
	var client *krpcgo.KRPCClient
	showUI := false
	launch := false
	connect := false

	flag.BoolFunc("ui", "show UI", func(_ string) error { setFlag(&showUI, true); return nil })
	flag.BoolFunc("launch", "launchify me, capn", func(_ string) error { setFlag(&launch, true); return nil })
	flag.BoolFunc("connect", "connect to server", func(_ string) error { setFlag(&connect, true); return nil })
	flag.Parse()
	if launch == true || connect == true {
		// Connect to the kRPC server with all default parameters.
		client = krpcgo.DefaultKRPCClient()
		err := client.Connect(context.Background())
		if err != nil {
			panic(err)
		}
		defer client.Close()
	}

	if launch == true {
		sc := spacecenter.New(client)

		vessel, _ := sc.ActiveVessel()
		control, _ := vessel.Control()

		control.SetSAS(true)
		control.SetRCS(false)
		control.SetThrottle(1.0)
		control.ActivateNextStage()
	}

	if showUI == true {
		ui.Ui()
	}

}
