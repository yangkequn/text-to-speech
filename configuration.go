package main

import (
	"os"
)

var SPEECH_KEY string
var SPEECH_REGION string

func init() {
	SPEECH_KEY = os.Getenv("SPEECH_KEY")
	SPEECH_REGION = os.Getenv("SPEECH_REGION")
	//panic if not set
	if SPEECH_KEY == "" || SPEECH_REGION == "" {
		panic("SPEECH_KEY or SPEECH_REGION not set")
	}
}
