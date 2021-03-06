package main

import (
	"container/list"
	"github.com/daizc311/TidyFilesScript/config"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
	"syscall"
	"time"
)

// 获取相对路径
var currentPath, _ = os.Getwd()

var suffixList = func() *list.List {
	l := list.New()
	l.PushBack(".mp4")
	l.PushBack(".avi")
	l.PushBack(".mkv")
	return l
}()

func main() {
	config.InitLog()
	log.SetFormatter(config.JsonLogFormatter())
	scanDirList(currentPath)
}

func getSuffix(fileName string) string {
	for i := suffixList.Front(); i != nil; i = i.Next() {

		if strings.HasSuffix(strings.ToLower(fileName), i.Value.(string)) {
			return i.Value.(string)
		}
	}
	return ""
}

func scanDirList(dirPath string) {
	var dir, _ = os.ReadDir(dirPath)

	for dirIndex := range dir {
		if dir[dirIndex].IsDir() {
			currentForeachDirName := dir[dirIndex].Name()
			path := dirPath + string(os.PathSeparator) + currentForeachDirName
			entries, _ := os.ReadDir(path)
			i := len(entries)
			if i == 1 {
				entry := entries[0]
				if entry.IsDir() {
					continue
				}
				entrySuffix := getSuffix(entry.Name())
				if entrySuffix != "" {

					var originalFilePath = path + string(os.PathSeparator) + entry.Name()
					var targetDirPath, err = getTargetDirByDate(currentPath, time.Now())
					var targetFilePath = targetDirPath + string(os.PathSeparator) + currentForeachDirName + entrySuffix
					if err != nil {
						log.Errorf("获取今日文件夹失败,cause:%s", err.Error())
						continue
					}
					moveFile(originalFilePath, targetFilePath)
					rmDir(path)
				}
			}
		}
	}
}

func rmDir(path string) {
	log.Infof("开始删除目录=>\n目标路径: %s", path)
	err := syscall.Rmdir(path)
	if err != nil {
		log.Errorf("删除目录失败: %s", err.Error())
	}
}

func moveFile(originalFilePath string, targetFilePath string) {
	log.Infof("开始移动文件=>\n当前文件路径: %s \n目标文件路径: %s", originalFilePath, targetFilePath)
	err := os.Rename(originalFilePath, targetFilePath)
	if err != nil {
		log.Errorf("移动文件失败: %s", err.Error())
	}
}

func getTargetDirByDate(path string, now time.Time) (string, error) {

	dateStr := now.Format("20060102")
	dirPath := path + string(os.PathSeparator) + dateStr

	var _, err = os.Stat(dirPath)
	if err != nil {
		if os.IsNotExist(err) {
			log.Infof("目录不存在，创建输出目录=>\n目标路径: %s", dirPath)
			err := os.Mkdir(dirPath, os.ModePerm)
			if err != nil {
				return dirPath, err
			}
		} else {
			return dirPath, err
		}
	}
	return dirPath, nil

}
