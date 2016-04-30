package base

import (
	"github.com/spacexnice/nice/pkg/util"
	"os"
	"bufio"
	"k8s.io/kubernetes/Godeps/_workspace/src/github.com/golang/glog"
	"io"
	"strings"
	"strconv"
	"sort"
	"math"
	"fmt"
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
	K3		string
	K6 		string
	Hole    MaxHole
}

type MaxHole struct {
	Start 	int
	End 	int
	Middle 	int
	Len     int
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
		fmt.Printf("%+v\n",b)
	}
}

func (b *Ball) keyPartition(pat int)string{
	secs := int(math.Ceil(float64(33)/float64(pat)))
	m := make(map[int]int)
	for i:=0;i<secs;i++{
		m[i] = 0
	}
	for _,i := range b.Reds{
		x := (i-1)/pat
		m[x] = m[x] +1
	}
	var rs = ""
	for i:=0;i<secs;i++{
		rs += fmt.Sprintf("%d-",m[i])
		//fmt.Println(fmt.Sprintf("%d:",m[i]))
	}
	return rs[0:len(rs)-1]
}

func (b *Ball) maxHole() MaxHole{
	pre,start,end,len := 0,0,0,0
	for i := range b.Reds{
		if (i - pre) > len{
			start,end,len = pre,i,(i - start)
		}
		pre = i
	}
	return MaxHole{
		Start:		start,
		End:		end,
		Middle:     (end-start)>>2,
		Len:        (end-start),
	}
}

func (b *Ball) String() string{

	return fmt.Sprintf("DATE:%*s  IDX:%s  REDS:%+v   BLUE:%s  [K3: %s   K6: %s   MAXHOLE:%+v]",
						b.Date,b.Index,b.Reds,b.Blue,b.Attr.K3,b.Attr.K6,b.Attr.Hole)
}

func Intersection(b1,b2 Ball) []int{
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

		ball.Attr = Attribute{
			K3:   	ball.keyPartition(11),
			K6: 	ball.keyPartition(6),
			Hole:	ball.maxHole(),
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

