package waka_layout

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	wzlib_utils "github.com/infra-whizz/wzlib/utils"
	"gopkg.in/yaml.v2"
)

type WkLayoutConfPartition struct {
	Size          int64
	Label         string
	PartitionCode string
	GUID          string
	Mountpoint    string
	Type          string

	device  string
	mounted string
}

// SetDevice to which partition is mounted on. This is mount-time use.
func (partition *WkLayoutConfPartition) SetDevice(device string) {
	partition.device = device
}

// UnsetDevice to which partition is mounted on. This is umount-time use.
func (partition *WkLayoutConfPartition) UnsetDevice() {
	partition.device = ""
}

// GetDevice of the partition
func (partition *WkLayoutConfPartition) GetDevice() string {
	return partition.device
}

// SetMountpoint to the partition
func (partition *WkLayoutConfPartition) SetMountpoint(mpt string) {
	partition.mounted = mpt
}

// RemoveMountpoint of the partition
func (partition *WkLayoutConfPartition) RemoveMountpoint() {
	partition.mounted = ""
}

// GetMountedAs location on the current layout
func (partition *WkLayoutConfPartition) GetMountedAs() string {
	return partition.mounted
}

type WkLayoutConf struct {
	Path           string
	Version        string
	Name           string
	Os             string
	Size           int64
	PackageManager string
	Packages       []string
	Partitions     []*WkLayoutConfPartition
	Repositories   []string
}

type WkImageLayout struct {
	conf                    *WkLayoutConf
	partitionTypeCodes      map[string][]string
	partitionFSTypeOverride map[string]string
	partitionMountpoints    map[string]string
}

func NewWkImageLayout(layoutPath string) *WkImageLayout {
	imglt := new(WkImageLayout)

	imglt.partitionMountpoints = make(map[string]string)
	imglt.partitionMountpoints["linux_bios"] = ""
	imglt.partitionMountpoints["linux_efi"] = "/boot/efi"
	imglt.partitionMountpoints["linux_boot"] = "/boot"
	imglt.partitionMountpoints["linux_root"] = "/"
	imglt.partitionMountpoints["linux_home"] = "/home"

	imglt.partitionFSTypeOverride = make(map[string]string)
	imglt.partitionFSTypeOverride["linux_bios"] = "vfat"
	imglt.partitionFSTypeOverride["linux_efi"] = "vfat"

	imglt.partitionTypeCodes = make(map[string][]string)
	// Linux partitions support
	imglt.partitionTypeCodes["linux_bios"] = []string{"ef02", "21686148-6449-6E6F-744E-656564454649"}    // BIOS boot partition
	imglt.partitionTypeCodes["linux_efi"] = []string{"ef00", "C12A7328-F81F-11D2-BA4B-00A0C93EC93B"}     // EFI system partition
	imglt.partitionTypeCodes["linux_boot"] = []string{"ea00", "BC13C2FF-59E6-4262-A352-B275FD6F7172"}    // Freedesktop $BOOT
	imglt.partitionTypeCodes["linux_swap"] = []string{"8200", "0657FD6D-A4AB-43C4-84E5-0933C84B4F4F"}    // Linux swap
	imglt.partitionTypeCodes["linux_root"] = []string{"8300", "0FC63DAF-8483-4772-8E79-3D69D8477DE4"}    // Linux filesystem
	imglt.partitionTypeCodes["linux_data"] = []string{"8300", "0FC63DAF-8483-4772-8E79-3D69D8477DE4"}    // Linux filesystem data
	imglt.partitionTypeCodes["linux_home"] = []string{"8302", "933AC7E1-2EB4-4F13-B844-0E14E2AEF915"}    // Linux /home
	imglt.partitionTypeCodes["linux_dmcrypt"] = []string{"8308", "7FFEC5C9-2D00-49B7-8941-3EA10A5586B7"} // Linux dm-crypt

	// Add your OS partitions below in format: <osname>_<type> :-)

	imglt.loadAndParse(layoutPath)

	return imglt
}

func (imglt *WkImageLayout) loadAndParse(schemaPath string) {
	var err error
	imglt.conf = new(WkLayoutConf)
	imglt.conf.Path, err = filepath.Abs(path.Join(schemaPath, "layout.conf"))
	if err != nil {
		panic(err)
	}

	lcpf, err := os.Stat(imglt.conf.Path)
	if err != nil {
		panic(err)
	}
	if lcpf.IsDir() {
		panic("Layout config cannot be a directory")
	}

	buff, err := ioutil.ReadFile(imglt.conf.Path)
	if err != nil {
		panic(err)
	}

	var layoutBuff map[string]interface{}

	if err := yaml.Unmarshal(buff, &layoutBuff); err != nil {
		panic(err)
	}

	imglt.setMainData(layoutBuff)
	imglt.setPackageList(layoutBuff)
	imglt.setPartitioningMap(layoutBuff)
	imglt.setRepoList(layoutBuff)
	imglt.verifyPartitionConfiguration()
}

// Set main data of the image
func (imglt *WkImageLayout) setMainData(buff map[string]interface{}) {
	name, ex := buff["name"]
	if !ex || name.(string) == "" {
		name = "untitled"
	}
	imglt.conf.Name = name.(string)

	version, ex := buff["version"]
	if !ex {
		fmt.Fprintln(os.Stderr, "Error: Version is not specified.")
		os.Exit(wzlib_utils.EX_GENERIC)
	}
	imglt.conf.Version = version.(string)

	osdata, ex := buff["os"]
	if !ex {
		fmt.Fprintln(os.Stderr, "Error: OS is not defined in the configuration")
		os.Exit(wzlib_utils.EX_GENERIC)
	}
	imglt.conf.Os = strings.ToLower(osdata.(string))

	imgSize, ex := buff["size"]
	if !ex {
		fmt.Fprintln(os.Stderr, "Error: Image size is not defined")
		os.Exit(wzlib_utils.EX_GENERIC)
	}
	imglt.conf.Size = int64(imgSize.(int))

	packman, ex := buff["package-manager"]
	if !ex {
		fmt.Fprintln(os.Stderr, "Error: unknown package manager")
		os.Exit(wzlib_utils.EX_GENERIC)
	}
	imglt.conf.PackageManager = packman.(string)
}

