package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const (
	url = "http://srv.msk01.gigacorp.local:8080/stat"
)

func main() {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Unable to fetch server statistic")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Unable to fetch server statistic")
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Unable to fetch server statistic")
		return
	}

	reader := csv.NewReader(strings.NewReader(string(body)))
	reader.Comma = ','

	record, err := reader.Read()
	if err != nil || len(record) != 7 {
		fmt.Println("Unable to fetch server statistic")
		return
	}

	// поля по заданию
	cpu, _ := strconv.Atoi(record[0])
	load, _ := strconv.ParseFloat(record[1], 64)
	memUsed, _ := strconv.ParseInt(record[2], 10, 64)
	memTotal, _ := strconv.ParseInt(record[3], 10, 64)
	diskUsed, _ := strconv.ParseInt(record[4], 10, 64)
	diskTotal, _ := strconv.ParseInt(record[5], 10, 64)
	net, _ := strconv.ParseInt(record[6], 10, 64)

	// CPU
	if cpu > 80 {
		fmt.Printf("CPU usage too high: %d%%\n", cpu)
	}

	// Load Average
	if load > 1 {
		fmt.Printf("Load average too high: %.2f\n", load)
	}

	// Memory (целочисленно!)
	if memTotal > 0 {
		memPercent := memUsed * 100 / memTotal
		if memPercent > 80 {
			fmt.Printf("Memory usage too high: %d%%\n", memPercent)
		}
	}

	// Disk
	if diskTotal > 0 {
		diskPercent := diskUsed * 100 / diskTotal
		if diskPercent > 90 {
			fmt.Printf("Disk usage too high: %d%%\n", diskPercent)
		}
	}

	// Network (без умножения на 8, как писали в чате)
	if net > 100000000 {
		fmt.Printf("Network traffic too high: %d\n", net)
	}

	// программа ДОХОДИТ ДО КОНЦА и ЗАВЕРШАЕТСЯ
	os.Exit(0)
}