package gbfs

import (
	"strconv"
	"strings"
)

type Bool bool

func (b *Bool) UnmarshalJSON(data []byte) (err error) {
	switch str := strings.ToLower(strings.Trim(string(data), `"`)); str {
	case "true":
		*b = true
	case "false":
		*b = false
	default:
		var n float64
		n, err = strconv.ParseFloat(str, 64)
		if n > 0 {
			*(*bool)(b) = true
		}
	}
	return err
}
