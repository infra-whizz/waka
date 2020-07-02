package waka_diskman

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"

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

func (dm *WkDiskManager) cleanup() {
	dm.tmpDir = ""
}

func (dm *WkDiskManager) createTemporarySpace() error {
	var err error
	dm.tmpDir, err = ioutil.TempDir(dm.tmpRoot, "waka-build")
	if err != nil {
		return err
	}
	return nil
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
