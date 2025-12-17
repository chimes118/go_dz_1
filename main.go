package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

const statsURL = "http://srv.msk01.gigacorp.local/_stats"

func main() {
	errorCount := 0

	resp, err := http.Get(statsURL)
	if err != nil {
		errorCount++
		printErrorIfNeeded(errorCount)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errorCount++
		printErrorIfNeeded(errorCount)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		errorCount++
		printErrorIfNeeded(errorCount)
		return
	}

	values := strings.Split(strings.TrimSpace(string(body)), ",")
	if len(values) != 7 {
		errorCount++
		printErrorIfNeeded(errorCount)
		return
	}

	num := make([]float64, 7)
	for i, v := range values {
		n, err := strconv.ParseFloat(v, 64)
		if err != nil {
			errorCount++
			printErrorIfNeeded(errorCount)
			return
		}
		num[i] = n
	}

	// 1. Load Average
	if num[0] > 30 {
		fmt.Printf("Load Average is too high: %.0f\n", num[0])
	}

	// 2. Memory usage
		memUsage := int(num[2] * 100 / num[1])
	if memUsage > 80 {
  		fmt.Printf("Memory usage too high: %d%%\n", memUsage)
}
	// 3. Disk space
		freeDisk := num[3] - num[4]
	freeMb := int(freeDisk / 1024 / 1024)
	if freeDisk*100/num[3] < 10 {
		fmt.Printf("Free disk space is too low: %d Mb left\n", freeMb)
}


	// 4. Network bandwidth
	freeNet := num[5] - num[6]
	freeMbit := int(freeNet*8/1024/1024)
	if freeNet*100/num[5] < 10 {
   		fmt.Printf("Network bandwidth usage high: %d Mbit/s available\n", freeMbit)
	}
}

func printErrorIfNeeded(count int) {
	if count >= 3 {
		fmt.Println("Unable to fetch server statistic")
	}
}