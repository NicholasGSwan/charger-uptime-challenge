package station

import (
	"bufio"
	"cmp"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

// honestly could just be a map of station ids and uptimes, but putting it in a type helps me organize it in my head a little better
type Station struct {
	stationId  uint
	chargerIds []uint
	uptime     [][2]uint
}

const (
	stationsHeader    = "[Stations]"
	availReportHeader = "[Charger Availability Reports]"
)

func check(err error) bool {
	return err == nil
}

func parseuint(s string) uint {
	num, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		panic(err)
	}
	return uint(num)
}

func sortCharge(a, b [2]uint) int {
	if n := cmp.Compare(a[0], b[0]); n != 0 {
		return n
	}

	return cmp.Compare(a[1], b[1])
}

func sortStation(a, b Station) int {
	return cmp.Compare(a.stationId, b.stationId)
}
func UptimeReport(filename string) {
	fmt.Println("Hello, world")
	fs, err := os.Open(filename)
	if err == nil {
		fmt.Println("File opened")
	}
	reader := bufio.NewReader(fs)
	// start read of file line by line
	for line, err := reader.ReadString('\n'); check(err); {
		//map of Stations, keys are the station id
		stats := make(map[uint]Station)
		// charge map; the key is the charger id and the value is the station id it belongs to
		chargeMap := make(map[uint]uint)
		line := strings.TrimSpace(line)
		for line != availReportHeader && check(err) {

			if line != stationsHeader && len(line) != 0 {
				fmt.Println(line)
				fields := strings.Fields(line)
				stat := Station{}
				var statId uint
				for i, field := range fields {
					if i == 0 {
						stat.stationId = parseuint(field)
						statId = parseuint(field)

					} else {
						chargeMap[parseuint(field)] = statId
						stat.chargerIds = append(stat.chargerIds, parseuint(field))
					}

					///fmt.Println(field)
				}
				stats[stat.stationId] = stat

			}

			line, err = reader.ReadString('\n')
			line = strings.TrimSpace(line)
		}

		for _, stat := range stats {
			fmt.Printf("Station Id: %v, Chargers in Station: %v \n", stat.stationId, stat.chargerIds)
		}

		for check(err) {
			if strings.TrimSpace(line) == availReportHeader {
				fmt.Println("printing charger reads: ")
			} else {
				fields := strings.Fields(line)
				if len(fields) == 3 || (len(fields) == 4 && fields[3] == "true") {
					//will modifying this stat modify the one in the map?
					//it does not
					stat := stats[chargeMap[parseuint(fields[0])]]
					stat.uptime = append(stat.uptime, [2]uint{parseuint(fields[1]), parseuint(fields[2])})
					stats[stat.stationId] = stat
				}

				fmt.Println(line)
			}
			line, err = reader.ReadString('\n')
			line = strings.TrimSpace(line)
		}
		fmt.Println(err)

		if err.Error() == "EOF" {
			fields := strings.Fields(line)
			if len(fields) == 3 || (len(fields) == 4 && fields[3] == "true") {
				//will modifying this stat modify the one in the map?
				//it does not
				stat := stats[chargeMap[parseuint(fields[0])]]
				stat.uptime = append(stat.uptime, [2]uint{parseuint(fields[1]), parseuint(fields[2])})
				stats[stat.stationId] = stat
			}
		}
		// create station slice to be sorted
		statsSlice := []Station{}
		// need to sort stations before output
		for _, stat := range stats {
			slices.SortFunc(stat.uptime, sortCharge)
			mergeUptimes(&stat)
			statsSlice = append(statsSlice, stat)

		}
		//sort stations
		slices.SortFunc(statsSlice, sortStation)
		for _, stat := range statsSlice {
			fmt.Printf("The uptimes for station : %v are: %v   length of uptime slice is: %v\n", stat.stationId, stat.uptime, len(stat.uptime))
		}

		//line, err = reader.ReadString('\n')
	}
}

func mergeUptimes(stat *Station) {
	L, R := 0, 1
	uptime := stat.uptime
	if len(uptime) > 1 {
		for len(uptime) > 1 && R < len(uptime) {
			if uptime[L][0] == uptime[R][0] {
				//this assumes I have sorted all the elements smallest to largest, including the upper bound of the uptime
				uptime = slices.Delete(uptime, L, L+1)
				// check if there is overlap
			} else if uptime[R][0] <= uptime[L][1] {
				//make sure L element contains all uptime from L and R
				if uptime[L][1] < uptime[R][1] {
					uptime[L][1] = uptime[R][1]
				}
				//deletes R element since it is redundant
				uptime = slices.Delete(uptime, R, R+1)
			} else {
				//this should only happen if uptimes do not overlap
				L++
				R++
			}
		}
	}
}
