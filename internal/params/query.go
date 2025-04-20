package params

import (
	"net/url"
	"strconv"
	"strings"
)

type QueryParams struct {
	Search string
	SortBy string
	Order  string
	Page   int
	Limit  int
}

func ParseQuery(values url.Values) QueryParams {
	q := QueryParams{
		Search: values.Get("search"),
		SortBy: "id",  // default
		Order:  "asc", // default
		Page:   1,
		Limit:  10,
	}

	// Parse sort
	if sort := values.Get("sort"); sort != "" {
		parts := strings.Split(sort, ".")
		if len(parts) == 2 {
			q.SortBy = parts[0]
			q.Order = parts[1]
		}
	}

	// Parse page
	if p, err := strconv.Atoi(values.Get("page")); err == nil && p > 0 {
		q.Page = p
	}

	// Parse limit
	if l, err := strconv.Atoi(values.Get("limit")); err == nil && l > 0 {
		q.Limit = l
	}

	return q
}
