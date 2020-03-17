package decoder

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

func IncludeText(decode Decoder) Decoder {
	return func(s *Status, config interface{}) {
		var decoded string
		String(&decoded)(s, config)
                s.LoadYaml(decode, []byte(decoded))
	}
}

func IncludeFile(decode Decoder) Decoder {
	return func(s *Status, config interface{}) {
		var decoded string
		String(&decoded)(s, config)

		bytes, err := readFile(decoded)
		if err != nil {
			s.Error("could not read file")
			return
		}

                sub := New(decoded)
                sub.LoadYaml(decode, bytes)
                s.errors = append(s.errors, sub.errors...)
	}
}

func readFile(filename string) ([]byte, error) {
	if !filepath.IsAbs(filename) {
		directory, err := os.Getwd()
		if err != nil {
			return []byte{}, err
		}
		filename, err = filepath.Abs(filepath.Join(directory, filename))
		if err != nil {
			return []byte{}, err
		}
	}
	return ioutil.ReadFile(filename)
}
