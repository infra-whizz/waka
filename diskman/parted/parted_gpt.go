package waka_parted

import (
	"fmt"

	waka_layout "github.com/infra-whizz/waka/layout"
	wzlib_subprocess "github.com/infra-whizz/wzlib/subprocess"
)

type WakaPartitionerGPT struct {
	parted string
	device string
}

func NewWakaPartitionerGPT(device string) *WakaPartitionerGPT {
	gpt := new(WakaPartitionerGPT)
	gpt.parted = "sgdisk"
	gpt.device = device

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

// Create a next available partition
func (gpt *WakaPartitionerGPT) Create(partition *waka_layout.WkLayoutConfPartition) error {
	fmt.Printf("Creating partition \"%s\" with size %d Mb\n", partition.Label, partition.Size)

	var psize string
	if partition.Size == 0 {
		psize = "0" // The rest of the space
	} else {
		psize = fmt.Sprintf("+%dM", partition.Size)
	}

	cmd, err := wzlib_subprocess.BufferedExec(gpt.parted, gpt.device,
		"-n", fmt.Sprintf("0:0:%s", psize),
		"-t", fmt.Sprintf("0:%s", partition.PartitionCode),
		"-c", fmt.Sprintf("0:\"%s\"", partition.Label))
	if err != nil {
		return err
	}
	out := cmd.StdoutString()
	fmt.Println("DEBUG:", out)
	if err := cmd.Wait(); err != nil {
		return err
	}

	return nil
}

// GetDiskDevice name
func (gpt *WakaPartitionerGPT) GetDiskDevice() string {
	return gpt.device
}
