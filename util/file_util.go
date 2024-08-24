package util

import (
	"os"
	"path/filepath"
)

func ChangeExtension(filename string, newExt string) string {
	name := filename[:len(filename)-len(filepath.Ext(filename))]
	return name + "." + newExt
}
func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return err == nil
}
