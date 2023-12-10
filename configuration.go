package main

import (
	"fmt"
	"os"
)

var SPEECH_KEY string
var SPEECH_REGION string

func init() {
	SPEECH_KEY = os.Getenv("SPEECH_KEY")
	SPEECH_REGION = os.Getenv("SPEECH_REGION")
	//panic if not set
	if SPEECH_KEY == "" || SPEECH_REGION == "" {
		panic("SPEECH_KEY or SPEECH_REGION are failed to be loaded from environment variables.")
	} else if len(SPEECH_KEY) < 10 {
		fmt.Println("invalid  SPEECH_KEY, too short", SPEECH_KEY)
		return
	} else if len(SPEECH_REGION) < 3 {
		fmt.Println("invalid  SPEECH_REGION, too short", SPEECH_REGION)
		return
	}
	fmt.Println("SPEECH_KEY and SPEECH_REGION are loaded from environment variables.")
	//display SPEECH_KEY and SPEECH_REGION,  but only leading and trailing 3 characters, with the rest replaced by asterisks
	fmt.Println("SPEECH_KEY:", SPEECH_KEY[:3]+"********"+SPEECH_KEY[len(SPEECH_KEY)-3:])
	fmt.Println("SPEECH_REGION:", SPEECH_REGION[:3]+"********"+SPEECH_REGION[len(SPEECH_REGION)-3:])

}
