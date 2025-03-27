package station

import (
	"bufio"
	"cmp"
	"fmt"
	"math"
	"os"
	"slices"
	"strconv"
	"strings"
)

// honestly could just be a map of station ids and uptimes, but putting it in a type helps me organize it in my head a little better
// added the min and max and minset after above comment.  Probably a better way to do this, but it works.
type Station struct {
	stationId  uint
	chargerIds []uint
	uptime     [][2]uint
	min        uint
	max        uint
	minSet     bool
}

const (
	stationsHeader    = "[Stations]"
	availReportHeader = "[Charger Availability Reports]"
	ERROR             = "ERROR"
	EOF               = "EOF"
)

func UptimeReport(filename string) {

	fs, err := os.Open(filename)
	var reader *bufio.Reader
	if check(err) {
		reader = bufio.NewReader(fs)
	} else {
		fmt.Println(ERROR)
		os.Exit(0)
	}
	// start read of file line by line
	for line, err := reader.ReadString('\n'); check(err); {
		//map of Stations, keys are the station id
		stats := make(map[uint]Station)
		// charge map; the key is the charger id and the value is the station id it belongs to
		chargeMap := make(map[uint]uint)
		line := strings.TrimSpace(line)
		for line != availReportHeader && check(err) {

			if line != stationsHeader && len(line) != 0 {
				// fmt.Println(line)
				fields := strings.Fields(line)
				stat := Station{}
				var statId uint
				for i, field := range fields {
					if i == 0 {
						statId = parseid(field)
						stat.stationId = statId
					} else {
						chargeMap[parseid(field)] = statId
						stat.chargerIds = append(stat.chargerIds, parseid(field))
					}

					///fmt.Println(field)
				}
				stats[stat.stationId] = stat

			}

			line, err = reader.ReadString('\n')
			line = strings.TrimSpace(line)
		}

		// for _, stat := range stats {
		// 	fmt.Printf("Station Id: %v, Chargers in Station: %v \n", stat.stationId, stat.chargerIds)
		// }

		for check(err) {
			if strings.TrimSpace(line) == availReportHeader {
				//fmt.Println("printing charger reads: ")
			} else {
				fields := strings.Fields(line)
				if len(fields) == 4 {
					stat := stats[chargeMap[parseid(fields[0])]]
					//will modifying this stat modify the one in the map?
					//it does not
					parseChargeLines(parsereads(fields[1]), parsereads(fields[2]), &stat, fields[3])
					stats[stat.stationId] = stat

				} else if len(fields) > 0 {
					fmt.Println(ERROR)
					//panic(ERROR)
					os.Exit(0)
				}

				//fmt.Println(line)
			}
			line, err = reader.ReadString('\n')
			line = strings.TrimSpace(line)
		}
		//fmt.Println(err)

		if err.Error() == EOF {
			fields := strings.Fields(line)
			if len(fields) == 4 {
				stat := stats[chargeMap[parseid(fields[0])]]
				//will modifying this stat modify the one in the map?
				//it does not
				parseChargeLines(parsereads(fields[1]), parsereads(fields[2]), &stat, fields[3])
				stats[stat.stationId] = stat
			} else if len(fields) > 0 {
				fmt.Println(ERROR)
				//panic(ERROR)
				os.Exit(0)
			}
		} else {
			fmt.Println(ERROR)
			//panic(ERROR)
			os.Exit(0)
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
			//fmt.Printf("The uptimes for station : %v are: %v   length of uptime slice is: %v\n", stat.stationId, stat.uptime, len(stat.uptime))
			fmt.Print(calcUptimePercent(stat))
		}

		//line, err = reader.ReadString('\n')
	}
}

func parseChargeLines(start uint, end uint, stat *Station, isUptime string) {
	//will modifying this stat modify the one in the map?
	//it does not
	if isUptime == "true" {
		stat.uptime = append(stat.uptime, [2]uint{start, end})

	}
	if !stat.minSet {
		stat.min = start
		stat.minSet = true
	} else if stat.min > start {
		stat.min = start
	}
	if stat.max < end {
		stat.max = end
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
	stat.uptime = uptime

}

func calcUptimePercent(stat Station) string {
	if len(stat.uptime) == 0 {
		fmt.Printf("%v %v\n", stat.stationId, 0)
	}
	var ret string
	if len(stat.uptime) > 0 {
		uptime := stat.uptime
		min := stat.min
		max := stat.max
		// fmt.Printf("the min time is: %v and the max time is: %v\n", min, max)
		//setting downtime to the entire time period
		downtime := max - min
		for _, times := range uptime {
			// fmt.Printf("downtime: %v\n ", downtime)
			//subtracting periods of uptime from downtime
			downtime = downtime - (times[1] - times[0])
		}
		// fmt.Printf("downtime: %v\n ", downtime)
		total := float64(max - min)
		percent := 100 - 100*(float64(downtime)/total)

		//fmt.Printf("min: %v max: %v total: %v  downtime: %v  downtime percent: %v\n", min, max, total, downtime, float64(downtime)/total)
		ret = fmt.Sprintf("%v %v\n", stat.stationId, math.Trunc(percent))

	}
	return ret
}

func check(err error) bool {
	return err == nil
}

func parseuint(s string, bitsize int) uint {
	s = strings.ReplaceAll(s, ",", "")
	num, err := strconv.ParseUint(s, 10, bitsize)
	if err != nil {
		fmt.Println(ERROR)
		//panic(err)
		os.Exit(0)
	}
	return uint(num)
}

func parseid(s string) uint {
	return parseuint(s, 32)
}
func parsereads(s string) uint {
	return parseuint(s, 64)
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
