package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	humanize "github.com/dustin/go-humanize"
	"github.com/lannonbr/MirrorBandwidthStats/models"
)

func cleanupBytes(str string) string {
	return strings.TrimSuffix(str, ".0")
}

func loadBandwidthCSV(filename string) []models.BandwidthEntry {
	csvfile, _ := os.Open(filename)
	reader := csv.NewReader(bufio.NewReader(csvfile))
	reader.Comment = '"'

	entries := []models.BandwidthEntry{}

	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalln("Error reading csv line", err)
		}

		recv, _ := strconv.ParseUint(cleanupBytes(line[0]), 10, 64)
		send, _ := strconv.ParseUint(cleanupBytes(line[1]), 10, 64)

		entries = append(entries, models.BandwidthEntry{
			Recv:      recv,
			Send:      send,
			Timestamp: line[2],
		})
	}

	return entries
}

func main() {

	filename := os.Args[1]

	fmt.Println("Preparing to read CSV File:", filename)

	entries := loadBandwidthCSV(filename)

	fmt.Printf("Number of entries: %d\n", len(entries))

	for _, entry := range entries {
		fmt.Println(entry.ToJSON())
	}

	var totalRecv, totalSend uint64

	for _, entry := range entries {
		totalRecv += entry.Recv
		totalSend += entry.Send
	}

	fmt.Println("Total Received:", humanize.Bytes(totalRecv))
	fmt.Println("Total Sent:", humanize.Bytes(totalSend))
}
