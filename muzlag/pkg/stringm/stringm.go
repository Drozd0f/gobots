package stringm

import (
	"fmt"
	"strconv"
)

func ToInt64(s string) (int64, error) {
	count, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("strconv parse int: %w", err)
	}

	return count, nil
}
