package main

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/pili-video-server/streamserver/def"
)

func sendErrorResponse(w http.ResponseWriter, errResp def.ErrorResponse) {
	w.WriteHeader(errResp.HttpSC)
	resStr, _ := json.Marshal(&errResp.Error)
	io.WriteString(w, string(resStr))
}

func sendNormalResponse(w http.ResponseWriter, resp string, sc int) {
	w.WriteHeader(sc)
	io.WriteString(w, resp)
}
