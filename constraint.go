package gobubble

type (
	ConstraintType string

	Constraint struct {
		Key            string         `json:"key"`
		ConstraintType ConstraintType `json:"constraint_type"`
		Value          interface{}    `json:"value"`
	}
)

const (
	KeyID                   = "_id"
	Equal    ConstraintType = "equals"
	NotEqual ConstraintType = "not equal"
	In       ConstraintType = "in"
)
