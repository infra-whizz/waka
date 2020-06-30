package waka

type WakaMounter struct {
	imgname string
}

func NewWakaMounter(imgname string) *WakaMounter {
	devmap := new(WakaMounter)
	devmap.imgname = imgname
	return devmap
}

func (devmap *WakaMounter) Mount() error {
	return nil
}

func (devmap *WakaMounter) Umount() error {
	return nil
}

func (devmap *WakaMounter) Device() string {
	return ""
}

func (devmap *WakaMounter) Geometry() (map[string]interface{}, error) {
	return nil, nil
}
