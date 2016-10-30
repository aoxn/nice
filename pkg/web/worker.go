package web

import (
	"github.com/jinzhu/gorm"
	"github.com/spacexnice/nice/pkg/base"
	"time"
	"github.com/spacexnice/nice/pkg/util"
	"encoding/json"
	//"fmt"
	"github.com/golang/glog"
	"strings"
	"fmt"
)

type Worker struct {
	Period time.Duration
	Stop   chan struct{}
	DB     *gorm.DB
}

const (
	WORKER_PERIOD = 1 * time.Hour
)

func NewWorker(db *gorm.DB) *Worker {

	return &Worker{
		DB:     db,
		Period: WORKER_PERIOD,
		Stop:   make(chan struct{}),
	}
}

func (w *Worker) Run() {
	// IDX 编号从0开始算
	go util.Until(func() {
		w.Work(base.NewBucket(false, -1))
	}, w.Period, w.Stop)
}

func (w *Worker) FillDatabaseTest() {

	for i := 1950; i < 1951; i ++ {
		w.Work(base.NewBucket(false, i))
	}
}

func (w *Worker) Work(bkt *base.Bucket) {
	result := bkt.Nice(-1)

	bkt.Statistic()

	glog.Infoln("++++++++++++++++++++++++ [RESULT] ++++++++++++++++++++++++++\nAfter PdtGroup Merge: ")
	result.NicePrint()
}

func (w *Worker) addResult(r *Record) {
	err := w.DB.Create(r).Error
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			glog.Warningf("预测值已经存在: %s, %+v \n", err.Error(), r)
			return
		}
		panic(err)
	}
	return
}

func (w *Worker) updatePreviousResult(bkt *base.Bucket) {
	r := Record{
		IDX:    bkt.TargetIdx - 1,
	}
	err := w.DB.First(&r).Error
	if err != nil {
		if err.Error() != "record not found" {
			glog.Warningln("UNKOWN DB ERROR:", err.Error())
			panic(err)
		} else {
			glog.Warningln("Record not found:", err.Error())
			return
		}
	}
	var nice base.RankList
	e := json.Unmarshal([]byte(r.NiceJson), &nice)
	if e != nil {
		panic(e)
	}
	for _, n := range nice {
		it := bkt.Balls[bkt.TargetIdx - 1].Intersection(n.SplitKey())
		n.Cross = fmt.Sprintf("%v|len=[%d]", it, len(it))
	}
	bt, _ := json.Marshal(nice)
	r.NiceJson = string(bt)
	b, e := json.Marshal(bkt.Balls[bkt.TargetIdx - 1])
	if e != nil {
		panic(e.Error())
	}
	r.BallJson = string(b)
	err = w.DB.Save(&r).Error
	if err != nil {
		panic(err)
	}
}

func (w *Worker) Exist(idx int) bool {
	r := &Record{
		IDX:    idx,
	}
	err := w.DB.Where(r).Find(r).Error;
	if err == nil {
		glog.Warningln("INDEX:[%d] Exist!\n", idx)
		return true
	}

	return false
}