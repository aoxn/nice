package util

import (
	"os"
	"net/http"
	"io"
	"github.com/jinzhu/gorm"
	"github.com/golang/glog"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"bytes"
	"encoding/gob"
	"bufio"
)

const (
	SSQ_FILE = "ssq.txt"
	SSQ_URL = "http://www.17500.cn/getData/ssq.TXT"
)

var DEFAULT_DB = "gorm.db"

func OpenInit(path string) *gorm.DB {
	db, err := gorm.Open("sqlite3", fmt.Sprintf("%s/%s", path, DEFAULT_DB))
	if err != nil {
		panic(err)
		return nil
	}
	db.LogMode(true)
	glog.Infof("DATABASE INIT: database [%s] init database and tables... ", fmt.Sprintf("%s/%s", path, DEFAULT_DB))
	return db
}
// 检查文件或目录是否存在
// 如果由 filename 指定的文件或目录存在则返回 true，否则返回 false
func Exist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

func LoadFile(force bool) {
	if force {
		os.Remove(SSQ_FILE)
	}
	if !Exist(SSQ_FILE) {
		res, e := http.Get(SSQ_URL)
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

func DeepCopy(dst, src interface{}) error {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(src); err != nil {
		return err
	}
	return gob.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(dst)
}

func ForEachLine(oper func(line string, idx int)) {

	f, err := os.OpenFile(SSQ_FILE, os.O_RDONLY, 0666)
	if err != nil {
		glog.Errorf("Error open file[%s],reason[%s]", SSQ_FILE, err.Error())
		panic(err)
	}
	defer f.Close()

	bf := bufio.NewReader(f)
	for i := 1; ; i++ {
		line, _, err := bf.ReadLine()
		if err == io.EOF {
			break
		}
		oper(string(line), i)
	}
}

func SetBall(array []int) *[33]int {
	ball := [33]int{}
	for _, b := range array {
		ball[b - 1] = 1
	}
	return &ball
}