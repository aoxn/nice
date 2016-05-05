package algorithm

import (
    "github.com/jinzhu/gorm"
    "github.com/spacexnice/nice/pkg/base"
    "github.com/spacexnice/ctlplane/pro/util"
    "time"
    "github.com/golang/glog"
)

type Worker struct {
    Period  time.Duration
    Stop    chan struct{}
    DB *    gorm.DB
}

const (
    WORKER_PERIOD = 1 * time.Hour
)

func NewWorker(db * gorm.DB) *Worker{

    return &Worker{
        DB:     db,
        Period: WORKER_PERIOD,
        Stop:   make(chan struct{}),
    }
}

func (w *Worker) Run() {
    go util.Until(func(){
        bkt := base.NewBucket(true)
        if w.Exist(len(bkt.Balls)){
            return
        }
        s := w.nice(bkt)
        r := Result{
            IDX:len(bkt.Balls),
            K3: s,
            K3S:s.ToJson(),
        }
        w.store(&r)
    },w.Period,w.Stop)
}

func (w *Worker) nice(bkt *base.Bucket) ScoreList{
    return NewPredicator(bkt).PKey3(len(bkt.Balls)-1)
}


func (w *Worker) donice(idx int) ScoreList{
    bkt := base.NewBucket(true)
    //prd.PKey3(idx).NicePrint()
    prd := NewPredicator(bkt)
    return prd.PKey3(idx)
}

func (w *Worker) store(r *Result){
    err := w.DB.Create(r).Error
    if err != nil {
        panic(err)
    }
    return
}

func (w *Worker) Exist(idx int)bool {
    err := w.DB.Find(&Result{IDX:idx}).Error;
    if err == nil{
        glog.Infof("INDEX:[%s] Exist!\n",idx)
        return true
    }
    return false
}