package storage

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type Duration time.Duration

func (d Duration) Value() (driver.Value, error) {
	val := time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC).Add(time.Duration(d))

	return driver.Value(
		fmt.Sprintf("%d:%d:%d", val.Hour(), val.Minute(), val.Second())), nil
}

func (d *Duration) Scan(raw interface{}) error {
	switch v := raw.(type) {
	case int64:
		*d = Duration(time.Duration(v))
	case nil:
		*d = Duration(time.Duration(0))
	case time.Time:
		zero := time.Date(0, time.January, 1, 0, 0, 0, 0, time.UTC)
		*d = (Duration)(v.Sub(zero))
	default:
		return fmt.Errorf("cannot sql.Scan() strfmt.Duration from: %#v", v)
	}
	return nil
}
