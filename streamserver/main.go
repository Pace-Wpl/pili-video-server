package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/pili-video-server/streamserver/def"
)

type middleWareHandler struct {
	r  *httprouter.Router
	cl *ConnLimiter
}

func (m middleWareHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//middle ware handler
	if !m.cl.GetConn() {
		sendErrorResponse(w, def.ErrorConnectLimit)
		return
	}

	//允许访问域
	w.Header().Set("Access-Control-Allow-Origin", "http://172.19.21.3:3000")
	w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type,Depth, User-Agent, X-Session-Id, X-User-Name, If-Modified-Since, Cache-Control, Origin")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("access-control-allow-methods", "GET, POST, OPTIONS, PUT, DELETE")
	m.r.ServeHTTP(w, r)
	//releaseConnect
	defer m.cl.ReleaseConn()

}

func NewMiddleWareHandler(r *httprouter.Router, cc int) http.Handler {
	m := middleWareHandler{}
	m.r = r
	m.cl = NewConnLimiter(cc)
	return m
}

func RegisterHandler() *httprouter.Router {
	router := httprouter.New()

	router.GET("/test", testPageHandler)

	router.GET("/video/:vid_id/:mat", streamHandler)
	router.POST("/video/:vid_id", uploadHandler)

	router.GET("/icon/:vid_id", GetIcon)

	return router
}

func main() {
	r := RegisterHandler()
	mh := NewMiddleWareHandler(r, def.STREAM_LIMIT)
	http.ListenAndServe(":9000", mh)
}
