package gbfs

import (
	"strings"
	"time"
)

type Clock time.Time

const clockFormat = "15:04:05"

func (clock *Clock) UnmarshalJSON(input []byte) error {
	newTime, err := time.Parse(clockFormat, strings.Trim(string(input), `"`))
	if err != nil {
		return err
	}

	*(*time.Time)(clock) = newTime
	return nil
}

func (clock Clock) MarshalJSON() ([]byte, error) {
	return []byte(clock.String()), nil
}

func (clock *Clock) String() string {
	return "\"" + time.Time(*clock).Format(clockFormat) + "\""
}
