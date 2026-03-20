// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package mongo

import (
	"regexp"
	"strings"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/database"
)

const regexOperator = "$regex"

var logicOperatorMapping = map[database.LogicOperator]string{
	database.AND: "$and",
	database.OR:  "$or",
}

var compareOperatorMapping = map[database.ComparsionOperator]string{
	database.EQ:    "$eq",
	database.EQI:   "$regex",
	database.NE:    "$ne",
	database.GT:    "$gt",
	database.GTE:   "$gte",
	database.LT:    "$lt",
	database.LTE:   "$lte",
	database.LIKE:  "$regex",
	database.LIKEI: "$regex",
}

func buildQuery(qc *database.QueryConfig) (bson.D, *options.FindOptionsBuilder) {
	filter := bson.D{}
	if qc.Matcher != nil {
		addGroupToLevel(*qc.Matcher, &filter)
	}

	// TODO: findOne if limit == 1 ?
	opts := options.Find()
	if qc.SortConfig != nil {
		opts.SetSort(renderSort(*qc.SortConfig))
	}

	if len(qc.KeepAttributes) > 0 {
		opts.SetProjection(renderKeep(qc.KeepAttributes))
	}

	if len(qc.UnsetAttributes) > 0 {
		opts.SetProjection(renderUnset(qc.UnsetAttributes))
	}

	if qc.Limit != nil {
		opts.SetLimit(int64(qc.Limit.Count))
		opts.SetSkip(int64(qc.Limit.Offset))
	}

	return filter, opts
}

func renderKeep(attributes []string) bson.D {
	var res bson.D
	for _, attr := range attributes {
		res = append(res, bson.E{
			Key:   strings.ToLower(attr),
			Value: 1,
		})
	}
	return res
}

func renderUnset(attributes []string) bson.D {
	var res bson.D
	for _, attr := range attributes {
		res = append(res, bson.E{
			Key:   strings.ToLower(attr),
			Value: 0,
		})
	}
	return res
}

func renderSort(config database.SortConfig) bson.D {
	var res bson.D
	for _, attrConfig := range config {
		orderNum := 1
		if attrConfig.Order == database.ASC {
			orderNum = -1
		}
		res = append(res, bson.E{
			Key:   strings.ToLower(attrConfig.Name),
			Value: orderNum,
		})
	}
	return res
}

func addGroupToLevel(config database.MatchGroup, level any) {
	if len(config.Chain) == 1 {
		if _, ok := config.Chain[0].(database.MatchGroup); !ok {
			addRandomMatcherToLevel(config.Chain[0], level)
			return
		}
	}
	var groupEntries bson.A
	for _, matcherConfig := range config.Chain {
		addGroupToLevel(matcherConfig.(database.MatchGroup), &groupEntries)
	}
	switch level := level.(type) {
	case *bson.D:
		*level = append(*level, bson.E{
			Key:   logicOperatorMapping[config.Operator],
			Value: groupEntries,
		},
		)
	case *bson.A:
		*level = append(*level, bson.D{
			bson.E{
				Key:   logicOperatorMapping[config.Operator],
				Value: groupEntries,
			},
		},
		)
	}
}

func addRandomMatcherToLevel(matcher, level any) {
	switch matcher := matcher.(type) {
	case database.AttributeMatcherInfo:
		addCondToLevel(
			renderAttributeMatcher(matcher),
			level,
		)
	case database.ArrayElemMatcherInfo:
		addCondToLevel(
			renderArrayElemMatcher(matcher),
			level,
		)
	case database.ArrayElemSubfieldMatcherInfo:
		addCondToLevel(
			renderArrayElemSubfieldMatcher(matcher),
			level,
		)
	}
}

func renderArrayElemSubfieldMatcher(info database.ArrayElemSubfieldMatcherInfo) bson.E {
	return bson.E{
		Key: strings.ToLower(info.Name + "." + info.SubField),
		Value: bson.D{
			bson.E{
				Key:   compareOperatorMapping[info.Operator],
				Value: info.Compare,
			},
		},
	}
}

func renderArrayElemMatcher(info database.ArrayElemMatcherInfo) bson.E {
	return bson.E{
		Key: strings.ToLower(info.Name),
		Value: bson.D{
			bson.E{
				Key:   compareOperatorMapping[info.Operator],
				Value: info.Compare,
			},
		},
	}
}

func renderAttributeMatcher(info database.AttributeMatcherInfo) bson.E {
	switch info.Operator {
	case database.EQI:
		return bson.E{
			Key: strings.ToLower(info.Name),
			Value: bson.D{
				bson.E{
					Key:   compareOperatorMapping[info.Operator],
					Value: "^" + regexp.QuoteMeta(info.Compare.(string)) + "$",
				},
				bson.E{
					Key:   "$options",
					Value: "i",
				},
			},
		}
	case database.LIKE:
		prep := strings.ReplaceAll(
			regexp.QuoteMeta(info.Compare.(string)),
			"%",
			".*",
		)
		return bson.E{
			Key: strings.ToLower(info.Name),
			Value: bson.D{
				bson.E{
					Key:   compareOperatorMapping[info.Operator],
					Value: "^" + prep + "$",
				},
			},
		}
	case database.LIKEI:
		prep := strings.ReplaceAll(
			regexp.QuoteMeta(info.Compare.(string)),
			"%",
			".*",
		)
		return bson.E{
			Key: strings.ToLower(info.Name),
			Value: bson.D{
				bson.E{
					Key:   compareOperatorMapping[info.Operator],
					Value: "^" + prep + "$",
				},
				bson.E{
					Key:   "$options",
					Value: "i",
				},
			},
		}
	}

	return bson.E{
		Key: strings.ToLower(info.Name),
		Value: bson.D{
			bson.E{
				Key:   compareOperatorMapping[info.Operator],
				Value: info.Compare,
			},
		},
	}
}

func addCondToLevel(c bson.E, level any) {
	switch l := level.(type) {
	case *bson.D:
		*l = append(*l, c)
	case *bson.A:
		*l = append(*l, bson.D{
			c,
		})
	}
}
