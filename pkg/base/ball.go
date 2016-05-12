package base

import (
	"github.com/spacexnice/nice/pkg/util"
	"os"
	"bufio"
	"io"
	"strings"
	"strconv"
	"sort"
	"math"
	"fmt"
    "github.com/golang/glog"
	"encoding/json"
)

const (
	 K3 = "11:11:11"
	 K6 = "6:5:6:5:6:5"
	 KEY_PREFIX = "GRP"
)

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
	Policy map[string]*EstimatePolicy

	Hole      MaxHole
}

type EstimatePolicy struct {
	PatKey    string

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
	Next     int

	//被评估对象出现频次的标准差
	Std      float64
	//修正标准差
	FixStd   float64
	//累积标准差
	AccStd   float64

	//期望值
	Expect   float64
}
type Group struct {
	//like 11:11:11
	Pattern  string

	//like 2:1:3
	Key      string

	//like 1, which part of the group
	Index    int

	//Count like 2
	Count    int

	//group start num
	Start    int
	End      int

	List     ScoreList
}


func NewGroups(pat,key string) *[]Group{
	var pts *[]Group
	pats := strings.Split(pat,":")
	keys := strings.Split(key,":")
	start:= 1
	for k,v := range pats {
		if v == ":"{
			continue
		}
		i,_ := strconv.Atoi(v)
		j,_ := strconv.Atoi(keys[k])
		pts = append(pts,&Group{
			Pattern: pat,
			Key:key,
			Index: k,
			Count: j,
			Start: start,
			End: start + i,
		})
		start = start + i
	}
	return pts
}


type MaxHole struct {
	Start 	  int
	End 	  int
	Middle 	  int
	Len       int
}

type Pattern struct {
	Pat       int
	Cnt       int
}

type Bucket struct {
	Balls 	  []*Ball
	NextIDX   int
}

func NewBucket(force bool) *Bucket{
	util.LoadFile(force)
	return LoadBucket()
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

func (b *Ball) foreach1(f func(key string)){
	for _,v :=range b.Reds{

		f(fmt.Sprintf("%d",v))
	}
}

func (b *Ball) foreach2(f func(key string)){
	for k1,_ :=range b.Reds{
		for k2,_ :=range b.Reds{
			if k2 < k1 {
				continue
			}
			f(fmt.Sprintf("%d:%d",b.Reds[k1],b.Reds[k2]))
		}
	}
}

func (b *Ball) foreach3(f func(key string)){
	for k1,_ :=range b.Reds{
		for k2,_ :=range b.Reds{
			if k2 < k1{
				continue
			}
			for k3,_ :=range b.Reds{
				if k3<k2{
					continue
				}
				f(fmt.Sprintf("%d:%d:%d",b.Reds[k1],b.Reds[k2],b.Reds[k3]))
			}
		}
	}
}

func (b *Ball) EstimatePolicy(balls []*Ball,pat string) *EstimatePolicy{
	estimates := map[string]*Estimate{}

	loop := func(key string){
		var cnt,total,accstd = 1,1,0.0
		for i:=len(balls)-1;i>=0;i--{
			if v,e := balls[i].Policy[pat].Estimates[key]; !e {
				cnt ++
			}else {
				total = v.AccCount + 1
				accstd= v.AccStd
				break
			}
		}
		avg   := float64(len(balls))/float64(total)
		accstd = accstd + math.Abs(float64(avg - cnt))
		estimates[key] = &Estimate{
			Num:    len(balls)-1,
			Key: 	key,
			Last:   cnt,
			AccCount: 	total,
			Avg:   	avg,
			Next:   2*avg - cnt,
			AccStd: accstd,
			Std:    accstd/float64(total),
		}
	}
	key := b.KeyPartition(pat)
	b.foreach1(loop)
	b.foreach2(loop)
	b.foreach3(loop)
	loop(key)
	return &EstimatePolicy{
		PatKey:       key,
		Estimates: estimates,
	}
}

func (b *Ball) KeyPartition(pat string) string{
	var pts []Pattern
	s := strings.Split(pat,":")
	for _,v := range s{
		if v == ":"{
			continue
		}
		i,_ := strconv.Atoi(v)
		pts = append(pts,Pattern{Pat:i})
	}
	for _,v := range b.Reds {
		rs := v
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
		rts += fmt.Sprintf("%d-",v.Cnt)
	}
	return fmt.Sprintf("%s/%s",KEY_PREFIX,rts[0:len(rts)-1])
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
	return len(b.intersection(r))>0
}

func (b * Ball) intersection(r []int)[]int{
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

func LoadBucket() *Bucket {

	var balls []*Ball

	ForEachLine(func(line string){
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
		}
		sort.Ints(ball.Reds)

		//pre := Ball{}
		//if len(balls)>0{
		//	pre = balls[len(balls)-1]
		//}
		ball.Hole   = ball.maxHole()

		ball.Policy = map[string]*EstimatePolicy{
			K3:   	ball.EstimatePolicy(balls,K3),
		}

		balls = append(balls,ball)
		return
	})

	return &Bucket{
		Balls:	 balls,
		NextIDX: len(balls),
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

type KeyScore struct {
	Pattern string
	Key     string
	Behind  int
	// 标准差.度量期望组合的可信度
	Std     float64
	// 修正绝对标准差.假设当前出现期望的组合时的标准差
	FixStd  float64
	//分数分级
	Expect  float64

	Ball    *Ball
}

type ScoreList []*KeyScore


func (s KeyScore) String()string{
	return fmt.Sprintf("Key:%s, Score:%4d,ScoreExponent:%10f, Std:%10f,FixStd:%10f, Ball:%s",
		s.Key,s.Behind,s.Expect,s.Std,s.FixStd,s.Ball)
}


func (l ScoreList) Len()int{
	return len(l)
}

func (l ScoreList) Less(i,j int)bool{
	if l[i].Expect >= l[j].Expect {
		return true
	}
	return false
}

func (l ScoreList) Swap(i,j int){
	t   := l[i]
	l[i] = l[j]
	l[j] = t
}
func (l ScoreList) NicePrint(){
	for _,v := range l{
		if v.FixStd >= 10{
			continue
		}
		fmt.Println(v)
	}
	return
}

func (l ScoreList) ToJson() string{
	b,e := json.Marshal(l)
	if e != nil {
		panic(e)
	}
	return string(b)
}
