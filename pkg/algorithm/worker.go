package algorithm

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
        w.Work(base.NewBucket(true,0))
    },w.Period,w.Stop)
}

func (w *Worker) FillDatabaseTest(){

    for i:= 1950;i< 1951;i ++{
        w.Work(base.NewBucket(false,i))
    }
}


func (w *Worker) Work(bkt *base.Bucket){
    bkt.AddPolicy(base.UK34)

    if bkt.Product && w.Exist(bkt.NextIDX){
        return
    }
    s := w.partition(bkt)

    for _,v := range s {
        offset := 1
        unionkey:= []*base.UnionKey{}
        for _,skey:=range v.SplitKey(){
            u := &base.UnionKey{
                // 5:6/1/6
                Pattern:    base.K11,
                Offset :  	offset,
                Count  :    skey,
                Length :    11,
                Full   :    true,
            }
            bkt.AddPolicy(u)

            glog.Infoln("MMMM:",skey,"  ",u,"   ",v)
            haha := NewPartitionNicer(bkt).predicate(bkt.NextIDX,u)

            //wac := NewRelateNicer(bkt).Predicate(bkt.NextIDX,haha)
            for _,m := range haha{
                glog.Infoln("HAHA:",m)
            }
            offset += 11
            unionkey = append(unionkey,u)
        }
        v.Uk = unionkey
    }

    for _,v := range s{
        glog.Infof("base34UK 11:11:11/1/6::  %s, %+v\n",v.Pattern,v)
    }
    //p := w.relate(bkt,s)
    //glog.Infoln("XP:::",s.ToJson())
    //r := Record{
    //    //本Record属于预测期,所以编号IDX应该为第len(bkt.Balls)
    //    IDX:    bkt.NextIDX,
    //    Index:  bkt.Balls[bkt.NextIDX-1].Index + 1,
    //    K3Json: s.ToJson(),
    //    NiceJson:p.ToJson(),
    //}
    ////glog.Infoln("PRINT:",p.ToJson())
    //w.addResult(&r)
    //w.updatePreviousResult(bkt)
}


func (w *Worker) partition(bkt *base.Bucket) base.ScoreList{
    //return NewPredicator(bkt).PKey3(len(bkt.Balls)-1)

    return NewPartitionNicer(bkt).PKey3(bkt.NextIDX)
}
func (w *Worker) relate(bkt *base.Bucket,list base.ScoreList) base.ScoreList{
    //return NewPredicator(bkt).PKey3(len(bkt.Balls)-1)

    return NewRelateNicer(bkt).Predicate(bkt.NextIDX,list)
}

func (w *Worker) donice(idx int) base.ScoreList{
    bkt := base.NewBucket(true,0)
    //prd.PKey3(idx).NicePrint()
    prd := NewPartitionNicer(bkt)
    return prd.PKey3(idx)
}

func (w *Worker) addResult(r *Record){
    err := w.DB.Create(r).Error
    if err != nil {
        if strings.Contains(err.Error() , "UNIQUE constraint failed"){
            glog.Warningf("预测值已经存在: %s, %+v \n",err.Error(),r)
            return
        }
        panic(err)
    }
    return
}

func (w *Worker) updatePreviousResult(bkt *base.Bucket) {
    r := Record{
        IDX:    bkt.NextIDX - 1,
    }
    err := w.DB.First(&r).Error
    if err != nil{
        if err.Error() != "record not found"{
            glog.Warningln("UNKOWN DB ERROR:",err.Error())
            panic(err)
        }else {
            glog.Warningln("Record not found:",err.Error())
            return
        }
    }
    var nice base.ScoreList
    e  := json.Unmarshal([]byte(r.NiceJson),&nice)
    if e != nil {
        panic(e)
    }
    for _,n := range nice{
        it := bkt.Balls[bkt.NextIDX - 1].Intersection(n.SplitKey())
        n.Cross = fmt.Sprintf("%v|len=[%d]",it,len(it))
    }
    bt,_ := json.Marshal(nice)
    r.NiceJson = string(bt)
    b ,e  := json.Marshal(bkt.Balls[bkt.NextIDX - 1])
    if e != nil{
        panic(e.Error())
    }
    r.BallJson= string(b)
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