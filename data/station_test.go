package station

import (
	"errors"
	"testing"
)

func TestCheck(t *testing.T) {
	var err error

	b := check(err)

	if b != true {
		t.Errorf(`check() should return true when err is nil, instead returned %v`, b)
	}
}

func TestCheck_False(t *testing.T) {
	err := errors.New("New Error")

	b := check(err)

	if b != false {
		t.Errorf(`check() should return false when err is NOT nil, instead returned %v`, b)
	}
}

func TestSortStation(t *testing.T) {
	s1 := Station{}
	s2 := Station{}

	s1.stationId = 1
	s2.stationId = 2

	ret := sortStation(s1, s2)

	if ret != -1 {
		t.Errorf(`Station 1 should be less than Station 2 and return -1, intead returned %v`, ret)
	}

	ret = sortStation(s2, s1)

	if ret != 1 {
		t.Errorf(`Station 2 should be greater than Station 1 and return 1, instead returned %v`, ret)
	}
}

func TestParseuint(t *testing.T) {
	stringint := "1234"

	actualint := parseid(stringint)

	if actualint != 1234 {
		t.Errorf(`parseid("1234") should equal 1234, instead returns %v`, actualint)
	}

	actualint = parsereads(stringint)

	if actualint != 1234 {
		t.Errorf(`parsereads("1234") should equal 1234, instead returns %v`, actualint)
	}

	maxUint64String := "18446744073709551615"
	maxUint64StringWithCommas := "18,446,744,073,709,551,615"

	actualint = parsereads(maxUint64String)
	if actualint != 18446744073709551615 {
		t.Errorf(`parsereads() should parse a 64 bit uint correctly`)
	}

	//can parseuint handle commas
	actualint = parsereads(maxUint64StringWithCommas)
	if actualint != 18446744073709551615 {
		t.Errorf(`parsereads() should parse a 64 bit uint correctly`)
	}
}

func TestCalcUptimePercent(t *testing.T) {
	stat := Station{}
	stat.stationId = 1
	stat.uptime = append(stat.uptime, [2]uint{0, 10})
	stat.min = 0
	stat.max = 10

	ret := calcUptimePercent(stat)

	if ret != "1 100\n" {
		t.Errorf(`calcUptimePercent should return "1 100" instead returned %v`, ret)
	}
}

func TestMergeUptimes(t *testing.T) {
	stat := Station{}
	stat.stationId = 1
	stat.uptime = append(stat.uptime, [2]uint{0, 10})
	stat.uptime = append(stat.uptime, [2]uint{0, 30})
	stat.uptime = append(stat.uptime, [2]uint{0, 50})

	if len(stat.uptime) != 3 {
		t.Errorf(`stat %v uptimes should have a length of 3, but has a length of %v`, stat.stationId, len(stat.uptime))
	}

	mergeUptimes(&stat)

	if len(stat.uptime) != 1 {
		t.Errorf(`stat %v uptime should have merged to one`, stat.stationId)
	}

	if stat.uptime[0][0] != 0 || stat.uptime[0][1] != 50 {
		t.Errorf(`The reads should be %v and %v but are %v and %v`, 0, 50, stat.uptime[0][0], stat.uptime[0][1])
	}

}

func TestParseChargeLines(t *testing.T) {
	stat := Station{}

	parseChargeLines(0, 10, &stat, "true")

	if stat.max != 10 || len(stat.uptime) == 0 {
		t.Error("parseChargeLines() did not function properly")
	}
}
