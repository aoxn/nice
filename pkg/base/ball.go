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
	"k8s.io/kubernetes/third_party/golang/go/doc/testdata"
)

const (
	 K3 = "11:11:11"
	 K6 = "6:5:6:5:6:5"
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

	Attr   *Attribute
}

type EstimatePolicy struct {
	Key       string

	Estimates map[string]*Estimate
}

type Attribute struct {
	ParKey    map[string]*Estimate
	Hole      MaxHole
	CoRelate2 *[34]*[34]Estimate            `json:"-"`
	CoRelate3 *[34]*[34]*[34]Estimate        `json:"-"`
	CoRelate1 *[34]Estimate                `json:"-"`
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
	Avg      int
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
	Balls 	  []Ball
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

//func (b *Ball) keyPartition(pat int)string{
//	secs := int(math.Ceil(float64(33)/float64(pat)))
//	m := make(map[int]int)
//	for i:=0;i<secs;i++{
//		m[i] = 0
//	}
//	for _,i := range b.Reds{
//		x := (i-1)/pat
//		m[x] = m[x] +1
//	}
//	var rs = ""
//	for i:=0;i<secs;i++{
//		rs += fmt.Sprintf("%d-",m[i])
//		//fmt.Println(fmt.Sprintf("%d:",m[i]))
//	}
//
//	return rs[0:len(rs)-1]
//}

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
				f(fmt.Sprintf("%d:%d:5d",b.Reds[k1],b.Reds[k2],b.Reds[k3]))
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
		avg   := len(balls)/total
		accstd = accstd + math.Abs(float64(avg - cnt))
		estimates[key] = &Estimate{
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
		Key:       key,
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
	return rts[0:len(rts)-1]
}

func (b *Ball) PartitionGroup(balls []*Ball,pat string) *Estimate {
	var cnt,total,found,accstd = 1,1,false,0.0
	kNum := b.KeyPartition(pat)
	for i:=len(balls)-1;i>=0;i--{
		if e,v := balls[i].Policy[pat].Estimates[kNum]; !e {
			if !found {
				cnt ++
			}
		}else {
			total ++
			if !found{
				found = true
				accstd = balls[i].Attr.ParKey[pat].AccStd
			}
		}
	}
	avg   := len(balls)/total
	accstd = accstd + math.Abs(float64(avg - cnt))
	return &Estimate{
		Key: 	kNum,
		Last:   cnt,
		AccCount: 	total,
		Avg:   	avg,
		Next:   2*avg - cnt,
		AccStd: accstd,
		Std:    accstd/float64(total),
	}
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
func (b *Ball) corelation1(balls []Ball,pre *Ball)*[34]Estimate {
	if pre.Attr.CoRelate1 == nil{
		pre.Attr.CoRelate1 = &[34]Estimate{}
	}
	co1 := *pre.Attr.CoRelate1

	for _, i := range b.Reds{
		co1[i].AccCount = co1[i].AccCount + 1
	}

	for k,_ := range co1{
		if k == 0 {
			continue
		}
		co,cnt,accstd := &co1[k],0,float64(0)

		for i:=len(balls)-1;i>=0;i--{
			if !balls[i].contains([]int{k}){
				cnt ++
			}else {
				accstd = (*balls[i].Attr.CoRelate1)[k].AccStd
				break
			}
		}
		if co.AccCount <= 0 {
			co.AccCount = 1
		}
		avg := len(balls)/co.AccCount
		co.Avg = avg
		co.Key = fmt.Sprintf("%d",k)
		co.Last= cnt
		co.Next= 2*avg - cnt
		co.AccStd = accstd + math.Abs(float64(avg - cnt))
		co.Std = co.AccStd/float64(co.AccCount)
		co.FixStd = (math.Abs(float64(co.Next))+float64(co.AccCount) * co.Std)/(float64(co.AccCount)+1)
		co.Expect = float64(co.Next)/co.FixStd
	}
	b.Attr.CoRelate1 = &co1
	return b.Attr.CoRelate1
}

func (b *Ball) corelation2(balls []Ball,pre *Ball) *[34][34]Estimate {
	if pre.Attr.CoRelate2 == nil{
		pre.Attr.CoRelate2 = &[34]*[34]Estimate{}
		for k,_ := range pre.Attr.CoRelate2{
			pre.Attr.CoRelate2[k] = &[34]Estimate{}
		}
	}
	co2 := *pre.Attr.CoRelate2
	util.DeepCopy()
	for _, bfirst := range b.Reds {
		for _, bsecond := range b.Reds {
			(*co2[bfirst])[bsecond].AccCount = (*co2[bfirst])[bsecond].AccCount + 1
		}
	}
	for k1,_ := range co2{
		if k1 == 0 {
			continue
		}
		for k2,_ := range &co2[k1]{
			if k2 == 0 {
				continue
			}

			co,cnt,accstd := &co2[k1][k2],0,float64(0)

			for i:=len(balls)-1;i>=0;i--{
				if !balls[i].contains([]int{k1,k2}){
					cnt ++
				}else {
					accstd = (*balls[i].Attr.CoRelate2)[k1][k2].AccStd
					break
				}
			}
			if co.Total <= 0 {
				co.Total = 1
			}
			avg := len(balls)/co.Total
			co.Avg = avg
			co.Key = fmt.Sprintf("%d:%d",k1,k2)
			co.Last= cnt
			co.Next= 2*avg - cnt
			co.AccStd = accstd + math.Abs(float64(avg - cnt))
			co.Std = co.AccStd/float64(co.Total)
		}
	}
	b.Attr.CoRelate2 = &co2
	return b.Attr.CoRelate2
}

func (b *Ball) corelation3(balls []Ball,pre *Ball) *[34][34][34]Estimate {
	if pre.Attr.CoRelate3 == nil{
		pre.Attr.CoRelate3 = &[34][34][34]Estimate{}
	}
	co3 := *pre.Attr.CoRelate3
	for _, bfirst := range b.Reds {
		for _, bsecond := range b.Reds {
			for _, bthird := range b.Reds{
				co3[bfirst][bsecond][bthird].Total = co3[bfirst][bsecond][bthird].Total + 1
			}
		}
	}

	for k1,_ := range co3{
		if k1 == 0 {
			continue
		}
		for k2,_ := range co3[k1]{
			if k2 == 0 {
				continue
			}
			for k3,_ :=range co3[k2]{
				if k3 == 0 {
					continue
				}

				co,cnt,accstd := &co3[k1][k2][k3],0,float64(0)

				for i:=len(balls)-1;i>=0;i--{
					if !balls[i].contains([]int{k1,k2,k3}){
						cnt ++
					}else {
						accstd = (*balls[i].Attr.CoRelate3)[k1][k2][k3].AccStd
						break
					}
				}
				if co.Total <= 0 {
					co.Total = 1
				}
				avg := len(balls)/co.Total
				co.Avg = avg
				co.Key = fmt.Sprintf("%d:%d:%d",k1,k2,k3)
				co.Last= cnt
				co.Next= 2*avg - cnt
				co.AccStd = accstd + math.Abs(float64(avg - cnt))
				co.Std = co.AccStd/float64(co.Total)



			}
		}
	}
	b.Attr.CoRelate3 = &co3

	return b.Attr.CoRelate3
}

func (b Ball) String() string{
	var m string = ""
	for k,v := range b.Attr.ParKey{
		m += fmt.Sprintf("%2s %+v",k,v) +" ## "
	}
	return fmt.Sprintf("DATE:%s  IDX:%d  REDS:%+2v   BLUE:%2d  [PARKEY: %95s   MAXHOLE: %+2v  FREQENCY: %v]",
						b.Date,b.Index,b.Reds,b.Blue,m[0:len(m)-3],b.Attr.Hole,b.Attr.CoRelate1)
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

	var balls []Ball

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
		ball := Ball{
			Date:	l[1],
			Index:  idx,
			Reds:	[]int{r1,r2,r3,r4,r5,r6},
			Blue:	b1,
		}
		sort.Ints(ball.Reds)

		pre := Ball{}
		if len(balls)>0{
			pre = balls[len(balls)-1]
		}
		ball.Attr = Attribute{
			ParKey:map[string]*Estimate{
				K3:ball.PartitionGroup(balls,K3),
				K6:ball.PartitionGroup(balls,K6),
			},
			Hole:	 ball.maxHole(),
		}
		ball.corelation1(balls,&pre),
		ball.corelation2(balls,&pre),
		ball.corelation3(balls,&pre),

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

