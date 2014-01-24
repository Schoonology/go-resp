package resp

import (
	"fmt"
	"strings"
)

func Unmarshall(str string) (interface{}, error) {
	s := NewScanner(strings.NewReader(str))
	s.Scan()
	return s.Obj(), s.Err()
}

func Marshall(v interface{}) (string, error) {
	switch v := v.(type) {
	case bool:
		if v {
			return Marshall(1)
		} else {
			return Marshall(0)
		}
	case int, int8, int16, int32, int64:
		return fmt.Sprintf(":%d\r\n", v), nil
	case string:
		return fmt.Sprintf("$%d\r\n%s\r\n", len(v), v), nil
	case []byte:
		return Marshall(string(v))
	case Status:
		return fmt.Sprintf("+%s\r\n", v.Message), nil
	case Error:
		return fmt.Sprintf("-%s %s\r\n", v.Type, v.Message), nil
	case []interface{}:
		var err error
		arr := make([]string, len(v)+1)
		arr[0] = fmt.Sprintf("*%d\r\n", len(v))
		for k, v := range v {
			arr[k+1], err = Marshall(v)
			if err != nil {
				return "", err
			}
		}
		return strings.Join(arr, ""), nil
	default:
		return "", fmt.Errorf("Invalid type %T", v)
	}
	return "", nil
}
