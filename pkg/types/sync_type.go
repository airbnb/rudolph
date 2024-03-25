package types

import "fmt"

type SyncType string

const (
	SyncTypeNormal   SyncType = "normal"
	SyncTypeClean    SyncType = "clean"
	SyncTypeCleanAll SyncType = "clean_all"
)

// UnmarshalText
func (s *SyncType) UnmarshalText(text []byte) error {
	switch syncType := string(text); syncType {
	case "normal":
		fallthrough
	case "NORMAL":
		*s = SyncTypeNormal
	case "clean":
		fallthrough
	case "CLEAN":
		*s = SyncTypeClean
	case "clean_all":
		fallthrough
	case "CLEAN_ALL":
		*s = SyncTypeCleanAll
	default:
		return fmt.Errorf("unknown sync_type value %q", syncType)
	}
	return nil
}

// MarshalText
func (s SyncType) MarshalText() ([]byte, error) {
	switch s {
	case SyncTypeNormal:
		return []byte(SyncTypeNormal), nil
	case SyncTypeClean:
		return []byte(SyncTypeClean), nil
	case SyncTypeCleanAll:
		return []byte(SyncTypeCleanAll), nil
	default:
		return nil, fmt.Errorf("unknown sync_type %s", s)
	}
}
