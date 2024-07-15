package botkit

import (
	"encoding/json"
	"fmt"
)

func ParseJSON[T any](src string) (T, error) {
	const op = "botkit.ParseJSON"

	var args T

	if err := json.Unmarshal([]byte(src), &args); err != nil {
		return *(new(T)), fmt.Errorf("%s: %w", op, err)
	}

	return args, nil
}
