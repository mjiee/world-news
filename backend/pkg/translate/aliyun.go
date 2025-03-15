package translate

import (
	"context"
	"strings"

	"github.com/mjiee/world-news/backend/pkg/errorx"

	alimt20181012 "github.com/alibabacloud-go/alimt-20181012/v2/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
)

// aliyunTranslator is the translator for aliyun
type aliyunTranslator struct {
	client *alimt20181012.Client
}

// newAliyunTranslater creates a new aliyun translator
func newAliyunTranslater(accessKey, accessSecret string) (*aliyunTranslator, error) {
	config := &openapi.Config{
		AccessKeyId:     tea.String(accessKey),
		AccessKeySecret: tea.String(accessSecret),
		Endpoint:        tea.String("mt.aliyuncs.com"),
	}

	client, err := alimt20181012.NewClient(config)
	if err != nil {
		return nil, err
	}

	return &aliyunTranslator{
		client: client,
	}, nil
}

// Translate translates the given texts to the target language
func (a *aliyunTranslator) Translate(ctx context.Context, toLang string, data ...string) ([]string, error) {
	var (
		results = make([]string, 0, len(data))
		texts   = make([]string, 0)
		flag    = 0
	)

	for idx, text := range data {
		if flag+len(text) < 5000 {
			texts = append(texts, text)
			flag += len(text)

			if idx != len(data)-1 {
				continue
			}
		}

		trans, err := a.translate(ctx, toLang, texts...)
		if err != nil {
			return nil, err
		}

		results = append(results, trans...)
		texts = []string{text}
		flag = len(text)
	}

	return results, nil
}

// translate translates the given texts to the target language
func (t *aliyunTranslator) translate(ctx context.Context, toLang string, data ...string) ([]string, error) {
	var (
		sourceText              = strings.Join(data, "\n")
		translateGeneralRequest = &alimt20181012.TranslateGeneralRequest{
			FormatType:     tea.String("text"),
			SourceLanguage: tea.String("auto"),
			TargetLanguage: tea.String(toLang),
			SourceText:     tea.String(sourceText),
			Scene:          tea.String("general"),
		}
		runtime = &util.RuntimeOptions{}
	)

	resp, err := t.client.TranslateGeneralWithOptions(translateGeneralRequest, runtime)
	if err != nil {
		if sdkErr, ok := err.(*tea.SDKError); ok {
			return nil, errorx.InternalError.SetErr(err).SetMessage(tea.StringValue(sdkErr.Message))
		}

		return nil, err
	}

	return strings.Split(tea.StringValue(resp.Body.Data.Translated), "\n"), nil
}
