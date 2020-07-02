package waka

import (
	"fmt"
	"os"

	waka_diskman "github.com/infra-whizz/waka/diskman"
	waka_layout "github.com/infra-whizz/waka/layout"
	wzlib_util "github.com/infra-whizz/wzlib/utils"
)

type Waka struct {
	imageLayout *waka_layout.WkImageLayout
	diskman     *waka_diskman.WkDiskManager
}

func NewWaka() *Waka {
	w := new(Waka)
	return w
}

// SetSchemaPath to the image description and layout schema
func (w *Waka) LoadSchema(schemaPath string) *Waka {
	w.imageLayout = waka_layout.NewWkImageLayout(schemaPath)
	return w
}

// Prepare environment, make disks, mount
func (w *Waka) prepare(force bool) {
	w.diskman = waka_diskman.NewWkDiskManager(w.imageLayout).SetTempDir(".")
	if force {
		if err := w.diskman.Remove(); err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err.Error())
			os.Exit(wzlib_util.EX_GENERIC)
		}
	}
	if err := w.diskman.Create(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err.Error())
		os.Exit(wzlib_util.EX_GENERIC)
	}
	if err := w.diskman.Loop(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err.Error())
		os.Exit(wzlib_util.EX_GENERIC)
	}
}

func (w *Waka) partitioning() {
	fmt.Println("Partitioning")
	if err := w.diskman.MakePartitions(); err != nil {
		fmt.Fprintln(os.Stderr, "Partitioning Error:", err.Error())
		os.Exit(wzlib_util.EX_GENERIC)
	}
}

// Bootstrap basic components, setup CMS system
func (w *Waka) bootstrap() {
}

// RunCMS on prepared image mount
func (w *Waka) runCMS() {
}

// Cleanup image, devices, data etc
func (w *Waka) cleanup() {
	if err := w.diskman.LoopOff(); err != nil {
		fmt.Fprintln(os.Stderr, "Removing loop device error:", err.Error())
		os.Exit(wzlib_util.EX_GENERIC)
	}
}

// Build the image
func (w *Waka) Build(force bool) {
	w.prepare(force)
	w.partitioning()
	w.bootstrap()
	w.runCMS()
	w.cleanup()
}
