package sf_config

type Config_Value_Type int

const (
	CONF_VAL_TYPE_NUM  Config_Value_Type = 1
	CONF_VAL_TYPE_STR  Config_Value_Type = 2
	CONF_VAL_TYPE_NIL  Config_Value_Type = 3
	CONF_VAL_TYPE_BOOL Config_Value_Type = 4

	CONF_VAL_TYPE_OBJ Config_Value_Type = 10
	CONF_VAL_TYPE_ARR Config_Value_Type = 11

	CONF_VAL_TYPE_UNKNOWN Config_Value_Type = 20
)

type Config interface {
	Read(data []byte) (err error)

	Get(key string) (c Config, err error)
	Type() (cvt Config_Value_Type, err error)
	//	Arr_type() (cvt Config_Value_Type, err error)

	Bool() (v bool, err error)
	Float() (v float64, err error)
	String() (v string, err error)

	Arr() (v []Config)
}

func get_type(v interface{}) Config_Value_Type {
	switch v.(type) {
	case float64:
		return CONF_VAL_TYPE_NUM
	case bool:
		return CONF_VAL_TYPE_BOOL
	case string:
		return CONF_VAL_TYPE_STR
	case []interface{}:
		return CONF_VAL_TYPE_ARR
	case map[string]interface{}:
		return CONF_VAL_TYPE_OBJ
	case nil:
		return CONF_VAL_TYPE_NIL
	}

	return CONF_VAL_TYPE_UNKNOWN
}
