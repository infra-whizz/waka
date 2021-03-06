package waka

import (
	"encoding/json"
	"fmt"
	"path"
	"strings"

	waka_diskman "github.com/infra-whizz/waka/diskman"
	wzlib_subprocess "github.com/infra-whizz/wzlib/subprocess"
)

// Prepare environment, make disks, mount
func (w *Waka) prepare(force bool) {
	w.diskman = waka_diskman.NewWkDiskManager(w.imageLayout).SetTempDir(".").SetBuildOutput(w.outputPath)
	w.imageSetup = NewWakaImageSetup(w.diskman, w.imageLayout)
	if force {
		ExitOnError(w.diskman.Remove())
	}
	ExitOnError(w.diskman.Create())
	ExitOnError(w.diskman.Loop())
}

func (w *Waka) partitioning() {
	fmt.Println("Partitioning")
	ExitOnErrorPreamble(w.diskman.MakePartitions(), "Partitioning did not suceeded -")
}

// Remount all partitions that has been just created
func (w *Waka) mount() {
	ExitOnErrorPreamble(w.diskman.Loop(), "Remount failed -")
}

// Format all partitions of the current device
func (w *Waka) format() {
	ExitOnErrorPreamble(w.diskman.FormatAllPartitions(), "Formatting failure -")
}

func (w *Waka) listRootFs() {
	root, _ := w.imageSetup.getRootPartitionMountpoint()
	cmd, _ := wzlib_subprocess.BufferedExec("find", root)
	content := cmd.StdoutString()
	cmd.Wait()
	fmt.Println("---------\nContent:\n", content, "---------")
}

// Bootstrap basic components, setup CMS system
func (w *Waka) preProvision() {
	ExitOnErrorPreamble(w.imageSetup.CopyPayload(), "Unable to copy payload")
	ExitOnErrorPreamble(w.imageSetup.CopyPreRootFs(), "Unable to copy pre-install rootfs layout")
	ExitOnErrorPreamble(w.imageSetup.CopyWhizzLocal(), "Unable to prepare whizz daemon")
}

func (w *Waka) postProvision() {
	ExitOnErrorPreamble(w.imageSetup.CopyPostRootFs(), "Unable to copy post-install rootfs layout")
}

// RunCMS on prepared image mount
func (w *Waka) runCMS() {
	rootfs, _ := w.imageSetup.getRootPartitionMountpoint()
	wzdBinPath := path.Join(rootfs, "tmp", ".waka", "bin", "wzd")
	collection := path.Join(rootfs, "tmp", ".waka", "collection")

	caller := wzlib_subprocess.NewEnvBufferedExec().SetEnv("WAKA_MOUNT", rootfs)
	cmd, err := caller.Exec(wzdBinPath, "-f", "json", "local", "-r", rootfs, "-d", collection, "-s", "init")
	ExitOnErrorPreamble(err, "Unable to run Whizz locally")
	stout := strings.Split(cmd.StdoutString(), "\n")
	sterr := strings.Split(cmd.StderrString(), "\n")
	cmd.Wait()

	// Each line expected to be a JSON. Otherise passed "as is".
	var buff map[string]string
	for _, out := range [][]string{stout, sterr} {
		for _, line := range out {
			line = strings.TrimSpace(line)
			err = json.Unmarshal([]byte(line), &buff)
			if err != nil {
				w.GetLogger().Info(line)
			} else {
				switch buff["level"] {
				case "debug":
					w.GetLogger().Debug(buff["msg"])
				case "error":
					w.GetLogger().Error(buff["msg"])
				case "warning":
					w.GetLogger().Warning(buff["msg"])
				default:
					w.GetLogger().Info(buff["msg"])
				}
			}
		}
	}
}

// Cleanup image, devices, data etc
func (w *Waka) cleanup() {
	// w.listRootFs() // debug
	ExitOnErrorPreamble(w.diskman.Umount(), "Unable to umount devices")
	ExitOnErrorPreamble(w.diskman.LoopOff(),
		fmt.Sprintf("Umount loop device %s failed -", w.diskman.GetDiskImageDevice()))
	ExitOnErrorPreamble(w.diskman.Cleanup(), "Unable to cleanup mountpoints")
}
