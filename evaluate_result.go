package togglr

import (
	"fmt"
	"strconv"
	"strings"
)

type EvalResult struct {
	featureKey string
	rawValue   string
	enabled    bool
	found      bool
	err        error
}

func (r *EvalResult) Err() error {
	return r.err
}

func (r *EvalResult) Found() bool {
	return r.found
}

func (r *EvalResult) Enabled() bool {
	return r.enabled
}

func (r *EvalResult) Value() string {
	if !r.found || !r.enabled {
		return ""
	}

	if r.err != nil {
		return ""
	}

	return r.rawValue
}

func (r *EvalResult) Result() (string, error) {
	if r.err != nil {
		return "", r.err
	}

	return r.Value(), nil
}

func (r *EvalResult) Bool() (bool, error) {
	if r.err != nil {
		return false, r.err
	}

	if !r.found || !r.enabled {
		return false, nil
	}

	switch strings.ToLower(r.rawValue) {
	case "true", "1", "yes", "on":
		return true, nil
	case "false", "0", "no", "off":
		return false, nil
	default:
		return false, fmt.Errorf("cannot convert %q to bool", r.rawValue)
	}
}

func (r *EvalResult) Int32() (int32, error) {
	if r.err != nil {
		return 0, r.err
	}

	if !r.found || !r.enabled {
		return 0, nil
	}

	val, err := strconv.ParseInt(r.rawValue, 10, 32)
	if err != nil {
		return 0, err
	}

	return int32(val), nil
}

func (r *EvalResult) UInt32() (uint32, error) {
	if r.err != nil {
		return 0, r.err
	}

	if !r.found || !r.enabled {
		return 0, nil
	}

	val, err := strconv.ParseUint(r.rawValue, 10, 64)
	if err != nil {
		return 0, err
	}

	return uint32(val), nil
}

func (r *EvalResult) Float32() (float32, error) {
	if r.err != nil {
		return 0, r.err
	}

	if !r.found || !r.enabled {
		return 0, nil
	}

	val, err := strconv.ParseFloat(r.rawValue, 64)
	if err != nil {
		return 0, err
	}

	return float32(val), nil
}

func (r *EvalResult) Int64() (int64, error) {
	if r.err != nil {
		return 0, r.err
	}

	if !r.found || !r.enabled {
		return 0, nil
	}

	return strconv.ParseInt(r.rawValue, 10, 64)
}

func (r *EvalResult) UInt64() (uint64, error) {
	if r.err != nil {
		return 0, r.err
	}

	if !r.found || !r.enabled {
		return 0, nil
	}

	return strconv.ParseUint(r.rawValue, 10, 64)
}

func (r *EvalResult) Float64() (float64, error) {
	if r.err != nil {
		return 0, r.err
	}

	if !r.found || !r.enabled {
		return 0, nil
	}

	return strconv.ParseFloat(r.rawValue, 64)
}
