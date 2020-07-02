package waka_diskman

import (
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

func (dm *WkDiskManager) createRawDisk() {
	cmd, err := wzlib_subprocess.BufferedExec("dd", "if=/dev/zero",
		fmt.Sprintf("of=%s", dm.getDiskPath()), "bs=1024K",
		fmt.Sprintf("seek=%d", dm.getDiskLayoutConfig().Size), "count=0")
	if err != nil {
		fmt.Println("Call init error:", err.Error())
	}
	fmt.Println(cmd.StdoutString())
	cmd.Wait()
}

func (dm *WkDiskManager) cleanup() {
	dm.tmpDir = ""
}

func (dm *WkDiskManager) createTemporarySpace() error {
	var err error
	dm.tmpDir, err = ioutil.TempDir(dm.tmpRoot, "waka-build") // Config!
	if err != nil {
		return err
	}
	return nil
}