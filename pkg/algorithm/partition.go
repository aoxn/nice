package algorithm

import (
	"sort"
	"fmt"
	"math"
	"github.com/spacexnice/nice/pkg/base"
)

type KeyNicer struct {
	Bucket 		*base.Bucket
}



func NewPartitionNicer(bucket *base.Bucket) *KeyNicer {

	return &KeyNicer{
		Bucket: 		bucket,
	}
}

func (p *KeyNicer) PKey3(idx int)base.ScoreList{
	return p.predicate(idx,base.K3)
}

func (p *KeyNicer) PKey6(idx int)base.ScoreList{
	return p.predicate(idx,base.K6)
}

func (p *KeyNicer) predicate(idx int,key string) base.ScoreList{
	cnt,rt := 0,make(map[string]*base.KeyScore)
	for i := idx - 1;i>=0;i--{
		cnt ++
		pk := p.Bucket.Balls[i].Policy[key]
		if _,e := rt[pk.PatKey];e{
			continue
		}
		est := pk.Estimates[pk.PatKey]
		score  := cnt - est.Next
		fixStd := (math.Abs(float64(score))+float64(est.AccCount) * est.Std)/(float64(est.AccCount)+1)
		rt[pk.PatKey] = &base.KeyScore{
			Key:	pk.PatKey,
			Std:	est.Std,
			FixStd: fixStd,
			Expect: float64(score)/fixStd,
			Behind:	score,
  			Ball:	p.Bucket.Balls[i],
		}
	}

	return p.Prune(rt)
}

func (p *KeyNicer) Prune(ls map[string]*base.KeyScore) base.ScoreList {
	var ss,res []*base.KeyScore
	for _,v := range ls{
		ss = append(ss,v)
	}
	sort.Sort(base.ScoreList(ss))
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

func (p *KeyNicer) Show(idx int){
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
