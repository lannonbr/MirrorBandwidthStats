package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

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

		entries = append(entries, models.BandwidthEntry{
			Recv:      line[0],
			Send:      line[1],
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

	for i := range entries {
		entries[i].Recv = cleanupBytes(entries[i].Recv)
		entries[i].Send = cleanupBytes(entries[i].Send)
	}

	for _, entry := range entries {
		fmt.Println(entry.ToJSON())
	}

}
