package api

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func afterNow(date, now time.Time) bool {
	dateComp := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	nowComp := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	return dateComp.After(nowComp)
}

func NextDate(now time.Time, dstart string, repeat string) (string, error) {

	if strings.TrimSpace(repeat) == "" {
		return "", errors.New("в параметре repeat — пустая строка")
	}

	dStart, err := time.Parse(dateFormat, dstart)
	if err != nil {
		return "", fmt.Errorf("некорректная дата dstart: %w", err)
	}
	repeat = strings.TrimSpace(repeat)
	ruleRepeat := strings.Split(repeat, " ")

	switch ruleRepeat[0] {
	case "d":
		if len(ruleRepeat) != 2 {
			return "", errors.New("правило требует указания двух параметров в формате d N")
		}
		days, err := strconv.Atoi(ruleRepeat[1])
		if err != nil {
			return "", errors.New("недопустимый символ")
		}
		if days <= 0 {
			return "", errors.New("не указан интервал в днях")
		}
		if days > 400 {
			return "", errors.New("превышен максимально допустимый интервал")
		}
		for {
			dStart = dStart.AddDate(0, 0, days)
			if afterNow(dStart, now) {
				break
			}
		}

	case "y":
		if len(ruleRepeat) != 1 {
			return "", errors.New("для правила указаны лишние параметры")
		}

		for {
			dStart = dStart.AddDate(1, 0, 0)
			if afterNow(dStart, now) {
				break
			}
		}

	default:
		return "", errors.New("unsupported format")
	}

	return dStart.Format(dateFormat), nil
}
