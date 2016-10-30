package api

import (
	"strconv"
	"github.com/spacexnice/nice/pkg/util"
	"github.com/golang/glog"
	"strings"
	"fmt"
	"math"
)

type Ball struct {
	//publish date
	Date  string

	//publish index
	Index int

	// Red Balls
	Reds  *[33]int

	// Blue Balls
	Blue  int
}

//每11个数一个分区,共3个分区,计算每个分区的出现的数字的个数, 如 3:2:1
func (ball *Ball) K3() string {
	result, cnt := "", 0
	for i, r := range ball.Reds {
		if r == 1 {
			cnt += 1
		}
		if i == 10 || i == 21 || i == 32 {
			result += fmt.Sprintf("%d:", cnt)
			cnt = 0
		}
	}
	return result[0:len(result) - 1]
}

//返回红球,如[1,4,6,11,24,32]
func (ball *Ball) Red() []int {
	red := []int{}
	for i, r := range ball.Reds {
		if r == 1 {
			red = append(red, i + 1)
		}
	}
	return red
}

type Bucket struct {
	Balls     []*Ball
	TargetIdx int
	Product   bool
}

func found() bool {
	return true
}


//		|
//		|~~~~~~~~~~~~~~~~~ <-- ProbaleLine (Pr线)     	当C线越过这条线后则当前模式出现的概率非常高
//		|
//		|----------------- <-- AverageLine (P0线)		当C线到达这条线后出现概率比较高
//		|
//		|@@@@@@@@@@@@@@@@@@@@@@@   <-- CurrentShowingUp (C线)
//		|
//		|
// --200|________________________  <-- LastShowingUp (L线)
//      |
//		|-.-.-.-.-.-.-.-.- <-- PreviousAveragerLine (P线)
//     TotoalCount  Interval
//
func (bucket *Bucket) Estimate(p Policy) *Estimator {
	first := true
	totalCount := 0.0
	lastLine := 0.0
	currentLine := float64(len(bucket.Balls) - 1)
	for k := len(bucket.Balls) - 1; k > 0; k-- {
		if !p.Found(bucket,k) {
			continue
		}
		if first {
			lastLine = float64(k)
			first = false
		}
		totalCount ++
	}

	interval := lastLine / (totalCount - 1)
	possibleLine := interval * totalCount

	// =================== std error =======================
	previous, cnt, sum := -1.0, 0.0, 0.0
	//std error
	for k := 0; k < len(bucket.Balls) - 1; k ++ {
		if !p.Found(bucket,k) {
			continue
		}
		cnt ++
		if previous < 0 {
			previous = float64(k)
			continue
		}
		i := float64(k) - previous
		//fmt.Printf("Curr: %f, Previous:%f , sub: %f \n",float64(k),previous,i)
		sum += i * i
		previous = float64(k)
	}
	std := math.Sqrt(sum / cnt)
	// ======================================================

	//probableLine2 := possibleLine + std // possibale + std
	probableLine := lastLine + std
	//fmt.Printf("std:%f,currentLine:%f,last:%f, totalCount: %f,  interval:%f, Possible:%f, Probable:%f, Probable2:%f\n",
	//	std, currentLine, lastLine, totalCount, interval, possibleLine, probableLine, probableLine2)
	return &Estimator{
		Key         :    p.Key(),
		AVG         :    interval,
		ProbableLine:    probableLine,
		PossibleLine:    possibleLine,
		CurrentLine :    currentLine,
		LastLine    :   lastLine,
	}
}

func LoadBucket(idx int, force bool) *Bucket {

	util.LoadFile(force)

	var balls []*Ball
	util.ForEachLine(func(line string, i int) {
		glog.Infof("处理第 %d 个\n", i)
		l := strings.Split(line, " ")
		idx, _ := strconv.Atoi(l[0])
		r1, _ := strconv.Atoi(l[2])
		r2, _ := strconv.Atoi(l[3])
		r3, _ := strconv.Atoi(l[4])
		r4, _ := strconv.Atoi(l[5])
		r5, _ := strconv.Atoi(l[6])
		r6, _ := strconv.Atoi(l[7])
		b1, _ := strconv.Atoi(l[8])
		ball := &Ball{
			Date:    l[1],
			Index:  idx,
			Reds:    util.SetBall([]int{r1, r2, r3, r4, r5, r6}),
			Blue:    b1,
		}

		//ball.Hole   = ball.maxHole()
		balls = append(balls, ball)
		return
	})

	bkt := &Bucket{
		Balls:     balls,
	}
	return bkt
}

