// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package search

type SortableOptions interface {
	ShouldOrder() bool
	GetFilterBy() map[string][]string
	HasColumnFilters() bool
	HasFilter() bool
	GetFilterString() string
	HasPaginationActive() bool
	GetPage() int
	GetItemsPerPage() int
	IsSortAsc() bool
	GetSortKey() string
}

type RequestSearchOptions struct {
	FilterString string              `json:"filterString" validate:"gte=0,lte=80"`
	GroupBy      []string            `json:"groupBy" validate:"dive,gte=0,lte=10"`
	GroupDesc    []string            `json:"groupDesc" validate:"dive,gte=0,lte=10"`
	ItemsPerPage int64               `json:"itemsPerPage"`
	MultiSort    bool                `json:"multiSort"`
	MustSort     bool                `json:"mustSort"`
	Page         int64               `json:"page"`
	SortDesc     []bool              `json:"sortDesc"`
	SortBy       []string            `json:"sortBy" validate:"dive,gte=0,lte=100"`
	FilterBy     map[string][]string `json:"filterBy"`
}

func (options *RequestSearchOptions) GetPage() int {
	return int(options.Page)
}

func (options *RequestSearchOptions) GetItemsPerPage() int {
	return int(options.ItemsPerPage)
}

func (options *RequestSearchOptions) HasFilter() bool {
	return len(options.FilterString) > 0
}

func (options *RequestSearchOptions) HasPaginationActive() bool {
	return options.ItemsPerPage > -1
}

func (options *RequestSearchOptions) GetFilterBy() map[string][]string {
	return options.FilterBy
}

func (options *RequestSearchOptions) GetFilterString() string {
	return options.FilterString
}

func (options *RequestSearchOptionsNew) GetFilterString() string {
	return options.FilterString
}

func (options *RequestSearchOptions) ShouldOrder() bool {
	return len(options.SortBy) > 0
}

func (options *RequestSearchOptions) HasColumnFilters() bool {
	for _, values := range options.FilterBy {
		if len(values) > 0 {
			return true
		}
	}
	return false
}

func (options *RequestSearchOptions) GetSortKey() string {
	return options.SortBy[0]
}

func (options *RequestSearchOptions) IsSortAsc() bool {
	return !options.SortDesc[0]
}

type SortBy struct {
	Key   string `json:"key"`
	Order string `json:"order"`
}

func (s SortBy) IsDesc() bool {
	return s.Order == "desc"
}

func (s SortBy) IsAsc() bool {
	return s.Order != "desc"
}

type RequestSearchOptionsNew struct {
	FilterString string              `json:"filterString" validate:"gte=0,lte=80"`
	GroupBy      []string            `json:"groupBy" validate:"dive,gte=0,lte=10"`
	GroupDesc    []string            `json:"groupDesc" validate:"dive,gte=0,lte=10"`
	ItemsPerPage int64               `json:"itemsPerPage"`
	MultiSort    bool                `json:"multiSort"`
	MustSort     bool                `json:"mustSort"`
	Page         int64               `json:"page"`
	FilterBy     map[string][]string `json:"filterBy"`

	SortBy []SortBy `json:"sortBy"`
}

func (options *RequestSearchOptionsNew) GetItemsPerPage() int {
	return int(options.ItemsPerPage)
}

func (options *RequestSearchOptionsNew) HasFilter() bool {
	return len(options.FilterString) > 0
}

func (options *RequestSearchOptionsNew) GetFilterBy() map[string][]string {
	return options.FilterBy
}

func (options *RequestSearchOptionsNew) HasPaginationActive() bool {
	return options.ItemsPerPage > -1
}

func (options *RequestSearchOptionsNew) ShouldOrder() bool {
	return len(options.SortBy) > 0
}

func (options *RequestSearchOptionsNew) HasColumnFilters() bool {
	for _, values := range options.FilterBy {
		if len(values) > 0 {
			return true
		}
	}
	return false
}

func (options *RequestSearchOptionsNew) GetPage() int {
	return int(options.Page)
}

func (options *RequestSearchOptionsNew) GetSortKey() string {
	return options.SortBy[0].Key
}

func (options *RequestSearchOptionsNew) IsSortAsc() bool {
	return options.SortBy[0].IsAsc()
}
