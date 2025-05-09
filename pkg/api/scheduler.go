package api

import (
	"errors"
	"fmt"
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
	ErrInvalidDayOfWeek   = errors.New("invalid day of week")
	ErrInvalidMonth       = errors.New("invalid month")

	validWeekDays  = make(map[int]bool)
	validMonths    = make(map[int]bool)
	validMonthDays = make(map[int]bool)
)

const (
	maxDayInterval = 400
	maxDayOfMonth  = 31
	maxMonth       = 12
	maxDayInWeek   = 7
)

// Инициализирующая функция
// создает мапы со списками валидных дней недели, месяцев и дней месяца
func init() {
	for i := 1; i <= maxDayInWeek; i++ {
		validWeekDays[i] = true
	}

	for i := 1; i <= maxMonth; i++ {
		validMonths[i] = true
	}

	validMonthDays[-2] = true
	validMonthDays[-1] = true
	for i := 1; i <= maxDayOfMonth; i++ {
		validMonthDays[i] = true
	}
}

// NewDate возвращает дату исходя из правила repeat
//
// Параметры:
//
//	now - время от которого происходит расчет
//	dateStr - дата в формате "20060102", исходное время
//	repeat - правило
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
	// d - day, m - month, y - year
	switch sliceRepeat[0] {
	case "d":
		if len(sliceRepeat) != 2 {
			return "", ErrInvalidFormat
		}
		return handleDayRule(now, date, sliceRepeat[1])
	case "y":
		return handleYearRule(now, date)
	case "m":
		if len(sliceRepeat) < 2 && len(sliceRepeat) > 4 {
			return "", ErrInvalidFormat
		}
		return handleMonthRule(now, date, strings.Join(sliceRepeat[1:], " "))
	case "w":
		if len(sliceRepeat) != 2 {
			return "", ErrInvalidFormat
		}
		return handleWeekRule(now, date, sliceRepeat[1])
	default:
		return "", ErrInvalidFormat
	}
}

// handleDayRule функция обработчик для правила d
//
// Параметры:
//
//	now - время от которого происходит расчет
//	date - дата в формате "20060102", исходное время
//	interval - интервал в днях
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
		if afterNow(result, now) {
			break
		}
	}

	return result.Format("20060102"), nil
}

// handleYearRule функция обработчик для правила y
//
// Параметры:
//
//	now, date - время от которого происходит расчет и дата в формате "20060102"
//	interval - интервал в днях
func handleYearRule(now, date time.Time) (string, error) {
	result := date
	for {
		result = result.AddDate(1, 0, 0)
		if afterNow(result, now) {
			break
		}
	}
	// Обработка високосного года
	// TODO: возможно есть более локаничное решение
	if date.Month() == 2 && date.Day() == 29 {
		if !isLeap(result.Year()) {
			result = time.Date(result.Year(), 3, 1, 0, 0, 0, 0, time.UTC)
		}
	}

	return result.Format("20060102"), nil
}

// extra TASK

// parseAndValidate парсит строку и валидирует ее
// c значениями сформированными в validMap
//
// Параметры:
//
//	s string - строка для парсинга
//	validMap map[int]bool - map с валидными значениями
//
// Возвращает:
//
//	map[int]bool - map с валидными значениями
//	error - ошибка
func parseAndValidate(s string, validMap map[int]bool) (map[int]bool, error) {
	parts := strings.Split(s, ",")
	result := make(map[int]bool)

	// Итерируемся по parts, проверяем каждый элемент на валидность
	// Формируем map result с валидными значениями
	// если что-то не валидно, возвращаем ошибку
	for _, part := range parts {
		num, err := strconv.Atoi(part)
		if err != nil {
			return nil, ErrInvalidFormat
		}

		if !validMap[num] {
			return nil, fmt.Errorf("invalid value %d", num)
		}

		result[num] = true
	}

	if len(result) == 0 {
		return nil, ErrInvalidFormat
	}

	return result, nil
}

// handleWeekRule функция обработчик для правила w
//
// Параметры:
//
//	now, date time.Time - время от которого происходит расчет и дата в формате "20060102"
//	daysStr string - строка с днями недели
func handleWeekRule(now, date time.Time, daysStr string) (string, error) {
	days, err := parseAndValidate(daysStr, validWeekDays)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrInvalidDayOfWeek, err)
	}

	result := date

	for {
		result = result.AddDate(0, 0, 1)
		if afterNow(result, now) {
			weekday := int(result.Weekday())
			// обработка ВС
			if weekday == 0 {
				weekday = 7
			}
			if days[weekday] {
				break
			}
		}
	}

	return result.Format("20060102"), nil
}

// handleMonthRule функция обработчик для правила m
//
// Параметры:
//
//	now, date time.Time - время от которого происходит расчет и дата в формате "20060102"
//	rule string - правило
func handleMonthRule(now, date time.Time, rule string) (string, error) {
	sliceRule := strings.Split(rule, " ")
	if len(sliceRule) > 2 {
		return "", ErrInvalidFormat
	}

	days, err := parseAndValidate(sliceRule[0], validMonthDays)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrInvalidDayOfMonth, err)
	}

	var months map[int]bool
	if len(sliceRule) == 2 {
		months, err = parseAndValidate(sliceRule[1], validMonths)
		if err != nil {
			return "", fmt.Errorf("%w: %v", ErrInvalidMonth, err)
		}
	}

	result := date
	for {
		result = result.AddDate(0, 0, 1)
		if afterNow(result, now) {
			currentMonth := int(result.Month())
			currentDay := result.Day()
			lastDay := lastDayOfMonth(result)

			if months != nil && !months[currentMonth] {
				continue
			}

			if (days[-1] && currentDay == lastDay) ||
				(days[-2] && currentDay == lastDay-1) ||
				(days[currentDay]) {
				break
			}
		}
	}

	return result.Format("20060102"), nil
}
