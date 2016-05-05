package util

import (
	"os"
	"net/http"
	"io"
    "github.com/jinzhu/gorm"
    "github.com/golang/glog"
    "fmt"
    _ "github.com/mattn/go-sqlite3"
)

const (
	SSQ_FILE = "ssq.txt"
	SSQ_URL  = "http://www.17500.cn/getData/ssq.TXT"
)
var DEFAULT_DB = "gorm.db"

func OpenInit(path string) *gorm.DB{
    db, err := gorm.Open("sqlite3", fmt.Sprintf("%s/%s",path,DEFAULT_DB))
    if err != nil {
        panic(err)
        return nil
    }
    db.LogMode(true)
    glog.Infof("DATABASE INIT: database [%s] init database and tables... ",fmt.Sprintf("%s/%s",path,DEFAULT_DB))
    return db
}
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