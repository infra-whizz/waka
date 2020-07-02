package waka

import (
	"fmt"
)

/*
	Prepares the image file system.
	The whole purpose of the fsmake is just create an image,
	formatted with the supported filesystem.

	Design of this class is subject to change in a future.
*/

type WakaDiskFormatter struct {
	sizeMb      int64
	rawFilename string
}

func NewWakaDiskFormatter() *WakaDiskFormatter {
	fm := new(WakaDiskFormatter)
	return fm
}

func (fm *WakaDiskFormatter) SetSizeMb(size int64) *WakaDiskFormatter {
	fm.sizeMb = size
	return fm
}

func (fm *WakaDiskFormatter) SetOutputFile(fname string) *WakaDiskFormatter {
	fm.rawFilename = fname
	return fm
}

func (fm *WakaDiskFormatter) Format(fstype string) error {
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
func (fm *WakaDiskFormatter) ext2() {
	fm.createRawImage()
}

// Format image with ext3 FS
func (fm *WakaDiskFormatter) ext3() {
	fm.createRawImage()
}

// Format image with ext4 FS
func (fm *WakaDiskFormatter) ext4() {
	fmt.Println("Formatting raw image with ext4")
	fm.createRawImage()
}

// Format image with xfs FS
func (fm *WakaDiskFormatter) xfs() {
	fm.createRawImage()
}

// Format image with cramfs FS
func (fm *WakaDiskFormatter) cramfs() {
	fm.createRawImage()
}

func (fm *WakaDiskFormatter) createRawImage() {
}
