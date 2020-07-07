package waka

import (
	"fmt"

	wzlib_subprocess "github.com/infra-whizz/wzlib/subprocess"

	waka_diskman "github.com/infra-whizz/waka/diskman"
)

// Prepare environment, make disks, mount
func (w *Waka) prepare(force bool) {
	w.diskman = waka_diskman.NewWkDiskManager(w.imageLayout).SetTempDir(".")
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

// Bootstrap basic components, setup CMS system
func (w *Waka) bootstrap() {
	ExitOnErrorPreamble(w.imageSetup.CopyPayload(), "Unable to copy payload")
	ExitOnErrorPreamble(w.imageSetup.CopyPreRootFs(), "Unable to copy pre-install rootfs layout")
}

func (w *Waka) postProvision() {
	ExitOnErrorPreamble(w.imageSetup.CopyPostRootFs(), "Unable to copy post-install rootfs layout")
}

// RunCMS on prepared image mount
func (w *Waka) runCMS() {
}

// Cleanup image, devices, data etc
func (w *Waka) cleanup() {
	ExitOnErrorPreamble(w.diskman.Umount(), "Unable to umount devices")
	ExitOnErrorPreamble(w.diskman.LoopOff(),
		fmt.Sprintf("Umount loop device %s failed -", w.diskman.GetDiskImageDevice()))
	ExitOnErrorPreamble(w.diskman.Cleanup(), "Unable to cleanup mountpoints")
}
