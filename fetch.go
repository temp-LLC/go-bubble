package gobubble

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
)

const FetchLimitMax = 100

type (
	FetchRequest struct {
		URL         string
		Token       string
		Target      string
		Constraints []Constraint
	}

	payload struct {
		Response struct {
			Results   json.RawMessage `json:"results"`
			Count     int             `json:"count"`
			Remaining int             `json:"remaining"`
		} `json:"response"`
	}

	parsedResponse[T any] struct {
		data      []T
		remaining int
	}
)

func NewFetchRequest(url, token, target string, constants []Constraint) FetchRequest {
	return FetchRequest{
		URL:         url,
		Token:       token,
		Target:      target,
		Constraints: constants,
	}
}

func newHttpRequest(req FetchRequest, cursor int) (*http.Request, error) {
	u, err := url.Parse(req.URL)
	if err != nil {
		return nil, fmt.Errorf("invalid url: %w", err)
	}
	u.Path = path.Join(u.Path, "api/1.1/obj", req.Target)

	qcs, err := json.Marshal(req.Constraints)
	if err != nil {
		return nil, fmt.Errorf("json.Marshal: %w", err)
	}

	q := u.Query()
	q.Set("constraints", string(qcs))
	q.Set("cursor", strconv.Itoa(cursor))
	u.RawQuery = q.Encode()

	r, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("http.NewRequest: %w", err)
	}
	r.Header.Add("Authorization", "Bearer"+" "+req.Token)
	return r, nil
}

func parseResponse[T any](res *http.Response) (*parsedResponse[T], error) {
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("io.ReadAll: %w", err)
	}
	var p payload
	if err := json.Unmarshal(body, &p); err != nil {
		return nil, fmt.Errorf("json.Unmarshal: %w", err)
	}
	var dest []T
	if err = json.Unmarshal(p.Response.Results, &dest); err != nil {
		return nil, fmt.Errorf("json.Unmarshal: %w", err)
	}
	return &parsedResponse[T]{
		data:      dest,
		remaining: p.Response.Remaining,
	}, nil
}

func fetch[T any](fetchRequest FetchRequest, cursor int) (parsedResponse *parsedResponse[T], err error) {
	httpReq, err := newHttpRequest(fetchRequest, cursor)
	if err != nil {
		return nil, fmt.Errorf("generateRequest: %w", err)
	}

	rawResponse, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("http.Get: %w", err)
	}
	if rawResponse.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", rawResponse.StatusCode)
	}
	defer func() {
		if cerr := rawResponse.Body.Close(); cerr != nil {
			err = fmt.Errorf("original: %w, close: %w", err, cerr)
		}
	}()

	parsedResponse, err = parseResponse[T](rawResponse)
	if err != nil {
		return nil, fmt.Errorf("parseResponse: %w", err)
	}

	return
}

func Fetch[T any](req FetchRequest) ([]T, error) {
	var collected []T
	for cursor := 0; ; cursor += FetchLimitMax {
		parsedResponse, err := fetch[T](req, cursor)
		if err != nil {
			return nil, fmt.Errorf("fetch: %w", err)
		}

		collected = append(collected, parsedResponse.data...)
		if parsedResponse.remaining == 0 {
			return collected, nil
		}
	}
}
