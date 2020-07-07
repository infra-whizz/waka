package waka

import (
	"fmt"
	"path"

	"github.com/isbm/go-shutil"

	waka_diskman "github.com/infra-whizz/waka/diskman"
	waka_layout "github.com/infra-whizz/waka/layout"
)

type WakaImageSetup struct {
	diskman               *waka_diskman.WkDiskManager
	layout                *waka_layout.WkImageLayout
	preInstallRootfsPath  string
	postInstallRootfsPath string
}

func NewWakaImageSetup(diskman *waka_diskman.WkDiskManager, layout *waka_layout.WkImageLayout) *WakaImageSetup {
	ims := new(WakaImageSetup)
	ims.diskman = diskman
	ims.layout = layout
	ims.preInstallRootfsPath = path.Join(path.Base(ims.layout.GetConfig().Path), "rootfs", "pre-install")
	ims.postInstallRootfsPath = path.Join(path.Base(ims.layout.GetConfig().Path), "rootfs", "post-install")
	return ims
}

func (ims *WakaImageSetup) getRootPartitionMountpoint() (string, error) {
	partition := ims.layout.GetPartitionByPType("root")
	if partition == nil {
		return "", fmt.Errorf("Root partition is not found")
	}

	if partition.GetMountedAs() == "" {
		return "", fmt.Errorf("Root partition %s is not mounted", partition.GetDevice())
	}

	return partition.GetMountedAs(), nil
}

// CopyPayload of all bits required to run inside the image (chroot)
func (ims *WakaImageSetup) CopyPayload() error {
	dst, err := ims.getRootPartitionMountpoint()
	if err != nil {
		return err
	}
	src := path.Join(path.Dir(ims.layout.GetConfig().Path), "collection")
	dst = path.Join(dst, "tmp", "waka-build", "collection")
	return shutil.CopyTree(src, dst, &shutil.CopyTreeOptions{
		Symlinks:               false,
		Ignore:                 nil,
		OverlayMode:            true,
		CopyFunction:           shutil.Copy,
		IgnoreDanglingSymlinks: false})
}

// Copy root FS.
func (ims *WakaImageSetup) copyRootFs(schema string) error {
	dst, err := ims.getRootPartitionMountpoint()
	if err != nil {
		return err
	}
	src := path.Join(path.Dir(ims.layout.GetConfig().Path), "rootfs", schema)
	return shutil.CopyTree(src, dst, &shutil.CopyTreeOptions{
		Symlinks:               false,
		Ignore:                 nil,
		OverlayMode:            true,
		CopyFunction:           shutil.Copy,
		IgnoreDanglingSymlinks: false})
}

// CopyPostRootFs is copying the entire root FS from "post-install" directory
func (ims *WakaImageSetup) CopyPostRootFs() error {
	return ims.copyRootFs("post-install")
}

// CopyPreRootFs is copying the entire root FS from "pre-install" directory
func (ims *WakaImageSetup) CopyPreRootFs() error {
	return ims.copyRootFs("pre-install")
}
