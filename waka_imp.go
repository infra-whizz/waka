package waka

import (
	"fmt"

	waka_diskman "github.com/infra-whizz/waka/diskman"
)

// Prepare environment, make disks, mount
func (w *Waka) prepare(force bool) {
	w.diskman = waka_diskman.NewWkDiskManager(w.imageLayout).SetTempDir(".")
	if force {
		ExitOnError(w.diskman.Remove())
	}
	ExitOnError(w.diskman.Create())
	ExitOnError(w.diskman.Loop())
}

func (w *Waka) partitioning() {
	fmt.Println("Partitioning")
	ExitOnErrorPreamble(w.diskman.MakePartitions(), "Partitioning failed")
}

// Remount all partitions that has been just created
func (w *Waka) remount() {
	w.cleanup()
	ExitOnErrorPreamble(w.diskman.Loop(), "Remount failed -")
}

// Format all partitions of the current device
func (w *Waka) format() {

}

// Bootstrap basic components, setup CMS system
func (w *Waka) bootstrap() {
}

// RunCMS on prepared image mount
func (w *Waka) runCMS() {
}

// Cleanup image, devices, data etc
func (w *Waka) cleanup() {
	ExitOnErrorPreamble(w.diskman.LoopOff(),
		fmt.Sprintf("Umount loop device %s failed -", w.diskman.GetDiskImageDevice()))
}
