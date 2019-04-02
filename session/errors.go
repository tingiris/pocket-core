package session

import "errors"

const (
	devid                   = "devid"
	blockhash               = "blockhash"
	nodelist                = "nodelist"
	requestedchain          = "requested chain"
	IsNotInSeed             = " is empty or nil "
	IsInvalidFormat         = " is not in the correct format"
	insufficientNodeString  = "not enough nodes to fulfill the session"
	incompleteSessionString = "invalid session, missing information needed for key generation"
)

var (
	NoDevID                = errors.New(devid + IsNotInSeed)
	NoBlockHash            = errors.New(blockhash + IsNotInSeed)
	NoNodeList             = errors.New(nodelist + IsNotInSeed)
	NoReqChain             = errors.New(requestedchain + IsNotInSeed)
	InvalidBlockHashFormat = errors.New(blockhash + IsInvalidFormat)
	InvalidDevIDFormat     = errors.New(devid + IsInvalidFormat)
	InsufficientNodes      = errors.New(insufficientNodeString)
	IncompleteSession      = errors.New(incompleteSessionString)
)
