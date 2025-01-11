package convertutils

import (
	"fmt"
	"strconv"
	"strings"
)

// TODO: delete
func IntSliceConvertIntoString(slice []int) string {
	if len(slice) == 0 {
		return ""
	}

	values := make([]string, len(slice))
	for i, v := range slice {
		values[i] = strconv.Itoa(v)
	}

	return "{" + strings.Join(values, ",") + "}"
}

func StringConvertIntoIntSlice(str string) ([]int, error) {
	str = strings.Trim(str, "{}")

	if str == "" {
		return []int{}, nil
	}

	strValues := strings.Split(str, ",")

	intValues := make([]int, len(strValues))
	for i, str := range strValues {
		val, err := strconv.Atoi(str)
		if err != nil {
			return nil, fmt.Errorf("Error to convert '%s' into int: %w", str, err)
		}
		intValues[i] = val
	}

	return intValues, nil
}
