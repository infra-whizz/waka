package waka_diskman

import (
	"fmt"
	"os"
	"path"

	waka_layout "github.com/infra-whizz/waka/layout"
	wzlib_utils "github.com/infra-whizz/wzlib/utils"
)

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
type WkDiskManager struct {
	imglt   *waka_layout.WkImageLayout
	tmpRoot string
	tmpDir  string
}

// NewWkDiskManager creates a new disk manager instance.
func NewWkDiskManager(imglt *waka_layout.WkImageLayout) *WkDiskManager {
	dm := new(WkDiskManager)
	dm.imglt = imglt
	dm.tmpRoot = "/tmp"

	return dm
}

//SetWorkingDir sets where all output is set. Default is /tmp
func (dm *WkDiskManager) SetTempDir(wdir string) *WkDiskManager {
	if wdir != "" {
		dm.tmpRoot = wdir
	}
	return dm
}

// Create disk, according to the layout.
// Output goes to the same directory is where layout.conf as build/ subdir.
func (dm *WkDiskManager) Create() {
	diskPath := dm.getDiskPath()
	nfo, _ := os.Stat(diskPath)
	if nfo != nil {
		fmt.Fprintf(os.Stderr, "Error: file %s exist\n", nfo.Name())
		os.Exit(wzlib_utils.EX_GENERIC)
	}

	fmt.Println("Creating image", dm.getDiskName())
	if err := os.MkdirAll(path.Dir(diskPath), 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error: unable to create disk path %s: %s\n", path.Dir(diskPath), err.Error())
		os.Exit(wzlib_utils.EX_GENERIC)
	}

	dm.createRawDisk()
}

func (dm *WkDiskManager) Mount() {
	dm.createTemporarySpace()
}

// Bind system (/proc, /dev, /sys etc)
func (dm *WkDiskManager) Bind() {
}

// GetPartitionMountpoint where particular partition is mounted from the current image.
func (dm *WkDiskManager) GetPartitionMountpoint(partname string) string {
	return ""
}

// Umount disks (all)
func (dm *WkDiskManager) Umount() {
	dm.cleanup()
}
