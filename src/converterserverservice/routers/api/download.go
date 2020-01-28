package api

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"bytes"
	"time"
	"rapGO.io/src/converterserverservice/pkg/setting"


)

type OutputReq struct {
	OutputUUID string `json:'outputUUID'`
}

func DownloadOutput(w http.ResponseWriter, r *http.Request) {
	//TODO : when the client has the right heartbeat, it will call this route to download the output file inside a component.
	decoder := json.NewDecoder(r.Body)
	var o OutputReq
	err := decoder.Decode(&o)
	if err != nil {
		panic(err)
	}
	filename := setting.TmpFolder()+setting.OutputPrefix()+"_"+o.OutputUUID+setting.OutputSuffix()
	outputdata, err := ioutil.ReadFile(filename)
	http.ServeContent(w, r, filename, time.Now(), bytes.NewReader(outputdata))
}