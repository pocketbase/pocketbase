package cron

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Moment represents a parsed single time moment.
type Moment struct {
	Second    int `json:"second"`
	Minute    int `json:"minute"`
	Hour      int `json:"hour"`
	Day       int `json:"day"`
	Month     int `json:"month"`
	DayOfWeek int `json:"dayOfWeek"`
}

// NewMoment creates a new Moment from the specified time.
func NewMoment(t time.Time) *Moment {
	return &Moment{
		Second:    t.Second(),
		Minute:    t.Minute(),
		Hour:      t.Hour(),
		Day:       t.Day(),
		Month:     int(t.Month()),
		DayOfWeek: int(t.Weekday()),
	}
}

// Schedule stores parsed information for each time component when a cron job should run.
type Schedule struct {
	Seconds    map[int]struct{} `json:"seconds"`
	Minutes    map[int]struct{} `json:"minutes"`
	Hours      map[int]struct{} `json:"hours"`
	Days       map[int]struct{} `json:"days"`
	Months     map[int]struct{} `json:"months"`
	DaysOfWeek map[int]struct{} `json:"daysOfWeek"`
}

// IsDue checks whether the provided Moment satisfies the current Schedule.
func (s *Schedule) IsDue(m *Moment) bool {

	if _, ok := s.Seconds[m.Second]; !ok {
		return false
	}

	if _, ok := s.Minutes[m.Minute]; !ok {
		return false
	}

	if _, ok := s.Hours[m.Hour]; !ok {
		return false
	}

	if _, ok := s.Days[m.Day]; !ok {
		return false
	}

	if _, ok := s.DaysOfWeek[m.DayOfWeek]; !ok {
		return false
	}

	if _, ok := s.Months[m.Month]; !ok {
		return false
	}
	return true
}

var macros = map[string]string{
	"@yearly":   "0 0 0 1 1 *",
	"@annually": "0 0 0 1 1 *",
	"@monthly":  "0 0 0 1 * *",
	"@weekly":   "0 0 0 * * 0",
	"@daily":    "0 0 0 * * *",
	"@midnight": "0 0 0 * * *",
	"@midday": 	 "0 0 12 * * *",
	"@hourly":   "0 0 * * * *",
	"@minutely": "0 * * * * *",
	"@secondly": "* * * * * *",
	"@tick":     "* * * * * *",
}

// NewSchedule creates a new Schedule from a cron expression.
//
// A cron expression could be a macro OR 5/6 segments separated by space,
// representing: minute, hour, day of the month, month and day of the week.
//
// The following segment formats are supported:
//   - wildcard: *
//   - range:    1-30
//   - step:     */n or 1-30/n
//   - list:     1,2,3,10-20/n
//
// The following macros are supported:
//   - @yearly (or @annually)
//   - @monthly
//   - @weekly
//   - @daily (or @midnight)
//   - @midday (each 12h)
//   - @hourly
//   - @minutely
//   - @secondly (or @tick)
func NewSchedule(cronExpr string) (*Schedule, error) {
	if v, ok := macros[cronExpr]; ok {
		cronExpr = v
	}

	segments := strings.Split(cronExpr, " ")

	if len(segments) == 5 {
		segments = strings.Split("0 "+cronExpr, " ")
	}
	if len(segments) != 6 {
		return nil, errors.New("invalid cron expression - must be a valid macro or to have exactly 5 or 6 space separated segments")
	}
	seconds, err := parseCronSegment(segments[0], 0, 59, map[string]string{
	})
	if err != nil {
		return nil, err
	}

	minutes, err := parseCronSegment(segments[1], 0, 59, map[string]string{
	})
	if err != nil {
		return nil, err
	}

	hours, err := parseCronSegment(segments[2], 0, 23, map[string]string{
	})
	if err != nil {
		return nil, err
	}

	days, err := parseCronSegment(segments[3], 1, 31, map[string]string{
	})
	if err != nil {
		return nil, err
	}

	months, err := parseCronSegment(segments[4], 1, 12, map[string]string{
		"JAN": "1",
		"FEB": "2",
		"MAR": "3",
		"APR": "4",
		"MAY": "5",
		"JUN": "6",
		"JUL": "7",
		"AUG": "8",
		"SEP": "9",
		"OCT": "10",
		"NOV": "11",
		"DEC": "12",
	})
	if err != nil {
		return nil, err
	}

	daysOfWeek, err := parseCronSegment(segments[5], 0, 6, map[string]string{
		"SAT": "6",
		"FRI": "5",
		"THU": "4",
		"WED": "3",
		"TUE": "2",
		"MON": "1",
		"SUN": "0",
		"7":   "0",
	})
	if err != nil {
		return nil, err
	}

	return &Schedule{
		Seconds:    seconds,
		Minutes:    minutes,
		Hours:      hours,
		Days:       days,
		Months:     months,
		DaysOfWeek: daysOfWeek,
	}, nil
}

