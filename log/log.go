package log

import (
	"fmt"
)

var (
	LoggingEnabled = false
	LogAdvertisements = false
	LogPackets = false
)

func Print(s string) {
	if LoggingEnabled {
		fmt.Println(s)
	}
}

func Printf(s string, i ...interface{}){
	if LoggingEnabled {
		fmt.Printf(s +"\n", i...)
	}
}

func Println(s string, i ...interface{}){
	if LoggingEnabled {
		fmt.Println(s)
	}
}

func AdvPrint(s string){
	if LogAdvertisements {
		Print(s)
	}
}

func AdvPrintf(s string, i ...interface{}){
	if LogAdvertisements {
		Printf(s,i...)
	}
}

func PacketPrint(s string){
	if LogPackets {
		Print(s)
	}
}

func PacketPrintf(s string, i ...interface{}){
	if LogPackets {
		Printf(s,i...)
	}
}
