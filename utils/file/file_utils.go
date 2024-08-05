package file

import (
	"fmt"
	"io"
	"os"
)

func ReadText(fileName string) (string, error) {

	file, err := os.Open(fileName)
	if err != nil {
		return "", fmt.Errorf("文件不存在")
	}
	defer file.Close()

	all, _ := io.ReadAll(file)
	return string(all), nil
}

func WriteText(fileName string, content string) error {
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("文件创建失败")
	}
	defer file.Close()
	_, err = file.WriteString(content + "\n")
	if err != nil {
		return fmt.Errorf("文件写入失败")
	}
	return nil
}

func WriteFile() {

}

func WriteFileCSV() {

}
