package web

import (
    "github.com/jinzhu/gorm"
    "github.com/gin-gonic/gin"
    "net/http"
    "fmt"
)

type WebHandler struct {
    DB       * gorm.DB
}

func NewWebHandler(db * gorm.DB) * WebHandler{

    return &WebHandler{
        DB:         db,
    }
}


func (h * WebHandler) Index(c *gin.Context) {
    res,rec := []Result{},[]algorithm.Record{}
    e := h.DB.Limit(10).Order("idx desc").Find(&rec).Error
    if e != nil {
        c.HTML(http.StatusOK,
            "index",
            gin.H{
                "title":        "NICE",
                "results" :     "",
                "error":        fmt.Sprintf("Error occured[%s]",e.Error()),
            },
        )
        return
    }
    //bkt := base.NewBucket(false)
    //fmt.Println("xxxx:",len(bkt.Balls))
    for k,_ := range rec{
        res = append(res,rec[k].LoadResult())
    }
    c.HTML(http.StatusOK,
        "index",
        gin.H{
            "title":       "NICE",
            "results" :    res,
            "error":       "",
        },
    )
}