// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package couch

import (
	"encoding/json"
	"regexp"
	"strings"

	"mercedes-benz.ghe.com/foss/disuko/infra/repository/database"
)

const (
	RegexOperator     = "$regex"
	ElemMatchOperator = "$elemMatch"
	ObjectKeyOperator = "$keyMapMatch"
)

var compareOperatorMapping = map[database.ComparsionOperator]string{
	database.EQ:  "$eq",
	database.NE:  "$ne",
	database.GT:  "$gt",
	database.GTE: "$gte",
	database.LT:  "$lt",
	database.LTE: "$lte",
}

var logicOperatorMapping = map[database.LogicOperator]string{
	database.AND: "$and",
	database.OR:  "$or",
}

type sortOrder string

const (
	ASC  sortOrder = "asc"
	DESC sortOrder = "desc"
)

type sort map[string]sortOrder

type selector map[string]interface{}
type logicChain []selector

type query struct {
	Selector selector `json:"selector"`
	Sort     []sort   `json:"sort,omitempty"`
	Skip     int      `json:"skip"`
	Limit    int      `json:"limit"`
	Fields   []string `json:"fields,omitempty"`
}

type condition map[string]interface{}

func BuildQuery(qc *database.QueryConfig) string {
	q := query{
		Selector: make(selector),
	}

	if qc.Matcher != nil {
		processMatcher(*qc.Matcher, q.Selector)
	}

	if qc.SortConfig != nil {
		q.addSort(*qc.SortConfig)
	}

	if len(qc.KeepAttributes) > 0 {
		q.Fields = append(q.Fields, qc.KeepAttributes...)
	}

	// TODO: there is no unset in couchDB?

	q.Limit = 100000
	if qc.Limit != nil {
		q.Limit = qc.Limit.Count
		q.Skip = qc.Limit.Offset
	}

	// d, err := json.MarshalIndent(q, "", "    ")
	// if err != nil {
	// 	fmt.Printf("error: %s\n", err)
	// } else {
	// 	fmt.Printf("Query:\n%s\n", d)
	// }
	d, _ := json.Marshal(q)
	return string(d)
}

func (q *query) addSort(i database.SortConfig) {
	for _, s := range i {
		order := ASC
		if s.Order == database.DESC {
			order = DESC
		}
		q.Sort = append(q.Sort, sort{
			s.Name: order,
		})
	}
}

func processMatcher(g database.MatchGroup, level interface{}) {
	// Is single matcher / has no children, dont add ( )
	if len(g.Chain) == 1 {
		if _, ok := g.Chain[0].(database.MatchGroup); !ok {
			processRandomMatcher(g.Chain[0], level)
			return
		}
	}

	n := addLogicChainToLevel(g, level)
	for _, m := range g.Chain {
		if g, ok := m.(database.MatchGroup); ok {
			processMatcher(g, n)
		} else {
			processRandomMatcher(g, level)
		}
	}
}

func processRandomMatcher(i interface{}, level interface{}) {
	switch m := i.(type) {
	case database.AttributeMatcherInfo:
		addCondToLevel(
			m.Name,
			newCondition(m.Operator, m.Compare),
			level,
		)
	case database.ArrayElemMatcherInfo:
		addCondToLevel(
			m.Name,
			newArrayElemSelector(m.Operator, m.Compare),
			level,
		)
	case database.ArrayElemSubfieldMatcherInfo:
		addCondToLevel(
			m.Name,
			newArrayElemSubfieldSelector(m.Operator, m.SubField, m.Compare),
			level,
		)
	case database.MapKeyMatcherInfo:
		addCondToLevel(
			m.Name,
			newMapKeySelector(m.Operator, m.Compare),
			level,
		)
	}
}

func addLogicChainToLevel(g database.MatchGroup, level interface{}) *logicChain {
	n := logicChain{}
	op := logicOperatorMapping[g.Operator]
	switch l := level.(type) {
	case selector:
		l[op] = &n
	case *logicChain:
		*l = append(*l, selector{
			op: &n,
		})
	}
	return &n
}

func addCondToLevel(name string, c condition, level interface{}) {
	switch l := level.(type) {
	case selector:
		l[name] = c
	case *logicChain:
		*l = append(*l, selector{
			name: c,
		})
	}
}

func newMapKeySelector(iOp database.ComparsionOperator, compare interface{}) condition {
	return condition{
		ObjectKeyOperator: newCondition(iOp, compare),
	}
}

func newArrayElemSelector(iOp database.ComparsionOperator, compare interface{}) condition {
	return condition{
		ElemMatchOperator: newCondition(iOp, compare),
	}
}

func newArrayElemSubfieldSelector(iOp database.ComparsionOperator, subfield string, compare interface{}) condition {
	return condition{
		ElemMatchOperator: selector{
			subfield: newCondition(iOp, compare),
		},
	}
}

func newCondition(iOp database.ComparsionOperator, compare interface{}) condition {
	switch iOp {
	case database.EQI:
		return condition{
			string(RegexOperator): "(?i)^" + regexp.QuoteMeta(compare.(string)) + "$",
		}
	case database.LIKE:
		prep := strings.ReplaceAll(
			regexp.QuoteMeta(compare.(string)),
			"%",
			".*",
		)
		return condition{
			string(RegexOperator): "^" + prep + "$",
		}
	case database.LIKEI:
		prep := strings.ReplaceAll(
			regexp.QuoteMeta(compare.(string)),
			"%",
			".*",
		)
		return condition{
			string(RegexOperator): "(?i)" + "^" + prep + "$",
		}
	}

	op := compareOperatorMapping[iOp]
	return condition{
		op: compare,
	}
}
