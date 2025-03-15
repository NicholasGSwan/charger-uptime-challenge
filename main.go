package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	station "github.com/NicholasGSwan/charger-uptime-challenge/data"
)

const stationsHeader = "[Stations]"
const availReportHeader = "[Charger Availability Reports]"

func check(err error) bool {
	if err == nil {
		return true
	}
	return false
}

func parse32uint(s string) uint32 {
	num, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		panic(err)
	}
	return uint32(num)
}

func main() {
	fmt.Println("Hello, world")
	fs, err := os.Open("/Users/uglygrayduck/Dev/charger-uptime/charger-uptime-challenge/the task/input_1.txt")
	if err == nil {
		fmt.Println("File opened")
	}
	reader := bufio.NewReader(fs)

	for line, err := reader.ReadString('\n'); check(err); {
		stats := []station.Station{}
		chargeMap := make(map[uint32]uint32)
		line := strings.TrimSpace(line)
		for line != availReportHeader && check(err) {

			if line != stationsHeader && len(line) != 0 {
				fmt.Println(line)
				fields := strings.Fields(line)
				stat := station.Station{}
				var statId uint32
				for i, field := range fields {
					if i == 0 {
						stat.StationId = parse32uint(field)
						statId = parse32uint(field)

					} else {
						chargeMap[parse32uint(field)] = statId
						stat.ChargerIds = append(stat.ChargerIds, parse32uint(field))
					}

					///fmt.Println(field)
				}
				stats = append(stats, stat)

			}

			line, err = reader.ReadString('\n')
			line = strings.TrimSpace(line)
		}

		for _, stat := range stats {
			fmt.Printf("Station Id: %v, Chargers in Station: %v \n", stat.StationId, stat.ChargerIds)
		}

		for check(err) {
			if strings.TrimSpace(line) == availReportHeader {
				fmt.Println("printing charger reads: ")
			} else {
				fmt.Println(line)
			}
			line, err = reader.ReadString('\n')
			line = strings.TrimSpace(line)
		}

		//line, err = reader.ReadString('\n')
	}

}
