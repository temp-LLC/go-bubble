package gobubble

import "fmt"

type (
	FetchByIDsRequest struct {
		URL         string
		Token       string
		Target      string
		Constraints []Constraint
		IDs         []string
	}
)

func NewFetchByIDsRequest(
	url, token, target string,
	constants []Constraint,
	ids []string,
) FetchByIDsRequest {
	return FetchByIDsRequest{
		URL:         url,
		Token:       token,
		Target:      target,
		Constraints: constants,
		IDs:         ids,
	}
}

func fetchCount(ids []string) int {
	if len(ids) < FetchLimitMax {
		return len(ids)
	}
	return FetchLimitMax
}

func FetchByIDs[T any](req FetchByIDsRequest) ([]T, error) {
	var collected []T
	ids := req.IDs
	for len(ids) > 0 {
		fetchCount := fetchCount(ids)
		fetched, err := Fetch[T](
			FetchRequest{
				URL:    req.URL,
				Token:  req.Token,
				Target: req.Target,
				Constraints: append(
					req.Constraints,
					Constraint{
						KeyID, In, ids[:fetchCount],
					}),
			})
		if err != nil {
			return nil, fmt.Errorf("fetch: %w", err)
		}

		collected = append(collected, fetched...)
		ids = ids[fetchCount:]
	}
	return collected, nil
}
