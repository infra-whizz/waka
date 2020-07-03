package waka_diskman

import (
	"fmt"
	"os"
	"path"

	waka_parted "github.com/infra-whizz/waka/diskman/parted"
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
	parted  waka_parted.WakaPartitioner
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

// MakePartitions on the disk
func (dm *WkDiskManager) MakePartitions() error {
	for _, partMeta := range dm.getDiskLayoutConfig().Partitions {
		if err := dm.addPartition(partMeta); err != nil {
			return err
		}
	}
	if err := dm.updateMountedDeviceMap(); err != nil {
		fmt.Fprintln(os.Stderr, "ERROR:", err.Error())
		os.Exit(wzlib_utils.EX_GENERIC)
	}
	return nil
}

// FormatAllPartitions that are present in the disk device, according to the configured layout.
func (dm *WkDiskManager) FormatAllPartitions() error {
	for _, partition := range dm.imglt.GetConfig().Partitions {
		if err := dm.formatPartition(partition); err != nil {
			fmt.Println("FAILED:", err.Error())
			return err
		} else {
			fmt.Println("FORMATTED")
		}
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

// Loop created image
func (dm *WkDiskManager) Loop() error {
	if err := dm.loopDiskImage(); err != nil {
		return err
	}
	diskDevice, err := dm.getDiskImageDevice()
	if err != nil {
		return err
	}
	dm.parted = waka_parted.NewWakaPartitionerGPT(diskDevice)
	if err := dm.updateMountedDeviceMap(); err != nil {
		// No partition table
		fmt.Println("ERROR update mounted device map:", err.Error())
	} else {
		// Mount everything
		for _, partition := range dm.imglt.GetConfig().Partitions {
			fmt.Println("Mounted partition", partition.GetDevice())
		}
	}
	fmt.Println("DEBUG: Mounted as", dm.parted.GetDiskDevice())
	return nil
}

// LoopOff turns looped image off
func (dm *WkDiskManager) LoopOff() error {
	return dm.unLoopDiskImage()
}

func (dm *WkDiskManager) Mount() error {
	return dm.createTemporarySpace()
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
	dm.flushDeviceMap()
	dm.cleanup()
}

// GetDiskImageDevice corresponding to the mounted disk loop
func (dm *WkDiskManager) GetDiskImageDevice() string {
	dev, _ := dm.getDiskImageDevice()
	return dev
}
