package sf_net

import "container/list"

type filter_result int

const (
	FILTER_MISS  filter_result = 0 + iota // not match rule
	FILTER_MATCH                          // match rule, and filter in
	FILTER_OUT                            // match rule, but filter out
)

type Type_filter interface {
	match(t int) filter_result
}

type type_filter struct {
	rules list.List
}
