// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package database

type LogicOperator string

const (
	AND LogicOperator = "&&"
	OR  LogicOperator = "||"
)

type ComparsionOperator string

const (
	EQ    ComparsionOperator = "=="
	EQI   ComparsionOperator = "EQI"
	NE    ComparsionOperator = "!="
	LT    ComparsionOperator = "<"
	GT    ComparsionOperator = ">"
	LTE   ComparsionOperator = "<="
	GTE   ComparsionOperator = ">="
	LIKE  ComparsionOperator = "LIKE"
	LIKEI ComparsionOperator = "LIKEI"
)

type SortOrder string

const (
	DESC SortOrder = "DESC"
	ASC  SortOrder = "ASC"
)

type SortAttribute struct {
	Name  string
	Order SortOrder
}

type SortConfig []SortAttribute

type singleValueMatch struct {
	Name     string
	Operator ComparsionOperator
	Compare  interface{}
}

type subFieldValueMatch struct {
	singleValueMatch
	SubField string
}

type AttributeMatcherInfo struct {
	singleValueMatch
}

type ArrayElemMatcherInfo struct {
	singleValueMatch
}

type ArrayElemSubfieldMatcherInfo struct {
	subFieldValueMatch
}

type MapKeyMatcherInfo struct {
	singleValueMatch
}

type MatchGroup struct {
	Chain    []any
	Operator LogicOperator
}

type LimitConfig struct {
	Count  int
	Offset int
}

type QueryConfig struct {
	SortConfig      *SortConfig
	Matcher         *MatchGroup
	Limit           *LimitConfig
	UnsetAttributes []string
	KeepAttributes  []string
}

type RevKeyHolder struct {
	Key string
	Rev string
}

func New() *QueryConfig {
	return &QueryConfig{}
}

func (qc *QueryConfig) SetMatcher(group MatchGroup) *QueryConfig {
	qc.Matcher = &group
	return qc
}

func (qc *QueryConfig) SetUnset(unset []string) *QueryConfig {
	qc.UnsetAttributes = unset
	return qc
}

func (qc *QueryConfig) SetKeep(keep []string) *QueryConfig {
	qc.KeepAttributes = keep
	return qc
}

func (qc *QueryConfig) SetSort(attributes SortConfig) *QueryConfig {
	qc.SortConfig = &attributes
	return qc
}

func (qc *QueryConfig) SetLimit(offset, count int) *QueryConfig {
	qc.Limit = &LimitConfig{
		Count:  count,
		Offset: offset,
	}
	return qc
}

func (qc *QueryConfig) DoesModify() bool {
	return len(qc.KeepAttributes) > 0 || len(qc.UnsetAttributes) > 0
}

func AndChain(groups ...MatchGroup) MatchGroup {
	res := MatchGroup{}
	res.Operator = AND
	for _, g := range groups {
		res.Chain = append(res.Chain, g)
	}
	return res
}

func OrChain(groups ...MatchGroup) MatchGroup {
	res := MatchGroup{}
	for _, g := range groups {
		res.Operator = OR
		res.Chain = append(res.Chain, g)
	}
	return res
}

func AttributeMatcher(name string, op ComparsionOperator, value interface{}) MatchGroup {
	return MatchGroup{
		Chain: []interface{}{
			AttributeMatcherInfo{
				singleValueMatch: singleValueMatch{
					Name:     name,
					Operator: op,
					Compare:  value,
				},
			},
		},
	}
}

func ArrayElemMatcher(name string, op ComparsionOperator, value interface{}) MatchGroup {
	return MatchGroup{
		Chain: []interface{}{
			ArrayElemMatcherInfo{
				singleValueMatch: singleValueMatch{
					Name:     name,
					Operator: op,
					Compare:  value,
				},
			},
		},
	}
}

func ArrayElemSubfieldMatcher(name string, subfield string, op ComparsionOperator, value interface{}) MatchGroup {
	return MatchGroup{
		Chain: []interface{}{
			ArrayElemSubfieldMatcherInfo{
				subFieldValueMatch: subFieldValueMatch{
					singleValueMatch: singleValueMatch{
						Name:     name,
						Operator: op,
						Compare:  value,
					},
					SubField: subfield,
				},
			},
		},
	}
}

func MapKeyMatcher(name string, op ComparsionOperator, value interface{}) MatchGroup {
	return MatchGroup{
		Chain: []interface{}{
			MapKeyMatcherInfo{
				singleValueMatch: singleValueMatch{
					Name:     name,
					Operator: op,
					Compare:  value,
				},
			},
		},
	}
}
