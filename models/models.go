package models

import (
	"fmt"
	"math/big"

	humanize "github.com/dustin/go-humanize"
)

// BandwidthEntry a single entry given from dstat
type BandwidthEntry struct {
	Recv      string // how many bytes mirror recieved
	Send      string // how many bytes mirror sent
	Timestamp string // timestamp
}

// RecvBytes Return a condensed version of the bytes
// ex: if Recv = 9546144, e.RecvBytes() returns "9.5 MB"
func (e *BandwidthEntry) RecvBytes() string {
	bigRepr, _ := big.NewInt(0).SetString(e.Recv, 10)
	return humanize.BigBytes(bigRepr)
}

// SendBytes Return a condensed version of the bytes
// ex: if Send = 9546144, e.SendBytes() returns "9.5 MB"
func (e *BandwidthEntry) SendBytes() string {
	bigRepr, _ := big.NewInt(0).SetString(e.Send, 10)
	return humanize.BigBytes(bigRepr)
}

// ToJSON return entry in a json format of { time, recv, send }
func (e *BandwidthEntry) ToJSON() string {
	return fmt.Sprintf("{\"time\": \"%s\", \"recv\": \"%s\", \"send\": \"%s\"}", e.Timestamp, e.RecvBytes(), e.SendBytes())
}
