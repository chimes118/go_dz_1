package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

const url = "http://srv.msk01.gigacorp.local/_stats"

func main() {
	for {
		resp, err := http.Get(url)
		if err != nil {
			return
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return
		}

		reader := csv.NewReader(strings.NewReader(string(body)))
		record, err := reader.Read()
		if err != nil || len(record) != 7 {
			return
		}

		load, _ := strconv.Atoi(strings.Split(record[0], ".")[0])

		memTotal, _ := strconv.ParseInt(record[1], 10, 64)
		memUsed, _ := strconv.ParseInt(record[2], 10, 64)

		diskTotal, _ := strconv.ParseInt(record[3], 10, 64)
		diskUsed, _ := strconv.ParseInt(record[4], 10, 64)

		netTotal, _ := strconv.ParseInt(record[5], 10, 64)
		netUsed, _ := strconv.ParseInt(record[6], 10, 64)

		// Load Average
		if load > 30 {
			fmt.Printf("Load Average is too high: %d\n", load)
		}

		// Memory
		if memTotal > 0 {
			memPercent := memUsed * 100 / memTotal
			if memPercent > 80 {
				fmt.Printf("Memory usage too high: %d%%\n", memPercent)
			}
		}

		// Disk — ВАЖНО!
		if diskTotal > 0 {
			freeBytes := diskTotal - diskUsed
			freePercent := freeBytes * 100 / diskTotal

			if freePercent < 10 {
				freeMb := freeBytes / (1024 * 1024)
				fmt.Printf("Free disk space is too low: %d Mb left\n", freeMb)
			}
		}

		// Network
		if netTotal > 0 {
			if netUsed*100/netTotal > 90 {
				freeMbit := (netTotal - netUsed) / 1_000_000
				fmt.Printf("Network bandwidth usage high: %d Mbit/s available\n", freeMbit)
			}
		}
	}
}