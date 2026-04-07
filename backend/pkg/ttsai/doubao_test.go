package ttsai

import (
	"context"
	"testing"

	"github.com/mjiee/world-news/backend/pkg/audio"
)

func TestDoubaoTextToSpeech(t *testing.T) {
	client, err := NewDoubaoTTSClient(&Config{
		ApiKey: "xxx",
		Model:  "seed-icl-2.0",
		Voices: []*Voice{
			{
				Id:    "zh_female_wanwanxiaohe_moon_bigtts",
				Model: "seed-tts-1.0",
			},
			{
				Id:    "zh_male_xudong_conversation_wvae_bigtts",
				Model: "seed-tts-1.0",
			},
		},
	})
	if err != nil {
		t.Error(err)
		return
	}

	ctx := context.Background()

	dataA, err := client.TextToSpeech(ctx, &TtsScript{
		Content:    "It helps you understand large codebases, automate tedious work, and ship faster.",
		Speaker:    "zh_male_xudong_conversation_wvae_bigtts",
		Emotion:    "radio",
		SpeechRate: 1.1,
		Volume:     100,
		Silence:    0.2,
	})
	if err != nil {
		t.Error(err)
		return
	}

	_, err = audio.WriteMp3sToFile("./axu_reference.mp3", dataA.AudioData)
	if err != nil {
		t.Error(err)
		return
	}

	dataB, err := client.TextToSpeech(ctx, &TtsScript{
		Content:    "It helps you understand large codebases, automate tedious work, and ship faster.",
		Speaker:    "zh_female_wanwanxiaohe_moon_bigtts",
		Emotion:    "neutral",
		SpeechRate: 1.1,
		Volume:     100,
		Silence:    0.2,
	})
	if err != nil {
		t.Error(err)
		return
	}

	_, err = audio.WriteMp3sToFile("./yuyuan_reference.mp3", dataB.AudioData)
	if err != nil {
		t.Error(err)
		return
	}
}
