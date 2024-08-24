package converter

import (
	"fmt"
	"os/exec"
)

func ConvertImage(image1 string, image2 string, resize string, quality string, crop string) error {
	args := []string{image1, "-auto-orient"}
	if quality != "" {
		args = append(args, "-quality", quality)
	}
	if resize != "" && crop == "" {
		args = append(args, "-resize", resize)
	}
	if crop != "" {
		args = append(args, "-crop", crop)
	}
	args = append(args, image2)
	cmd := exec.Command("magick", args...)
	output, err := cmd.Output()
	if err != nil {
		return err
	}
	fmt.Println(string(output))
	return nil
}
