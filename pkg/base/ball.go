package base

import (
	"github.com/spacexnice/nice/pkg/util"
	"os"
	"bufio"
	"io"
	"strings"
	"strconv"
	"sort"
	"fmt"
    "github.com/golang/glog"
	"encoding/json"
)

const (
	 K11= "5:6"
	 K34 = "11:11:11"
	 K6 = "6:5:6:5:6:5"
	 KEY_PREFIX = "GRP/"
)

var UK34 = &UnionKey{
	// 11:11:11/1/6
	Pattern:    K34,
	Offset :  	1,
	Count  :    6,
	Length :    33,
	Full   :    false,
}

var UK111 = &UnionKey{
	// 11:11:11/1/6
	Pattern:    K34,
	Offset :  	12,
	Count  :    6,
	Length :    33,
	Full   :    false,
}
var UK1112 = &UnionKey{
	// 11:11:11/1/6
	Pattern:    K34,
	Offset :  	1,
	Count  :    6,
	Length :    33,
	Full   :    false,
}
var UK1123 = &UnionKey{
	// 11:11:11/1/6
	Pattern:    K34,
	Offset :  	1,
	Count  :    6,
	Length :    33,
	Full   :    false,
}

type Ball struct {
	//publish date
	Date   string

	//publish index
	Index  int

	// Red Balls
	Reds   []int

	// Blue Balls
	Blue   int

	//like K3 => 2:3:1
	//kike 2  => 2
	Policy map[string]*Group

	Hole      MaxHole
}

type EstimatePolicy struct {
	PatKey    string

	UnionKey  *UnionKey
	Estimates map[string]*Estimate
}

type Estimate struct {
	Num      int

	//被评估对象
	Key      string

	//上次出现的地点
	Last     int
	//被评估对象累积出现次数
	AccCount int
	//被评估对象的平均出现次数
	Avg      float64
	//被评估对象的预计下一次出现次数
	Next     float64

	//被评估对象出现频次的标准差
	Std      float64
	//修正标准差
	FixStd   float64
	//累积标准差
	AccStd   float64

	//期望值
	Expect   float64

	Group    *Group
}

type UnionKey struct {
	// 11:11:11
	Pattern		string

	//
	Key 		string

	// 1
	Offset      int

	// 2
	Count  		int

	// 33
	Length      int

	Full        bool
}
func (u *UnionKey) PKey()string{
	// 11:11:11/1/2
	return fmt.Sprintf("%s/%d/%d",u.Pattern,u.Offset,u.Count)
}
type MaxHole struct {
	Start 	  int
	End 	  int
	Middle 	  int
	Len       int
}

type PatKey struct {
	Pat       int
	Cnt       int
}

type Bucket struct {
	Balls 	  []*Ball
	NextIDX   int
	Product   bool
}

func NewBucket(force bool,idx int) *Bucket{

	return LoadBucket(idx,force)
}


func (bkt *Bucket) RedBall(idx int) []int {

	return bkt.Balls[idx].Reds
}

func (bkt *Bucket) BlueBall(idx int) int{

	return bkt.Balls[idx].Blue
}


func (bkt *Bucket) NicePrint(){
	for _,b := range bkt.Balls{
		fmt.Println(b)
	}
}

func (b * Ball) ForEachInGroup(g *UnionKey,f func(key string)) {
	for i := 1;i<=3;i ++ {
		switch i {
		case 1:
			//glog.Info("Group Case 1")
			for _, v := range b.Reds {
				if v < g.Offset || v > (g.Offset + g.Length-1){
					continue
				}
				f(fmt.Sprintf("%d", v))
			}
		case 2:
			//glog.Info("Group Case 2")
			for k1, v1 := range b.Reds {
				if v1 < g.Offset || v1 > (g.Offset + g.Length-1){
					continue
				}
				for k2, v2 := range b.Reds {
					if k2 < k1 {
						continue
					}
					if v2 < g.Offset || v2 > (g.Offset + g.Length-1){
						continue
					}
					f(fmt.Sprintf("%d:%d", v1, v2))
				}
			}

			//glog.Info("Group Case 2 end")
		case 3:
			//glog.Info("Group Case 3")
			for k1, v1 := range b.Reds {
				if v1 < g.Offset || v1 > (g.Offset + g.Length-1){
					continue
				}
				for k2, v2 := range b.Reds {
					if k2 < k1 {
						continue
					}
					if v2 < g.Offset || v2 > (g.Offset + g.Length-1){
						continue
					}
					for k3, v3 := range b.Reds {
						if k3 < k2 {
							continue
						}
						if v3 < g.Offset || v3 > (g.Offset + g.Length-1){
							continue
						}
						f(fmt.Sprintf("%d:%d:%d", v1, v2, v3))
					}
				}
			}
		}
	}
}



