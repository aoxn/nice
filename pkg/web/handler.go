package web

import (
    "github.com/jinzhu/gorm"
    "github.com/gin-gonic/gin"
    "net/http"
    "fmt"
    "github.com/spacexnice/nice/pkg/algorithm"
)

type WebHandler struct {
    RegURL   string
    DB       * gorm.DB
}

func NewWebHandler(db * gorm.DB) * WebHandler{

    return &WebHandler{
        DB:         db,
    }
}


func (h * WebHandler) Index(c *gin.Context) {
    r := []algorithm.Result{}
    e := h.DB.Limit(15).Find(&r).Error
    if e != nil {
        c.HTML(http.StatusOK,
            "index",
            gin.H{
                "title":      "NICE",
                "results" :     "",
                "error":  fmt.Sprintf("Error occured[%s]",e.Error()),
            },
        )
        return
    }
    c.HTML(http.StatusOK,
        "index",
        gin.H{
            "title":      "NICE",
            "results" :    r,
            "error":  "",
        },
    )
}