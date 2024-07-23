package query

import "strings"

// Query represents an individual portion of a user's query
type Query interface {
	Match(target string) bool
	Key() string
}

type queryKey struct {
	key string
}

type queryWildcard struct{}

func (c *queryKey) Match(target string) bool {
	return target == c.key
}

func (c *queryKey) Key() string {
	return c.key
}

func (c *queryWildcard) Match(target string) bool {
	return true
}

func (c *queryWildcard) Key() string {
	return "*"
}

func Build(src string, querySeparator rune) ([]Query, error) {
	var queries []Query
	for _, q := range strings.Split(src, string([]rune{querySeparator})) {
		if q == "*" {
			queries = append(queries, &queryWildcard{})
		} else {
			queries = append(queries, &queryKey{
				key: q,
			})
		}
	}
	return queries, nil
}
