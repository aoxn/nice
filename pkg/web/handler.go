package web

import (
    "github.com/jinzhu/gorm"
    "github.com/gin-gonic/gin"
    "net/http"
    "fmt"
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

    c.HTML(http.StatusOK,
        "index",
        gin.H{
            "title":     "NICE",
            "balls" :     []string{"x","y"},
            "has":       false,
            "currtag":   "",
            "errorinfo": fmt.Sprintf("List Repository Error with search"),
        },
    )
}