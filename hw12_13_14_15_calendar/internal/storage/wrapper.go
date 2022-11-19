package storage

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type Duration time.Duration

func (d *Duration) Value() (driver.Value, error) {
	if d == nil {
		return nil, nil
	}

	zero := time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC)

	return driver.Value(zero.Add(time.Duration(*d)).Format("15:04:05")), nil
}

func (d *Duration) Scan(raw interface{}) error {
	switch v := raw.(type) {
	case int64:
		*d = Duration(time.Duration(v))
	case nil:
		*d = Duration(time.Duration(0))
	case time.Time:
		h, m, s := v.Clock()
		hd, md, sd := time.Duration(h)*time.Hour, time.Duration(m)*time.Minute, time.Duration(s)*time.Second

		*d = (Duration)(hd + md + sd)
	default:
		return fmt.Errorf("cannot sql.Scan() strfmt.Duration from: %#v", v)
	}
	return nil
}
