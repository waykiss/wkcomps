package date

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGetTempo(t *testing.T) {
	//teste com mais de 1 dia de diferença
	start := time.Date(2022, '2', 17, 11, 0, 23, 0, time.UTC)
	end := time.Date(2022, '2', 21, 12, 30, 23, 0, time.UTC)
	retorno := GetTempo(end.Sub(start))
	assert.Equal(t, "4 dias, 01h30m", retorno, "Retorno do tempo decorrido está diferente do esperado")

	//teste com apenas horas e minutos de diferença
	start = time.Date(2022, '2', 17, 11, 0, 23, 0, time.UTC)
	end = time.Date(2022, '2', 17, 12, 35, 23, 0, time.UTC)
	retorno = GetTempo(end.Sub(start))
	assert.Equal(t, "01h35m", retorno, "Retorno do tempo decorrido está diferente do esperado")

	//teste com apenas minutos de diferença
	start = time.Date(2022, '2', 17, 11, 0, 23, 0, time.UTC)
	end = time.Date(2022, '2', 17, 11, 35, 23, 0, time.UTC)
	retorno = GetTempo(end.Sub(start))
	assert.Equal(t, "00h35m", retorno, "Retorno do tempo decorrido está diferente do esperado")
}

func TestGetDaysBetweenDates(t *testing.T) {
	baseDate := time.Date(2022, 7, 27, 0, 0, 0, 0, time.Local)
	testCases := []struct {
		startDate, endDate time.Time
		days               int
	}{
		{baseDate, baseDate, 0},
		{baseDate, baseDate, 0},
		{baseDate, baseDate.AddDate(0, 0, -5), -5},
		{baseDate, baseDate.AddDate(0, 0, 3), 3},
		{baseDate, baseDate.AddDate(1, 0, 0), 365},
		{baseDate.AddDate(0, 0, 3), baseDate.AddDate(0, 0, 0), -3},
		{baseDate.AddDate(1, 0, 0), baseDate.AddDate(0, 0, 0), -365},
	}
	for _, v := range testCases {
		got := GetDaysBetweenDates(v.startDate, v.endDate)
		assert.Equal(t, v.days, got)
	}
}

func TestDateEqual(t *testing.T) {
	baseDate := time.Date(2022, 7, 27, 0, 0, 0, 0, time.Local)
	testCases := []struct {
		date1, date2 time.Time
		expected     bool
	}{
		{baseDate, baseDate, true},
		{baseDate, baseDate.Add(time.Minute * 60), true},
		{baseDate, baseDate.Add(time.Hour * 24).Add(time.Second * -1), true},
		{baseDate, baseDate.Add(time.Hour * 24), false},
		{baseDate, baseDate.Add(time.Hour * -24), false},
	}
	for _, v := range testCases {
		got := DateEqual(v.date1, v.date2)
		assert.Equal(t, v.expected, got)
	}
}
