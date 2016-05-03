package app

import (
    "net"
    "github.com/jinzhu/gorm"
    "os"
    "flag"
    "github.com/gin-gonic/gin"
    "github.com/spacexnice/nice/pkg/util"
    "github.com/spacexnice/nice/pkg/web"
)


var (
    ENV_DB_DATA_PATH= "DB_DATA_PATH"
)

type Config struct {
    Port              int
    Address           net.IP

    DataPath          string

    EnableProfiling   bool
    BindPodsQPS       float32
    BindPodsBurst     int
}

type NiceServer struct {
    Cnf        * Config
    // interface to access sqlite db
    DB         * gorm.DB
    eng        * gin.Engine
    Handler    * web.WebHandler
}

func NewNiceServer() *NiceServer {
    cnf := createConfig()
    db  := util.OpenInit(cnf.DataPath)

    return &NiceServer{
        DB:      db,
        Cnf:     cnf,
        eng:     gin.Default(),
        Handler: web.NewWebHandler(db),
    }
}



func createConfig()* Config{
    dpath := os.Getenv(ENV_DB_DATA_PATH)
    if dpath == ""{
        p,_ := os.Getwd()
        dpath = p
    }
    return &Config{
        DataPath:   dpath,
        Port:       8080,
        Address:    net.IPv4(0,0,0,0),
    }
}

func (s *NiceServer) AddFlags(){
    flag.Set("logtostderr", "true")
}
func (s *NiceServer) route(){
    r := s.eng
    r.Static("/js", "js")
    r.Static("/css", "css")
    r.Static("/fonts", "fonts")
    r.LoadHTMLGlob("pages/*.html")

    r.GET("/",s.Handler.Index)


    //// Authorization group
    //// authorized := r.Group("/", AuthRequired())
    //// exactly the same as:
    //authorized := r.Group("/")
    //// per group middleware! in this case we use the custom created
    //// AuthRequired() middleware just in the "authorized" group.
    //authorized.Use(AuthRequired())
    //{
    //    authorized.POST("/login", loginEndpoint)
    //    authorized.POST("/submit", submitEndpoint)
    //    authorized.POST("/read", readEndpoint)
    //
    //    // nested group
    //    testing := authorized.Group("testing")
    //    testing.GET("/analytics", analyticsEndpoint)
    //}
}


func (s *NiceServer) Run(){
    s.route()
    s.eng.Run(":8000")
}