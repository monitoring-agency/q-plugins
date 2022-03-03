package formatter

import (
	"fmt"
	"strings"
	"time"
)

func FormatTimeDuration(duration *time.Duration) string {
	r := duration.Round(time.Second)
	d := r / (time.Hour * 24)
	r -= d * time.Hour * 24
	h := r / time.Hour
	r -= h * time.Hour
	m := r / time.Minute
	r -= m * time.Minute
	s := r / time.Second

	var f []string

	if d == 1 {
		f = append(f, fmt.Sprintf("%d day", d))
	} else if d > 1 {
		f = append(f, fmt.Sprintf("%d days", d))
	}

	if h == 1 {
		f = append(f, fmt.Sprintf("%d hour", h))
	} else if h > 1 {
		f = append(f, fmt.Sprintf("%d hours", h))
	}

	if m > 0 {
		f = append(f, fmt.Sprintf("%d min", m))
	}

	if s > 0 {
		f = append(f, fmt.Sprintf("%ds", s))
	}
	return fmt.Sprint(strings.Join(f, ", "))
}
