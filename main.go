package main

import (
	"github.com/BurntSushi/toml"
	"net/http"
	"github.com/urfave/negroni"
	"strconv"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/yanzay/log"
	"github.com/cool2645/youtube-to-netdisk/model"
	"github.com/cool2645/youtube-to-netdisk/handler"
	. "github.com/cool2645/youtube-to-netdisk/config"
	"github.com/cool2645/youtube-to-netdisk/broadcaster"
)

var mux = httprouter.New()

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "app/index.html")
}

func main() {

	_, err := toml.DecodeFile("config.toml", &GlobCfg)
	if err != nil {
		panic(err)
	}

	db, err := gorm.Open("mysql", ParseDSN(GlobCfg))
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("Database init done")
	defer db.Close()

	db.AutoMigrate(&model.Keyword{}, &model.Task{}, &model.Subscriber{})
	model.Db = db

	if GlobCfg.RIRI_ENABLE {
		go broadcaster.ServeRiri(model.Db, GlobCfg.RIRI_ADDR, GlobCfg.RIRI_KEY)
	}

	mux.GET("/api", handler.Pong)

	mux.POST("/api/trigger", handler.TriggerTask)
	mux.POST("/api/kill/:id", handler.KillTask)
	mux.GET("/api/retry/:id", handler.Retry)
	mux.POST("/api/retry/:id", handler.Retry)
	mux.GET("/api/running", handler.GetRunningTaskStatus)
	mux.GET("/api/task", handler.GetTasks)
	mux.GET("/api/task/:id", handler.GetTask)
	mux.GET("/api/task/:id/log", handler.GetTaskLog)
	mux.GET("/api/keyword", handler.GetKeywords)
	mux.POST("/api/keyword", handler.NewKeyword)
	mux.DELETE("/api/keyword", handler.DeleteKeyword)

	mux.ServeFiles("/static/*filepath", http.Dir("static"))

	mux.NotFound = http.HandlerFunc(NotFoundHandler)

	c := cors.New(cors.Options{
		AllowedOrigins:   GlobCfg.ALLOW_ORIGIN,
		AllowedMethods:   []string{"GET", "POST", "OPTIONS", "PUT", "DELETE"},
		AllowCredentials: true,
		//AllowedHeaders: []string{""},
	})
	h := c.Handler(mux)

	n := negroni.New()
	n.Use(negroni.NewStatic(http.Dir("app")))
	n.UseHandler(h)

	http.ListenAndServe(":"+strconv.FormatInt(GlobCfg.PORT, 10), n)

}
