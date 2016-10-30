package api

import (
	"fmt"
	"math"
	//"github.com/golang/glog"
)

type Estimator struct {
	Key          string
	//AVG
	AVG          float64
	//非常可能出现的位置
	ProbableLine float64

	//很可能出现的位置
	PossibleLine float64

	//按平均算法预测的下一期出现位置
	AverageLine  float64

	//当前期数所处位置
	CurrentLine  float64

	//最近一次出现的位置
	LastLine     float64
}

func (e *Estimator) Po() float64 {
	return e.CurrentLine - e.PossibleLine
}

func (e *Estimator) Pr() float64 {
	return e.CurrentLine - e.ProbableLine
}

//无量纲化的非常可能的位置
func (e *Estimator) DimensionlessPo() float64 {
	return e.Po()/e.AVG
}

//无量纲化的可能的位置
func (e *Estimator) DimensionlessPr() float64 {
	return e.Pr()/e.AVG
}

func (e *Estimator) String() string {
	return fmt.Sprintf("Key:%5s,   上次位置:%12f,    当前位置:%12f,    可能的位置:%12f (%10f [%10f]),    非常可能的位置:%12f (%10f [%10f]),  平均值:%12f\n",
		e.Key,e.LastLine,e.CurrentLine, e.PossibleLine,e.Po(),e.DimensionlessPo(), e.ProbableLine,e.Pr(),e.DimensionlessPr(),e.AVG )
}

type Estimators []*Estimator

// Len is the number of elements in the collection.
func (e Estimators) Len() int {
	return len(e)
}

const DIFF=0.1

func (e Estimators) ZeroLine() Estimators{
	result := Estimators{}
	for _,m := range e{
		if m.DimensionlessPo() + DIFF > 0 {
			result = append(result,m)
		}
	}
	return result
}

func (e Estimators) Cross(et Estimators) Estimators{
	result := Estimators{}
	for _,v1 := range e{
		for _,v2 := range et{
			if v1.Key == v2.Key{
				result = append(result, v1)
			}
		}
	}
	return result
}

func (e Estimators) UminusN(et Estimators) Estimators{
	result := Estimators{}
	for _,v1 := range e{
		found := false
		for _,v2 := range et{
			if v1.Key == v2.Key{
				found = true
				break
			}
		}
		if ! found {
			result = append(result,v1)
		}
	}

	for _,v1 := range et{
		found := false
		for _,v2 := range e{
			if v1.Key == v2.Key{
				found = true
				break
			}
		}
		if ! found {
			result = append(result,v1)
		}
	}
	return result
}

// Less reports whether the element with
// index i should sort before the element with index j.
func (e Estimators) Less(i, j int) bool {
	ai,aj := e[i],e[j]
	if math.Abs(ai.ProbableLine - aj.ProbableLine) < DIFF{
		// equal
		if ai.PossibleLine < aj.PossibleLine{
			return true
		}
		return false
	}else if ai.ProbableLine < aj.ProbableLine{
		return true
	}
	return false
}
// Swap swaps the elements with indexes i and j.
func (e Estimators) Swap(i, j int) {
	e[i],e[j] = e[j],e[i]
}

func (e Estimators) String() string {
	result := ""
	for _, v := range e{
		result += v.String()
	}
	return result
}
