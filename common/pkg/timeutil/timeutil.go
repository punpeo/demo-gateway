package timeutil

import (
	"math"
	"net/http"
	"time"
)

var (
	cst *time.Location
)

// CSTLayout China Standard Time Layout
const CSTLayout = "2006-01-02 15:04:05"
const CSTDateLayout = "2006-01-02"

func init() {
	var err error
	if cst, err = time.LoadLocation("Asia/Shanghai"); err != nil {
		panic(err)
	}

	// 默认设置为中国时区
	time.Local = cst
}

// RFC3339ToCSTLayout convert rfc3339 value to china standard time layout
// 2020-11-08T08:18:46+08:00 => 2020-11-08 08:18:46
func RFC3339ToCSTLayout(value string) (string, error) {
	ts, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return "", err
	}

	return ts.In(cst).Format(CSTLayout), nil
}

// CSTLayoutString 格式化时间
// 返回 "2006-01-02 15:04:05" 格式的时间
func CSTLayoutString() string {
	ts := time.Now()
	return ts.In(cst).Format(CSTLayout)
}

// ParseCSTInLocation 格式化时间
func ParseCSTInLocation(date string) (time.Time, error) {
	return time.ParseInLocation(CSTLayout, date, cst)
}

// CSTLayoutStringToUnix 返回 unix 时间戳
// 2020-01-24 21:11:11 => 1579871471
func CSTLayoutStringToUnix(cstLayoutString string) (int64, error) {
	stamp, err := time.ParseInLocation(CSTLayout, cstLayoutString, cst)
	if err != nil {
		return 0, err
	}
	return stamp.Unix(), nil
}

// CSTLayoutStringToUnix 返回 unix 时间戳
// 2020-01-24 => 1579871471
func CSTDateLayoutStringToUnix(cstDateLayoutString string) (int64, error) {
	stamp, err := time.ParseInLocation(CSTDateLayout, cstDateLayoutString, cst)
	if err != nil {
		return 0, err
	}
	return stamp.Unix(), nil
}

// GMTLayoutString 格式化时间
// 返回 "Mon, 02 Jan 2006 15:04:05 GMT" 格式的时间
func GMTLayoutString() string {
	return time.Now().In(cst).Format(http.TimeFormat)
}

// ParseGMTInLocation 格式化时间
func ParseGMTInLocation(date string) (time.Time, error) {
	return time.ParseInLocation(http.TimeFormat, date, cst)
}

// SubInLocation 计算时间差
func SubInLocation(ts time.Time) float64 {
	return math.Abs(time.Now().In(cst).Sub(ts).Seconds())
}

func GetTodayStr() string {
	todayTime := time.Now().Format("2006-01-02")

	return todayTime
}

func GetTodayTimeStr() string {
	todayTime := time.Now().Format("2006-01-02 00:00:00")

	return todayTime
}

func GetYesterdayStr() string {
	nTime := time.Now()
	yesTime := nTime.AddDate(0, 0, -1)

	return yesTime.Format("2006-01-02")
}

func GetBeforeDayStr(days int) string {
	nTime := time.Now()
	yesTime := nTime.AddDate(0, 0, days)

	return yesTime.Format("2006-01-02")
}

func GetMonthStr(month int) string {
	nTime := time.Now()
	yesTime := nTime.AddDate(0, month, 0)

	return yesTime.Format("2006-01-02")
}

/*
获取日期当天0点和23点59分
*/
func GetDateStartAndEndTime(date string) (int64, int64) {
	//日期当天0点时间戳(拼接字符串)
	startDate := date + " 00:00:00"
	startTime, _ := time.ParseInLocation("2006-01-02 15:04:05", startDate, cst)

	//日期当天23时59分时间戳
	endDate := date + " 23:59:59"
	end, _ := time.ParseInLocation("2006-01-02 15:04:05", endDate, cst)

	//返回当天0点和23点59分的时间戳
	return startTime.Unix(), end.Unix()
}

func UnixToString(unixTime int64, format string) string {
	timeStr := time.Unix(int64(unixTime), 0).Format(format)
	return timeStr
}

// 时间戳转日期
func UnixToCSTDateLayout(timeTemp int64) string {
	formatTimeStr := time.Unix(timeTemp, 0).Format(CSTDateLayout)

	return formatTimeStr
}

// 时间戳转日期时间
func UnixToCSTLayout(timeTemp int64) string {
	formatTimeStr := time.Unix(timeTemp, 0).Format(CSTLayout)

	return formatTimeStr
}

// CSTLayoutToUnix 日期转时间戳
func CSTLayoutToUnix(timeStr string) (int64, error) {
	t, err := time.ParseInLocation(CSTLayout, timeStr, cst)
	if err != nil {
		return 0, err
	}
	return t.Unix(), nil
}
