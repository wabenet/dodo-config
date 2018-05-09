package config

import (
	"reflect"
)

// BuildConfig represents the build configuration for a docker image
type BuildConfig struct {
	Context      string
	Dockerfile   string
	Steps        []string
	Args         KeyValueList
	NoCache      bool
	ForceRebuild bool
}

func DecodeBuild(name string, config interface{}) (BuildConfig, error) {
	var result BuildConfig
	switch t := reflect.ValueOf(config); t.Kind() {
	case reflect.String:
		decoded, err := DecodeString(name, config)
		if err != nil {
			return result, err
		}
		result.Context = decoded
	case reflect.Map:
		for k, v := range t.Interface().(map[interface{}]interface{}) {
			switch key := k.(string); key {
			case "context":
				decoded, err := DecodeString(key, v)
				if err != nil {
					return result, err
				}
				result.Context = decoded
			case "dockerfile":
				decoded, err := DecodeString(key, v)
				if err != nil {
					return result, err
				}
				result.Dockerfile = decoded
			case "steps":
				decoded, err := DecodeStringSlice(key, v)
				if err != nil {
					return result, err
				}
				result.Steps = decoded
			case "args":
				decoded, err := DecodeKeyValueList(key, v)
				if err != nil {
					return result, err
				}
				result.Args = decoded
			case "no_cache":
				decoded, err := DecodeBool(key, v)
				if err != nil {
					return result, err
				}
				result.NoCache = decoded
			case "force_rebuild":
				decoded, err := DecodeBool(key, v)
				if err != nil {
					return result, err
				}
				result.ForceRebuild = decoded
			default:
				return result, errorUnsupportedKey(name, key)
			}
		}
	default:
		return result, errorUnsupportedType(name, t.Kind())
	}
	return result, nil
}