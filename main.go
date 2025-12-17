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
			// сервер завершился — корректно выходим
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

		load, _ := strconv.Atoi(strings.Split(record[0], ".")[0])

		memTotal, _ := strconv.ParseInt(record[1], 10, 64)
		memUsed, _ := strconv.ParseInt(record[2], 10, 64)

		diskTotal, _ := strconv.ParseInt(record[3], 10, 64)
		diskUsed, _ := strconv.ParseInt(record[4], 10, 64)

		netTotal, _ := strconv.ParseInt(record[5], 10, 64)
		netUsed, _ := strconv.ParseInt(record[6], 10, 64)

		if load > 30 {
			fmt.Printf("Load Average is too high: %d\n", load)
		}

		if memTotal > 0 {
			memPercent := memUsed * 100 / memTotal
			if memPercent > 80 {
				fmt.Printf("Memory usage too high: %d%%\n", memPercent)
			}
		}

		if diskTotal > 0 {
			freeMb := (diskTotal - diskUsed) / (1024 * 1024)
			if freeMb*100/diskTotal < 10 {
				fmt.Printf("Free disk space is too low: %d Mb left\n", freeMb)
			}
		}

		if netTotal > 0 {
			freeMbit := (netTotal - netUsed) / 1_000_000
			if netUsed*100/netTotal > 90 {
				fmt.Printf("Network bandwidth usage high: %d Mbit/s available\n", freeMbit)
			}
		}
	}
}
