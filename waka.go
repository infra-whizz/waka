package waka

import (
	"os"

	wzlib_util "github.com/infra-whizz/wzlib/utils"
	"github.com/isbm/go-nanoconf"
)

type Waka struct {
	conf *nanoconf.Config
}

func NewWaka() *Waka {
	w := new(Waka)
	return w
}

// SetSchema for the image to be built
func (w *Waka) SetSchemaConfig(conf *nanoconf.Config) *Waka {
	w.conf = conf
	return w
}

// Prepare environment, make disks, mount
func (w *Waka) prepare() {
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
	w.bootstrap()
	w.runCMS()
	w.cleanup()
}
