package converter

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os/exec"
)

func ConvertImage(image1 string, image2 string, resize string, quality string, crop string, blur string) error {
	args := []string{image1, "-auto-orient", "-strip"}
	if quality != "" {
		args = append(args, "-quality", quality)
	}
	if resize != "" && crop == "" {
		args = append(args, "-resize", resize)
	}
	if crop != "" {
		args = append(args, "-crop", crop)
	}
	if blur != "" {
		args = append(args, "-blur", blur)
	}
	args = append(args, image2)
	cmd := exec.Command("magick", args...)
	logrus.Info("Converting image: ", cmd.Args)
	output, err := cmd.Output()
	if err != nil {
		return err
	}
	fmt.Println(string(output))
	return nil
}
