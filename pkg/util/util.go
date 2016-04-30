package util

import (
	"os"
	"net/http"
	"io"
)

const (
	SSQ_FILE = "ssq.txt"
	SSQ_URL  = "http://www.17500.cn/getData/ssq.TXT"
)

// 检查文件或目录是否存在
// 如果由 filename 指定的文件或目录存在则返回 true，否则返回 false
func Exist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}


func LoadFile(force bool){
	if force{
		os.Remove(SSQ_FILE)
	}
	if !Exist(SSQ_FILE){
		res, e  := http.Get(SSQ_URL)
		if e != nil {
			panic(e)
		}

		file, e := os.Create(SSQ_FILE)
		if e != nil {
			panic(e)
		}
		defer file.Close()

		io.Copy(file, res.Body)
	}
}