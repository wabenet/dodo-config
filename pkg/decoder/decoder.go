package decoder

import (
	"reflect"
	"strconv"
)

type (
	Decoder  func(*Status, interface{})
	Producer func() (interface{}, Decoder)
)

func Kinds(lookup map[reflect.Kind]Decoder) Decoder {
	return func(s *Status, config interface{}) {
		kind := reflect.ValueOf(config).Kind()
		if decode, ok := lookup[kind]; ok {
			decode(s, config)
		} else {
			s.Error("invalid type")
		}
	}
}

func Keys(lookup map[string]Decoder) Decoder {
	return func(s *Status, config interface{}) {
		decoded, ok := reflect.ValueOf(config).Interface().(map[interface{}]interface{})
		if !ok {
			s.Error("not a map")
			return
		}
		for k, v := range decoded {
			key := k.(string)
			if decode, ok := lookup[key]; ok {
				s.Run(key, decode, v)
			} else {
				s.Error("unexpected key")
			}
		}
	}
}

func Slice(produce Producer, target interface{}) Decoder {
	return func(s *Status, config interface{}) {
		decoded, ok := reflect.ValueOf(config).Interface().([]interface{})
		if !ok {
			s.Error("not a slice")
			return
		}
		items := reflect.ValueOf(target).Elem()
		for i, item := range decoded {
			ptr, decode := produce()
			s.Run(strconv.Itoa(i), decode, item)
			items.Set(reflect.Append(items, reflect.ValueOf(ptr).Elem()))
		}
	}
}

func Singleton(produce Producer, target interface{}) Decoder {
	return func(s *Status, config interface{}) {
		items := reflect.ValueOf(target).Elem()
		ptr, decode := produce()
		s.Run("", decode, config)
		items.Set(reflect.Append(items, reflect.ValueOf(ptr).Elem()))
	}
}

func Map(produce Producer, target interface{}) Decoder {
	return func(s *Status, config interface{}) {
		decoded, ok := reflect.ValueOf(config).Interface().(map[interface{}]interface{})
		if !ok {
			s.Error("not a map")
			return
		}
		items := reflect.ValueOf(target).Elem()
		for key, value := range decoded {
			ptr, decode := produce()
			s.Run(key.(string), decode, value)
			items.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(ptr).Elem())
		}
	}
}

func String(target interface{}) Decoder {
	return func(s *Status, config interface{}) {
		decoded, ok := config.(string)
		if !ok {
			s.Error("not a string")
			return
		}
		templated, err := ApplyTemplate(s, decoded)
		if err != nil {
			s.Error("invalid templating")
			return
		}
		reflect.ValueOf(target).Elem().SetString(templated)
	}
}

func NewString() Producer {
	return func() (interface{}, Decoder) {
		var target string
		return &target, String(&target)
	}
}

func Bool(target interface{}) Decoder {
	return func(s *Status, config interface{}) {
		decoded, ok := config.(bool)
		if !ok {
			s.Error("not a boolean")
			return
		}
		reflect.ValueOf(target).Elem().SetBool(decoded)
	}
}

func NewBool() Producer {
	return func() (interface{}, Decoder) {
		var target bool
		return &target, Bool(&target)
	}
}

func Int(target interface{}) Decoder {
	return func(s *Status, config interface{}) {
		decoded, ok := config.(int64)
		if !ok {
			s.Error("not an integer")
			return
		}
		reflect.ValueOf(target).Elem().SetInt(decoded)
	}
}

func NewInt() Producer {
	return func() (interface{}, Decoder) {
		var target int64
		return &target, Int(&target)
	}
}
