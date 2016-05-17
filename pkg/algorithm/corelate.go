package algorithm

import (
	"github.com/spacexnice/nice/pkg/base"
	"math"
	"fmt"
	"sort"
	"github.com/golang/glog"
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

func (p *RelateNicer) Haha()base.ScoreList{
	return p.PreKey(430,base.ScoreList([]*base.KeyScore{
		&base.KeyScore{
			Key: "2:1:3",
			Pattern:"11:11:11",
		},
	}))
}

func (p *RelateNicer) Predicate(idx int,list base.ScoreList)base.ScoreList{
	return p.PreKey(idx,list)
}


func (p *RelateNicer) PreKey(idx int, list base.ScoreList) base.ScoreList{
	var rtp base.ScoreList

	for _,v := range list{
		glog.Infoln("PMXY:",v)
		grps := base.NewGroups(v.Pattern,v.Key)

		for m,g := range grps{

			//glog.Infof("GROUP: %+v\n",g)
			rt := make(map[string]*base.KeyScore)
			p.ForEachInGroup(g,func(key string){

				cnt := 0.0
				for i := idx - 1;i>=1;i--{
					cnt ++
					pk := p.Bucket.Balls[i].Policy[g.Pattern]

					est,e := pk.Estimates[key]
					if !e {
						continue
					}
					score  := cnt - est.Next
					fixStd := (math.Abs(float64(score))+float64(est.AccCount) * est.Std)/(float64(est.AccCount)+1)
					rt[key] = &base.KeyScore{
						Key:	key,
						PatKey: pk.PatKey,
						Pattern: v.Pattern,
						Std:	est.Std,
						FixStd: fixStd,
						Expect: float64(score)/fixStd,
						Behind:	score,
						Ball:	p.Bucket.Balls[i],
					}
					//glog.Infof("RT[key]: %+v\n",rt[key])
					//find and then end loop
					return
				}
			})
			grps[m].List = p.Prune(rt)
		}

		// here we get all the result into the grps
		// do Merge
		rtp = append(rtp,p.MergeGroups(grps)...)

	}
	sort.Sort(rtp)
	var res base.ScoreList
	for _,v := range rtp {
		//if v.FixStd >= 10{
		//	//过滤掉修正方差大于10的
		//	continue
		//}
		if len(res) <= 20 {
			//只取前4个
			res=append(res,v)
			glog.Infoln(v)
		}
	}

	glog.Warningln("预测值长度",len(res))

	return res
}

func (p *RelateNicer) MergeGroups(grps []*base.Group) base.ScoreList{
	//尾递归

	return p.merge(0,grps,base.ScoreList{})
}

func (p *RelateNicer) merge(i int,grps []*base.Group,list base.ScoreList)base.ScoreList{
	if i >= len(grps){
		//收割的时候
		return base.ScoreList{list.Merge()}
	}
	if (grps)[i].List.Len() == 0 {
		return p.merge(i+1,grps,list)
	}
	var rtp base.ScoreList
	for _,v := range (grps)[i].List{
		l := append(list,v)
		rtp = append(rtp,p.merge(i+1,grps,l)...)
	}
	return rtp
}

func (p *RelateNicer) Prune(ls map[string]*base.KeyScore) base.ScoreList {
	//至少保留2个
	var ss,res []*base.KeyScore
	for _,v := range ls{
		ss = append(ss,v)
	}
	sort.Sort(base.ScoreList(ss))
	for _,v := range ss {
		//if v.FixStd >= 10{
		//	//过滤掉修正方差大于10的
		//	continue
		//}
		if len(res) <= 10 {
			//只取前4个
			res=append(res,v)
		}
	}
	return res
}

func (p * RelateNicer) ForEachInGroup(g *base.Group,f func(key string)) {
	switch g.Count {
	case 1:
		glog.Info("Group Case 1")
		for i:=g.Start;i< g.End;i++{
			f(fmt.Sprintf("%d",i))
		}
	case 2:
		glog.Info("Group Case 2")
		for i:=g.Start;i< g.End;i++{
			for j:=g.Start;j< g.End;j++ {
				if j <= i{
					continue
				}
				f(fmt.Sprintf("%d:%d", i,j))
			}
		}

		glog.Info("Group Case 2 end")
	case 3:
		glog.Info("Group Case 3")
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
