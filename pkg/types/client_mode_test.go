package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClientModeTypes_MarshallText(t *testing.T) {
	tests := []struct {
		name       string
		clientMode ClientMode
		want       []byte
		wantErr    bool
	}{
		{"MONITOR", Monitor, []byte("MONITOR"), false},
		{"LOCKDOWN", Lockdown, []byte("LOCKDOWN"), false},
		{"MISSPELLED", ClientMode(0), nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.clientMode.MarshalText()
			if (err != nil) != tt.wantErr {
				t.Errorf("ClientMode.MarshalText() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestClientModeTypes_UnmarshallText(t *testing.T) {
	tests := []struct {
		name    string
		text    []byte
		want    ClientMode
		wantErr bool
	}{
		{"MONITOR", []byte("MONITOR"), Monitor, false},
		{"LOCKDOWN", []byte("LOCKDOWN"), Lockdown, false},
		{"MISSPELLED", nil, ClientMode(0), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got ClientMode
			err := got.UnmarshalText(tt.text)
			if (err != nil) != tt.wantErr {
				t.Errorf("ClientMode.UnmarshalText() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
