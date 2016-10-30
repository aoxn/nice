package api

import (
	"sort"
	"fmt"
	"strconv"
	"strings"
	"time"
	"os"
	"io"
)

type Picker struct {

}

func AreaPicker(bkt *Bucket)*Estimators{
	e1 := bkt.Estimate(&K3Policy{Target:"2:2:2"})
	e6 := bkt.Estimate(&K3Policy{Target:"2:1:3"})
	e7 := bkt.Estimate(&K3Policy{Target:"2:3:1"})
	e2 := bkt.Estimate(&K3Policy{Target:"1:2:3"})
	e3 := bkt.Estimate(&K3Policy{Target:"1:3:2"})
	e4 := bkt.Estimate(&K3Policy{Target:"3:2:1"})
	e5 := bkt.Estimate(&K3Policy{Target:"3:1:2"})
	e8 := bkt.Estimate(&K3Policy{Target:"3:3:0"})
	k3 := Estimators{e1,e2,e3,e4,e5,e6,e7,e8}
	sort.Sort(k3)
	return &k3
}

func SimplePicker(bkt *Bucket)*Estimators{

	simple := Estimators{}
	for i:=1;i<=33;i++ {
		e := bkt.Estimate(&SimplePolicy{Target: i})
		simple = append(simple,e)
	}
	sort.Sort(simple)
	return &simple
}

func ContinuousPicker(bkt *Bucket)*Estimators{

	simple := Estimators{}
	for i:=1;i<=33;i++ {
		e := bkt.Estimate(&ContinuousPolicy{Target: i})
		simple = append(simple,e)
	}
	sort.Sort(simple)
	return &simple
}

func ContinuousIntervalPicker(bkt * Bucket) * Estimators{

	return &Estimators{
		bkt.Estimate(&ContinuousIntervalPolicy{Target:"持续间隔"}),
	}
}

func CoLivePicker(bkt *Bucket) *Estimators{
	r := Estimators{}
	for i := 1;i <=33;i++{
		for j:=i+1;j<=33;j++{
			r = append(r,bkt.Estimate(&CoLivePolicy{Target:fmt.Sprintf("%d:%d",i,j)}))
		}
	}
	sort.Sort(r)
	return &r
}

func Output(desc string,target interface{}) {
	fmt.Println(desc)
	fmt.Println(target)
}

func Pick(bkt *Bucket) {

	//分区项集
	area := AreaPicker(bkt)

	//单独频繁项集
	simp := SimplePicker(bkt)

	//连续频繁项集
	cont := ContinuousPicker(bkt)

	//连续频繁间隔项集
	cti  := ContinuousIntervalPicker(bkt)

	//协同频繁项集
	colive:= CoLivePicker(bkt)

	Output("协同频繁项集",colive)
	Output("分区项集", area)
	Output("单独频繁项集",simp)
	Output("连续频繁项集",cont)
	Output("连续频繁间隔项集",cti)

	Output("单独频繁项集 [预测>0]",simp.ZeroLine())
	Output("连续频繁项集 [预测>0]",cont.ZeroLine())

	//连续出现的期望点与单独出现的期望点的交集
	Output("[连续频繁项集] 与 [单独频繁项集] 的交集",
		simp.ZeroLine().Cross(cont.ZeroLine()),
	)

	ver := Vertical(bkt.Balls[len(bkt.Balls)-1],cont.ZeroLine())
	Output("[连续频繁项集] 预测",ver)

	a1 := append(simp.ZeroLine().Cross(cont.ZeroLine()),
		simp.ZeroLine().UminusN(cont.ZeroLine())...,
	)

	Hint("(单独频繁项集) ## (连续频繁项集) WithOut Predict", area.ZeroLine(), a1)

	Hint("(单独频繁项集) ## (连续频繁项集) With    Predict", area.ZeroLine(), append(ver,a1...))

	for _,v := range ver{
		Hint("(协同频繁项集) With    Predict",area.ZeroLine(),append(Estimators{v},Horizontal(v,colive)...))
	}

	Write2File("\n")
}

// 上一期[a,b,c,d,e,f]  ==> 连续频繁项集
func Vertical(ball *Ball,est Estimators) Estimators{
	result := Estimators{}
	red := ball.Red()
	for _,v := range est{
		i,_ := strconv.Atoi(v.Key)
		for _,n := range red{
			if i == n {
				result = append(result, v)
			}
		}
	}
	return result
}

// 协同频繁项列表
func Horizontal(e *Estimator,est *Estimators) Estimators{
	result := Estimators{}
	for _,v := range *est{
		found,key := false,""
		tmp := strings.Split(v.Key,":")
		for _,n := range tmp{
			if e.Key == n{
				found = true
			}else {
				key = n
			}
		}
		if !found{
			continue
		}
		v.Key = key
		result = append(result,v)
	}
	return result
}

func Hint(desc string,area Estimators, target Estimators) []string{
	response := []string{}
	for _,ar := range area{
		result := ""
		tmp := strings.Split(ar.Key,":")
		t1,_ := strconv.Atoi(tmp[0])
		t2,_ := strconv.Atoi(tmp[1])
		t3,_ := strconv.Atoi(tmp[2])
		for _,b := range target{
			i,_ := strconv.Atoi(b.Key)
			if i<=11 && i>=1 && t1 >0 {
				t1--
				result += fmt.Sprintf("%d:",i)
			}
			if i<=22 && i>=12 && t2 >0 {
				t2--
				result += fmt.Sprintf("%d:",i)
			}
			if i<=33 && i>=23 && t3 >0 {
				t3--
				result += fmt.Sprintf("%d:",i)
			}
		}
		msg := fmt.Sprintf("Hint:[%s] [%s] [%d] %s => %s\n",
			time.Now().Format("2006-01-02 15:04:05"),desc,int64(ar.CurrentLine),ar.Key,result[0:len(result)-1],
		)
		fmt.Printf(msg)
		Write2File(msg)
		response = append(response, result[0:len(result)-1])
		//break
	}
	return response
}

var RESULT_FILE = "PREDICT.TXT"

func Write2File(msg string) {
	file,err := os.OpenFile(RESULT_FILE,os.O_CREATE|os.O_APPEND|os.O_RDWR,0660)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	io.WriteString(file,msg)
}