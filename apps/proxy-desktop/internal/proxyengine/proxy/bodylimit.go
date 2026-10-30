package proxy

import (
	"errors"
	"fmt"
	"io"
)

var errRequestBodyTooLarge = errors.New("request body too large")

func readLimited(r io.Reader, max int64) ([]byte, error) {
	if max < 1 {
		return nil, fmt.Errorf("invalid max body size")
	}
	data, err := io.ReadAll(io.LimitReader(r, max+1))
	if err != nil {
		return nil, err
	}
	if int64(len(data)) > max {
		return nil, fmt.Errorf("%w (max %d bytes)", errRequestBodyTooLarge, max)
	}
	return data, nil
}
