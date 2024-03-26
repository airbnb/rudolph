package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSyncTypeTypes_MarshallText(t *testing.T) {
	tests := []struct {
		name     string
		syncType SyncType
		want     []byte
		wantErr  bool
	}{
		{"normal", SyncTypeNormal, []byte("normal"), false},
		{"clean", SyncTypeClean, []byte("clean"), false},
		{"clean_all", SyncTypeCleanAll, []byte("clean_all"), false},
		{"MISSPELLED", SyncType(""), nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.syncType.MarshalText()
			if (err != nil) != tt.wantErr {
				t.Errorf("SyncType.MarshalText() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSyncTypeTypes_UnmarshallText(t *testing.T) {
	tests := []struct {
		name    string
		text    []byte
		want    SyncType
		wantErr bool
	}{
		{"normal", []byte("normal"), SyncTypeNormal, false},
		{"clean", []byte("clean"), SyncTypeClean, false},
		{"clean_all", []byte("clean_all"), SyncTypeCleanAll, false},
		{"MISSPELLED", nil, SyncType(""), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got SyncType
			err := got.UnmarshalText(tt.text)
			if (err != nil) != tt.wantErr {
				t.Errorf("SyncType.UnmarshalText() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
