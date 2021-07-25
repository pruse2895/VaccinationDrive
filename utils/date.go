package utils

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

const (
	//DFddMMMyyyy denotes date format dd-MMM-yyyy
	DFddMMMyyyy = "02-Jan-2006"
	//DFyyyyMMdd denotes date format yyyy-MM-dd
	DFyyyyMMdd = "2006-01-02"
	//DFddMMyyyy denotes date format dd/MM/yyyy
	DFddMMyyyy = "02/01/2006"
	//DFMMMyyyy denotes date format MMM/yyyy
	DFMMMyyyy = "Jan/2006"
	//DFyyyyMM denotes date format yyyy-MM
	DFyyyyMM = "2006-01"
	//DFMMMyyyy denotes date format MMMyyyy
	DFMMMyyyyNoSep = "Jan2006"
)

//StringyyyyMMddToDate is func to parse string as date
func StringyyyyMMddToDate(dateStr string) (time.Time, error) {
	// log.Printf("string to date %+v", dateStr)
	// d, err := time.ParseInLocation(DFyyyyMMdd, dateStr, time.UTC)
	// log.Printf("string to date %v", d)
	// log.Printf("string to date err %v", err)
	return time.ParseInLocation(DFyyyyMMdd, dateStr, time.UTC)
}

//StringddMMyyyyToDate is func to parse string as date
func StringddMMyyyyToDate(dateStr string) (time.Time, error) {
	// log.Printf("string to date %+v", dateStr)
	// d, err := time.ParseInLocation(DFyyyyMMdd, dateStr, time.UTC)
	// log.Printf("string to date %v", d)
	// log.Printf("string to date err %v", err)
	return time.ParseInLocation(DFddMMyyyy, dateStr, time.UTC)
}

//DateToStringDFddMMMyyyy is func to format date as string
func DateToStringDFddMMMyyyy(date time.Time) string {
	return date.Format(DFddMMMyyyy)
}

//DateToStringDFyyyyMMdd is func to format date as string
func DateToStringDFyyyyMMdd(date time.Time) string {
	return date.Format(DFyyyyMMdd)
}

//DateToStringDFMMMyyyy is func to format date as string
func DateToStringDFMMMyyyy(date time.Time) string {
	return date.Format(DFMMMyyyy)
}

//FormatDateTime parse string to time format
func FormatDateTime(timeStampString string) (time.Time, error) {
	layOut := "2006-01-02 15:04"
	timeStamp, err := time.Parse(layOut, timeStampString)
	if err != nil {
		log.Println(err)
	}
	// hr, min, sec := timeStamp.Clock()
	// log.Println("hrs,mins,secs", hr, min, sec)
	return timeStamp, nil
}

// CurrentTimeWithZone will return the present time in zone
func CurrentTimeWithZone(zone string) (time.Time, error) {
	if zone == "" {
		zone = "Asia/Kolkata"
	}
	loc, err := time.LoadLocation(zone)
	if err != nil {
		loc, err = time.LoadLocation("Asia/Kolkata")
		if err != nil {
			return time.Now(), err
		}
	}
	currentTime := time.Now().In(loc)
	return currentTime, nil
}

func TimeByZoneInt(zone string, timeObj time.Time) time.Time {
	hr, _ := strconv.Atoi(fmt.Sprintf("%s%s", zone[:1], zone[1:3]))
	min, _ := strconv.Atoi(fmt.Sprintf("%s%s", zone[:1], zone[3:5]))
	timeObj = timeObj.Add(time.Hour * time.Duration(hr))
	timeObj = timeObj.Add(time.Minute * time.Duration(min))
	return timeObj
}

//DurationToString convert time.duration to string
func DurationToString(duration time.Duration) string {
	durationStr := duration.String()
	if strings.HasSuffix(durationStr, "0s") {
		durationStr = durationStr[:len(durationStr)-2]
	}
	return durationStr
}

//DurationToHourMins convert hours and mins to int64
func DurationToHourMins(duration time.Duration) (int32, int32) {
	totalMins := duration.Minutes()
	totalMinsInt := int32(totalMins)
	hrs := totalMinsInt / 60
	mins := totalMinsInt % 60
	return hrs, mins
}

//TimeWithZoneToDate convert date format to date in string
func TimeWithZoneToDate(date time.Time) string {
	return date.Format(DFyyyyMMdd)
}

//GetMonthYearFromMMMYYYY get month year
func GetMonthYearFromMMMYYYY(timeStampString string) (time.Time, error) {
	layOut := "Jan-2006"
	timeStamp, err := time.Parse(layOut, timeStampString)
	if err != nil {
		log.Println(err)
	}
	return timeStamp, nil
}

// MonthInterval for get month start and end date
func MonthInterval(y int, m time.Month) (beginningOfMonth, enOfMonth time.Time) {
	beginningOfMonth = time.Date(y, m, 1, 0, 0, 0, 0, time.UTC)
	enOfMonth = time.Date(y, m+1, 1, 0, 0, 0, -1, time.UTC)
	return beginningOfMonth, enOfMonth
}

// DurationIntToHHMMStr function is use to change duration total minutes to hours and minutes string
func DurationIntToHHMMStr(duration int) string {

	hrs := duration / 60
	mins := duration % 60

	str := ""
	if duration > 0 {
		str = str + fmt.Sprintf("%dh", hrs)
	}

	if duration > 0 {
		str = str + fmt.Sprintf("%dm", mins)
	}
	return str
}

//CheckDateFormatyyyyMMdd function is use to check given date is valid format or not
func CheckDateFormatyyyyMMdd(dateStr string) bool {
	_, err := time.ParseInLocation(DFyyyyMMdd, dateStr, time.UTC)
	if err != nil {
		return false
	}
	return true
}

//StringyyyyMMddToddMMyyyy is func to parse string as date
func StringyyyyMMddToddMMyyyy(dateStr string) string {
	d, _ := time.ParseInLocation(DFyyyyMMdd, dateStr, time.UTC)
	return d.Format(DFddMMMyyyy)
}

//StringyyyyMMddToddMMyyyy is func to parse string as date
func ConvertFormatToddMMyyyyWithSlash(dateStr string) string {
	d, _ := time.ParseInLocation(DFyyyyMMdd, dateStr, time.UTC)
	return d.Format(DFddMMyyyy)
}

//StringyyyyMMddToddMMyyyy is func to parse string as date
func ConvertyyyyMMToMMMyyyy(dateStr string) string {
	d, _ := time.ParseInLocation(DFyyyyMM, dateStr, time.UTC)
	return d.Format(DFMMMyyyyNoSep)
}

//ConvertDateToStringDFyyyymmdd is func to format date as string
func ConvertDateToStringDFyyyyMMdd(date time.Time) string {
	return date.Format(DFyyyyMMdd)
}

func ConvertDateToStringWithDate(date time.Time) string {
	layOut := "2006-01-02 15:04"
	return date.Format(layOut)
}

func IsNumeric(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}
