package waka_diskman

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"sort"

	waka_layout "github.com/infra-whizz/waka/layout"
	wzlib_subprocess "github.com/infra-whizz/wzlib/subprocess"
)

// Return disk layout confnigurtaion
func (dm *WkDiskManager) getDiskLayoutConfig() *waka_layout.WkLayoutConf {
	return dm.imglt.GetConfig()
}

// Get disk name
func (dm *WkDiskManager) getDiskName() string {
	return fmt.Sprintf("%s-%s.raw", dm.getDiskLayoutConfig().Name, dm.getDiskLayoutConfig().Version)
}

// Return disk path
func (dm *WkDiskManager) getDiskPath() string {
	return path.Join(path.Dir(dm.getDiskLayoutConfig().Path), "build", dm.getDiskName())
}

func (dm *WkDiskManager) createRawDisk() error {
	cmd, err := wzlib_subprocess.BufferedExec("dd", "if=/dev/zero",
		fmt.Sprintf("of=%s", dm.getDiskPath()), "bs=1024K",
		fmt.Sprintf("seek=%d", dm.getDiskLayoutConfig().Size), "count=0")
	if err != nil {
		return err
	}
	out := cmd.StderrString()
	fmt.Println("DEBUG:", out)

	return cmd.Wait()
}

// Connect disk image to the loop
func (dm *WkDiskManager) loopDiskImage() error {
	if dm.parted == nil {
		cmd, err := wzlib_subprocess.BufferedExec("losetup", "-fP", dm.getDiskPath())
		if err != nil {
			return err
		}
		out := cmd.StdoutString()
		fmt.Println(out)
		return cmd.Wait()
	}
	return nil
}

// Disconnect disk image from the loop
func (dm *WkDiskManager) unLoopDiskImage() error {
	if dm.parted != nil {
		cmd, err := wzlib_subprocess.BufferedExec("losetup", "-d", dm.parted.GetDiskDevice())
		if err != nil {
			return err
		}
		out := cmd.StdoutString()
		if err := cmd.Wait(); err != nil {
			return err
		}
		fmt.Println("DEBUG:", out)
		dm.parted = nil
	}
	return nil
}

func (dm *WkDiskManager) getDiskImageDevice() (string, error) {
	cmd, err := wzlib_subprocess.BufferedExec("losetup", "-lJ")
	if err != nil {
		return "", err
	}
	var buff map[string][]map[string]interface{}
	if err := json.Unmarshal([]byte(cmd.StdoutString()), &buff); err != nil {
		return "", err
	}

	if err := cmd.Wait(); err != nil {
		return "", err
	}

	for _, deviceLoopMap := range buff["loopdevices"] {
		if deviceLoopMap["back-file"].(string) == dm.getDiskPath() {
			return deviceLoopMap["name"].(string), nil
		}
	}

	return "", fmt.Errorf("No device found for %s disk", dm.getDiskName())
}

func (dm *WkDiskManager) addPartition(partition *waka_layout.WkLayoutConfPartition) error {
	return dm.parted.Create(partition)
}

func (dm *WkDiskManager) updateMountedDeviceMap() error {
	device, err := dm.getDiskImageDevice()
	if err != nil {
		return err
	}

	cmd, err := wzlib_subprocess.BufferedExec("sfdisk", "-lJ", device)
	if err != nil {
		return err
	}
	out := cmd.StdoutString()
	if err := cmd.Wait(); err != nil {
		return err
	}
	var partTable map[string]map[string]interface{}
	if err := json.Unmarshal([]byte(out), &partTable); err != nil {
		return err
	}

	for partIdx, partInfo := range partTable["partitiontable"]["partitions"].([]interface{}) {
		dm.imglt.GetConfig().Partitions[partIdx].SetDevice(partInfo.(map[string]interface{})["node"].(string))
	}
	return nil
}

func (dm *WkDiskManager) flushDeviceMap() {
	for _, partition := range dm.imglt.GetConfig().Partitions {
		partition.UnsetDevice()
	}
}

// Mount partition
func (dm *WkDiskManager) mountPartition(partition *waka_layout.WkLayoutConfPartition) error {
	return dm._partitionMounter(partition, true)
}

// Umount partition
func (dm *WkDiskManager) umountPartition(partition *waka_layout.WkLayoutConfPartition) error {
	return dm._partitionMounter(partition, false)
}

// Partition mounter
func (dm *WkDiskManager) _partitionMounter(partition *waka_layout.WkLayoutConfPartition, mount bool) error {
	if dm.tmpDir == "" {
		return fmt.Errorf("Temp directory is not defined")
	}
	mountpoint := path.Join(dm.tmpDir, partition.Mountpoint)

	var cmd *wzlib_subprocess.BufferedCmd
	var err error
	if mount {
		fmt.Println("Mounting partition", partition.GetDevice(), "to", mountpoint)
		if err := os.MkdirAll(mountpoint, 0700); err != nil {
			fmt.Println("ERROR creating directory:", err.Error())
		}
		cmd, err = wzlib_subprocess.BufferedExec("mount", partition.GetDevice(), mountpoint)
	} else {
		fmt.Println("Umounting partition", partition.GetDevice(), "from", mountpoint)
		cmd, err = wzlib_subprocess.BufferedExec("umount", partition.GetDevice())
	}
	sout := cmd.StdoutString()
	serr := cmd.StderrString()

	if err != nil {
		fmt.Println("DEBUG ERROR:", serr)
		return err
	}

	err = cmd.Wait()
	if err != nil {
		fmt.Println("DEBUG ERROR:", serr)
		return err
	}
	fmt.Println("DEBUG:", sout)

	if mount {
		partition.SetMountpoint(mountpoint)
	} else {
		partition.RemoveMountpoint()
	}

	return nil
}

// Format given partition
func (dm *WkDiskManager) formatPartition(partition *waka_layout.WkLayoutConfPartition) error {
	fmt.Println("Formatting partition:", partition.Label, "available at", partition.GetDevice(), "as", partition.FsType)
	var args []string
	switch partition.FsType {
	case "vfat":
		args = []string{"-t", partition.FsType, "-F", "32", partition.GetDevice()}
	case "ext2", "ext3", "ext4", "xfs":
		// XFS support is minimal default at the moment
		args = []string{"-t", partition.FsType, "-L", partition.Label, partition.GetDevice()}
	default:
		return fmt.Errorf("Unsupported filesystem: %s", partition.FsType)
	}

	cmd, err := wzlib_subprocess.BufferedExec("mkfs", args...)
	out := cmd.StdoutString()
	sterr := cmd.StderrString()

	if err != nil {
		fmt.Println("DEBUG ERROR:", sterr)
		return err
	}

	if err := cmd.Wait(); err != nil {
		fmt.Println("DEBUG ERROR:", sterr)
		return err
	}

	fmt.Println("DEBUG:", out)
	return nil
}

// Get mount points
func (dm *WkDiskManager) getOrderedMountPoints(reversed bool) []string {
	mountPoints := make([]string, 0)
	for _, partition := range dm.imglt.GetConfig().Partitions {
		mountPoints = append(mountPoints, partition.Mountpoint)
	}
	if reversed {
		sort.Sort(sort.Reverse(sort.StringSlice(mountPoints)))
	} else {
		sort.Strings(mountPoints)
	}
	return mountPoints
}
