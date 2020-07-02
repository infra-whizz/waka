package waka_diskman

/*
	Disk Manager.
	Purpose:
	  - Create disk images in general
	  - Partitioning
	  - Formatting partitions
	  - Mounting/umounting partitions
	  - Accessing partitions
	  - Converting formats
*/

// WkDiskManager object
type WkDiskManager struct{}

// NewWkDiskManager creates a new disk manager instance.
func NewWkDiskManager() *WkDiskManager {
	dm := new(WkDiskManager)
	return dm
}

func (dm *WkDiskManager) Create() {}
func (dm *WkDiskManager) Mount()  {}

// Bind system (/proc, /dev, /sys etc)
func (dm *WkDiskManager) Bind()                        {}
func (dm *WkDiskManager) GetPartition(partname string) {}
func (dm *WkDiskManager) Umount()                      {}
