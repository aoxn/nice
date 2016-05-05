package web

import (
    "github.com/jinzhu/gorm"
    "github.com/gin-gonic/gin"
    "net/http"
    "fmt"
    "github.com/spacexnice/nice/pkg/algorithm"
    "github.com/spacexnice/nice/pkg/base"
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
    r := []algorithm.Result{}
    e := h.DB.Limit(15).Order("idx desc").Find(&r).Error
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
    bkt := base.NewBucket(false)
    fmt.Println("xxxx:",len(bkt.Balls))
    for k,_ := range r{
        r[k].LoadJson()
        if r[k].IDX >= len(bkt.Balls)-1 && k == 0 {
            r[k].Ball = base.Ball{Reds:[]int{0,0,0,0,0,0}}
            continue
        }
        r[k].Ball = bkt.Balls[r[k].IDX]
    }
    c.HTML(http.StatusOK,
        "index",
        gin.H{
            "title":       "NICE",
            "results" :    r,
            "error":       "",
        },
    )
}