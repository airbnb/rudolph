package clock

import "time"

func Unixtimestamp(t time.Time) int64 {
	return t.UTC().Unix()
}

func FromUnixtimestamp(i int64) time.Time {
	return time.Unix(i, 0)
}

func RFC3339(t time.Time) string {
	return t.UTC().Format(time.RFC3339)
}

func ParseRFC3339(ts string) (time.Time, error) {
	return time.Parse(time.RFC3339, ts)
}
