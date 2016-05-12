package algorithm

import (
	"github.com/spacexnice/nice/pkg/base"
	"fmt"
    "github.com/jfrazelle/go/canonical/json"
	"github.com/golang/glog"
)


type Record struct {
	//IDX编号从0开始算起
	IDX         int       `gorm:"primary_key"`
	Index       int

	BallJson    string    `gorm:"size:4096"`
	K3Json      string    `gorm:"size:65535"`

	NiceJson    string    `gorm:"size:65535"`
}

type Result struct {
	Record      Record

	ParKey3     string
	K3          base.ScoreList
	Nice        base.ScoreList
	Ball        base.Ball

	Date        string
}


func (rec Record) LoadResult() Result{
	res := Result{
		Record:  rec,
	}
	e := json.Unmarshal([]byte(rec.K3Json),&res.K3)
	if e != nil {
		panic(e)
	}
	e  = json.Unmarshal([]byte(rec.NiceJson),&res.Nice)
	if e != nil {
		panic(e)
	}

	res.Ball = base.Ball{Reds:[]int{0,0,0,0,0,0}}
	if rec.BallJson == "" {
		glog.Warningf("It seems that this is the latest Predict,which Ball is unkonw,[%d][%d]\n",rec.IDX,rec.Index)
		return res
	}

	e  = json.Unmarshal([]byte(rec.BallJson),&res.Ball)
	if e != nil{
		fmt.Println("Could be nil:",e)
	}
	//res.ParKey3 = res.Ball.Attr.ParKey[base.K3].Key
	res.Date    = res.Ball.Date
	return res
}
