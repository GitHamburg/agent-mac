// +build darwin
// +build !cgo

package disk

import "../internal/common"

func IOCounters(names ...string) (map[string]IOCountersStat, error) {
	return nil, common.ErrNotImplementedError
}
