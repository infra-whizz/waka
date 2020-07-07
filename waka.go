package waka

import (
	waka_diskman "github.com/infra-whizz/waka/diskman"
	waka_layout "github.com/infra-whizz/waka/layout"
)

type Waka struct {
	imageLayout *waka_layout.WkImageLayout
	diskman     *waka_diskman.WkDiskManager
	imageSetup  *WakaImageSetup
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

// Build the image
func (w *Waka) Build(force bool) {
	w.prepare(force)
	w.partitioning()
	w.format()
	w.mount()
	w.bootstrap()
	w.runCMS()
	w.postProvision()
	w.cleanup()
}
