package net

import (
	"test/sf/config"
)

type Strategy int

const (
	_ Strategy = iota
	STRATEGY_LOOSE
	STRATEGY_STRICT

	CONF_STRATEGY_LOOSE  string = "loose"
	CONF_STRATEGY_STRICT string = "strict"

	FILTER_TYPE_ALL int = 0 // 0 represent all type
)

type Filter interface {
	// return true: means this msg need be filter or drop
	// return false: means this msg can be recv/send
	Filter(mt MsgType) bool
}

type filter struct {
	rules    map[int]Matcher
	strategy Strategy
}

// match type of msg matcher first,
// if miss and strategy is STRATEGY_LOOSE then match all_type
func (f *filter) Filter(mt MsgType) bool {
	if f == nil {
		return false
	}

	t := mt.GetType()
	if t != FILTER_TYPE_ALL {
		m, ok := f.rules[t]
		if ok {
			result := m.Match(mt.GetProto())
			switch {
			case result == MATCHRESULT_MATCH:
				return true
			case result == MATCHRESULT_UNMATCH:
				return false
			}
		}
	}

	m, ok := f.rules[FILTER_TYPE_ALL]
	if ok {
		result := m.Match(mt.GetProto())
		switch {
		case result == MATCHRESULT_MATCH:
			return true
		case result == MATCHRESULT_UNMATCH:
			return false
		case f.strategy == CONF_STRATEGY_STRICT:
			return true
		case f.strategy == CONF_STRATEGY_LOOSE:
			return false
		}
	}

	return false
}

const (
	CONF_FILTER_TYPE        string = "type"
	CONF_FILTER_STRATEGY    string = "strategy"
	CONF_FILTER_RULES       string = "rules"
	CONF_FILTER_RULES_TYPE  string = "type"
	CONF_FILTER_RULES_VALUE string = "value"

	CONF_FILTER_TYPE_DEFAULT string = "default"
)

func NewFilter(conf config.Config) Filter {
	t, _ := conf.Get(CONF_FILTER_TYPE)

	switch v, _ := t.String(); v {
	case CONF_FILTER_TYPE_DEFAULT:
		fallthrough
	default:
		return newDefaultFilter(conf)
	}
}

// filter config demo:
/*
{
	"type" : "default"
	, "strategy" : "loose"
	, "rules" : [
		{ "type" : 0, "value" : "matcher-value" }
		, { "type" : 1, "value" : "matcher-value" }
	]
}
*/
func newDefaultFilter(conf config.Config) Filter {
	f := new(filter)
	f.rules = make(map[int]Matcher)
	f.strategy = STRATEGY_LOOSE

	c, _ := conf.Get(CONF_FILTER_STRATEGY)
	switch s, _ := c.String(); s {
	case CONF_STRATEGY_STRICT:
		f.strategy = STRATEGY_STRICT
	case CONF_STRATEGY_LOOSE:
		fallthrough
	default:
		f.strategy = STRATEGY_LOOSE
	}

	c, _ = conf.Get(CONF_FILTER_RULES)
	switch t, _ := c.Type(); t {
	case config.CONF_VAL_TYPE_ARR:
		for _, e := range c.Arr() {
			readRule(f, e)
		}
	case config.CONF_VAL_TYPE_OBJ:
		readRule(f, t)
	}

	return f
}

func readRule(f *filter, conf config.Config) {
	if f == nil || conf == nil {
		return
	}

	c, _ := conf.Get(CONF_FILTER_RULES_TYPE)
	t, e := c.Float()
	if e != nil {
		panic("the type of rule of filter must be int")
	}

	m, ok := f.rules[int(t)]
	if !ok {
		m = newDefaultMatcher()
	}
	c, _ = conf.Get(CONF_FILTER_RULES_VALUE)
	v, e := c.String()
	if e != nil {
		return
	}
	m.RegisterRule(v)
}
