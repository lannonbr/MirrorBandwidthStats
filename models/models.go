package models

import (
	"fmt"

	humanize "github.com/dustin/go-humanize"
)

// BandwidthEntry a single entry given from dstat
type BandwidthEntry struct {
	Recv      uint64 // how many bytes mirror recieved
	Send      uint64 // how many bytes mirror sent
	Timestamp string // timestamp
}

// ToJSON return entry in a json format of { time, recv, send }
func (e *BandwidthEntry) ToJSON() string {
	return fmt.Sprintf("{\"time\": \"%s\", \"recv\": \"%s\", \"send\": \"%s\"}", e.Timestamp, humanize.Bytes(e.Recv), humanize.Bytes(e.Send))
}
