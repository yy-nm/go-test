package net

import "container/list"

type MatchResult int

const (
	_                   MatchResult = iota
	MATCHRESULT_MATCH               // 匹配
	MATCHRESULT_MISS                // 没有符合的匹配项
	MATCHRESULT_UNMATCH             // 不匹配, 有明确的项表示不匹配
)

type Matcher interface {
	Match(int) MatchResult
	RegisterRule(string)
}

type matcherUnit struct {
	t     int
	value int
}

type mather struct {
	rules list.List
}

func (m *mather) RegisterRule(s string) {
	if m == nil {
		return
	}

}

func newDefaultMatcher() Matcher {
	return new(mather)
}
