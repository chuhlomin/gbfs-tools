package gbfs

import (
	"strings"
	"time"
)

type Date time.Time

const dateFormat = "2006-01-02"

func (d *Date) UnmarshalJSON(input []byte) error {
	newTime, err := time.Parse(dateFormat, strings.Trim(string(input), `"`))
	if err != nil {
		return err
	}

	*(*time.Time)(d) = newTime
	return nil
}

func (d Date) MarshalJSON() ([]byte, error) {
	return []byte(d.String()), nil
}

func (d Date) Unix() int64 {
	return time.Time(d).Unix()
}

func (d Date) Time() time.Time {
	return time.Time(d).UTC()
}

func (d *Date) String() string {
	return "\"" + time.Time(*d).Format(dateFormat) + "\""
}
