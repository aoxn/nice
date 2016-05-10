package algorithm

import (
    "github.com/jinzhu/gorm"
    "github.com/spacexnice/nice/pkg/base"
    "time"
    "github.com/spacexnice/nice/pkg/util"
    "encoding/json"
    "fmt"
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
    // IDX 编号从0开始算
    go util.Until(func(){
        bkt := base.NewBucket(true)
        //if w.Exist(len(bkt.Balls)){
        //    return
        //}
        s := w.nice(bkt)
        r := Record{
            //本Record属于预测期,所以编号IDX应该为第len(bkt.Balls)
            IDX:    len(bkt.Balls),
            Index:  bkt.Balls[len(bkt.Balls)-1].Index + 1,
            K3Json: s.ToJson(),
        }
        w.addResult(&r)
        w.updatePreviousResult(bkt)
    },w.Period,w.Stop)
}

func (w *Worker) nice(bkt *base.Bucket) ScoreList{
    //return NewPredicator(bkt).PKey3(len(bkt.Balls)-1)

    return NewPredicator(bkt).PKey3(len(bkt.Balls)-1)
}


func (w *Worker) donice(idx int) ScoreList{
    bkt := base.NewBucket(true)
    //prd.PKey3(idx).NicePrint()
    prd := NewPredicator(bkt)
    return prd.PKey3(idx)
}

func (w *Worker) addResult(r *Record){
    err := w.DB.Create(r).Error
    if err != nil {
        panic(err)
    }
    return
}

func (w *Worker) updatePreviousResult(bkt *base.Bucket) {
    r := Record{
        IDX:    len(bkt.Balls)-1,
    }
    err := w.DB.First(&r).Error
    if err != nil{
        if err.Error() != "record not found"{
            glog.Warningln("UNKOWN DB ERROR:",err.Error())
        }
        panic(err)
    }
    glog.Warningln("rrrr: ",r,"  PPPPPPPP:")
    b ,e  := json.Marshal(bkt.Balls[len(bkt.Balls)-1])
    if e != nil{
        panic(e.Error())
    }
    r.BallJson= string(b)
    fmt.Println("BALLJSON:",r.BallJson)
    err = w.DB.Save(&r).Error
    if err != nil {
        panic(err)
    }
}

func (w *Worker) Exist(idx int)bool {
    r := &Record{
        IDX:    idx,
    }
    err := w.DB.Where(r).Find(r).Error;
    if err == nil{
        glog.Warningln("INDEX:[%d] Exist!\n",idx)
        return true
    }

    return false
}