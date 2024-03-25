package types

import "fmt"

// ClientMode specifies which mode the Santa client will evaluate rules in.
type ClientMode int

const (
	Monitor ClientMode = iota + 1
	Lockdown
)

// UnmarshalText yes
func (c *ClientMode) UnmarshalText(text []byte) error {
	switch mode := string(text); mode {
	case "MONITOR":
		*c = Monitor
	case "LOCKDOWN":
		*c = Lockdown
	default:
		return fmt.Errorf("unknown client_mode value %q", mode)
	}
	return nil
}

// MarshalText yes
func (c ClientMode) MarshalText() ([]byte, error) {
	switch c {
	case Monitor:
		return []byte("MONITOR"), nil
	case Lockdown:
		return []byte("LOCKDOWN"), nil
	default:
		return nil, fmt.Errorf("unknown client_mode %d", c)
	}
}
