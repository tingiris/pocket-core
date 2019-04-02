package common

import (
	"github.com/pokt-network/pocket-core/common/fbs"
	"strconv"
)

func UnmarshalBlockchain(flatBuffer []byte) Blockchain {
	res := fbs.GetRootAsBlockchain(flatBuffer, 0)
	return Blockchain{string(res.NameBytes()), strconv.Itoa(int(res.Netid())), strconv.Itoa(int(res.Version()))}
}
