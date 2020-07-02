package waka_parted

import (
	waka_layout "github.com/infra-whizz/waka/layout"
)

type WakaPartitioner interface {
	Flush()
	Create(partition *waka_layout.WkLayoutConfPartition) error
	GetDiskDevice() string
}