//offset start from 1
func (b *Ball) KeyArea(grp * Group) string{
	var pts []PatKey
	t_cnt := 0
	for _,v := range grp.Pattern.Key{
		t_cnt += v
		pts = append(pts, PatKey{Pat:v})
	}
	for _,v := range b.Reds {
		if v < grp.Start || v > grp.End {
			continue
		}
		rs := v - grp.Start + 1
		for j,p := range pts {
			rs = rs - p.Pat
			if rs <= 0 {
				pts[j].Cnt += 1
				break
			}
		}
	}
	var rts string = ""
	for _,v := range pts {
		rts += fmt.Sprintf("%d:",v.Cnt)
	}
	return fmt.Sprintf("%s",rts[0:len(rts)-1])
}

//offset start from 1
func (b *Ball) KeyPartition(pat string,offset int) string{
	var pts []PatKey
	s,t_cnt := strings.Split(pat,":"),0
	for _,v := range s{
		if v == ":"{
			continue
		}
		i,_ := strconv.Atoi(v)
		t_cnt += i
		pts = append(pts, PatKey{Pat:i})
	}
	for _,v := range b.Reds {
		if v < offset || v >= offset + t_cnt{
			continue
		}
		rs := v - offset + 1
		for j,p := range pts {
			rs = rs - p.Pat
			if rs <= 0 {
				pts[j].Cnt += 1
				break
			}
		}
	}
	var rts string = ""
	for _,v := range pts {
		rts += fmt.Sprintf("%d:",v.Cnt)
	}
	return fmt.Sprintf("%s%s",KEY_PREFIX,rts[0:len(rts)-1])
}

func (b *Ball) maxHole() MaxHole{
	pre,start,end,len := 0,0,0,0
	for _,i := range b.Reds{
		if (i - pre) > len{
			start,end,len = pre,i,(i - start)
		}
		pre = i
	}
	return MaxHole{
		Start:		start,
		End:		end,
		Middle:     (end-start)>>1 + start,
		Len:        (end-start),
	}
}

func (b Ball) String() string{
	var m string = ""
	for k,v := range b.Policy{
		m += fmt.Sprintf("%2s %+v",k,v) +" ## "
	}
	return fmt.Sprintf("DATE:%s  IDX:%d  REDS:%+2v   BLUE:%2d  [PARKEY: %95s   MAXHOLE: %+2v  FREQENCY: ]",
						b.Date,b.Index,b.Reds,b.Blue,m[0:len(m)-3],b.Hole)
}

func (b * Ball) contains(r []int) bool {
	sort.Ints(r)
	return len(b.Intersection(r))>0
}

func (b * Ball) Intersection(r []int)[]int{
	var rt []int
	for i,j:=0,0;i<len(b.Reds)&&j<len(r);{
		if b.Reds[i] < r[j]{
			i++
			continue
		}
		if b.Reds[i] == r[j]{
			rt = append(rt,r[j])
			i++
			j++
			continue
		}
		if b.Reds[i] > r[j]{
			j++
			continue
		}
	}
	return rt
}

func (bkt *Bucket) Intersection(b1,b2 Ball) []int{
	var rt []int
	for i,j:=0,0;i<len(b1.Reds)&&j<len(b2.Reds);{
		if b1.Reds[i] < b2.Reds[j]{
			i++
			continue
		}
		if b1.Reds[i] == b2.Reds[j]{
			i++
			j++
			rt = append(rt,b1.Reds[i])
			continue
		}
		if b1.Reds[i] > b2.Reds[j]{
			j++
			continue
		}
	}
	return rt
}

