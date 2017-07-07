package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/net/context"

	"github.com/currantlabs/ble"
	"github.com/currantlabs/ble/examples/lib/dev"
)

var (
	device = flag.String("device", "default", "implementation of ble")
	du     = flag.Duration("du", 5*time.Second, "scanning duration")
	dup    = flag.Bool("dup", true, "allow duplicate reported")
	bredr  = flag.Bool("bredr", false, "scan fro BR/EDR devices")
)

var ctx context.Context

var address = ble.NewAddr("00:a0:96:14:18:5b")

func main() {
	flag.Parse()

	d, err := dev.NewDevice(*device)
	if err != nil {
		log.Fatalf("can't new device : %s", err)
	}
	ble.SetDefaultDevice(d)

	// Scan for specified durantion, or until interrupted by user.
	fmt.Printf("Scanning for %s...\n", *du)
	ctx = ble.WithSigHandler(context.WithTimeout(context.Background(), *du))

	if *bredr {
		chkErr(ble.Inquire(ctx, 255, inqHandler))
	} else {
		chkErr(ble.Scan(ctx, *dup, advHandler, nil))
	}
}

func inqHandler(i ble.Inquiry) {
	fmt.Printf("[%s] %3d\n", i.Address(), i.RSSI())

	if i.Address().String() == address.String() {
		fmt.Printf("Found medtracker, dialing\n")
		cli, err := ble.DialRFCOMM(ctx, address)
		if err != nil {
			fmt.Printf("Error dialing: %s\n", err)
			return
		}
		_, err = cli.Write([]byte("status\n"))
		if err != nil {
			fmt.Printf("Error writing: %s\n", err)
		}
		bs := make([]byte, 1024)
		n, err := cli.Read(bs)
		if err != nil {
			fmt.Printf("Error reading: %s\n", err)
		}
		fmt.Printf("Read %s", bs[:n])
		cli.CancelConnection()
	}
}

func advHandler(a ble.Advertisement) {
	if a.Connectable() {
		fmt.Printf("[%s] C %3d:", a.Address(), a.RSSI())
	} else {
		fmt.Printf("[%s] N %3d:", a.Address(), a.RSSI())
	}
	comma := ""
	if len(a.LocalName()) > 0 {
		fmt.Printf(" Name: %s", a.LocalName())
		comma = ","
	}
	if len(a.Services()) > 0 {
		fmt.Printf("%s Svcs: %v", comma, a.Services())
		comma = ","
	}
	if len(a.ManufacturerData()) > 0 {
		fmt.Printf("%s MD: %X", comma, a.ManufacturerData())
	}
	fmt.Printf("\n")
}

func chkErr(err error) {
	switch errors.Cause(err) {
	case nil:
	case context.DeadlineExceeded:
		fmt.Printf("done\n")
	case context.Canceled:
		fmt.Printf("canceled\n")
	default:
		log.Fatalf(err.Error())
	}
}
