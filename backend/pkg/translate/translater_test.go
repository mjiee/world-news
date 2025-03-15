package translate

import (
	"context"
	"fmt"
	"testing"
)

func testTranslate(cfg *Config) {
	translater, err := NewTranslater(cfg)
	if err != nil {
		fmt.Println(err)

		return
	}

	results, err := translater.Translate(context.Background(), "zh", "apple", "banana")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(results)
}

func TestAliyunTranslater(t *testing.T) {
	cfg := &Config{
		Platform:  aliyunPlatform,
		AppId:     "Access_Key_ID",
		AppSecret: "Access_Key_Secret",
	}

	testTranslate(cfg)
}

func TestBaiduTranslater(t *testing.T) {
	cfg := &Config{
		Platform:  baiduPlatform,
		AppId:     "app_id",
		AppSecret: "app_secret",
	}

	testTranslate(cfg)
}

func TestMicrosoft(t *testing.T) {
	cfg := &Config{
		Platform: microsoftPlatform,
		AppId:    "subscription_key",
	}

	testTranslate(cfg)
}

func TestGoogleTranslater(t *testing.T) {
	cfg := &Config{
		Platform: googlePlatform,
		AppId:    "project_id",
	}

	testTranslate(cfg)
}
