package linux

import (
	"github.com/thomascriley/ble"
	"github.com/thomascriley/ble/linux/gatt"
	"github.com/thomascriley/ble/linux/hci"
	"github.com/thomascriley/ble/linux/rfcomm"
)

var testDevice ble.Device = &Device{}
var testConn ble.Conn = &hci.Conn{}
var testClient ble.ClientBLE = &gatt.Client{}
var testRFCOMMClient ble.ClientRFCOMM = &rfcomm.Client{}
var testAdv ble.Advertisement = &hci.Advertisement{}