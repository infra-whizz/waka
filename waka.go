package waka

import (
	"fmt"
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
	err := NewWakaFSMake().SetSizeMb(500).Format("ext4")
	if err != nil {
		fmt.Println(err.Error())
	}
	os.Exit(wzlib_util.EX_GENERIC)
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
