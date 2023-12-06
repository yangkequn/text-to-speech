package main

import (
	"fmt"
	"time"

	"github.com/Microsoft/cognitive-services-speech-sdk-go/audio"
	"github.com/Microsoft/cognitive-services-speech-sdk-go/common"
	"github.com/Microsoft/cognitive-services-speech-sdk-go/speech"
)

func synthesizeStartedHandler(event speech.SpeechSynthesisEventArgs) {
	defer event.Close()
	fmt.Println("Synthesis started.")
}

func synthesizingHandler(event speech.SpeechSynthesisEventArgs) {
	defer event.Close()
	//fmt.Printf("Synthesizing, audio chunk size %d.\n", len(event.Result.AudioData))
}

func synthesizedHandler(event speech.SpeechSynthesisEventArgs) {
	defer event.Close()
	//fmt.Printf("Synthesized, audio length %d.\n", len(event.Result.AudioData))
}

func cancelledHandler(event speech.SpeechSynthesisEventArgs) {
	defer event.Close()
	fmt.Println("Received a cancellation.")
}

type TTSResult struct {
	// AudioData presents the synthesized audio.
	AudioData []byte

	// AudioDuration presents the time duration of synthesized audio.
	AudioDurationSecn time.Duration
}

func TTSInfosToSpeech(ssml string) (result *TTSResult, err error) {
	var (
		audioConfig       *audio.AudioConfig
		speechConfig      *speech.SpeechConfig
		speechSynthesizer *speech.SpeechSynthesizer
	)
	if audioConfig, err = audio.NewAudioConfigFromDefaultSpeakerOutput(); err != nil {
		fmt.Println("NewAudioConfigFromDefaultSpeakerOutput Got an error: ", err)
		return
	}
	defer audioConfig.Close()
	if speechConfig, err = speech.NewSpeechConfigFromSubscription(SPEECH_KEY, SPEECH_REGION); err != nil {
		fmt.Println("NewSpeechConfigFromSubscription Got an error: ", err)
		return
	}
	defer speechConfig.Close()

	//speechConfig.SetPropertyByString("OPENSSL_DISABLE_CRL_CHECK", "false")
	//speechConfig.SetProperty(common.SpeechLogFilename, "SpeechLog.txt")
	//16khz ogg
	speechConfig.SetSpeechSynthesisOutputFormat(common.Ogg16Khz16BitMonoOpus)
	if speechSynthesizer, err = speech.NewSpeechSynthesizerFromConfig(speechConfig, audioConfig); err != nil {
		fmt.Println("NewSpeechSynthesizerFromConfig Got an error: ", err)
		return
	}
	defer speechSynthesizer.Close()

	speechSynthesizer.SynthesisStarted(synthesizeStartedHandler)
	speechSynthesizer.Synthesizing(synthesizingHandler)
	speechSynthesizer.SynthesisCompleted(synthesizedHandler)
	speechSynthesizer.SynthesisCanceled(cancelledHandler)

	//dangerous: 如果TTSInfosLeft 被错误得重写成TTSInfos，那么完成转换的语音将会消失
	task := speechSynthesizer.SpeakSsmlAsync(ssml)

	var outcome speech.SpeechSynthesisOutcome
	select {
	case outcome = <-task:
	case <-time.After(60 * time.Second):
		return nil, fmt.Errorf("TTS Timed out")
	}
	defer outcome.Close()
	if outcome.Error != nil {
		fmt.Println("Synthesis error: ", outcome.Error, " ssml: ", ssml)
		return nil, outcome.Error
	} else {
		fmt.Println("Synthesis success! audio length: ", len(outcome.Result.AudioData), "ssml: ", ssml)
	}

	if outcome.Result.Reason == common.SynthesizingAudioCompleted {
		result = &TTSResult{AudioData: outcome.Result.AudioData, AudioDurationSecn: outcome.Result.AudioDuration}
		return result, nil
	} else {
		fmt.Printf("TTS CANCELED: Reason=%d.\n", outcome.Result.Reason)
		return nil, fmt.Errorf("TTS CANCELED: Reason=%d.\n", outcome.Result.Reason)
	}
}
