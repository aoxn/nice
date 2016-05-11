package algorithm

import (
	"github.com/spacexnice/nice/pkg/base"
	"sort"
	"github.com/golang/glog"
)

type RelateNicer struct {
	Bucket       *base.Bucket
}

func NewRelateNicer(bkt *base.Bucket) *RelateNicer{
	return &RelateNicer{
		Bucket: 	bkt,
	}
}

func (p *RelateNicer) Predicate(idx int)ScoreList{
	glog.Info("XYYYYYYYY")
	return p.predicate(idx,1,12)
}

func (p *RelateNicer) predicate(idx int,start,end int) ScoreList{
	var res []*KeyScore
	for k := start;k <end;k++ {
		v := (*p.Bucket.Balls[p.Bucket.NextIDX-1].Attr.CoRelate1)[k]
		glog.Info("MMMMMM:",v)
		if k==0 {
			continue
		}
		res = append(res,&KeyScore{
			Key:	v.Key,
			Std:    v.Std,
			FixStd: v.FixStd,
			Expect: v.Expect,
			Behind: v.Next,
		})
	}
	sort.Sort(ScoreList(res))
	glog.Info("HHHHHHH:",res)
	for i,v := range res {
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

