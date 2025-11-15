package ttsai

import (
	"bufio"
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/pkg/errors"

	"github.com/mjiee/gokit"

	"github.com/mjiee/world-news/backend/pkg/logx"
)

const (
	doubaoCreateTTsUrl = "https://openspeech.bytedance.com/api/v3/tts/unidirectional"
)

var (
	AudioEmotions = []string{
		"happy", "sad", "angry", "surprised", "fear", "hate",
		"excited", "coldness", "neutral", "depressed", "lovey-dovey", "shy",
		"comfort", "tension", "tender", "storytelling", "radio", "magnetic",
		"advertising", "vocal - fry", "ASMR", "news", "entertainment",
	}
)

// Doubao TTS client
type DoubaoTTSClient struct {
	apiKey string
	model  string
	voices []*Voice
	client *http.Client
}

func NewDoubaoTTSClient(config *Config) (*DoubaoTTSClient, error) {
	if config.Model == "" {
		config.Model = "seed-icl-2.0"
	}

	return &DoubaoTTSClient{
		apiKey: config.ApiKey,
		model:  config.Model,
		voices: config.Voices,
		client: http.DefaultClient,
	}, nil
}

// doubaoAudioParams audio params
type doubaoAudioParams struct {
	Format       string `json:"format"`        // mp3/ogg_opus/pcm
	SampleRate   int    `json:"sample_rate"`   // [8000,16000,22050,24000,32000,44100,48000]
	Emotion      string `json:"emotion"`       // https://www.volcengine.com/docs/6561/1257544
	EmotionScale int    `json:"emotion_scale"` // [1,5]
	SpeechRate   int    `json:"speech_rate"`   // [-50,100]
	LoudnessRate int    `json:"loudness_rate"` // [-50,100]
}

func newDoubaoAudioParams(script *TtsScript) *doubaoAudioParams {
	data := &doubaoAudioParams{Format: "mp3", SampleRate: 16000}

	if script.Emotion != "" {
		data.Emotion = script.Emotion
		data.EmotionScale = 5
	}

	if script.SpeechRate < 1 {
		data.SpeechRate = int((script.SpeechRate - float32(1)) * 50)
	} else if script.SpeechRate > 1 {
		data.SpeechRate = int((script.SpeechRate - float32(1)) * 50)
	}

	if script.Volume < 50 {
		data.LoudnessRate = script.Volume - 50
	} else if script.Volume > 50 {
		data.LoudnessRate = (100 - script.Volume) * 2
	}

	return data
}

// doubaoAudioReqParams req params
type doubaoAudioReqParams struct {
	Text        string             `json:"text"`
	Speaker     string             `json:"speaker"`
	AudioParams *doubaoAudioParams `json:"audio_params"`
}

// TextToSpeech synthesizes text into audio using the Doubao TTS API.
func (c *DoubaoTTSClient) TextToSpeech(ctx context.Context, script *TtsScript) (*TtsTask, error) {
	var (
		reqBody = map[string]interface{}{
			"user": map[string]string{"uid": "world_news"},
			"req_params": &doubaoAudioReqParams{
				Text:        script.Content,
				Speaker:     script.Speaker,
				AudioParams: newDoubaoAudioParams(script),
			},
		}
		header = map[string]string{"X-Api-Resource-Id": c.model}
		voice  = gokit.SliceFind(c.voices, func(v *Voice) bool { return v.Id == script.Speaker })
	)

	if voice != nil && voice.Model != "" {
		header["X-Api-Resource-Id"] = voice.Model
	}

	data, err := c.do(ctx, http.MethodPost, doubaoCreateTTsUrl, header, reqBody)
	if err != nil {
		return nil, err
	}

	audioData, err := c.parseResponse(data)
	if err != nil {
		return nil, err
	}

	return &TtsTask{
		AudioId:   data.Header.Get("X-Tt-Logid"),
		AudioData: audioData,
		Type:      "mp3",
		Status:    ProcessStatusCompleted,
	}, nil
}

// do a request
func (d *DoubaoTTSClient) do(ctx context.Context, method string, url string, header map[string]string, body any,
) (resp *http.Response, err error) {
	var req *http.Request

	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		req, err = http.NewRequest(method, url, bytes.NewBuffer(bodyBytes))
	} else {
		req, err = http.NewRequest(method, url, nil)
	}
	if err != nil {
		return nil, errors.WithStack(err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", d.apiKey)

	for k, v := range header {
		req.Header.Set(k, v)
	}

	return d.client.Do(req)
}

// parseResponse parse tts response
func (d *DoubaoTTSClient) parseResponse(resp *http.Response) ([]byte, error) {
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)

		return nil, fmt.Errorf("http error, status: %d, body: %s", resp.StatusCode, string(body))
	}

	var audioData []byte
	scanner := bufio.NewScanner(resp.Body)

	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}

		var apiResp struct {
			Code    int     `json:"code"`
			Message string  `json:"message"`
			Data    *string `json:"data"`
		}

		if err := json.Unmarshal(line, &apiResp); err != nil {
			continue
		}

		if apiResp.Code != 0 && apiResp.Code != 20000000 {
			return nil, fmt.Errorf("code: %d, message: %s", apiResp.Code, apiResp.Message)
		}

		if apiResp.Data == nil {
			continue
		}

		chunk, err := base64.StdEncoding.DecodeString(*apiResp.Data)
		if err != nil {
			return nil, fmt.Errorf("unmarshal audio data, %v", err.Error())
		}

		audioData = append(audioData, chunk...)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if len(audioData) == 0 {
		logx.Error("parseResponse", errors.New("No audios data found"))
	}

	return audioData, nil
}
