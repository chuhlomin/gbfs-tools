package gbfs

import (
	"strconv"
	"time"
)

type Timestamp time.Time

func (t *Timestamp) UnmarshalJSON(b []byte) (err error) {
	q, err := strconv.ParseInt(string(b), 10, 64)
	if err != nil {
		return err
	}
	*(*time.Time)(t) = time.Unix(q, 0)
	return nil
}

func (t Timestamp) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatInt(time.Time(t).Unix(), 10)), nil
}

func (t Timestamp) Unix() int64 {
	return time.Time(t).Unix()
}

func (t Timestamp) Time() time.Time {
	return time.Time(t).UTC()
}

func (t Timestamp) String() string {
	return t.Time().String()
}
