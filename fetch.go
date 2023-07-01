package gobubble

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
)

type (
	Constraint struct {
		Key            string      `json:"key"`
		ConstraintType string      `json:"constraint_type"`
		Value          interface{} `json:"value"`
	}

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
)

func generateRequest(req FetchRequest)(*http.Request, error) {
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
	// q.Set("cursor", strconv.Itoa(cursor))
	u.RawQuery = q.Encode()

	r, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("http.NewRequest: %w", err)
	}
	r.Header.Add("Authorization", "Bearer"+" "+req.Token)
  return r, nil
}

func parseResponse[T any](res *http.Response) ([]T, error) {
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("io.ReadAll: %w", err)
	}
	var p payload
	if err := json.Unmarshal(body, &p); err != nil {
		return nil, fmt.Errorf("json.Unmarshal: %w", err)
	}
  var dest []T
  if err = json.Unmarshal(p.Response.Results, &dest);err != nil {
    return nil, fmt.Errorf("json.Unmarshal: %w", err)
  }
  return dest, nil
}

func Fetch[T any](req FetchRequest) ([]T, error) {
  r, err := generateRequest(req)
  if err != nil {
    return nil, fmt.Errorf("generateRequest: %w", err)
  }

	res, err := http.DefaultClient.Do(r)
	if err != nil {
		return nil, fmt.Errorf("http.Get: %w", err)
	}
	defer func() {
		err := res.Body.Close()
		if err != nil {
			// TODO
			fmt.Println(err)
		}
	}()

  ret, err := parseResponse[T](res)
  if err != nil {
    return nil, fmt.Errorf("parseResponse: %w", err)
  }
	return ret, nil
}
