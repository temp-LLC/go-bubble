package gobubble

import "fmt"

const KeyID = "_id"

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
		ret, err := Fetch[T](
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

		collected = append(collected, ret...)
		ids = ids[fetchCount:]
	}
	return collected, nil
}

/*
	var allDateApplications []entity.DateApplication
	for len(ids) > 0 {
		lenIDs := len(ids)
		var n int
		if lenIDs < bubble.MaxLimit {
			n = lenIDs
		} else {
			n = bubble.MaxLimit
		}

		c := append(constraints, bubble.GenerateConstraint(bubble.KeyID, bubble.ConstraintTypeIn, ids[:n]))
		dateApplications, err := d.Fetch(c)
		if err != nil {
			return nil, err
		}

		allDateApplications = append(allDateApplications, dateApplications...)
		ids = ids[n:]
	}

	return allDateApplications, nil

*/
