package api

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

type validatable interface {
	Valid() bool
}

func setScalar(f reflect.Value, name, raw string) error {
	switch f.Kind() {
	case reflect.Bool:
		b, err := strconv.ParseBool(raw)
		if err != nil {
			return fmt.Errorf("%s must be a boolean (got %q)", name, raw)
		}
		f.SetBool(b)
	case reflect.Int, reflect.Int64:
		n, err := strconv.Atoi(raw)
		if err != nil {
			return fmt.Errorf("%s must be an integer (got %q)", name, raw)
		}
		f.SetInt(int64(n))
	case reflect.String:
		f.SetString(raw)
	default:
		return fmt.Errorf("unsupported query field type %s for %s", f.Kind(), name)
	}
	if chk, ok := f.Interface().(validatable); ok && !chk.Valid() {
		return fmt.Errorf("invalid value for %s: %q", name, raw)
	}
	return nil
}

func BindQuery[T any](r *http.Request) (T, error) {
	var out T
	vals := r.URL.Query()
	v := reflect.ValueOf(&out).Elem()
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		name := field.Tag.Get("query")
		if name == "" {
			continue
		}
		f := v.Field(i)

		if f.Kind() == reflect.Slice {
			raws, ok := vals[name]
			if !ok || len(raws) == 0 {
				continue
			}
			var parts []string
			for _, raw := range raws {
				parts = append(parts, strings.Split(raw, ",")...)
			}
			slice := reflect.MakeSlice(f.Type(), 0, len(parts))
			for _, raw := range parts {
				if raw == "" {
					continue
				}
				elem := reflect.New(f.Type().Elem()).Elem()
				if err := setScalar(elem, name, raw); err != nil {
					return out, err
				}
				slice = reflect.Append(slice, elem)
			}
			f.Set(slice)
			continue
		}

		raw := vals.Get(name)
		if raw == "" {
			continue
		}
		if err := setScalar(f, name, raw); err != nil {
			return out, err
		}
	}
	return out, nil
}
