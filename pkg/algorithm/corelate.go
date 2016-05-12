package algorithm

import (
	"github.com/spacexnice/nice/pkg/base"
	"math"
	"fmt"
	"sort"
)

type RelateNicer struct {
	Bucket       *base.Bucket

	Group        map[string]*[]M
}

type M struct {
	Key     string

}

func NewRelateNicer(bkt *base.Bucket) *RelateNicer{
	return &RelateNicer{
		Bucket: 	bkt,
	}
}

func (p *RelateNicer) Predicate(idx int)base.ScoreList{
	return p.PreKey(idx,base.ScoreList{})
}

func (p *RelateNicer) PreKey(idx int, list base.ScoreList) base.ScoreList{
	var rtp base.ScoreList
	for _,v := range list{

		grps := base.NewGroups(v.Pattern,v.Key)

		for m,g := range grps{
			rt := make(map[string]*base.KeyScore)
			p.ForEachInGroup(g,func(key string){

				cnt := 0
				for i := idx - 1;i>=0;i--{
					cnt ++
					pk := p.Bucket.Balls[i].Policy[g.Pattern]

					est,e := pk.Estimates[key]
					if !e {
						continue
					}
					score  := cnt - est.Next
					fixStd := (math.Abs(float64(score))+float64(est.AccCount) * est.Std)/(float64(est.AccCount)+1)
					rt[key] = &base.KeyScore{
						Key:	pk.PatKey,
						Std:	est.Std,
						FixStd: fixStd,
						Expect: float64(score)/fixStd,
						Behind:	score,
						Ball:	p.Bucket.Balls[i],
					}
					//find and then end loop
					return
				}

			})
			(*grps)[m].List = p.Prune(rt)
		}

		// here we get all the result into the grps
		// do Merge
		rtp = append(rtp,p.MergeGroups(grps)...)

	}
	return rtp
}

func (p *RelateNicer) MergeGroups(grps *[]base.Group) *base.ScoreList{

	return
}

func (p *RelateNicer) Prune(ls map[string]*base.KeyScore) base.ScoreList {
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

func (p * RelateNicer) ForEachInGroup(g *base.Group,f func(key string)) {
	switch g.Count {
	case 1:
		for i:=g.Start;i< g.End;i++{
			f(fmt.Sprintf("%d",i))
		}
	case 2:
		for i:=g.Start;i< g.End;i++{
			for j:=g.Start;j< g.End;i++ {
				if j <= i{
					continue
				}
				f(fmt.Sprintf("%d:%d", i,j))
			}
		}
	case 3:
		for i:=g.Start;i< g.End;i++{
			for j:=g.Start;j< g.End;j++ {
				if j <= i{
					continue
				}
				for k:=g.Start;k< g.End;k++ {
					if k <= j {
						continue
					}
					f(fmt.Sprintf("%d:%d:%d", i,j,k))
				}
			}
		}
	}
}
