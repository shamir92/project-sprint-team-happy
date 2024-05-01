package internal

import "time"

func ToISO8601Format(t time.Time) string {
	layout := "2006-01-02T15:04:05Z07:00" //ISO 8601
	return t.Format(layout)
}
