package algorithm

import (
	"sort"
	"fmt"
	"math"
	"github.com/spacexnice/nice/pkg/base"
	"strings"
	"strconv"
	//"github.com/golang/glog"
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
	return p.predicate(idx,base.UK34)
}

func (p *KeyNicer) keystart(key string) bool{
	if strings.Index(key,base.KEY_PREFIX) != -1{
		return true
	}
	return false
}

func (p *KeyNicer) keysum(key string) int{
	sum := 0
	for _,v := range strings.Split(strings.Replace(key,base.KEY_PREFIX,"",-1),":"){
		k,_ := strconv.Atoi(v)
		sum += k
	}
	return sum
}

func (p *KeyNicer) predicate(idx int, uk *base.UnionKey) base.ScoreList{
	cnt,rt := 0.0,make(map[string]*base.KeyScore)
	for i := idx - 1;i>=1;i--{
		pk := p.Bucket.Balls[i].Policy[uk.PKey()]
		//glog.Infoln("AB:",pk.PatKey,"  ",uk.Count)
		if p.keystart(pk.PatKey) && p.keysum(pk.PatKey) != uk.Count{
			continue
		}
		cnt ++
		if _,e := rt[pk.PatKey];e{
			continue
		}
		est := pk.Estimates[pk.PatKey]
		score  := cnt - est.Next
		fixStd := (math.Abs(float64(score))+float64(est.AccCount) * est.Std)/(float64(est.AccCount)+1)
		rt[pk.PatKey] = &base.KeyScore{
			Key:	pk.PatKey,
			Pattern: uk.Pattern,
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
	for _,v := range ss {
		if v.FixStd >= 10{
			//过滤掉修正方差大于10的
			continue
		}
		if len(res) <= 2 {
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
