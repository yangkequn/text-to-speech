package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", handleRequest)
	fmt.Println("Listening on port 80...")
	fmt.Println("Please open http://ip:port/ssml=xxx in your browser.")
	fmt.Println("Visit https://learn.microsoft.com/en-us/azure/ai-services/speech-service/speech-synthesis-markup for more information.")

	log.Fatal(http.ListenAndServe(":80", nil))
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	ssml := r.FormValue("ssml")
	if ssml == "" {
		DemoUsage := `
		<!DOCTYPE html>
		<html>
		<head>
		<title>SSML to Audio Info</title>
		</head>
		<body>
		<h1>SSML to Audio Info</h1>
		<p>Enter SSML to get audio info.</p>
		<p>For more information, please refer to the <a href="https://learn.microsoft.com/en-us/azure/ai-services/speech-service/speech-synthesis-markup">documentation</a>.</p>
		<form action="/" method="GET">
		<input type="text" name="ssml" />
		<input type="submit" value="Submit" />
		</form>
		</body>
		</html>
		`
		fmt.Fprintf(w, DemoUsage)
		return
	}
	if audioInfo, err := TTSInfosToSpeech(ssml); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else if audioInfo != nil {
		//send audioInfo back to client
		dataJson, err := json.Marshal(audioInfo)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, string(dataJson))
	} else {
		http.Error(w, "No audio info found for the given SSML.", http.StatusNotFound)
	}
	return
}
