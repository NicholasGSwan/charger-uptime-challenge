package main

import (
	"os"

	station "github.com/NicholasGSwan/charger-uptime-challenge/data"
)

func main() {
	var filename string
	if len(os.Args) > 0 {
		//it seems the command itself is the first 'arg'
		filename = os.Args[1]
	}
	//fmt.Println(filename)
	station.UptimeReport(filename)

}
