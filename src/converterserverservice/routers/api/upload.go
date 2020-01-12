package api

import (
	"net/http"
	"fmt"
	"io"
	"os"
	"encoding/json"

	"rapGO.io/src/converterserverservice/pkg/setting"
	"rapGO.io/src/converterserverservice/pkg/uuid"

)

type OutputUUID struct {
	Status int `json:'status'`
	OutputUIID string `json:'outputUUID'`
}

func UploadInputBLOB(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20) // limit your max input length!
    // in your case file would be fileupload
    file, _, err := r.FormFile("file")
    if err != nil {
        panic(err)
    }
    defer file.Close()
	filename := uuid.NewVoiceUUID() //use UUID here
	out, err := os.Create(setting.TmpFolder()+filename)
	if err != nil {
		fmt.Printf("Unable to create the file for writing. Check your write access privilege")
		return
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		fmt.Println(err)
	}
	res := OutputUUID{Status: 200, OutputUIID: filename}
	js, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type","application/json")
	w.Write(js)

}