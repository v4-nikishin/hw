package hw09structvalidator

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	errStr := "validation failed:"
	for _, e := range v {
		errStr += " " + e.Field + ": " + e.Err.Error()
	}
	return errStr
}

func validate(k, v string, rule string) (bool, ValidationErrors, error) {
	rs := strings.Split(rule, "|")
	ves := make(ValidationErrors, 0, len(rs))
	var ok bool
	var customErr ValidationError
	var err error
	for _, rv := range rs {
		ok, customErr, err = doValidate(k, v, rv)
		if ok {
			continue
		}
		if err != nil {
			return false, ValidationErrors{}, err
		}
		ves = append(ves, customErr)
	}
	if len(ves) > 0 {
		return false, ves, nil
	}
	return ok, ValidationErrors{}, err
}

func doValidate(k, v string, rule string) (bool, ValidationError, error) {
	r := strings.Split(rule, ":")
	if len(r) < 2 {
		return false, ValidationError{}, fmt.Errorf("incorrect validation rule: %s", rule)
	}
	switch r[0] {
	case "len":
		var n int
		fmt.Sscan(r[1], &n)
		if len(v) != n {
			return false, ValidationError{Field: k, Err: fmt.Errorf("incorrect string size")}, nil
		}
	case "in":
		if !strings.Contains(r[1], v) {
			return false, ValidationError{Field: k, Err: fmt.Errorf("value %s is not in %s", v, r[1])}, nil
		}
	case "min":
		min, err := strconv.Atoi(r[1])
		if err != nil {
			return false, ValidationError{}, fmt.Errorf("bad value type")
		}
		val, err := strconv.Atoi(v)
		if err != nil {
			return false, ValidationError{}, fmt.Errorf("bad value type")
		}
		if val < min {
			return false, ValidationError{Field: k, Err: fmt.Errorf("value %d is less than min %d", val, min)}, nil
		}
	case "max":
		max, err := strconv.Atoi(r[1])
		if err != nil {
			return false, ValidationError{}, fmt.Errorf("bad value type")
		}
		val, err := strconv.Atoi(v)
		if err != nil {
			return false, ValidationError{}, fmt.Errorf("bad value type")
		}
		if val > max {
			return false, ValidationError{Field: k, Err: fmt.Errorf("value %d is more than max %d", val, max)}, nil
		}
	case "regexp":
		reg, _ := regexp.Compile(r[1])
		if !reg.MatchString(v) {
			return false, ValidationError{Field: k, Err: fmt.Errorf("value %s does not match format %s", v, r[1])}, nil
		}
	}
	return true, ValidationError{}, nil
}

func validateSlice(fieldName string, sv reflect.Value, rule string) (bool, ValidationErrors, error) {
	var ves ValidationErrors
	for i := 0; i < sv.Len(); i++ {
		var strValue string
		v := sv.Index(i)
		switch v.Kind() { //nolint:exhaustive
		case reflect.String:
			strValue = v.String()
		case reflect.Int:
			strValue = strconv.Itoa(int(v.Int()))
		}
		ok, customErrs, err := validate(fieldName+"["+strconv.Itoa(i)+"]", strValue, rule)
		if ok {
			continue
		}
		if err != nil {
			return false, ValidationErrors{}, err
		}
		ves = append(ves, customErrs...)
	}
	if len(ves) > 0 {
		return false, ves, nil
	}
	return true, ValidationErrors{}, nil
}

func Validate(iv interface{}) error {
	v := reflect.ValueOf(iv)
	if v.Kind() != reflect.Struct {
		return fmt.Errorf("expected a struct, but received %T", iv)
	}
	var ves ValidationErrors
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if !field.IsExported() {
			continue
		}
		tagValue, ok := field.Tag.Lookup("validate")
		if !ok {
			continue
		}
		var strValue string
		fv := v.Field(i)
		switch field.Type.Kind() { //nolint:exhaustive
		case reflect.String:
			strValue = fv.String()
		case reflect.Int:
			strValue = strconv.Itoa(int(fv.Int()))
		case reflect.Slice:
			if fv.Len() == 0 {
				continue
			}
			ok, customErrs, err := validateSlice(field.Name, fv, tagValue)
			if ok {
				continue
			}
			if err != nil {
				return err
			}
			ves = append(ves, customErrs...)

			continue
		}
		ok, customErrs, err := validate(field.Name, strValue, tagValue)
		if ok {
			continue
		}
		if err != nil {
			return err
		}
		ves = append(ves, customErrs...)
	}
	if len(ves) > 0 {
		return fmt.Errorf(ves.Error())
	}
	return nil
}
