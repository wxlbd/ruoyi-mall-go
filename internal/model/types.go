package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"strconv"
	"strings"
)

// BitBool is a boolean that maps to BIT(1) in database
type BitBool bool

// Scan implements the Scanner interface.
func (b *BitBool) Scan(value interface{}) error {
	if value == nil {
		*b = false
		return nil
	}

	switch v := value.(type) {
	case []uint8:
		if len(v) > 0 {
			*b = BitBool(v[0] == 1)
		} else {
			*b = false
		}
	case int64:
		*b = BitBool(v == 1)
	case bool:
		*b = BitBool(v)
	default:
		return errors.New("incompatible type for BitBool")
	}
	return nil
}

// Value implements the driver Valuer interface.
func (b BitBool) Value() (driver.Value, error) {
	if b {
		return []byte{1}, nil
	}
	return []byte{0}, nil
}

func NewBitBool(b bool) BitBool {
	return BitBool(b)
}

// IntListFromCSV handles comma-separated integer lists from MyBatis IntegerListTypeHandler.
// Supports both "1,2,3" format and JSON "[1,2,3]" format.
type IntListFromCSV []int

func (l *IntListFromCSV) Scan(value interface{}) error {
	if value == nil {
		*l = nil
		return nil
	}

	var data []byte
	switch v := value.(type) {
	case []byte:
		data = v
	case string:
		data = []byte(v)
	default:
		return errors.New("incompatible type for IntListFromCSV")
	}

	if len(data) == 0 {
		*l = nil
		return nil
	}

	str := strings.TrimSpace(string(data))
	if str == "" {
		*l = nil
		return nil
	}

	// Try JSON format first
	if strings.HasPrefix(str, "[") {
		var result []int
		if err := json.Unmarshal(data, &result); err == nil {
			*l = result
			return nil
		}
	}

	// Parse as comma-separated
	parts := strings.Split(str, ",")
	result := make([]int, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		i, err := strconv.Atoi(p)
		if err != nil {
			return err
		}
		result = append(result, i)
	}
	*l = result
	return nil
}

func (l IntListFromCSV) Value() (driver.Value, error) {
	if len(l) == 0 {
		return "", nil
	}
	parts := make([]string, len(l))
	for i, v := range l {
		parts[i] = strconv.Itoa(v)
	}
	return strings.Join(parts, ","), nil
}

func (l IntListFromCSV) MarshalJSON() ([]byte, error) {
	return json.Marshal([]int(l))
}

// StringListFromCSV handles comma-separated string lists.
type StringListFromCSV []string

func (l *StringListFromCSV) Scan(value interface{}) error {
	if value == nil {
		*l = nil
		return nil
	}

	var data []byte
	switch v := value.(type) {
	case []byte:
		data = v
	case string:
		data = []byte(v)
	default:
		return errors.New("incompatible type for StringListFromCSV")
	}

	if len(data) == 0 {
		*l = nil
		return nil
	}

	str := strings.TrimSpace(string(data))
	if str == "" {
		*l = nil
		return nil
	}

	// Try JSON format first
	if strings.HasPrefix(str, "[") {
		var result []string
		if err := json.Unmarshal(data, &result); err == nil {
			*l = result
			return nil
		}
	}

	// Parse as comma-separated
	parts := strings.Split(str, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			result = append(result, p)
		}
	}
	*l = result
	return nil
}

func (l StringListFromCSV) Value() (driver.Value, error) {
	if len(l) == 0 {
		return "", nil
	}
	return strings.Join(l, ","), nil
}

func (l StringListFromCSV) MarshalJSON() ([]byte, error) {
	return json.Marshal([]string(l))
}
