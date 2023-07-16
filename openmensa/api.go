package openmensa

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

const (
	OPENMENSA_API_ENDPOINT = "https://openmensa.org/api/v2"
)

var (
	ErrExpectedStatusCodeOk = errors.New("expected status code ok (200)")
)

type OpenmensaApi struct {
	endpoint string
	client   *http.Client
}

func NewOpenmensaApi(endpoint string) OpenmensaApi {
	return OpenmensaApi{
		endpoint: endpoint,
		client:   &http.Client{},
	}
}

func makeJsonRequest[T interface{}](client *http.Client, req *http.Request, data *T) error {
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return ErrExpectedStatusCodeOk
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, data)
}

func (api OpenmensaApi) ListCanteens(page, limit int) ([]Canteen, error) {
	req, _ := http.NewRequest(http.MethodGet, api.endpoint+"/canteens", nil)
	req.Header.Add("Accept", "application/json")
	q := req.URL.Query()
	q.Add("page", strconv.Itoa(page))
	q.Add("limit", strconv.Itoa(limit))
	req.URL.RawQuery = q.Encode()

	var canteens []Canteen
	err := makeJsonRequest(api.client, req, &canteens)
	return canteens, err
}

func (api OpenmensaApi) ListDaysOfCanteen(id, page, limit int) ([]CanteenDay, error) {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/canteens/%d/days", api.endpoint, id), nil)
	req.Header.Add("Accept", "application/json")
	q := req.URL.Query()
	q.Add("page", strconv.Itoa(page))
	q.Add("limit", strconv.Itoa(limit))
	req.URL.RawQuery = q.Encode()

	var canteenDays []CanteenDay
	err := makeJsonRequest(api.client, req, &canteenDays)
	return canteenDays, err
}

func (api OpenmensaApi) ListMealsForADay(id int, date time.Time, page, limit int) ([]CanteenMeal, error) {
	fmtDate := date.Format("2006-01-02")
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/canteens/%d/days/%s/meals", api.endpoint, id, fmtDate), nil)
	req.Header.Add("Accept", "application/json")
	q := req.URL.Query()
	q.Add("page", strconv.Itoa(page))
	q.Add("limit", strconv.Itoa(limit))
	req.URL.RawQuery = q.Encode()

	var meals []CanteenMeal
	err := makeJsonRequest(api.client, req, &meals)
	return meals, err
}
