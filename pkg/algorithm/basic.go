package algorithm

import (
	"github.com/spacexnice/nice/pkg/base"
	"sort"
	"fmt"
	"math"
    "github.com/jfrazelle/go/canonical/json"
)

type GroupPredicator struct {
	Bucket 		*base.Bucket
}

type Result struct {
    IDX         int

    K3          ScoreList
    K3S         string `gorm:"size:65535"`
}
type Score struct {
	Key 		string
	Rank    	int
	// 标准差.度量期望组合的可信度
	Std     	float64
	// 修正绝对标准差.假设当前出现期望的组合时的标准差
	FixStd  	float64
	//分数分级
	ScoreExp 	float64

	Ball    	base.Ball
}

func (s Score) String()string{
	return fmt.Sprintf("Key:%s, Score:%4d,ScoreExponent:%10f, Std:%10f,FixStd:%10f, Ball:%s",
		s.Key,s.Rank,s.ScoreExp,s.Std,s.FixStd,s.Ball)
}

type ScoreList []*Score

func (l ScoreList) Len()int{
	return len(l)
}

func (l ScoreList) Less(i,j int)bool{
	if l[i].ScoreExp <= l[j].ScoreExp{
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

func NewPredicator(bucket *base.Bucket) *GroupPredicator{

	return &GroupPredicator{

		Bucket: 		bucket,
	}
}

func (p *GroupPredicator) PKey3(idx int)ScoreList{
	return p.predicate(idx,base.K3)
}

func (p *GroupPredicator) PKey6(idx int)ScoreList{
	return p.predicate(idx,base.K6)
}

func (p *GroupPredicator) predicate(idx int,key string) ScoreList{
	cnt,rt := 0,make(map[string]*Score)
	for i := idx - 1;i>=0;i--{
		cnt ++
		pk := p.Bucket.Balls[i].Attr.ParKey[key]
		if _,e := rt[pk.Key];e{
			continue
		}
		score  := cnt - pk.Next
		fixStd := (math.Abs(float64(score))+float64(pk.Total) * pk.Std)/(float64(pk.Total)+1)
		rt[pk.Key] = &Score{
			Key:	pk.Key,
			Std:	pk.Std,
			FixStd: fixStd,
			ScoreExp: float64(score)/fixStd,
			Rank:	score,
			Ball:	p.Bucket.Balls[i],
		}
	}
	var ss []*Score
	for _,v := range rt{
		ss = append(ss,v)
	}
	sort.Sort(ScoreList(ss))
	return ss
}

func (p *GroupPredicator) Show(idx int){
	le := len(p.Bucket.Balls)
	if idx > le{
		return
	}
	cnt := 0
	for i := idx ;i<le;i++{
		fmt.Println(p.Bucket.Balls[i])
		if cnt >= 3{
			break
		}
		cnt ++
	}
	return
}