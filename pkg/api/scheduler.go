package sheduler

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

var (
	ErrInvalidFormat      = errors.New("invalid repeat format")
	ErrUnsupportedFormat  = errors.New("unsupported repeat format")
	ErrInvalidDate        = errors.New("invalid date")
	ErrInvalidDayInterval = errors.New("invalid day interval")
	ErrInvalidDayOfMonth  = errors.New("invalid day of month")
	ErrInvalidMonth       = errors.New("invalid month")
)

const maxDayInterval = 400

func isLeap(year int) bool {
	return year%400 == 0 || (year%100 != 0 && year%4 == 0)
}

func afterNow(date, now time.Time) bool {
	return date.Year() > now.Year() ||
		(date.Year() == now.Year() && date.Month() > now.Month()) ||
		(date.Year() == now.Year() && date.Month() == now.Month() && date.Day() > now.Day())
}

func NewDate(now time.Time, dateStr string, repeat string) (string, error) {
	if repeat == "" {
		return "", ErrInvalidFormat
	}

	date, err := time.Parse("20060102", dateStr)
	if err != nil {
		return "", ErrInvalidDate
	}

	sliceRepeat := strings.Fields(repeat)
	if len(sliceRepeat) < 1 {
		return "", ErrInvalidFormat
	}

	switch sliceRepeat[0] {
	case "d":
		if len(sliceRepeat) != 2 {
			return "", ErrInvalidFormat
		}
		return handleDayRule(now, date, sliceRepeat[1])
	case "y":
		return handleYearRule(now, date)
	case "m":
		return "", ErrUnsupportedFormat //tmp
	case "w":
		return "", ErrUnsupportedFormat //tmp
	default:
		return "", ErrInvalidFormat
	}
}

func handleDayRule(now, date time.Time, interval string) (string, error) {
	days, err := strconv.Atoi(interval)
	if err != nil {
		return "", ErrInvalidDayInterval
	}

	if days < 1 || days > maxDayInterval {
		return "", ErrInvalidDayInterval
	}

	result := date
	for {
		result = result.AddDate(0, 0, days)
		if afterNow(date, now) {
			break
		}
	}

	return result.Format("20060102"), nil
}

func handleYearRule(now, date time.Time) (string, error) {
	result := date
	for {
		result = result.AddDate(1, 0, 0)
		if afterNow(result, now) {
			break
		}
	}
	if date.Month() == 2 && date.Day() == 29 {
		if !isLeap(result.Year()) {
			result = time.Date(result.Year(), 3, 1, 0, 0, 0, 0, time.UTC)
		}
	}

	return result.Format("20060102"), nil
}
