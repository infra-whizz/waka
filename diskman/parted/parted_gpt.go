package waka

import (
	"fmt"

	wzlib_subprocess "github.com/infra-whizz/wzlib/subprocess"
)

type WakaPartitionerGPT struct {
	parted  string
	device  string
	ostype  string
	partmap map[string]string
}

func NewWakaPartitionerGPT(ostype string, device string) *WakaPartitionerGPT {
	gpt := new(WakaPartitionerGPT)
	gpt.parted = "sgdisk"
	gpt.device = device
	gpt.ostype = ostype

	// To view all: sgdisk -L
	gpt.partmap = make(map[string]string)
	gpt.partmap["bios"] = "ef02"          // BIOS boot partition
	gpt.partmap["efi"] = "ef00"           // EFI system partition
	gpt.partmap["linux_boot"] = "ea00"    // Freedesktop $BOOT
	gpt.partmap["linux_swap"] = "8200"    // Linux swap
	gpt.partmap["linux_root"] = "8300"    // Linux filesystem
	gpt.partmap["linux_home"] = "8302"    // Linux /home
	gpt.partmap["linux_dmcrypt"] = "8308" // Linux dm-crypt

	return gpt
}

// Flush everything from the image
func (gpt *WakaPartitionerGPT) Flush() {
	fmt.Println("Creating a new GPT structure on device", gpt.device)
	cmd, err := wzlib_subprocess.BufferedExec(gpt.parted, "-og", gpt.device)
	if err != nil {
		panic(err)
	}
	fmt.Println(cmd.StdoutString())
}

// Get partition code type (see partmap)
func (gpt *WakaPartitionerGPT) getTypeName(name string) (string, error) {
	val, ex := gpt.partmap[fmt.Sprintf("%s_%s", gpt.ostype, name)]
	if !ex {
		return "", fmt.Errorf("Unsupported partition type %s on OS %s", name, gpt.ostype)
	}
	return val, nil
}

// Create a next available partition
func (gpt *WakaPartitionerGPT) Create(name string, mbsize int, typeName string) {
	fmt.Printf("Creating %s partition \"%s\" with size %d Mb\n", typeName, name, mbsize)
	partCode, err := gpt.getTypeName(typeName)
	if err != nil {
		panic(err)
	}

	var psize string
	if mbsize == 0 {
		psize = "0" // The rest of the space
	} else {
		psize = fmt.Sprintf("+%dM", mbsize)
	}

	layout := fmt.Sprintf("0:0:%s", psize)
	partition := fmt.Sprintf("0:%s", partCode)
	label := fmt.Sprintf("0:\"%s\"", name)
	cmd, err := wzlib_subprocess.BufferedExec(gpt.parted, "-n", layout, "-t", partition, "-c", label)
	if err != nil {
		panic(err)
	}
	fmt.Println(cmd.StdoutString())

// GetDiskDevice name
func (gpt *WakaPartitionerGPT) GetDiskDevice() string {
	return gpt.device
}
