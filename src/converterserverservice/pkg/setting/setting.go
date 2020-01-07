package setting

import (
	"log"
	"time"
	"os"
	"github.com/go-ini/ini"

	"rapGO.io/src/audioconverterservice/pkg/bucket"
)

type App struct {
	TmpFolder string //path to tmp storage
	InputPrefix string //name convention for the input filename (i.e : input_)
	InputSuffix string // base file extension (i.e : mp3)
}

var AppSetting = &App{}

type Server struct {
	RunMode string
	HttpPort string
	ReadTimeout time.Duration
	WriteTimeout time.Duration
}

var ServerSetting = &Server{}

type Storage struct {
	ProjectID string
	BucketName string
}

var StorageSetting = &Storage{}

var cfg *ini.File
var bi *bucket.BucketInterface

func Setup() {
	if !environmentDetected() {
		//if at least one env variable is not set, we use the app.ini file (meant to be used in development mode)
		var err error
		cfg, err = ini.Load("conf/app.ini")
		if err != nil {
			log.Fatalf("setting.Setup, fail to parse 'conf/app.ini': %v", err)
		}
		mapTo("app", AppSetting)
		mapTo("server", ServerSetting)
		mapTo("storage", StorageSetting)
	}
	// create the buket interface for google storage
	bi, err = bucket.NewBucketInterface(StorageSetting.ProjectID, StorageSetting.BucketName)
	if err != nil {
		log.Fatalf("setting.Setup, fail to create BucketInterface: %v", err)
	}
	
}

func mapTo(section string, v interface{}) {
	err := cfg.Section(section).MapTo(v)
	if err != nil {
		log.Fatalf("Cfg.MapTo %s err: %v", section, err)
	}
}

func environmentDetected() bool {
	//App related variables
	if v, ok := os.LookupEnv("TMP_FOLDER"); ok {
		AppSetting.TmpFolder = v
	} else {
		return false
	}

	if v, ok := os.LookupEnv("INPUT_PREFIX"); ok {
		AppSetting.InputPrefix = v
	} else {
		return false
	}
	///////////////////////
	// Server related variables
	if v, ok := os.LookupEnv("SERVER_RUN_MODE"); ok {
		ServerSetting.RunMode = v
	} else {
		return false
	}
	
	if v, ok := os.LookupEnv("SERVER_HTTP_PORT"); ok {
		ServerSetting.HttpPort = v
	} else {
		return false
	}

	if v, ok := os.LookupEnv("SERVER_READ_TIMEOUT"); ok {
		ServerSetting.ReadTimeout = v
	} else {
		return false
	}

	if v, ok := os.LookupEnv("SERVER_WRITE_TIMEOUT"); ok {
		ServerSetting.WriteTimeout = v
	} else {
		return false
	}
	//////////////////////
	// Storage(google bucket) related variables
	if v, ok := os.LookupEnv("STORAGE_PROJECT_ID"); ok {
		StorageSetting.ProjectID = v
	} else {
		return false
	}

	if v, ok := os.LookupEnv("STORAGE_BUCKET_NAME"); ok {
		StorageSetting.BucketName = v
	} else {
		return false
	}
	/////////////////////

	//finally, return true (meaning there were no errors)
	return true
}

func GetVoiceUUIDPrefix() string {
	return AppSetting.InputPrefix
}

func GetVoiceUUIDSuffix() string {
	return AppSetting.InputSuffix
}