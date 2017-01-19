package crontab

import (
	"crypto/sha256"
	"fmt"
	"strconv"
	"strings"

	"github.com/robfig/cron"
)

const (
	minMinute = 0
	maxMinute = 59
	minHour   = 0
	maxHour   = 23
	minDom    = 1
	maxDom    = 31
	minMonth  = 1
	maxMonth  = 12
	minDow    = 0
	maxDow    = 6
)

var (
	weekDays = []string{
		"Sun",
		"Mon",
		"Tue",
		"Wed",
		"Thu",
		"Fri",
		"Sat",
	}
)

// Schedule represents crontab spec and command
type Schedule struct {
	Spec    string
	Command string
}

// Parse parses crontab file and return a list of Schedule
func Parse(crontab string) ([]*Schedule, error) {
	schedules := []*Schedule{}
	lines := strings.Split(crontab, "\n")

	for _, line := range lines {
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "#") {
			continue
		}

		ss := strings.SplitN(line, " ", 6)
		if len(ss) < 6 {
			return []*Schedule{}, fmt.Errorf("Invalid format. line: %s", line)
		}

		schedules = append(schedules, &Schedule{
			Spec:    strings.Join(ss[0:5], " "),
			Command: ss[5],
		})
	}

	return schedules, nil
}

// ConvertToSystemdCalendar converts crontab spec format to Systemd Timer format
//   crontab:       https://en.wikipedia.org/wiki/Cron
//   Systemd Timer: https://www.freedesktop.org/software/systemd/man/systemd.time.html
func (s *Schedule) ConvertToSystemdCalendar() (string, error) {
	schedule, err := cron.ParseStandard(s.Spec)
	if err != nil {
		return "", err
	}

	specSchedule, ok := schedule.(*cron.SpecSchedule)
	if !ok {
		return "", fmt.Errorf("Unable to convert Schedule to SpecSchedule")
	}

	minutes := parseBits(specSchedule.Minute, minMinute, maxMinute)
	hours := parseBits(specSchedule.Hour, minHour, maxHour)
	doms := parseBits(specSchedule.Dom, minDom, maxDom)
	months := parseBits(specSchedule.Month, minMonth, maxMonth)
	dows := parseBits(specSchedule.Dow, minDow, maxDow)

	fields := []string{}

	if dows != "*" {
		weekdays, err := convertDowsToWeekdays(dows)
		if err != nil {
			return "", err
		}
		fields = append(fields, weekdays)
	}

	if months != "*" || doms != "*" {
		fields = append(fields, fmt.Sprintf("%s-%s", months, doms))
	}

	fields = append(fields, fmt.Sprintf("%s:%s", hours, minutes))

	return strings.Join(fields, " "), nil
}

// SHA256Sum generates SHA-256 checksum of schedule
func (s *Schedule) SHA256Sum() string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(fmt.Sprintf("%s;%s", s.Spec, s.Command))))
}

func convertDowsToWeekdays(bits string) (string, error) {
	dows := []string{}

	for _, bit := range strings.Split(bits, ",") {
		b, err := strconv.Atoi(bit)
		if err != nil {
			return "", err
		}
		dows = append(dows, weekDays[b])
	}

	return strings.Join(dows, ","), nil
}

func parseBits(n uint64, min, max int) string {
	var all1 uint64

	for i := min; i <= max; i++ {
		all1 |= 1 << uint(i)
	}

	if n&all1 == all1 {
		return "*"
	}

	bits := []string{}

	for i := 0; i <= max; i++ {
		if n&(1<<uint(i)) > 0 {
			bits = append(bits, strconv.Itoa(i))
		}
	}

	return strings.Join(bits, ",")
}
