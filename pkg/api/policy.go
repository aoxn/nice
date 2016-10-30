package api

import (
	"strconv"
	"strings"
)

type Policy interface {
	Found(bkt *Bucket,k int) bool
	Key() string
}

//3区策略控制器
type K3Policy struct {
	Policy
	Target string
}

func (k3 *K3Policy) Found(bkt *Bucket,k int) bool {
	if bkt.Balls[k].K3() == k3.Key() {
		return true
	}
	return false
}

func (k3 *K3Policy) Key() string {
	return k3.Target
}

// 简单策略控制器
type SimplePolicy struct {
	Policy
	Target 		  int
}

func (s *SimplePolicy) Found(bkt *Bucket,k int) bool{

	return bkt.Balls[k].Reds[s.Target - 1] == 1
}

func (s *SimplePolicy) Key() string {

	return strconv.Itoa(s.Target)
}

type ContinuousPolicy struct{
	Policy
	Target        int
}

func (c *ContinuousPolicy) Found(bkt *Bucket,k int) bool {
	if k <= 0 {
		return false
	}
	return bkt.Balls[k].Reds[c.Target - 1] == 1 &&
				bkt.Balls[k - 1].Reds[c.Target - 1] == 1
}

func (s *ContinuousPolicy) Key() string {

	return strconv.Itoa(s.Target)
}

type ContinuousIntervalPolicy struct{
	Policy
	Target        string
}

func (c *ContinuousIntervalPolicy) Found(bkt *Bucket,k int) bool {
	if k <= 0 {
		return false
	}
	cnt := 0
	for _,bi := range bkt.Balls[k].Red(){
		if bkt.Balls[k-1].Reds[bi-1] == 1{
			cnt += 1
		}
	}
	return cnt > 1
}

func (s *ContinuousIntervalPolicy) Key() string {

	return s.Target
}

//两个数字同时出现
type CoLivePolicy struct{
	Policy
	Target        string
}

func (c *CoLivePolicy) Found(bkt *Bucket,k int) bool {
	if k <= 0 {
		return false
	}
	tmp := strings.Split(c.Target,":")
	i,_ := strconv.Atoi(tmp[0])
	j,_ := strconv.Atoi(tmp[1])

	return bkt.Balls[k].Reds[i-1] == 1 && bkt.Balls[k].Reds[j-1] == 1
}

func (s *CoLivePolicy) Key() string {

	return s.Target
}