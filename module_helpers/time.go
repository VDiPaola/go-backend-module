package module_helpers

import "time"

func NowJS() int64 {
	return time.Now().UTC().UnixNano() / 1e6
}

func UnixNanoToJS(un int64) int64 {
	return un / 1e6
}
