package jobs

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"imagu/internal/db/repo"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var DeletionJob = Job{
	Func: deleteImages,
	Rate: 3600,
}
var uploadsDir = "./data/uploads"

func deleteImages() {
	deletionRate, err := repo.GetAutomaticallyDeletionTime()
	if err != nil {
		return
	}
	deletionDuration := time.Duration(deletionRate) * time.Minute

	now := time.Now()
	err = filepath.Walk(uploadsDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		fileInfo, err := os.Stat(path)
		if err != nil {
			return err
		}
		if strings.Contains(path, "base") {
			logrus.Debug("Skip base file")
			return nil
		}
		lastAccessTime := fileInfo.ModTime()
		if now.Sub(lastAccessTime) > deletionDuration {
			err := os.Remove(path)
			if err != nil {
				logrus.Errorf("Error deleting file %s: %v\n", path, err)
			} else {
				logrus.Infof("Deleted file %s\n", path)
			}
			uuid := filepath.Base(filepath.Dir(path))
			err = repo.UpdateSizeAndCount(uuid, -fileInfo.Size(), -1)
			if err != nil {
				logrus.Errorf("Error updating image entry %s: %v\n", path, err)
			}
		}
		return nil
	})
	if err != nil {
		fmt.Printf("Error walking through the directory: %v\n", err)
	}
}
