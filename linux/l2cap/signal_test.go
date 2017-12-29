package l2cap

import (
	"testing"
)

var configurationResponseBytes []byte = []byte{0x40, 0x00, 0x00, 0x00, 0x00, 0x00}

func TestUnMarshalConfigurationResponse(t *testing.T) {
	rsp := &ConfigurationResponse{}
	if err := rsp.Unmarshal(configurationResponseBytes); err != nil {
		t.Fatalf("Error: %s", err)
	}
}