// parseCronSegment parses a single cron expression segment and
// returns its time schedule slots.
func parseCronSegment(segment string, min int, max int, extra map[string]string) (map[int]struct{}, error) {
	slots := map[int]struct{}{}
	segment = strings.ToUpper(strings.Trim(segment, " "))
	list := strings.Split(segment, ",")
	for _, p := range list {
		stepParts := strings.Split(p, "/")

		// step (*/n, 1-30/n)
		var step int
		switch len(stepParts) {
		case 1:
			step = 1
		case 2:
			parsedStep, err := strconv.Atoi(stepParts[1])
			if err != nil {
				return nil, err
			}

			if parsedStep < 1 || parsedStep > max {
				return nil, fmt.Errorf("invalid segment step boundary - the step must be between %d and the %d", 1, max)
			}
			step = parsedStep
		default:
			return nil, errors.New("invalid segment step format - must be in the format */n or 1-30/n")
		}

		// find the min and max range of the segment part
		var rangeMin, rangeMax int
		if stepParts[0] == "*" {
			rangeMin = min
			rangeMax = max
		} else {
			// single digit (1) or range (1-30)
			rangeParts := strings.Split(stepParts[0], "-")
			switch len(rangeParts) {
			case 1:
				if step != 1 {
					return nil, errors.New("invalid segement step - step > 1 could be used only with the wildcard or range format")
				}
				newValue, isExtraValue := extra[rangeParts[0]]
				if isExtraValue {
					rangeParts[0] = newValue
				}
				parsed, err := strconv.Atoi(rangeParts[0])
				if err != nil {
					return nil, err
				}
				if parsed < min || parsed > max {
					return nil, errors.New("invalid segment value - must be between the min and max of the segment")
				}
				rangeMin = parsed
				rangeMax = rangeMin
			case 2:

				newValue, isExtraValue := extra[rangeParts[0]]
				if isExtraValue {
					rangeParts[0] = newValue
				}

				parsedMin, err := strconv.Atoi(rangeParts[0])
				if err != nil {
					return nil, err
				}
				if parsedMin < min || parsedMin > max {
					return nil, fmt.Errorf("invalid segment range minimum - must be between %d and %d", min, max)
				}
				rangeMin = parsedMin

				newValue, isExtraValue = extra[rangeParts[1]]
				if isExtraValue {
					rangeParts[1] = newValue
				}
				parsedMax, err := strconv.Atoi(rangeParts[1])
				if err != nil {
					return nil, err
				}
				if parsedMax < parsedMin || parsedMax > max {
					return nil, fmt.Errorf("invalid segment range maximum - must be between %d and %d", rangeMin, max)
				}
				rangeMax = parsedMax
			default:
				return nil, errors.New("invalid segment range format - the range must have 1 or 2 parts")
			}
		}

		// fill the slots
		for i := rangeMin; i <= rangeMax; i += step {
			slots[i] = struct{}{}
		}
	}

	return slots, nil
}
