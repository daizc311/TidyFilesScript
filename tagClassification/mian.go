package main

import (
	"github.com/daizc311/TidyFilesScript/config"
	log "github.com/sirupsen/logrus"
	"os"
	"path"
	"regexp"
	"strings"
)

// 获取相对路径
var currentPath, _ = os.Getwd()

func main() {
	config.InitLog()
	log.SetFormatter(config.JsonLogFormatter())

	// 按照文件类型分类
	tagDirMap, fileList := scanDirTagAndFileList(currentPath)

	// 遍历文件夹，找到对应的文件
	for index := range fileList {
		fileInfo := fileList[index]
		fileName := fileInfo.Name()
		// 使用文件名循环匹配DirTag
		for tag := range tagDirMap {
			if strings.Contains(fileName, tag) {
				matchTagDirInfo := tagDirMap[tag]
				copyFile(fileInfo, matchTagDirInfo)
				//log.Infof("开始移动文件: %s ==[move]==> %s", fileName, newFileInfo.Name())
				break
			}
		}
	}

}

func scanDirTagAndFileList(dirPath string) (map[string]os.DirEntry, []os.DirEntry) {

	var dir, _ = os.ReadDir(dirPath)
	// 用于匹配 #[TAG] 的正则
	var dirRegexp, _ = regexp.Compile("#\\[{1}[^\\]]*\\]{1}")

	var tagDirMap = make(map[string]os.DirEntry)
	var fileList []os.DirEntry

	for i := range dir {
		dirEntry := dir[i]
		dirEntryName := dirEntry.Name()
		if dirEntry.IsDir() {
			allDirTags := dirRegexp.FindAllString(dirEntryName, -1)
			// 有TAG的文件夹放进map
			if len(allDirTags) > 0 {
				// tag-file反向map
				for tagIndex := range allDirTags {
					// 处理掉tag两边的方括号
					dirTag := allDirTags[tagIndex]
					cleanTag := dirTag[2 : len(dirTag)-1]
					// 先来先到，后来没了
					if tagDirMap[cleanTag] == nil {
						tagDirMap[cleanTag] = dirEntry
					}
				}
				log.Debugf("查找到目录TAG: %s ==DirTag==> %s", dirEntryName, strings.Join(allDirTags, ","))
			}
		} else {
			fileList = append(fileList, dirEntry)
		}
	}

	if len(tagDirMap) > 0 {
		tagDirMapKeys := make([]string, 0, len(tagDirMap))
		for tag := range tagDirMap {
			tagDirMapKeys = append(tagDirMapKeys, tag)
		}
		log.Infof("目录[%s]中查找到TAG: [%s]", dirPath, strings.Join(tagDirMapKeys, ","))
	} else {
		log.Infof("目录[%s]中未查找到TAG", dirPath)
	}

	if len(fileList) > 0 {
		log.Infof("目录[%s]中查找到文件%d个", dirPath, len(fileList))
	} else {
		log.Infof("目录[%s]中未查找到文件", dirPath)
	}

	return tagDirMap, fileList
}

func copyFile(fileInfo os.DirEntry, dirInfo os.DirEntry) {

	fileName := fileInfo.Name()

	var currentFilePath = path.Join(currentPath, fileName)
	var targetFilePath = path.Join(currentPath, dirInfo.Name(), fileName)

	log.Infof("当前文件路径: %s \n目标文件路径: %s", currentFilePath, targetFilePath)
	err := os.Rename(currentFilePath, targetFilePath)
	if err != nil {
		log.Errorf("移动文件失败: %s", err.Error())
	}
}
