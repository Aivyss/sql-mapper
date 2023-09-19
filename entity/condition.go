package entity

type Condition struct {
	PartName string
	CaseName string
}

type PredicateConditions func() []*Condition
