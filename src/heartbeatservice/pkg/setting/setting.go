package setting

import (
	"os"
	"errors"
)

func LastHeartbeatDesc() string {
	v, ok := os.LookupEnv("LAST_HEARTBEAT_DESC")
	if !ok {
		panic(errors.New("the last heartbeat description is not detected."))
	}
	return v
}