package translate

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/mjiee/world-news/backend/pkg/errorx"
	"github.com/pkg/errors"
)

// microsoftApiRul is microsoft translate api
const microsoftApiRul = "https://api.cognitive.microsofttranslator.com/translate?api-version=3.0"

// microsoftTranslationResult is the result of microsoft translation
type microsoftTranslationResult []struct {
	Translations []struct {
		Text string `json:"text"`
		To   string `json:"to"`
	} `json:"translations"`
}

// microsoftTranslater is the translator for microsoft
type microsoftTranslater struct {
	apiKey string
	client *http.Client
}

// newMicrosoftTranslater creates a new microsoft translater
func newMicrosoftTranslater(apiKey string) (*microsoftTranslater, error) {
	return &microsoftTranslater{
		apiKey: apiKey,
		client: &http.Client{},
	}, nil
}

// specialLangCode is the special language code for microsoft
var specialLangCode = map[string]string{
	"zh": "zh-Hans",
}

// Translate translates the given texts to the target language
func (m *microsoftTranslater) Translate(ctx context.Context, toLang string, data ...string) ([]string, error) {
	var (
		requestBody = make([]map[string]string, len(data))
	)

	if lang, ok := specialLangCode[toLang]; ok {
		toLang = lang
	}

	for i, text := range data {
		requestBody[i] = map[string]string{"Text": text}
	}

	bodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	url := microsoftApiRul + "&to=" + toLang

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	req.Header.Set("Ocp-Apim-Subscription-Key", m.apiKey)
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	resp, err := m.client.Do(req)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errorx.InternalError.SetMessage(fmt.Sprintf("microsoft translate failed, status code: %d", resp.StatusCode))
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var (
		results            []string
		translationResults microsoftTranslationResult
	)

	if err = json.Unmarshal(respBody, &translationResults); err != nil {
		return nil, errors.WithStack(err)
	}

	for _, result := range translationResults {
		for _, translation := range result.Translations {
			results = append(results, translation.Text)
		}
	}

	return results, nil
}