// Read repositories
func (imglt *WkImageLayout) setRepoList(buff map[string]interface{}) {
	imglt.conf.Repositories = make([]string, 0)
	repolist, ex := buff["repositories"]
	if !ex {
		fmt.Fprintln(os.Stderr, "No repository URLs specified.")
		os.Exit(wzlib_utils.EX_GENERIC)
	}
	for _, repoURL := range repolist.([]interface{}) {
		imglt.conf.Repositories = append(imglt.conf.Repositories, repoURL.(string))
	}
}

// Get package list to the configuration structure
func (imglt *WkImageLayout) setPackageList(buff map[string]interface{}) {
	imglt.conf.Packages = make([]string, 0)
	pkglist, ex := buff["packages"]
	if !ex {
		fmt.Fprintln(os.Stderr, "Error: no packages defined for the image")
		os.Exit(wzlib_utils.EX_GENERIC)
	}
	for _, pkgname := range pkglist.([]interface{}) {
		imglt.conf.Packages = append(imglt.conf.Packages, pkgname.(string))
	}
}

func (imglt *WkImageLayout) getPartitionTypeCode(ostype string, partname string) []string {
	if partname == "data" {
		partname = "root"
	}
	value, ex := imglt.partitionTypeCodes[fmt.Sprintf("%s_%s", ostype, partname)]
	if !ex {
		return []string{"", ""}
	}
	return value
}

// Check if user did not wrote a nonsense partition FS type. For example, EFI should be vfat.
func (imglt *WkImageLayout) overridePartitionFSType(ostype string, partname string, fstype string) string {
	value, ex := imglt.partitionFSTypeOverride[fmt.Sprintf("%s_%s", ostype, partname)]
	if !ex {
		return fstype
	}
	return value
}

// Get default mountpoint
func (imglt *WkImageLayout) getDefaultMountpoint(ostype string, partname string, mountpoint string) string {
	mpt, ex := imglt.partitionMountpoints[fmt.Sprintf("%s_%s", ostype, partname)]
	if !ex {
		return mountpoint
	}
	return mpt
}

// Get partition map and convert data
func (imglt *WkImageLayout) setPartitioningMap(buff map[string]interface{}) {
	// TODO: probably this should be automatically
	// acquired from the current machine, in case not defined?

	imglt.conf.Partitions = make([]*WkLayoutConfPartition, 0)
	partmap, ex := buff["partitions"]
	if !ex {
		fmt.Fprintln(os.Stderr, "Error: no partitions defined for this image")
		os.Exit(wzlib_utils.EX_GENERIC)
	}
	for _, partItf := range partmap.([]interface{}) {
		for partType := range partItf.(map[interface{}]interface{}) {
			partCode := imglt.getPartitionTypeCode(imglt.conf.Os, partType.(string))
			if partCode[0] == "" {
				fmt.Fprintln(os.Stderr, "Error: unsupported partition type:", partType)
				os.Exit(wzlib_utils.EX_GENERIC)
			}

			partMeta := partItf.(map[interface{}]interface{})[partType].(map[interface{}]interface{})
			partFSType := partMeta["type"]
			if partFSType == nil {
				partFSType = imglt.overridePartitionFSType(imglt.conf.Os, partType.(string), "ext4")
			} else {
				partFSType = imglt.overridePartitionFSType(imglt.conf.Os, partType.(string), partFSType.(string)) // FS might be not supported
			}

			partMountpoint := partMeta["mountpoint"]
			if partMountpoint == nil {
				partMountpoint = imglt.getDefaultMountpoint(imglt.conf.Os, partType.(string), "")
				if partMountpoint == "" {
					fmt.Fprintf(os.Stderr, "Error: no mountpoint specified for partition %s\n", partType)
				}
			} else {
				partMountpoint = imglt.getDefaultMountpoint(imglt.conf.Os, partType.(string), partMountpoint.(string))
			}

			partSize := partMeta["size"]
			if partSize == nil {
				partSize = 0
			}

			partLabel, ex := partMeta["label"]
			if !ex {
				partLabel = fmt.Sprintf("Partition %s", partType)
			}

			imglt.conf.Partitions = append(imglt.conf.Partitions,
				&WkLayoutConfPartition{
					PartitionCode: partCode[0],
					GUID:          partCode[1],
					Size:          int64(partSize.(int)),
					Label:         partLabel.(string),
					Type:          partFSType.(string),
					Mountpoint:    partMountpoint.(string),
				})
		}
	}
}

func (imglt *WkImageLayout) verifyPartitionConfiguration() {
	partitions := len(imglt.conf.Partitions)
	for idx, partition := range imglt.conf.Partitions {
		if partition.Size == 0 && idx+1 < partitions {
			fmt.Fprintf(os.Stderr, "Error: partition %s wants the rest of the disk space, but there are more partitions...\n", partition.Label)
			os.Exit(wzlib_utils.EX_GENERIC)
		}
	}
}

// GetConfig of the image layout
func (imglt *WkImageLayout) GetConfig() *WkLayoutConf {
	return imglt.conf
}
