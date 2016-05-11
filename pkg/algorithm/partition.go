package algorithm

import (
	"sort"
	"fmt"
	"math"
	"github.com/spacexnice/nice/pkg/base"
)

type PartitionNicer struct {
	Bucket 		*base.Bucket
}



func NewPartitionNicer(bucket *base.Bucket) *PartitionNicer {

	return &PartitionNicer{
		Bucket: 		bucket,
	}
}

func (p *PartitionNicer) PKey3(idx int)ScoreList{
	return p.predicate(idx,base.K3)
}

func (p *PartitionNicer) PKey6(idx int)ScoreList{
	return p.predicate(idx,base.K6)
}

func (p *PartitionNicer) predicate(idx int,key string) ScoreList{
	cnt,rt := 0,make(map[string]*KeyScore)
	for i := idx - 1;i>=0;i--{
		cnt ++
		pk := p.Bucket.Balls[i].Attr.ParKey[key]
		if _,e := rt[pk.Key];e{
			continue
		}
		score  := cnt - pk.Next
		fixStd := (math.Abs(float64(score))+float64(pk.AccCount) * pk.Std)/(float64(pk.AccCount)+1)
		rt[pk.Key] = &KeyScore{
			Key:	pk.Key,
			Std:	pk.Std,
			FixStd: fixStd,
			Expect: float64(score)/fixStd,
			Behind:	score,
			Ball:	p.Bucket.Balls[i],
		}
	}
	var ss,res []*KeyScore
	for _,v := range rt{
		ss = append(ss,v)
	}
	sort.Sort(ScoreList(ss))
	for i,v := range ss {
		if v.FixStd >= 10{
			//过滤掉修正方差大于10的
			continue
		}
		if i <= 4 {
			//只取前4个
			res=append(res,v)
		}
	}
	return res
}

func (p *PartitionNicer) Show(idx int){
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
