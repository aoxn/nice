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
)

const (
	 K3 = "11:11:11"
	 K6 = "6:5:6:5:6:5"
)

type Ball struct {
	//publish date
	Date      string

	//publish index
	Index     int

	// Red Balls
	Reds	  []int

	// Blue Balls
	Blue      int

	Attr      Attribute
}


type Attribute struct {
	ParKey	  map[string]*Partition
	Hole      MaxHole
	CoRelate2 [34][34]int			`json:"-"`
	CoRelate3 [34][34][34]int		`json:"-"`
	AccFreq   []int					`json:"-"`
}

type Partition struct {
	Key 	  string
	Last   	  int
	Total     int
	Avg       int
	Next      int
	Std       float64
	AccStd    float64
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

func (b *Ball) PartitionGroup(balls []Ball,pat string) *Partition{
	var cnt,total,found,accstd = 1,1,false,0.0
	kNum := b.KeyPartition(pat)
	for i:=len(balls)-1;i>=0;i--{
		if balls[i].Attr.ParKey[pat].Key != kNum {
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
	return &Partition{
		Key: 	kNum,
		Last:   cnt,
		Total: 	total,
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
func (b *Ball) frequency(pre Ball)[]int{
	if pre.Attr.AccFreq == nil {
		pre.Attr.AccFreq = make([]int,34,34)
	}
	freq := make([]int,34,34)
	for idx, v := range pre.Attr.AccFreq{
		freq[idx] = v
	}
	for _, i := range b.Reds{
		freq[i] = pre.Attr.AccFreq[i] + 1
	}
	return freq
}

func (b *Ball) corelation2(pre *Ball) [34][34]int {
	coRelate2 := pre.Attr.CoRelate2
	for _, bfirst := range b.Reds {
		for _, bsecond := range b.Reds {
			coRelate2[bfirst][bsecond] = pre.Attr.CoRelate2[bfirst][bsecond] + 1
		}
	}
	return coRelate2
}

func (b *Ball) corelation3(pre *Ball) [34][34][34]int {
	coRelate3 := pre.Attr.CoRelate3
	for _, bfirst := range b.Reds {
		for _, bsecond := range b.Reds {
			for _, bthird := range b.Reds{
				coRelate3[bfirst][bsecond][bthird] = pre.Attr.CoRelate3[bfirst][bsecond][bthird] + 1
			}
		}
	}
	return coRelate3
}

func (b Ball) String() string{
	var m string = ""
	for k,v := range b.Attr.ParKey{
		m += fmt.Sprintf("%2s %+v",k,v) +" ## "
	}
	return fmt.Sprintf("DATE:%s  IDX:%d  REDS:%+2v   BLUE:%2d  [PARKEY: %95s   MAXHOLE: %+2v  FREQENCY: %v]",
						b.Date,b.Index,b.Reds,b.Blue,m[0:len(m)-3],b.Attr.Hole,b.Attr.AccFreq)
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
			ParKey:map[string]*Partition{
				K3:ball.PartitionGroup(balls,K3),
				K6:ball.PartitionGroup(balls,K6),
			},
			Hole:	 ball.maxHole(),
			AccFreq: ball.frequency(pre),
			CoRelate2: ball.corelation2(&pre),
			CoRelate3: ball.corelation3(&pre),
		}
		balls = append(balls,ball)
		return
	})

	return &Bucket{
		Balls:	balls,
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

