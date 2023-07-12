package gobubble

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

func FetchByIDs[T any](ids string) ([]T, error) {
	var collected []T
	for len(ids) > 0 {
		fetchCount := fetchCount(ids)
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
