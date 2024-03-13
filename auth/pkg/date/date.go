package date

import (
	"encoding/json"
	"fmt"
	"time"
)

// Date represents a date on DateOnly format.
// This code was based by: https://www.willem.dev/articles/change-time-format-json/#skip-the-manual-json-wrangling
type Date struct {
	T time.Time
}

func (d Date) MarshalJSON() ([]byte, error) {
	formatted := d.T.Format(time.DateOnly)
	return json.Marshal(formatted)
}

func (d *Date) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return fmt.Errorf("failed to unmarshal to string: %w", err)
	}
	t, err := time.Parse(time.DateOnly, s)
	if err != nil {
		return fmt.Errorf("error parsing date: %w", err)
	}
	d.T = t
	return nil
}

func (d Date) String() string {
	return d.T.Format(time.DateOnly)
}

func New(t time.Time) Date {
	return Date{T: t}
}
