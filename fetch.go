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

	Request struct {
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

func NewConstraint(key string, constraintType string, value interface{}) *Constraint {
	return &Constraint{
		Key:            key,
		ConstraintType: constraintType,
		Value:          value,
	}
}

// func NewRequest(url string, token string, target string) *Request {
// 	return &Request{
// 		url:    url,
// 		token:  token,
// 		target: target,
// 	}
// }

func Fetch(req Request) error {
	u, err := url.Parse(req.URL)
	if err != nil {
		return fmt.Errorf("invalid url: %w", err)
	}
	u.Path = path.Join(u.Path, "api/1.1/obj", req.Target)

	qcs, err := json.Marshal(req.Constraints)
	if err != nil {
		return fmt.Errorf("json.Marshal: %w", err)
	}

	q := u.Query()
	q.Set("constraints", string(qcs))
	// q.Set("cursor", strconv.Itoa(cursor))
	u.RawQuery = q.Encode()

	r, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return fmt.Errorf("http.NewRequest: %w", err)
	}
	r.Header.Add("Authorization", "Bearer"+" "+req.Token)

  fmt.Printf("req: %+v", r)

	res, err := http.DefaultClient.Do(r)
	if err != nil {
		return fmt.Errorf("http.Get: %w", err)
	}
	defer func() {
		err := res.Body.Close()
		if err != nil {
			// TODO
			fmt.Println(err)
		}
	}()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("io.ReadAll: %w", err)
	}
  var p payload
  if err := json.Unmarshal(body, &p); err != nil {
    return fmt.Errorf("json.Unmarshal: %w", err)
  }
  fmt.Println(p.Response)

	return nil
}

func SampleGenerics[T any](a, b T) []T {
	return []T{a, b}
}
