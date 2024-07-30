package gbfs

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type ID string

func (id *ID) UnmarshalJSON(data []byte) (err error) {
	var i interface{}
	if err = json.Unmarshal(data, &i); err != nil {
		return err
	}

	switch v := i.(type) {
	case int:
		*(*string)(id) = strconv.Itoa(v)
	case float64:
		*(*string)(id) = fmt.Sprintf("%1.0f", v)
	case string:
		*(*string)(id) = v
	default:
		return fmt.Errorf("parse ID %T", v)
	}
	return nil
}
