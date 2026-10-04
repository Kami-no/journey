package date

import (
	"bytes"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var marchMayChecker = regexp.MustCompile("M([^a]|$)")

// GetCurrentTime - whenever we need time.Now(), we use this function instead so that we always use UTC in journey
func GetCurrentTime() time.Time {
	return time.Now().UTC()
}

func GenerateTimeAgo(date *time.Time) []byte {
	timeAgo := GetCurrentTime().Sub(*date)
	if timeAgo.Minutes() < 1 {
		return []byte("a few seconds ago")
	}
	if timeAgo.Minutes() < 2 {
		return []byte("a minute ago")
	}
	if timeAgo.Minutes() < 60 {
		var buffer bytes.Buffer
		buffer.WriteString(strconv.Itoa(int(timeAgo.Minutes())))
		buffer.WriteString(" minutes ago")
		return buffer.Bytes()
	}
	if timeAgo.Hours() < 2 {
		return []byte("an hour ago")
	}
	if timeAgo.Hours() < 24 {
		var buffer bytes.Buffer
		buffer.WriteString(strconv.Itoa(int(timeAgo.Hours())))
		buffer.WriteString(" hours ago")
		return buffer.Bytes()
	}
	if timeAgo.Hours() < 48 {
		return []byte("a day ago")
	}
	days := int(timeAgo.Hours() / 24)
	if days < 25 {
		var buffer bytes.Buffer
		buffer.WriteString(strconv.Itoa(days))
		buffer.WriteString(" days ago")
		return buffer.Bytes()
	}
	if days < 45 {
		return []byte("a month ago")
	}
	if days < 345 {
		months := days / 30
		if months < 2 {
			months = 2
		}
		var buffer bytes.Buffer
		buffer.WriteString(strconv.Itoa(months))
		buffer.WriteString(" months ago")
		return buffer.Bytes()
	}
	if days < 548 {
		return []byte("a year ago")
	}
	years := days / 365
	if years < 2 {
		years = 2
	}
	var buffer bytes.Buffer
	buffer.WriteString(strconv.Itoa(years))
	buffer.WriteString(" years ago")
	return buffer.Bytes()
}

func FormatDate(format string, date *time.Time) []byte {

	// Do these first (so they don't accidentally replace something the others insert)
	format = strings.ReplaceAll(format, "h", replaceh(date))
	format = strings.ReplaceAll(format, "s", strconv.Itoa(date.Second()))

	// Year, month, and day
	format = strings.ReplaceAll(format, "Do", replaceDo(date))
	format = strings.ReplaceAll(format, "YYYY", strconv.Itoa(date.Year()))
	if date.Year() > 99 {
		format = strings.ReplaceAll(format, "YY", strconv.Itoa(date.Year())[2:])
	}
	format = strings.ReplaceAll(format, "Q", strconv.Itoa(((int(date.Month())-1)/3)+1))
	format = strings.ReplaceAll(format, "DDDD", replaceDDDD(date))
	format = strings.ReplaceAll(format, "DDD", replaceDDD(date))
	format = strings.ReplaceAll(format, "DD", replaceDD(date))
	format = strings.ReplaceAll(format, "X", strconv.FormatInt(date.Unix(), 10))
	// Unix ms ('x') is not used by ghost. Excluding it for now.
	// format = strings.ReplaceAll(format, "x", strconv.FormatInt((date.UnixNano()/1000000), 10))

	// Locale formats. Not supported yet
	format = strings.ReplaceAll(format, "gggg", strconv.Itoa(date.Year()))
	if date.Year() > 99 {
		format = strings.ReplaceAll(format, "gg", strconv.Itoa(date.Year())[2:])
	}
	format = strings.ReplaceAll(format, "ww", replaceww(date))
	format = strings.ReplaceAll(format, "w", replacew(date))
	format = strings.ReplaceAll(format, "e", strconv.Itoa(int(date.Weekday())))

	// ISO week date formats. Not supported yet - https://en.wikipedia.org/wiki/ISO_week_date
	format = strings.ReplaceAll(format, "GGGG", strconv.Itoa(date.Year()))
	if date.Year() > 99 {
		format = strings.ReplaceAll(format, "GG", strconv.Itoa(date.Year())[2:])
	}
	format = strings.ReplaceAll(format, "WW", replaceww(date))
	format = strings.ReplaceAll(format, "W", replacew(date))
	format = strings.ReplaceAll(format, "E", strconv.Itoa(int(date.Weekday())))

	// Hour, minute, second, millisecond, and offset
	format = strings.ReplaceAll(format, "HH", replaceHH(date))
	format = strings.ReplaceAll(format, "H", strconv.Itoa(date.Hour()))
	format = strings.ReplaceAll(format, "hh", replacehh(date))
	format = strings.ReplaceAll(format, "a", replacea(date))
	format = strings.ReplaceAll(format, "A", replaceA(date))
	format = strings.ReplaceAll(format, "mm", replacemm(date))
	format = strings.ReplaceAll(format, "m", strconv.Itoa(date.Minute()))
	format = strings.ReplaceAll(format, "ss", replacess(date))
	format = strings.ReplaceAll(format, "SSS", strconv.Itoa(date.Nanosecond()/1000000))
	format = strings.ReplaceAll(format, "SS", strconv.Itoa(date.Nanosecond()/10000000))
	format = strings.ReplaceAll(format, "S", strconv.Itoa(date.Nanosecond()/100000000))
	format = strings.ReplaceAll(format, "ZZ", replaceZZ(date))
	format = strings.ReplaceAll(format, "Z", replaceZ(date))

	// Not documented for moment.js, but seems to be used by ghost themes
	format = strings.ReplaceAll(format, "dddd", date.Weekday().String())

	// This needs to be last so that month strings don't interfere with the other replace functions
	format = strings.ReplaceAll(format, "MMMM", date.Month().String())
	if len(date.Month().String()) > 2 {
		format = strings.ReplaceAll(format, "MMM", date.Month().String()[:3])
	}
	format = strings.ReplaceAll(format, "MM", replaceMM(date))
	// Replace M - make sure the Ms in March and May don't get replaced.
	// TODO: Regex could be improved, only recognizes 'M's that are not followed by 'a's.
	format = marchMayChecker.ReplaceAllString(format, strconv.Itoa(int(date.Month())))

	return []byte(format)
}

func replaceMM(date *time.Time) string {
	var buffer bytes.Buffer
	month := int(date.Month())
	if month < 10 {
		buffer.WriteString("0")
		buffer.WriteString(strconv.Itoa(month))
	} else {
		buffer.WriteString(strconv.Itoa(month))
	}
	return buffer.String()
}

func replaceDDDD(date *time.Time) string {
	var buffer bytes.Buffer
	startOfYear := time.Date((date.Year() - 1), time.December, 31, 0, 0, 0, 0, time.UTC)
	days := int(date.Sub(startOfYear) / (24 * time.Hour))
	if days < 10 {
		buffer.WriteString("00")
	} else if days < 100 {
		buffer.WriteString("0")
	}
	buffer.WriteString(strconv.Itoa(days))
	return buffer.String()
}

func replaceDDD(date *time.Time) string {
	startOfYear := time.Date((date.Year() - 1), time.December, 31, 0, 0, 0, 0, time.UTC)
	return strconv.Itoa(int(date.Sub(startOfYear) / (24 * time.Hour)))
}

func replaceDD(date *time.Time) string {
	var buffer bytes.Buffer
	if date.Day() < 10 {
		buffer.WriteString("0")
		buffer.WriteString(strconv.Itoa(date.Day()))
	} else {
		buffer.WriteString(strconv.Itoa(date.Day()))
	}
	return buffer.String()
}

func replaceDo(date *time.Time) string {
	var buffer bytes.Buffer
	buffer.WriteString(strconv.Itoa(date.Day()))
	if date.Day() == 1 || date.Day() == 21 || date.Day() == 31 {
		buffer.WriteString("st")
	} else if date.Day() == 2 || date.Day() == 22 {
		buffer.WriteString("nd")
	} else if date.Day() == 3 || date.Day() == 23 {
		buffer.WriteString("rd")
	} else {
		buffer.WriteString("th")
	}
	return buffer.String()
}

func replaceww(date *time.Time) string {
	var buffer bytes.Buffer
	startOfYear := time.Date((date.Year() - 1), time.December, 25, 0, 0, 0, 0, time.UTC)
	weeks := int(date.Sub(startOfYear) / (24 * time.Hour * 7))
	if weeks < 10 {
		buffer.WriteString("0")
	}
	buffer.WriteString(strconv.Itoa(weeks))
	return buffer.String()
}

func replacew(date *time.Time) string {
	startOfYear := time.Date((date.Year() - 1), time.December, 25, 0, 0, 0, 0, time.UTC)
	return strconv.Itoa(int(date.Sub(startOfYear) / (24 * time.Hour * 7)))
}

func replaceHH(date *time.Time) string {
	var buffer bytes.Buffer
	hour := date.Hour()
	if hour < 10 {
		buffer.WriteString("0")
	}
	buffer.WriteString(strconv.Itoa(hour))
	return buffer.String()
}

func replacehh(date *time.Time) string {
	var buffer bytes.Buffer
	hour := date.Hour()
	if hour == 0 {
		hour = 12
	} else if hour > 12 {
		hour = hour - 12
	}
	if hour < 10 {
		buffer.WriteString("0")
	}
	buffer.WriteString(strconv.Itoa(hour))
	return buffer.String()
}

func replaceh(date *time.Time) string {
	var buffer bytes.Buffer
	hour := date.Hour()
	if hour == 0 {
		hour = 12
	} else if hour > 12 {
		hour = hour - 12
	}
	buffer.WriteString(strconv.Itoa(hour))
	return buffer.String()
}

func replacea(date *time.Time) string {
	if date.Hour() < 12 {
		return "am"
	}
	return "pm"
}

func replaceA(date *time.Time) string {
	if date.Hour() < 12 {
		return "AM"
	}
	return "PM"
}

func replacemm(date *time.Time) string {
	var buffer bytes.Buffer
	minute := date.Minute()
	if minute < 10 {
		buffer.WriteString("0")
	}
	buffer.WriteString(strconv.Itoa(minute))
	return buffer.String()
}

func replacess(date *time.Time) string {
	var buffer bytes.Buffer
	second := date.Second()
	if second < 10 {
		buffer.WriteString("0")
	}
	buffer.WriteString(strconv.Itoa(second))
	return buffer.String()
}

func replaceZZ(date *time.Time) string {
	var buffer bytes.Buffer
	_, offset := date.In(time.Local).Zone()
	offset = offset / 3600
	if offset > 0 {
		buffer.WriteString("+")
	} else {
		buffer.WriteString("-")
	}
	if offset < 0 {
		offset = -offset
	}
	if offset < 10 {
		buffer.WriteString("0")
	}
	buffer.WriteString(strconv.Itoa(offset))
	buffer.WriteString("00")
	return buffer.String()
}

func replaceZ(date *time.Time) string {
	var buffer bytes.Buffer
	_, offset := date.In(time.Local).Zone()
	offset = offset / 3600
	if offset > 0 {
		buffer.WriteString("+")
	} else {
		buffer.WriteString("-")
	}
	if offset < 0 {
		offset = -offset
	}
	if offset < 10 {
		buffer.WriteString("0")
	}
	buffer.WriteString(strconv.Itoa(offset))
	buffer.WriteString(":00")
	return buffer.String()
}
