package translate

// Imports the Google Cloud Translation library
import (
	"context"
	"fmt"

	goTranslate "cloud.google.com/go/translate/apiv3"
	"cloud.google.com/go/translate/apiv3/translatepb"
	"github.com/pkg/errors"
)

// googleTranslater is the translator for google
type googleTranslater struct {
	projectID string
}

// newGoogleTranslater creates a new google translater
func newGoogleTranslater(projectID string) (*googleTranslater, error) {
	return &googleTranslater{
		projectID: projectID,
	}, nil
}

// Translate translates the given texts to the target language
func (g *googleTranslater) Translate(ctx context.Context, toLang string, data ...string) ([]string, error) {
	if len(data) == 0 {
		return data, nil
	}

	client, err := goTranslate.NewTranslationClient(ctx)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer client.Close()

	req := &translatepb.TranslateTextRequest{
		Parent:             fmt.Sprintf("projects/%s/locations/global", g.projectID),
		TargetLanguageCode: toLang,
		MimeType:           "text/plain",
		Contents:           data,
	}

	resp, err := client.TranslateText(ctx, req)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	results := make([]string, 0, len(resp.GetTranslations()))

	for idx, translation := range resp.GetTranslations() {
		results[idx] = translation.GetTranslatedText()
	}

	return results, nil
}
