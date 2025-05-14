package api

import (
	"fmt"
	"time"

	"github.com/vborisov12/go_final/pkg/storage"
)

// isLeap проверяет, является ли год високосным
func isLeap(year int) bool {
	return year%400 == 0 || (year%100 != 0 && year%4 == 0)
}

// afterNow проверяет, находится ли дата после now
func afterNow(date, now time.Time) bool {
	return date.Year() > now.Year() ||
		(date.Year() == now.Year() && date.Month() > now.Month()) ||
		(date.Year() == now.Year() && date.Month() == now.Month() && date.Day() > now.Day())
}

// lastDayOfMonth возвращает последний день месяца
//
// Параметры:
//
//	t time.Time - дата
func lastDayOfMonth(t time.Time) int {
	return time.Date(t.Year(), t.Month()+1, 0, 0, 0, 0, 0, time.UTC).Day()
}

// checkDate проверяет корректность даты
func checkDate(task *storage.Task) error {
	now := time.Now()

	if task.Date == "" {
		task.Date = now.Format("20060102")
		return nil
	}

	t, err := time.Parse("20060102", task.Date)
	if err != nil {
		return fmt.Errorf("Invalid date format: %w", err)
	}

	if t.Format("20060102") == now.Format("20060102") {
		return nil
	}

	if task.Repeat != "" {
		next, err := NewDate(now, task.Date, task.Repeat)
		if err != nil {
			return err
		}
		if !afterNow(t, now) {
			task.Date = next
			return nil
		}
	} else if !afterNow(t, now) {
		task.Date = now.Format("20060102")
		return nil
	}

	return nil
}
