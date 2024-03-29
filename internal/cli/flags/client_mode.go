package flags

import (
	"fmt"

	"github.com/airbnb/rudolph/pkg/types"
)

const (
	monitorMode       = "monitor"
	lockdownMode      = "lockdown"
	defaultClientMode = "monitor"
)

// configMode is a custom type for use as a CLI flag representing the type of config mode being applied
type ClientMode types.ClientMode

func newConfigTypeValue(val string, c *ClientMode) *ClientMode {
	err := c.Set(val)
	if err != nil {
		fmt.Println(`Warning: invalid value for client mode provided, default to using "monitor"`)
		*c = ClientMode(types.Monitor)
	}

	return c
}

func (i *ClientMode) AsClientMode() types.ClientMode {
	return types.ClientMode(*i)
}

func (i *ClientMode) Set(c string) error {
	switch c {
	case monitorMode:
		*i = ClientMode(types.Monitor)
	case lockdownMode:
		*i = ClientMode(types.Lockdown)
	default:
		return fmt.Errorf(`invalid client mode; must be "monitor" or "lockdown"`)
	}
	return nil
}

func (i *ClientMode) Type() string {
	return "string"
}

func (i *ClientMode) String() string {
	c := (types.ClientMode)(*i)
	// select `monitor` or `lockdown` and if invalid, always default to `defaultClientMode` ClientMode
	switch c {
	case types.Monitor:
		return monitorMode
	case types.Lockdown:
		return lockdownMode
	default:
		return defaultClientMode
	}
}
