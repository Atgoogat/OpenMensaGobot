package emojifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Atgoogat/openmensarobot/domain"
)

type EmojifierTextFormater struct {
	url    string
	client *http.Client
}

func NewEmojifierTextFormatter(url string) EmojifierTextFormater {
	return EmojifierTextFormater{
		url:    url,
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

type emojifyReq struct {
	Text string `json:"text"`
}

type emojifyRes struct {
	EmojifiedText string `json:"emojifiedText"`
}

func (etf EmojifierTextFormater) Format(text string) (string, error) {
	requestBody := emojifyReq{Text: text}
	reqJson, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	res, err := etf.client.Post(etf.url+"/emojify", "application/json", bytes.NewBuffer(reqJson))
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("emojifier returned status %d", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	var emojified emojifyRes
	if err := json.Unmarshal(body, &emojified); err != nil {
		return "", err
	}
	if emojified.EmojifiedText == "" {
		return "", fmt.Errorf("emojifier returned empty text")
	}
	return emojified.EmojifiedText, nil
}

var _ domain.TextFormatter = (*EmojifierTextFormater)(nil)
