package session

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/google/flatbuffers/go"
	"github.com/pokt-network/pocket-core/common"
	"github.com/pokt-network/pocket-core/logs"
	"github.com/pokt-network/pocket-core/types"
	"io/ioutil"
	"log"
	"os"
	"reflect"
)

type Seed struct {
	DevID          []byte
	BlockHash      []byte
	RequestedChain []byte
	NodeList       []Node
}

// "NewSeed" is the constructor of the sessionSeed
func NewSeed(devID []byte, nodePoolFilePath string, requestedBlockchain []byte, blockHash []byte) (Seed, error) {
	np, err := FileToNodes(nodePoolFilePath)
	return Seed{DevID: devID, BlockHash: blockHash, RequestedChain: requestedBlockchain, NodeList: np}, err
}

// "FileToNodes" converts the world state noodPool.json file into a slice of session.Node
func FileToNodes(nodePoolFilePath string) ([]Node, error) {
	nws := FileToNWSSlice(nodePoolFilePath)
	return nwsToNodes(nws)
}

// "FileToNWSSlice" converts a file to a slice of NodeWorldState Nodes
func FileToNWSSlice(nodePoolFilePath string) []common.NodeWorldState {
	jsonFile, err := os.Open(nodePoolFilePath)
	if err != nil {
		logs.NewLog(err.Error(), logs.ErrorLevel, logs.JSONLogFormat)
		fmt.Println(err.Error())
	}
	defer jsonFile.Close()
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		logs.NewLog(err.Error(), logs.ErrorLevel, logs.JSONLogFormat)
		fmt.Println(err.Error())
	}
	var nodes []common.NodeWorldState
	err = json.Unmarshal(byteValue, &nodes)
	if err != nil {
		log.Fatal(err.Error())
	}
	return nodes
}

// "nwsToNodes" converts the NodeWorldState slice of nodes from the json file
// into a []Node which is used in our session seed
func nwsToNodes(nws []common.NodeWorldState) ([]Node, error) {
	var nodeList []Node
	for _, node := range nws {
		if !node.Active {
			continue
		}
		n, err := nwsToNode(node)
		if err != nil {
			return nodeList, err
		}
		nodeList = append(nodeList, n)
	}
	return nodeList, nil
}

// "nwsToNode" is a helper function to NWSToNode which takes a NodeWorldState Node
// and converts it to a session.Node
func nwsToNode(nws common.NodeWorldState) (Node, error) {
	chains := types.NewSet()
	var role role
	gid, ip, port, _ := nws.EnodeSplit()
	for _, c := range nws.Chains {
		marshalChain, err := common.MarshalBlockchain(flatbuffers.NewBuilder(0), c)
		if err != nil {
			return Node{}, err
		}
		chains.Add(hex.EncodeToString(common.SHA256FromBytes(marshalChain)))
	}
	switch nws.IsVal {
	case true:
		role = VALIDATE
	case false:
		role = SERVICE
	}
	return Node{GID: gid, IP: ip, Port: port, Chains: *chains, Role: role}, nil
}

// "ErrorCheck()" checks all of the fields of a seed to ensure that it is considered initially valid
func (s *Seed) ErrorCheck() error {
	if s.DevID == nil || len(s.DevID) == 0 {
		return NoDevID
	}
	if s.BlockHash == nil || len(s.BlockHash) == 0 {
		return NoBlockHash
	}
	if s.NodeList == nil || len(s.NodeList) == 0 {
		return NoNodeList
	}
	if reflect.DeepEqual(s.NodeList[0], Node{}) {
		return InsufficientNodes
	}
	if s.RequestedChain == nil || len(s.RequestedChain) == 0 {
		return NoReqChain
	}
	if len(s.BlockHash) != 32 {
		return InvalidBlockHashFormat
	}
	if len(s.DevID) != 32 {
		return InvalidDevIDFormat
	}
	return nil
}
