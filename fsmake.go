package waka

/*
	Prepares the image file system.
	The whole purpose of the fsmake is just create an image,
	formatted with the supported filesystem.

	Design of this class is subject to change in a future.
*/

type WakaFSMake struct{}

func NewWakaFSMake() *WakaFSMake {
	fm := new(WakaFSMake)
	return fm
}

func (fm *WakaFSMake) Format(fstype string) {
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
	}
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
