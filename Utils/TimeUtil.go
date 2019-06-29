package Utils

import (
	"github.com/Magezeng/go-framework/Converters"
	"github.com/pkg/errors"
	"strconv"
	"time"
)

/*
 Get a epoch time
  eg: 1521221737.376
*/
func EpochTimeNow() string {
	millisecond := time.Now().UnixNano() / 1000000
	epoch := strconv.Itoa(int(millisecond))
	epochBytes := []byte(epoch)
	epoch = string(epochBytes[:10]) + "." + string(epochBytes[10:])
	return epoch
}

/*
 Get a iso time
  eg: 2018-03-16T18:02:48.284Z
*/
func IsoTimeNow() string {
	return TimeToIso(time.Now())
}

func TimeToIso(t time.Time) string {
	utcTime := t.UTC()
	iso := utcTime.String()
	isoBytes := []byte(iso)
	iso = string(isoBytes[:10]) + "T" + string(isoBytes[11:23]) + "Z"
	return iso
}

/*
 Get utc +8 -- 1540365300000 -> 2018-10-24 15:15:00 +0800 CST
*/
func LongTimeToUTC8(longTime int64) time.Time {
	timeString := Converters.Int64ToString(longTime)
	sec := timeString[0:10]
	nsec := timeString[10:len(timeString)]
	return time.Unix(Converters.StringToInt64(sec), Converters.StringToInt64(nsec))
}

/*
 1540365300000 -> 2018-10-24 15:15:00
*/
func LongTimeToUTC8Format(longTime int64) string {
	return LongTimeToUTC8(longTime).Format("2006-01-02 15:04:05")
}

/*
  iso time change to time.Time
  eg: "2018-11-18T16:51:55.933Z" -> 2018-11-18 16:51:55.000000933 +0000 UTC
*/
func IsoToTime(iso string) (time.Time, error) {
	nilTime := time.Now()
	if iso == "" {
		return nilTime, errors.New("illegal parameter")
	}
	// "2018-03-18T06:51:05.933Z"
	isoBytes := []byte(iso)
	year, err := strconv.Atoi(string(isoBytes[0:4]))
	if err != nil {
		return nilTime, errors.New("illegal year")
	}
	month, err := strconv.Atoi(string(isoBytes[5:7]))
	if err != nil {
		return nilTime, errors.New("illegal month")
	}
	day, err := strconv.Atoi(string(isoBytes[8:10]))
	if err != nil {
		return nilTime, errors.New("illegal day")
	}
	hour, err := strconv.Atoi(string(isoBytes[11:13]))
	if err != nil {
		return nilTime, errors.New("illegal hour")
	}
	min, err := strconv.Atoi(string(isoBytes[14:16]))
	if err != nil {
		return nilTime, errors.New("illegal min")
	}
	sec, err := strconv.Atoi(string(isoBytes[17:19]))
	if err != nil {
		return nilTime, errors.New("illegal sec")
	}
	nsec, err := strconv.Atoi(string(isoBytes[20 : len(isoBytes)-1]))
	if err != nil {
		return nilTime, errors.New("illegal nsec")
	}
	return time.Date(year, time.Month(month), day, hour, min, sec, nsec, time.UTC), nil
}
