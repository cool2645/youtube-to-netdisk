package http

import (
	"github.com/BurntSushi/toml"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
	"github.com/urfave/negroni"
	"net/http"
	"strconv"
)

var config Config

type WebInterface struct{}

func (wi WebInterface) Driver() string {
	return "WebInterface"
}

func (wi WebInterface) Start() {
	_, err := toml.DecodeFile("http.toml", &config)
	if err != nil {
		panic(err)
	}

	mux := httprouter.New()
	mux.GET(config.API_PREFIX+"/", pong)

	mux.POST(config.API_PREFIX+"/trigger", triggerTask)
	mux.POST(config.API_PREFIX+"/cancel/:id", cancelTask)
	mux.GET(config.API_PREFIX+"/retry/:id", retryTask)
	mux.POST(config.API_PREFIX+"/retry/:id", retryTask)
	mux.GET(config.API_PREFIX+"/task", getTasks)
	mux.GET(config.API_PREFIX+"/task/:id", getTask)
	mux.GET(config.API_PREFIX+"/task/:id/log", getTaskLog)
	mux.GET(config.API_PREFIX+"/keyword", getKeywords)
	mux.POST(config.API_PREFIX+"/keyword", newKeyword)
	mux.DELETE(config.API_PREFIX+"/keyword", deleteKeyword)

	mux.ServeFiles("/static/*filepath", http.Dir("static"))

	mux.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, config.VIEW_DIR+"/index.html")
	})

	c := cors.New(cors.Options{
		AllowedOrigins:   config.ALLOW_ORIGIN,
		AllowedMethods:   []string{"GET", "POST", "OPTIONS", "PUT", "DELETE"},
		AllowCredentials: true,
	})
	h := c.Handler(mux)

	n := negroni.New()
	n.Use(negroni.NewStatic(http.Dir("app")))
	n.UseHandler(h)

	http.ListenAndServe(":"+strconv.FormatInt(config.PORT, 10), n)
}