func LoadBucket(idx int,force bool) *Bucket {

	util.LoadFile(force)

	var balls []*Ball

	ForEachLine(func(line string){
		glog.Infof("处理第 %d 个\n",len(balls))
		l := strings.Split(line," ")
		idx,_ := strconv.Atoi(l[0])
		r1, _ := strconv.Atoi(l[2])
		r2, _ := strconv.Atoi(l[3])
		r3, _ := strconv.Atoi(l[4])
		r4, _ := strconv.Atoi(l[5])
		r5, _ := strconv.Atoi(l[6])
		r6, _ := strconv.Atoi(l[7])
		b1, _ := strconv.Atoi(l[8])
		ball := &Ball{
			Date:	l[1],
			Index:  idx,
			Reds:	[]int{r1,r2,r3,r4,r5,r6},
			Blue:	b1,
			Policy: map[string]*Group{},
		}
		sort.Ints(ball.Reds)

		//ball.Hole   = ball.maxHole()
		balls = append(balls,ball)
		g33 := &[]*Group{
			&Group{
				Index: 1,
				IndexInner: 0,
				Start: 1,
				End:   33,
				AreaLen: 33,
				Pattern: &Pattern{
					Key:	[]int{33},
					Value: 	map[int]RankList{
						6:RankList{
							&RankScore{
								Predict:[]int{6},
							},
						},
					},
				},
				Level:   1,
			},
		}
		add := func(g *Group){
			// 在这儿计算group的EstimateKey
			ball.Policy[g.GroupKey()] = g
			g.AddEstimate(&balls)
			for _,rlist := range g.Parent.Pattern.Value{
				for _,r := range rlist{
					// 计算Pattern的预测
					pdt := r.Predict[g.IndexInner]
					g.AddRankValue(&balls,pdt)
				}
			}
			//glog.Infoln("Group Test: ", " ",g)
		}
		//NewSubGroups(g33,3, add)
		NewSubGroups(NewSubGroups(g33,3, add),2, add)

		return
	})

	next := len(balls)
	if idx > 0 {
		next = idx
	}
	glog.Infoln("预测期数:",next,":",len(balls))
	return &Bucket{
		Balls:	 balls,
		NextIDX: next,
	}
}

func ForEachLine(oper func(line string) ){

	f, err := os.OpenFile(util.SSQ_FILE, os.O_RDONLY, 0666)
	if err != nil{
		glog.Errorf("Error open file[%s],reason[%s]",util.SSQ_FILE,err.Error())
		panic(err)
	}
	defer f.Close()

	bf := bufio.NewReader(f)

	for {
		line, _, err := bf.ReadLine()
		if err == io.EOF{
			break
		}
		oper(string(line))
	}
}

type RankScore struct {
	PatKey      string
	Pattern 	string
	Predict     []int
	Key     	string
	Behind  	float64
	// 标准差.度量期望组合的可信度
	Std     	float64
	// 修正绝对标准差.假设当前出现期望的组合时的标准差
	FixStd  	float64
	//分数分级
	Expect  	float64

	Cross   	string

	Ball    	*Ball

	Group   	*Group
}

func (k *RankScore) SplitKey()[]int{
	var balls []int
	for _,v := range strings.Split(strings.Replace(k.Key,KEY_PREFIX,"",-1),":"){
		i,_ := strconv.Atoi(v)
		balls = append(balls,i)
	}
	return balls
}

type RankList []*RankScore


func (s RankScore) String()string{
	return fmt.Sprintf("Key:%s, Score:%10f,ScoreExponent:%10f, Std:%10f,FixStd:%10f",
		s.Key,s.Behind,s.Expect,s.Std,s.FixStd)
}


func (l RankList) Len()int{
	return len(l)
}

func (l RankList) Less(i,j int)bool{
	if l[i].Expect >= l[j].Expect {
		return true
	}
	return false
}

func (l RankList) Swap(i,j int){
	t   := l[i]
	l[i] = l[j]
	l[j] = t
}
func (l RankList) Merge() *RankScore {
	key,std,fixstd,expect,behind,pattern,pkey := "",0.0,0.0,0.0,0.0,"",""
	for _,v := range l{
		key += fmt.Sprintf("%s:",v.Key)
		std += v.Std
		fixstd += v.FixStd
		expect += v.Expect
		behind += v.Behind
		pattern = v.Pattern
		pkey    = v.PatKey
		//glog.Infof("TEST : %+v\n",v)
	}
	//glog.Infoln("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF")
	le := float64(len(l))
	return &RankScore{
		Key:  key[0:len(key)-1],
		Std:  std/le,
		FixStd: fixstd/le,
		Expect: expect/le,
		Behind: behind/le,
		Pattern:pattern,
		PatKey: pkey,
	}
}
func (l RankList) NicePrint(){
	for _,v := range l{
		if v.FixStd >= 10{
			continue
		}
		fmt.Println(v)
	}
	return
}

func (l RankList) ToJson() string{
	b,e := json.Marshal(l)
	if e != nil {
		panic(e)
	}
	return string(b)
}
