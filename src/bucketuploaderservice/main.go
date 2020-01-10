package main

import (

	"rapGO.io/src/bucketuploaderservice/pkg/eventproc"
)

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {

	eventproc.ProcessEvents()

}


