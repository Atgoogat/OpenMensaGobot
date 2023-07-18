package emojifier

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/Atgoogat/openmensarobot/domain"
)

type EmojifierTextFormater struct {
	url string
}

func NewEmojifierTextFormatter(url string) EmojifierTextFormater {
	return EmojifierTextFormater{
		url: url,
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

	res, err := http.Post(etf.url+"/emojify", "application/json", bytes.NewBuffer(reqJson))
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	var emojified emojifyRes
	err = json.Unmarshal(body, &emojified)
	return emojified.EmojifiedText, err
}

var _ domain.TextFormatter = (*EmojifierTextFormater)(nil)
