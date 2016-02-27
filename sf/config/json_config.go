package config

import (
	"encoding/json"
	"test/sf/misc"
)

type json_Config struct {
	content interface{}
	is_read bool
}

func (c *json_Config) Read(data []byte) (err error) {
	if c == nil {
		err = misc.ErrNilPointer
		return
	}

	c.is_read = false

	err = json.Unmarshal(data, &c.content)
	if err == nil {
		c.is_read = true
	}
	return
}

func (c *json_Config) Type() (cvt Config_Value_Type, err error) {
	if c == nil {
		err = misc.ErrNilPointer
		return
	} else if !c.is_read {
		err = misc.ErrConfigNotInit
		return
	}

	cvt = get_type(c.content)
	return
}

func (c *json_Config) Arr_type() (cvt Config_Value_Type, err error) {
	if c == nil {
		err = misc.ErrNilPointer
		return
	} else if !c.is_read {
		err = misc.ErrConfigNotInit
		return
	}

	vt := get_type(c.content)
	if vt != CONF_VAL_TYPE_ARR {
		err = misc.ErrConfigNotArr
		return
	}

	av, ok := c.content.([]interface{})
	if !ok {
		err = misc.ErrConfigConvert
		return
	}

	if cap(av) > 0 {
		cvt = get_type(av[0])
		return
	} else {
		cvt = CONF_VAL_TYPE_OBJ
		return
	}
}

func (c *json_Config) Get(key string) (config Config, err error) {
	if c == nil {
		err = misc.ErrNilPointer
		return
	} else if !c.is_read {
		err = misc.ErrConfigNotInit
		return
	}

	vt := get_type(c.content)
	if vt != CONF_VAL_TYPE_OBJ {
		err = misc.ErrConfigTypeNotMatch
		return
	}

	v, ok := c.content.(map[string]interface{})
	if !ok {
		err = misc.ErrConfigConvert
		return
	}

	config = &json_Config{content: v[key], is_read: true}
	return
}

func (c *json_Config) Bool() (v bool, err error) {
	if c == nil {
		err = misc.ErrNilPointer
		return
	} else if !c.is_read {
		err = misc.ErrConfigNotInit
		return
	}

	vt := get_type(c.content)
	if vt != CONF_VAL_TYPE_BOOL {
		err = misc.ErrConfigTypeNotMatch
		return
	}

	v = c.content.(bool)
	return
}

func (c *json_Config) Float() (v float64, err error) {
	if c == nil {
		err = misc.ErrNilPointer
		return
	} else if !c.is_read {
		err = misc.ErrConfigNotInit
		return
	}

	vt := get_type(c.content)
	if vt != CONF_VAL_TYPE_NUM {
		err = misc.ErrConfigTypeNotMatch
		return
	}

	v = c.content.(float64)
	return
}

func (c *json_Config) String() (v string, err error) {
	if c == nil {
		err = misc.ErrNilPointer
		return
	} else if !c.is_read {
		err = misc.ErrConfigNotInit
		return
	}

	vt := get_type(c.content)
	if vt != CONF_VAL_TYPE_STR {
		err = misc.ErrConfigTypeNotMatch
		return
	}

	v = c.content.(string)
	return
}

func (c *json_Config) Arr() (v []Config) {
	if c == nil {
		return
	} else if !c.is_read {
		return
	}

	vt := get_type(c.content)
	if vt != CONF_VAL_TYPE_ARR {
		return
	}

	av := c.content.([]interface{})
	v = make([]Config, len(av))
	for i, vav := range av {
		v[i] = &json_Config{content: vav, is_read: true}
	}

	return
}

func New_json_config() (c Config) {
	j := new(json_Config)
	j.is_read = false
	j.content = nil

	c = j
	return
}
