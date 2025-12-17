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
	resp, err := http.Get(statsURL)
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

	values := strings.Split(strings.TrimSpace(string(body)), ",")
	if len(values) != 7 {
		fmt.Println("Unable to fetch server statistic")
		return
	}

	numbers := make([]float64, 7)
	for i, v := range values {
		n, err := strconv.ParseFloat(v, 64)
		if err != nil {
			fmt.Println("Unable to fetch server statistic")
			return
		}
		numbers[i] = n
	}

	// дальше будет логика проверок — добавим следующим шагом
}