package translate

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/mjiee/world-news/backend/pkg/errorx"
	"github.com/pkg/errors"
)

// baidu translate api
const baiduApiURL = "https://fanyi-api.baidu.com/api/trans/vip/translate"

// baiduTranslationResult is the result of baidu translation
type baiduTranslationResult struct {
	From        string `json:"from"`
	To          string `json:"to"`
	TransResult []struct {
		Src string `json:"src"`
		Dst string `json:"dst"`
	} `json:"trans_result"`

	ErrorCode string `json:"error_code"`
	ErrorMsg  string `json:"error_msg"`
}

// baiduTranslater is the translator for baidu
type baiduTranslater struct {
	appID     string
	secretKey string
	client    *http.Client
}

// newBaiduTranslater creates a new baidu translater
func newBaiduTranslater(appID, secretKey string) (*baiduTranslater, error) {
	return &baiduTranslater{
		appID:     appID,
		secretKey: secretKey,
		client:    http.DefaultClient,
	}, nil
}

// generateSign generates the sign for baidu translation, using md5(appID + q + salt + secretKey)
func (b *baiduTranslater) generateSign(q string, salt string) string {
	str := b.appID + q + salt + b.secretKey
	hash := md5.Sum([]byte(str))

	return hex.EncodeToString(hash[:])
}

// Translate translates the given texts to the target language
func (b *baiduTranslater) Translate(ctx context.Context, toLang string, data ...string) ([]string, error) {
	var (
		results = make([]string, 0, len(data))
		texts   = make([]string, 0)
		flag    = 0
	)

	for idx, text := range data {
		if flag+len(text) < 6000 {
			texts = append(texts, text)
			flag += len(text)

			if idx != len(data)-1 {
				continue
			}
		}

		trans, err := b.translate(ctx, toLang, texts...)
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
func (b *baiduTranslater) translate(ctx context.Context, toLang string, data ...string) ([]string, error) {
	var (
		q    = strings.Join(data, "\n")
		salt = strconv.FormatInt(time.Now().Unix(), 10)
		sign = b.generateSign(q, salt)
	)

	params := url.Values{}
	params.Set("q", q)
	params.Set("from", "auto")
	params.Set("to", toLang)
	params.Set("appid", b.appID)
	params.Set("salt", salt)
	params.Set("sign", sign)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baiduApiURL+"?"+params.Encode(), nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	resp, err := b.client.Do(req)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var transData baiduTranslationResult

	if err = json.Unmarshal(body, &transData); err != nil {
		return nil, errors.WithStack(err)
	}

	if transData.ErrorCode != "" && transData.ErrorCode != "52000" {
		return nil, errorx.InternalError.SetMessage(fmt.Sprintf("translate failed: %s", transData.ErrorMsg))
	}

	results := make([]string, 0, len(data))

	for _, item := range transData.TransResult {
		results = append(results, strings.Split(item.Dst, "\n")...)
	}

	return results, nil
}
