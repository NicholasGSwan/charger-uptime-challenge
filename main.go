package main

import (
	"fmt"
	"os"
)

func main() {
	var filename string
	if len(os.Args) > 0 {
		filename = os.Args[0]
	}
	fmt.Println(filename)
	//station.UptimeReport(filename)
	//station.UptimeReport("/Users/uglygrayduck/Dev/charger-uptime/charger-uptime-challenge/the task/input_1.txt")
	//station.UptimeReport("/Users/uglygrayduck/Dev/charger-uptime/charger-uptime-challenge/the task/input_2.txt")

}
