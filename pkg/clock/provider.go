package clock

import "time"

type TimeProvider interface {
	Now() time.Time
}

type ConcreteTimeProvider struct{}

func (c ConcreteTimeProvider) Now() time.Time {
	return time.Now()
}

// This time provider is useful for unit tests, where you can "freeze" time
// to have concrete values you can test in the database
type FrozenTimeProvider struct {
	Current time.Time
}

func (p FrozenTimeProvider) Now() time.Time {
	return p.Current
}

// This time provider is useful to just conveniently freeze time at the turn of the millenium
// and you never need fine-grained control of time
type Y2K struct {
	TimeProvider
}

func (p Y2K) Now() time.Time {
	return Y2KTime()
}

// This time provider is useful if you need to hop through time in your tests; for example, to
// test cache expiration or such
type TimeMachine struct {
	Current time.Time
}

func (p *TimeMachine) Now() time.Time {
	return p.Current
}
func (p *TimeMachine) Travel(newTime time.Time) {
	p.Current = newTime
}

func Y2KTime() time.Time {
	var theTime, _ = ParseRFC3339("2000-01-01T00:00:00Z")
	return theTime.UTC()
}
