package algorithm

import (
	"math"
	"sort"
	"github.com/spacexnice/nice/pkg/base"
)

type GroupNice struct {
	Bucket     *base.Bucket
}

func (p *GroupNice) Nice(idx int, grp *base.Group) base.RankList {
	cnt,rt := 0.0,make(map[string]*base.RankScore)
	for i := idx - 1;i>=1;i-- {
		pk := p.Bucket.Balls[i].Policy[grp.GroupKey()]
		//glog.Infoln("AB:",pk.PatKey,"  ",uk.Count)
		cnt ++
		est := pk.Estimate
		if _,e := rt[est.Key];e{
			continue
		}
		score  := cnt - est.Next
		fixStd := (math.Abs(float64(score))+float64(est.AccCount) * est.Std)/(float64(est.AccCount)+1)
		rt[est.Key] = &base.RankScore{
			Key:	est.Key,
			Pattern: grp.Pattern,
			Std:	est.Std,
			FixStd: fixStd,
			Expect: float64(score)/fixStd,
			Behind:	score,
			Ball:	p.Bucket.Balls[i],
			Group:  grp,
		}
	}
	var s []*base.RankScore
	for _,v := range rt{
		s    = append(s,v)
	}
	return p.Prune(rt)
}

func (p *GroupNice) Prune(ss base.RankList) base.RankList {
	var res base.RankList

	sort.Sort(ss)
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

