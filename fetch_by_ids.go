package gobubble

import "fmt"

type (
	FetchIDsRequest struct {
		URL         string
		Token       string
		Target      string
		Constraints []Constraint
		IDs         string
	}
)

func fetchCount(ids string) int {
	if len(ids) < FetchLimitMax {
		return len(ids)
	}
	return FetchLimitMax
}

func FetchByIDs[T any](req FetchIDsRequest) ([]T, error) {
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
