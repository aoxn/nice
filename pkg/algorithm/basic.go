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
	K3          ScoreList
	Nice        ScoreList
	Ball        base.Ball

	Date        string
}

type KeyScore struct {
	Key    string
	Behind int
	// 标准差.度量期望组合的可信度
	Std    float64
	// 修正绝对标准差.假设当前出现期望的组合时的标准差
	FixStd float64
	//分数分级
	Expect float64

	Ball   base.Ball
}

type ScoreList []*KeyScore


func (s KeyScore) String()string{
	return fmt.Sprintf("Key:%s, Score:%4d,ScoreExponent:%10f, Std:%10f,FixStd:%10f, Ball:%s",
		s.Key,s.Behind,s.Expect,s.Std,s.FixStd,s.Ball)
}


func (l ScoreList) Len()int{
	return len(l)
}

func (l ScoreList) Less(i,j int)bool{
	if l[i].Expect >= l[j].Expect {
		return true
	}
	return false
}

func (l ScoreList) Swap(i,j int){
	t   := l[i]
	l[i] = l[j]
	l[j] = t
}
func (l ScoreList) NicePrint(){
	for _,v := range l{
		if v.FixStd >= 10{
			continue
		}
		fmt.Println(v)
	}
	return
}

func (l ScoreList) ToJson() string{
    b,e := json.Marshal(l)
    if e != nil {
        panic(e)
    }
    return string(b)
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
	res.ParKey3 = res.Ball.Attr.ParKey[base.K3].Key
	res.Date    = res.Ball.Date
	return res
}
