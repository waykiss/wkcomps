package date

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"golang.org/x/text/language"
)

type TimeZones string

const (
	TimeZoneSaoPaulo       = "America/Sao_Paulo"
	DateTimeFormatDatabase = "2006-01-02T15:04:05"
	DateFormatDatabase     = "2006-01-02"
	DateTimeBr             = "02/01/2006 15:04:05"
	DateBr                 = "02/01/2006"
	Days                   = time.Hour * 24
)

func TimeStampToDateTimeString(timestamp *timestamp.Timestamp) string {
	if timestamp == nil {
		return ""
	}
	value := TimeStampToTime(timestamp)
	return value.Format("2006-01-02 15:04:05")
}

func TimeToStringPtBR(timestamp time.Time) string {
	return timestamp.Format(DateTimeBr)
}

func TimeToString(timestamp time.Time) string {
	return timestamp.Format(DateTimeBr)
}

func TimeStampToDateString(timestamp *timestamp.Timestamp) string {
	value := TimeStampToTime(timestamp)
	return value.Format("2006-01-02")
}

// ToTime convert interface to Time
func ToTime(v interface{}) time.Time {

	// if can convert to int, it means that can be a unix time
	if value, err := strconv.ParseInt(fmt.Sprintf(v.(string)), 10, 64); err == nil {
		return IntToTime(value)
	}

	result, err := time.Parse(time.RFC3339, fmt.Sprint(v))
	if err != nil {
		return time.Time{}
	}
	return result
}

// IntToTime Convert time in int64 format to time.Time type
func IntToTime(v int64) time.Time {
	return time.Unix(int64(v)/1000, int64(v)%1000*1000000)
}

func TimeToTimeStamp(time time.Time) *timestamp.Timestamp {
	return &timestamp.Timestamp{Seconds: int64(time.Unix())}
}

func TimeStampToTime(timestamp *timestamp.Timestamp) time.Time {
	t, _ := ptypes.Timestamp(timestamp)
	return t
}

func GetDateFromTime(v time.Time, timezone TimeZones) time.Time {
	l, err := time.LoadLocation(string(timezone))
	if err != nil {
		panic(err)
	}
	return time.Date(v.Year(), v.Month(), v.Day(), v.Hour(), 0, 0, 0, l)
}

// GetEndOfTheDay essa funcao altera um horario para o final do dia, 23:59:59
func GetEndOfTheDay(v time.Time) time.Time {
	return time.Date(v.Year(), v.Month(), v.Day(), 23, 59, 59, 0, time.Local)
}

// GetStartOfTheDay essa funcao altera um horario para o inicio do dia, 00:00:00
func GetStartOfTheDay(v time.Time) time.Time {
	return time.Date(v.Year(), v.Month(), v.Day(), 00, 00, 00, 0, time.Local)
}

// NextDateByDay retorna a proxima data baseado no dia
func NextDateByDay(day, maxDay int) time.Time {
	now := time.Now()
	currentDay := now.Day()
	currentMonth := now.Month() + 1
	currentYear := now.Year()
	if currentDay > maxDay {
		currentMonth++
	}
	if currentMonth == 13 {
		currentYear++
		currentMonth = 1
	}
	return time.Date(currentYear, currentMonth, day, 23, 59, 59, 0, time.Local)
}

// GetTempo retorna uma string informado o tempo em portugues que falta a partir de um objeto Duration
func GetTempo(duration time.Duration) (r string) {
	diasString := ""
	duration = duration.Round(time.Minute)
	h := duration / time.Hour
	duration -= h * time.Hour
	m := duration / time.Minute
	if h > 23 {
		dias := 0
		for h > 23 {
			dias++
			h = h - 24
		}
		diasString = fmt.Sprintf("%d dias, ", dias)
	}
	return fmt.Sprintf("%s%02dh%02dm", diasString, int(math.Abs(float64(h))), int(math.Abs(float64(m))))
}

// IsWeekend retorna se uma data e final de semana ou nao
func IsWeekend(t time.Time) bool {
	t = t.UTC()
	switch t.Weekday() {
	case time.Saturday, time.Sunday:
		return true
	}
	return false
}

// IsWeekend retorna o próximo dia da semana
func GetNextDayOfWeek(t time.Time) time.Time {
	for IsWeekend(t) {
		t = t.AddDate(0, 0, 1)
	}
	return t
}

// DateEqual compare two dates if they are equals, ignoring the time and considering only the date
// return true if they are equals
func DateEqual(date1, date2 time.Time) bool {
	y1, m1, d1 := date1.Date()
	y2, m2, d2 := date2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

// FirstDayOfMonth return the time and set day 1
func FirstDayOfMonth(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
}

// LastDayOfMonth return the time with last day of the month
func LastDayOfMonth(t time.Time) time.Time {
	return FirstDayOfMonth(t).AddDate(0, 1, 0).Add(-time.Second)
}

func NextDay(t time.Time) time.Time {
	return t.AddDate(0, 0, 1)
}

// GetMonthName retorna o nome do mês baseado na lingua passada como parametro, padrao é ingles
func GetMonthName(date time.Time, lang language.Tag) string {
	switch lang {
	case language.BrazilianPortuguese:
		return shortMonthNamesPtBR[date.Format("Jan")]
	default:
		return date.Format("Jan")
	}
}

var shortMonthNamesPtBR = map[string]string{
	"Jan": "jan",
	"Feb": "fev",
	"Mar": "mar",
	"Apr": "abr",
	"May": "mai",
	"Jun": "jun",
	"Jul": "jul",
	"Aug": "ago",
	"Sep": "set",
	"Oct": "out",
	"Nov": "nov",
	"Dec": "dez",
}

// variaveis com a mapeamento do nome de dadtas em ingles para portugues
var longDayNamesPtBR = map[string]string{
	"Sunday":    "domingo",
	"Monday":    "segunda-feira",
	"Tuesday":   "terça-feira",
	"Wednesday": "quarta-feira",
	"Thursday":  "quinta-feira",
	"Friday":    "sexta-feira",
	"Saturday":  "sábado",
}

var shortDayNamesPtBR = map[string]string{
	"Sun": "dom",
	"Mon": "seg",
	"Tue": "ter",
	"Wed": "qua",
	"Thu": "qui",
	"Fri": "sex",
	"Sat": "sáb",
}

var longMonthNamesPtBR = map[string]string{
	"January":   "janeiro",
	"February":  "fevereiro",
	"March":     "março",
	"April":     "abril",
	"May":       "maio",
	"June":      "junho",
	"July":      "julho",
	"August":    "agosto",
	"September": "setembro",
	"October":   "outubro",
	"November":  "novembro",
	"December":  "dezembro",
}

// FromString function to return a time.Time given a string, it automatically tries to convert from several and
// common formats, if it has error, return time zeroValue
func FromString(date string) (r time.Time) {
	formats := []string{
		"2006-01-02",
		"2006-01-02T15:04:05Z07:00",
		"2006-01-02T15:04:05Z",
		"2006-01-02 15:04:05",
		"02/01/2006 15:04:05", //BR format
		"02/01/2006",
	}
	var err error
	for _, f := range formats {
		r, err = time.Parse(f, date)
		if err == nil {
			return
		}
	}
	return
}

// GetDaysBetweenDates return the diference of the of days between two dates object
func GetDaysBetweenDates(startDate, endDate time.Time) int {
	endDate = time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 0, 0, 0, 0, time.Local)
	startDate = time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, time.Local)
	hoursDiff := endDate.Sub(startDate).Hours()
	daysDiff := int(hoursDiff) / 24
	return daysDiff
}
