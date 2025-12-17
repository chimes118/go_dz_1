package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

const url = "http://srv.msk01.gigacorp.local:8080/stat"

func main() {
	for {
		resp, err := http.Get(url)
		if err != nil {
			// сервер выключился — НОРМАЛЬНО выходим
			return
		}

		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
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

		cpu, _ := strconv.Atoi(record[0])
		load, _ := strconv.Atoi(strings.Split(record[1], ".")[0])
		memUsed, _ := strconv.ParseInt(record[2], 10, 64)
		memTotal, _ := strconv.ParseInt(record[3], 10, 64)
		diskUsed, _ := strconv.ParseInt(record[4], 10, 64)
		diskTotal, _ := strconv.ParseInt(record[5], 10, 64)
		net, _ := strconv.ParseInt(record[6], 10, 64)

		if load > 50 {
			fmt.Printf("Load Average is too high: %d\n", load)
		}

		if memTotal > 0 {
			memPercent := memUsed * 100 / memTotal
			if memPercent > 80 {
				fmt.Printf("Memory usage too high: %d%%\n", memPercent)
			}
		}

		if net > 0 {
			fmt.Printf("Network bandwidth usage high: %d Mbit/s available\n", net/1_000_000)
		}

		if diskTotal > 0 {
			free := (diskTotal - diskUsed) / (1024 * 1024)
			if free < 50000 {
				fmt.Printf("Free disk space is too low: %d Mb left\n", free)
			}
		}
	}
}
