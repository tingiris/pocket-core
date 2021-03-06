package unit

import (
	"testing"

	"github.com/pokt-network/pocket-core/db"
	"github.com/pokt-network/pocket-core/node"
)

func DummyNode() node.Node {
	chains := []node.Blockchain{{Name: "ethereum", NetID: "1"}}
	n := node.Node{
		GID:         "test",
		IP:          "ipfromintegrationtest",
		RelayPort:   "portfromintegrationtest",
		ClientID:    "0",
		CliVersion:  "0",
		Blockchains: chains,
	}
	return n
}

func TestPut(t *testing.T) {
	d := db.DB()
	_, err := d.Add(DummyNode())
	if err != nil {
		t.Fatalf(err.Error())
	}
}

func TestRemove(t *testing.T) {
	d := db.DB()
	_, err := d.Remove(DummyNode())
	if err != nil {
		t.Fatalf(err.Error())
	}
}
