package waka

import (
	"fmt"

	wzlib_subprocess "github.com/infra-whizz/wzlib/subprocess"
)

/*
	Prepares the image file system.
	The whole purpose of the fsmake is just create an image,
	formatted with the supported filesystem.

	Design of this class is subject to change in a future.
*/

type WakaFSMake struct {
	sizeMb      int64
	rawFilename string
}

func NewWakaFSMake() *WakaFSMake {
	fm := new(WakaFSMake)
	return fm
}

func (fm *WakaFSMake) SetSizeMb(size int64) *WakaFSMake {
	fm.sizeMb = size
	return fm
}

func (fm *WakaFSMake) SetOutputFile(fname string) *WakaFSMake {
	fm.rawFilename = fname
	return fm
}

func (fm *WakaFSMake) Format(fstype string) error {
	if fm.sizeMb == 0 {
		return fmt.Errorf("Unable to format image with %s to a zero size", fstype)
	}
	switch fstype {
	case "ext4":
		fm.ext4()
	case "xfs":
		fm.xfs()
	case "ext3":
		fm.ext3()
	case "ext2":
		fm.ext2()
	case "cramfs":
		fm.cramfs()
	default:
		return fmt.Errorf("Unknown FS type: %s", fstype)
	}
	return nil
}

// Format image with ext2 FS
func (fm *WakaFSMake) ext2() {
	fm.createRawImage()
}

// Format image with ext3 FS
func (fm *WakaFSMake) ext3() {
	fm.createRawImage()
}

// Format image with ext4 FS
func (fm *WakaFSMake) ext4() {
	fmt.Println("Formatting raw image with ext4")
	fm.createRawImage()
}

// Format image with xfs FS
func (fm *WakaFSMake) xfs() {
	fm.createRawImage()
}

// Format image with cramfs FS
func (fm *WakaFSMake) cramfs() {
	fm.createRawImage()
}

func (fm *WakaFSMake) createRawImage() {
}
