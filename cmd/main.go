package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
)

type stationData struct {
	min   float64
	max   float64
	sum   float64
	count int
}

type stationName [32]byte

func printStations(stations map[stationName]*stationData) {
	stationNames := make([]string, len(stations))
	i := 0
	for k := range stations {
		stationNames[i] = string(k[:])
		i++
	}
	sort.Strings(stationNames)
	fmt.Print("{")
	for _, k := range stationNames {
		v := stations[(stationName)([]byte(k))]
		fmt.Printf("%s=%.1f/%.1f/%.1f, ", k, v.min, v.sum/float64(v.count), v.max)
	}
	fmt.Print("}\n")
}

func parseFloat(str string) float64 {
	parsed, err := strconv.ParseFloat(str, 64)
	if err != nil {
		log.Fatal(err)
	}

	return parsed
}

func main() {
	file, err := os.Open("./data/measurements2.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	stations := make(map[stationName]*stationData)
	br := bufio.NewReader(file)
	for {
		bs, err := br.ReadBytes(10)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}

		parts := bytes.Split(bs, []byte{';'})
		parts[0] = append(parts[0], make([]byte, 32-len(parts[0]))...)
		name := (stationName)(parts[0])
		temperature := parseFloat(string(parts[1][:len(parts[1])-1]))

		if data, ok := stations[name]; !ok {
			stations[name] = &stationData{
				min:   temperature,
				max:   temperature,
				sum:   temperature,
				count: 1,
			}
		} else {
			if temperature < data.min {
				data.min = temperature
			}
			if temperature > data.max {
				data.min = temperature
			}

			data.sum += temperature
			data.count++
		}
	}

	printStations(stations)
}
