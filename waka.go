package waka

import (
	waka_diskman "github.com/infra-whizz/waka/diskman"
	waka_layout "github.com/infra-whizz/waka/layout"
	wzlib_logger "github.com/infra-whizz/wzlib/logger"
)

type Waka struct {
	imageLayout *waka_layout.WkImageLayout
	diskman     *waka_diskman.WkDiskManager
	imageSetup  *WakaImageSetup

	outputPath    string
	cleanupMounts bool

	wzlib_logger.WzLogger
}

func NewWaka() *Waka {
	w := new(Waka)
	w.cleanupMounts = true
	return w
}

// SetCleanupOnExit turns off or on cleaning up (if not crashed prematurely)
func (w *Waka) SetCleanupOnExit(cleanup bool) *Waka {
	w.cleanupMounts = cleanup
	return w
}

// SetBuildOutput other than $LAYOUT/build
func (w *Waka) SetBuildOutput(outputPath string) *Waka {
	w.outputPath = outputPath
	return w
}

// SetSchemaPath to the image description and layout schema
func (w *Waka) LoadSchema(schemaPath string) *Waka {
	w.imageLayout = waka_layout.NewWkImageLayout(schemaPath)
	return w
}

// Build the image
func (w *Waka) Build(force bool) {
	w.prepare(force)
	w.partitioning()
	w.format()
	w.mount()
	w.preProvision()
	w.runCMS()
	w.postProvision()

	if w.cleanupMounts {
		w.cleanup()
	} else {
		w.GetLogger().Info("Leaving mounts untouched for debugging purposes")
	}
}
