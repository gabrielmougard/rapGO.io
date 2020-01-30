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
	OutputUUID string `json:"outputUUID"`
}

func DownloadOutput(w http.ResponseWriter, r *http.Request) {
	//TODO : when the client has the right heartbeat, it will call this route to download the output file inside a component.
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	
	//Unmarshal
	var o OutputReq
	err = json.Unmarshal(b, &o)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	filename := setting.TmpFolder()+setting.OutputPrefix()+"_"+o.OutputUUID+setting.OutputSuffix()
	outputdata, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	//fmt.Printf("length of the generated output : %d\n", len(outputdata))
	http.ServeContent(w, r, filename, time.Now(), bytes.NewReader(outputdata))
	//http.ServeFile(w, r, filename)
}
