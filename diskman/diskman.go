package waka_diskman

import (
	"fmt"
	"os"
	"path"

	waka_layout "github.com/infra-whizz/waka/layout"
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

// Remove disk
func (dm *WkDiskManager) Remove() error {
	diskPath := dm.getDiskPath()
	nfo, _ := os.Stat(diskPath)
	if nfo != nil {
		return os.Remove(diskPath)
	}
	return nil
}

// Create disk, according to the layout.
// Output goes to the same directory is where layout.conf as build/ subdir.
func (dm *WkDiskManager) Create() error {
	diskPath := dm.getDiskPath()
	nfo, _ := os.Stat(diskPath)
	if nfo != nil {
		return fmt.Errorf("File %s already exists", nfo.Name())
	}

	fmt.Println("Creating image", dm.getDiskName())
	if err := os.MkdirAll(path.Dir(diskPath), 0755); err != nil {
		return fmt.Errorf("Unable to create disk path %s: %s", path.Dir(diskPath), err.Error())
	}

	return dm.createRawDisk()
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
