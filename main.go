package main

import (
	"context"
	"errors"
	"flag"
	"os"

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
	var testFilePtr *string

	setUITestingFlag := func(flagVal string) error {
		if flagVal == "" {
			testFilePtr = nil
			return errors.New("-ui_test flag requires a value")
		}
		if flagVal == "runtest" || flagVal == "unmarshalTest" {
			// pass, nothin, use this
		} else if _, err := os.Stat(flagVal); errors.Is(err, os.ErrNotExist) {
			panic(err)
		}

		testFilePtr = &flagVal

		return nil
	}

	flag.BoolFunc("ui", "show UI", func(_ string) error { setFlag(&showUI, true); return nil })
	flag.Func("ui_test", "Test configuration file for testing UI", setUITestingFlag)
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
		if testFilePtr != nil {
			panic(errors.New("can't specify test file and standard ui"))
		}
		quit := ui.StartMissionControl()
		defer quit()
	} else if testFilePtr != nil {
		if *testFilePtr == "runtest" {
			ui.UIOnTest()
		} else if *testFilePtr == "unmarshalTest" {
			ui.TestUI([]byte(`{"drawings":[],"pause":3,"pause_div":2}`))
		} else {
			ui.TestUIFromFile(testFilePtr)
		}
	}

}
