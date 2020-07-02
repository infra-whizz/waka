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
func (w *Waka) prepare() {
	disk := waka_diskman.NewWkDiskManager(w.imageLayout).SetTempDir(".")
	if err := disk.Create(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err.Error())
		os.Exit(wzlib_util.EX_GENERIC)
	}
}

func (w *Waka) partitioning() {
	os.Exit(wzlib_util.EX_GENERIC)
}

// Bootstrap basic components, setup CMS system
func (w *Waka) bootstrap() {
	os.Exit(wzlib_util.EX_GENERIC)
}

// RunCMS on prepared image mount
func (w *Waka) runCMS() {
	os.Exit(wzlib_util.EX_GENERIC)
}

// Cleanup image, devices, data etc
func (w *Waka) cleanup() {
	os.Exit(wzlib_util.EX_GENERIC)
}

// Build the image
func (w *Waka) Build() {
	w.prepare()
	w.partitioning()
	w.bootstrap()
	w.runCMS()
	w.cleanup()
}
