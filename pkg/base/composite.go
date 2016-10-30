package base

import (
	"fmt"
	"math"
	"strings"
	"strconv"
	"github.com/golang/glog"
	"sort"
)

type Area struct {
	// 11:11:11
	Pattern string

	//
	Key     string

	// 1
	Offset  int

	// 2
	Count   int

	// 33
	Length  int

	Full    bool
}

func (u *Area) PKey() string {
	// 11:11:11/1/2
	return fmt.Sprintf("%s/%d/%d", u.Pattern, u.Offset, u.Count)
}

func (u Area) String() string {
	return fmt.Sprintf("Pattern:%s,key:%s,offset:%d,Count:%d,Length:%d",
		u.Pattern, u.Key, u.Offset, u.Count, u.Length)
}

func (ball *Ball) AddComposite(balls *[]*Ball, unks *[]*Area) {
	var estimate *Estimate
	iball := len(*balls) - 1
	loop := func(key string) {
		var cnt, total, accstd = 1, 1, 0.0
		for i := iball - 1; i >= 0; i-- {
			v, e := (*balls)[i].Composite[key]

			if !e {
				cnt ++
			} else {
				total = v.AccCount + 1
				accstd = v.AccStd
				break
			}
		}
		avg := float64(iball + 1) / float64(total)
		accstd = accstd + math.Abs(float64(avg) - float64(cnt))
		estimate = &Estimate{
			Num:    iball,
			Key:    key,
			Last:   cnt,
			AccCount:    total,
			Avg:    avg,
			Next:   2 * avg - float64(cnt),
			AccStd: accstd,
			Std:    accstd / float64(total),
		}
		ball.Composite[key] = estimate
		//glog.Infof("KEY:%s,  %+v",key,estimate)
	}
	for _, v := range *unks {
		ball.ForEachInGroup(v, loop)
	}
}

func NewArea(pat string) *[]*Area {
	var unks []*Area
	offset := 1
	ps := strings.Split(pat, ":")
	for _, i := range ps {
		v, _ := strconv.Atoi(i)
		unks = append(unks, &Area{
			Pattern: pat,
			Key:     i,
			Offset:  offset,
			Count:   1,
			Length:  v,
		})
		offset += v
	}
	return &unks
}

func (u *Area) ForEachCount(f func(key string)) {
	switch u.Count {
	case 1:
		//glog.Info("Group Case 1")
		for i := u.Offset; i < u.Offset + u.Length; i++ {
			f(fmt.Sprintf("%d", i))
		}
	case 2:
		//glog.Info("Group Case 2")
		for i := u.Offset; i < u.Offset + u.Length; i++ {
			for j := u.Offset; j < u.Offset + u.Length; j++ {
				if j <= i {
					continue
				}
				f(fmt.Sprintf("%d:%d", i, j))
			}
		}

	//glog.Info("Group Case 2 end")
	case 3:
		//glog.Info("Group Case 3")
		for i := u.Offset; i < u.Offset + u.Length; i++ {
			for j := u.Offset; j < u.Offset + u.Length; j++ {
				if j <= i {
					continue
				}
				for k := u.Offset; k < u.Offset + u.Length; k++ {
					if k <= j {
						continue
					}
					f(fmt.Sprintf("%d:%d:%d", i, j, k))
				}
			}
		}
	}
}

func (rank *RankScore) NewAreas() *[]*Area {
	unks, offset := []*Area{}, 1
	for k, v := range rank.Pattern {
		unks = append(unks, &Area{
			Pattern: fmt.Sprintf("%d", v),
			Key:     fmt.Sprintf("%d", rank.Predict[k]),
			Offset:  offset,
			Count:   rank.Predict[k],
			Length:  v,
		})
		offset += v
	}
	return &unks
}

func (rank *RankScore) AddPredictGroup(bkt *Bucket, idx int) {

	for _, v := range *rank.NewAreas() {
		glog.Infoln("UNIONKEY: ", v)
		rt := make(map[string]*RankScore)
		v.ForEachCount(func(key string) {
			cnt := 0.0
			for i := idx; i >= 1; i-- {
				cnt ++
				pk, e := bkt.Balls[i].Composite[key]
				if !e {
					continue
				}
				score := cnt - pk.Next
				fixStd := (math.Abs(float64(score)) + float64(pk.AccCount) * pk.Std) / (float64(pk.AccCount) + 1)
				rt[key] = &RankScore{
					Key:    key,
					Predict: rank.Predict,
					Pattern: rank.Pattern,
					Std:    pk.Std,
					FixStd: fixStd,
					Expect: float64(score) / fixStd,
					Behind:    score,
				}
				//glog.Infof("RT[key]: %+v\n",rt[key])
				//find and then end loop
				return
			}
		})
		var list RankList
		for _, m := range rt {
			list = append(list, m)
		}
		list = list.prune()
		rank.Group = append(rank.Group, &PdtGroup{Key:v.Count, Value:list})
	}
}

func (rank *RankScore) MergePdtGroup() RankList {
	rlist := RankList{}
	for _, pg := range rank.Group {
		glog.Infoln("Before Merge: ", pg.Key, "  ", pg.Value)
		rlist = rlist.merge(pg.Value, false)
		glog.Infoln("After Mege  : ", pg.Key, "  ", rlist)
	}
	return rlist
}

func (rlist RankList) prune() RankList {
	res, cnt := RankList{}, 2
	sort.Sort(rlist)
	glog.Infoln("****************************************************************")
	glog.Infof("BEFORE PRUNE: len = [%d]\n", len(rlist))
	for _, v := range rlist {
		if len(res) < cnt {
			glog.Infoln("PRUNE KEEP: ", v)
			//只取前4个
			res = append(res, v)
		}
	}
	glog.Infof("PRUNE KEEP:  len = [%d]\n\n", len(res))
	return res
}