package exprimentalPrograms_test

import (
	"regexp"
	"strings"
	"testing"
)

const regexString = "[0-9]+[m|h|s|d]"

// ConvertTimeInterval to the desired prometheus format
func ConvertTimeInterval(monitorTimeInt string) (string, error) {
	timeInterval := strings.ToLower(strings.Join(strings.Fields(monitorTimeInt), ""))
	reg, err := regexp.Compile(regexString)

	if !reg.MatchString(timeInterval) {
		return "", err
	}
	timePeriod := reg.FindString(timeInterval)

	return timePeriod, nil
}

func BenchmarkConvertTime(b *testing.B) {
	time := "query=over_time=9h30m"
	for n := 0; n < b.N; n++ {
		if res, err := ConvertTimeInterval(time); err != nil {
			b.Fail()
		} else {
			b.Log(res)
		}
	}
}
