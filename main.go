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

func humanizeBits(bits uint64) string {
	str := humanize.Bytes(bits)
	strArr := strings.Split(str, " ")
	str = strArr[0] + string([]rune(strArr[1])[0]) + "b"
	return str
}

func analyzeFile(filename string, fullPrint bool) (uint64, uint64, uint64, uint64) {
	entries := loadBandwidthCSV(filename)

	var totalRecv, totalSend uint64

	for _, entry := range entries {
		totalRecv += entry.Recv
		totalSend += entry.Send
	}

	if fullPrint {
		fmt.Println("Filename:", filename)
	}

	humanRecv := humanize.Bytes(totalRecv)
	humanSend := humanize.Bytes(totalSend)

	totalOverall := totalRecv + totalSend
	humanOverall := humanize.Bytes(totalOverall)

	rate := ((totalRecv + totalSend) * 8) / 3550

	if fullPrint {
		fmt.Println("Total Received:", humanRecv)
		fmt.Println("Total Sent:", humanSend)
		fmt.Println("Total overall:", humanOverall)
		fmt.Println("Rate:", humanizeBits(rate)+"/sec")
		fmt.Println("------")
	}

	return totalRecv, totalSend, totalOverall, rate
}

func main() {

	files := os.Args[1:]

	var counter int

	var totalRecv, totalSend, totalOverall, totalRate, avgRate uint64

	for _, arg := range files {
		recv, send, overall, rate := analyzeFile(arg, false)
		counter++

		totalRecv += recv
		totalSend += send
		totalOverall += overall

		totalRate += rate
	}

	avgRate = totalRate / uint64(counter)

	fmt.Println("Recieved:", humanize.Bytes(totalRecv))
	fmt.Println("Sent:", humanize.Bytes(totalSend))
	fmt.Println("Overall:", humanize.Bytes(totalOverall))
	fmt.Println("Rate:", humanizeBits(avgRate)+"/sec")

}
