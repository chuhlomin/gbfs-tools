package gbfs

import (
	"fmt"
	"strings"
	"time"
)

type Weekday time.Weekday

func (w *Weekday) UnmarshalJSON(b []byte) (err error) {
	switch str := strings.ToLower(strings.Trim(string(b), `"`)); str {
	case "mon":
		*(*time.Weekday)(w) = time.Monday
	case "tue":
		*(*time.Weekday)(w) = time.Tuesday
	case "wed":
		*(*time.Weekday)(w) = time.Wednesday
	case "thu":
		*(*time.Weekday)(w) = time.Thursday
	case "fri":
		*(*time.Weekday)(w) = time.Friday
	case "sat":
		*(*time.Weekday)(w) = time.Saturday
	case "sun":
		*(*time.Weekday)(w) = time.Sunday
	default:
		return fmt.Errorf("parse %q as weekday", str)
	}

	return err
}
