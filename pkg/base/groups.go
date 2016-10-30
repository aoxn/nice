package base

import (
	"fmt"
	"math"
	"github.com/golang/glog"
	"sort"
	"strings"
	"strconv"
)

type Pattern struct {
	Key   []int
	Value map[int]RankList
}

type Group struct {
	// start with 1
	Index      int

	IndexInner int

	Parent     *Group

	Children   []*Group

	Pattern    *Pattern

	// start with 1
	Start      int
	End        int

	// 区间长度
	AreaLen    int

	Level      int
	//

	Estimate   *Estimate
}

func (g *Group) GroupKey() string {
	return fmt.Sprintf("Level:%d/Index:%d/Pattern:%v/Start:%d/End:%d/AreaLen:%d",
		g.Level, g.Index, g.Pattern.Key, g.Start, g.End, g.AreaLen)
}

func (g Group) String() string {
	return fmt.Sprintf("Level:%d/Index:%d/IndexInner:%d/Pattern:%+v/Start:%d/End:%d/AreaLen:%d/children:%d",
		g.Level, g.Index, g.IndexInner, *g.Pattern, g.Start, g.End, g.AreaLen, len(g.Children))
}

func NewSubGroups(grp *[]*Group, grpcount int, AddGroup func(g *Group)) *[]*Group {
	var grps []*Group
	index, start := 1, 1
	for _, v := range (*grp) {
		for k, m := range v.Pattern.Key {
			p := breakdown(m, grpcount)
			g := &Group{
				Index:        index - 1,
				IndexInner: k,
				Parent:     v,
				Pattern:    &Pattern{
					Key:    p,
					Value:    map[int]RankList{
						6:RankList{
							&RankScore{
								Predict:[]int{6},
							},
						},
					},
				},
				Start:      start,
				End:        start + m - 1,
				AreaLen:    m,
				Level:      v.Level + 1,
				//EstimateKey: v.KeyGroup[0][k],
			}
			v.Children = append(v.Children, g)
			AddGroup(g)
			start += m
			index += 1
			grps = append(grps, g)
		}
	}
	return &grps
}

func breakdown(length, cnt int) []int {

	var pat []int
	for i := 1; i <= length; i ++ {
		j := (i - 1) % cnt
		if len(pat) <= j {
			pat = append(pat, 1)
		} else {
			pat[j] += 1
		}
	}
	return pat
}

func Test() {
	g := NewSubGroups(
		&[]*Group{
			&Group{
				Index: 1,
				Start: 1,
				End:   33,
				AreaLen: 33,
				Pattern: &Pattern{
					Key:[]int{33},

				},
			},
		}, 3,
		func(g *Group) {

		},
	)
	for k, v := range *g {
		glog.Infoln("Group Test: ", k, " ", v)
	}
}

func (grp *Group) AddEstimate(balls *[]*Ball) {
	var estimate *Estimate
	iball := len(*balls) - 1
	loop := func(key string) {
		var cnt, total, accstd = 1, 1, 0.0
		for i := iball - 1; i >= 0; i-- {
			policy := (*balls)[i].Policy[grp.GroupKey()]
			if policy == nil {
				// This is the first one,return an arbitary estimates
				break
			}
			v := policy.Estimate
			if v.Key != key {
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
			Group:  grp,
		}
	}
	key := (*balls)[iball].KeyArea(grp)
	loop(key)
	grp.Estimate = estimate
}

func CountKey(key string) int {
	cnt := 0
	for _, v := range strings.Split(key, ":") {
		i, _ := strconv.Atoi(v)
		cnt += i
	}
	return cnt
}

func SplitKey(key string) []int {
	var keys []int
	for _, v := range strings.Split(key, ":") {
		i, _ := strconv.Atoi(v)
		keys = append(keys, i)
	}
	return keys
}

func (grp *Group) AddRankValue(balls *[]*Ball, pdt int) {

	rank, iball := RankList{}, len(*balls) - 1
	nice := func() {
		cnt, rt := 0.0, make(map[string]*RankScore)
		for i := iball; i >= 1; i-- {
			pk := (*balls)[i].Policy[grp.GroupKey()]
			//glog.Infoln("AB:",pk.PatKey,"  ",uk.Count)
			// 过滤掉 key count != predict
			est := pk.Estimate
			//glog.Infof("COUNTKEY:%d   pdt:%d   KEY:%s\n",CountKey(est.Key),pdt,est.Key)
			if CountKey(est.Key) != pdt {
				continue
			}
			cnt ++
			if _, e := rt[est.Key]; e {
				continue
			}
			score := cnt - est.Next
			fixStd := (math.Abs(float64(score)) + float64(est.AccCount) * est.Std) / (float64(est.AccCount) + 1)
			rt[est.Key] = &RankScore{
				Key:     est.Key,
				Pattern: grp.Pattern.Key,
				Predict: SplitKey(est.Key),
				Std:    est.Std,
				FixStd: fixStd,
				Expect: float64(score) / fixStd,
				Behind:    score,
				Ball:    (*balls)[i],
			}
		}
		for _, v := range rt {
			rank = append(rank, v)
		}
	}
	nice()
	grp.Pattern.Value[pdt] = Prune(rank, pdt)
}

func MergeGroups(root *Group) RankList {
	//尾递归

	return combine(root, 6)
}

func combine(root *Group, key int) RankList {
	glog.Infoln("Group LEVEL: ", root)
	if root.Children == nil {
		return root.Pattern.Value[key]
	}
	var res RankList
	for _, v := range root.Pattern.Value[key] {
		r := RankList{}
		for i, m := range v.Predict {
			r = r.merge(combine(root.Children[i], m), true)
		}
		res = append(res, r...)
	}
	return res
}

func Prune(list RankList, pdt int) RankList {
	res, cnt := RankList{}, 2
	//glog.Infoln("BEFORE:",list)
	sort.Sort(list)
	if pdt == 1 {
		cnt = 1
	}
	for _, v := range list {
		if v.FixStd >= 10 {
			//过滤掉修正方差大于10的
			continue
		}
		if len(res) < cnt {
			//只取前4个
			res = append(res, v)
		}
	}
	//glog.Infoln("AFTER:",res)
	return res
}
