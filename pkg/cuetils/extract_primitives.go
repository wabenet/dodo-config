package cuetils

import (
	"fmt"

	"cuelang.org/go/cue"
)

func String(name string, value cue.Value) (string, error) {
	out, err := value.String()
	if err != nil {
		return "", fmt.Errorf("value for %s is not a string: %w", name, err)
	}

	return out, nil
}

func Bool(name string, value cue.Value) (bool, error) {
	out, err := value.Bool()
	if err != nil {
		return false, fmt.Errorf("value for %s is not a boolean: %w", name, err)
	}

	return out, nil
}

func Int(name string, value cue.Value) (int64, error) {
	out, err := value.Int64()
	if err != nil {
		return 0, fmt.Errorf("value for %s is not an integer: %w", name, err)
	}

	return out, nil
}
