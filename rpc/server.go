// This package is for the RPC/REST API
package rpc

import (
	"log"
	"net/http"

	"github.com/pokt-network/pocket-core/config"
	"github.com/pokt-network/pocket-core/logs"
	"github.com/pokt-network/pocket-core/rpc/client"
	"github.com/pokt-network/pocket-core/rpc/relay"
	"github.com/pokt-network/pocket-core/rpc/shared"
)

// "StartServers" executes the specified configuration for the client.
func StartServers() {
	if config.Get().CRPC { // if flag set
		go StartClientRPC(config.Get().CRPCPort) // run the client rpc in a goroutine
	}
	if config.Get().RRPC { // if flag set
		go StartRelayRPC(config.Get().RRPCPort) // run the relay rpc in a goroutine
	}
}

// "startClientRPC" starts the client RPC/REST API server at a specific port.
func StartClientRPC(port string) {
	log.Fatal(http.ListenAndServe(":"+port, shared.Router(client.Routes()))) // This starts the client RPC API.
	logs.NewLog("Started client server", logs.InfoLevel, logs.JSONLogFormat)
}

// "startRelayRPC" starts the client RPC/REST API server at a specific port.
func StartRelayRPC(port string) {
	log.Fatal(http.ListenAndServe(":"+port, shared.Router(relay.Routes()))) // This starts the relay RPC API.
	logs.NewLog("Started relay server", logs.InfoLevel, logs.JSONLogFormat)
}
