package enum

import (
	"database/sql/driver"
	"fmt"
)

// LogStatus represents the type of user.
type LogStatus uint

const (
	Unknown LogStatus = 0
	Queued  LogStatus = 1
	Pending LogStatus = 2
	Succeed LogStatus = 3
	Failed  LogStatus = 4
)

func (m *LogStatus) IsValid() bool {
	switch *m {
	case Queued, Pending, Succeed, Failed:
		return true
	default:
		return false
	}
}

func (m *LogStatus) String() string {
	switch *m {
	case Queued:
		return "Queued"
	case Pending:
		return "Pending"
	case Succeed:
		return "Succeed"
	case Failed:
		return "Failed"
	default:
		return "Unknown"
	}
}

// Value implements driver.Valuer so GORM/DB can store the enum as integer.
func (m *LogStatus) Value() (driver.Value, error) {
	if !m.IsValid() {
		return nil, fmt.Errorf("invalid LogStatus: %d", m)
	}
	return int64(*m), nil
}

// Scan implements sql.Scanner so GORM/DB can read the integer into the enum.
func (m *LogStatus) Scan(value any) error {
	if value == nil {
		*m = Unknown
		return nil
	}
	switch v := value.(type) {
	case int:
		*m = LogStatus(v)
	case int32:
		*m = LogStatus(v)
	case int64:
		*m = LogStatus(v)
	case uint64:
		*m = LogStatus(v)
	case string:
		switch v {
		case "1", "Queued", "queued":
			*m = Queued
		case "2", "Pending", "pending":
			*m = Pending
		case "3", "Succeed", "succeed":
			*m = Succeed
		case "4", "Failed", "failed":
			*m = Failed
		default:
			return fmt.Errorf("cannot scan LogStatus from string: %s", v)
		}
	default:
		return fmt.Errorf("unsupported scan type for LogStatus: %T", value)
	}
	if !m.IsValid() {
		return fmt.Errorf("scanned invalid LogStatus value: %d", *m)
	}
	return nil
}
