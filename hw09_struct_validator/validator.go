package hw09structvalidator

import (
	"bytes"
	"errors"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"sync/atomic"
)

var (
	ErrFailedReadArgument = errors.New("failed to read argument")
	ErrNeedPassStruct     = errors.New("need to pass structure")
)

var (
	ErrStringNotMatchRegExp              = errors.New("string does not match regular expression")
	ErrStringIsNotIncludedInSet          = errors.New("the string is not included in the set")
	ErrNumberIsNotIncludedInSet          = errors.New("the number is not included in the set")
	ErrNumberIsGreaterThanMaximumAllowed = errors.New("the number is greater than the maximum allowed")
	ErrNumberIsLessThanMinimumAllowed    = errors.New("the number is less than the minimum allowed")
	ErrIncorrectLength                   = errors.New("incorrect length")
)

type ValidationError struct {
	Field string
	Err   error
}

type Fn func(interface{}, string) error

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	buff := bytes.NewBufferString("")
	for _, fieldErr := range v {
		buff.WriteString(fieldErr.Err.Error() + "\n")
	}

	return strings.TrimSpace(buff.String())
}

func (v ValidationErrors) Is(target error) bool {
	var vt ValidationErrors
	if errors.As(target, &vt) {
		if len(vt) != len(v) {
			return false
		}

		sort.Slice(vt, func(i, j int) bool {
			return v[i].Field == vt[j].Field
		})

		return reflect.DeepEqual(v, vt)
	}

	return false
}

type Validator struct {
	tag    string
	sep    string
	rules  map[reflect.Kind]map[uint32]Fn
	names  map[string]uint32
	nameID uint32
}

func NewEmpty() *Validator {
	return &Validator{
		tag:   "validate",
		sep:   "|",
		rules: make(map[reflect.Kind]map[uint32]Fn),
		names: make(map[string]uint32),
	}
}

func New() *Validator {
	validator := NewEmpty()

	validator.SetRule(reflect.String, "in", stringIn)
	validator.SetRule(reflect.String, "len", stringLength)
	validator.SetRule(reflect.String, "regexp", stringRegexp)

	validator.SetRule(reflect.Int, "in", intIn)
	validator.SetRule(reflect.Int, "min", intMin)
	validator.SetRule(reflect.Int, "max", intMax)

	return validator
}

func (r *Validator) SetRule(kind reflect.Kind, rule string, executor Fn) {
	id, ok := r.names[rule]
	if !ok {
		id = atomic.AddUint32(&r.nameID, 1)
		r.names[rule] = id
	}

	if _, ok = r.rules[kind]; !ok {
		r.rules[kind] = make(map[uint32]Fn)
	}

	r.rules[kind][id] = executor
}

func (r *Validator) checkRules(rules []string, fieldValue reflect.Value, path string) ValidationErrors {
	var results ValidationErrors
	for _, rule := range rules {
		splits := strings.SplitN(rule, ":", 2)
		if len(splits) < 2 {
			continue
		}

		name, args := splits[0], splits[1]
		nameID, ok := r.names[name]
		if !ok {
			continue
		}

		fn, ok := r.rules[fieldValue.Kind()][nameID]
		if !ok {
			continue
		}

		if err := fn(fieldValue.Interface(), args); err != nil {
			results = append(results, ValidationError{
				Field: path,
				Err:   err,
			})
		}
	}

	return results
}

func (r *Validator) Validate(v interface{}) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	if rv.Kind() != reflect.Struct {
		return ErrNeedPassStruct
	}

	var results ValidationErrors
	for i := 0; i < rv.NumField(); i++ {
		structField := rv.Type().Field(i)
		tag := structField.Tag.Get(r.tag)
		if !structField.IsExported() || tag == "" || tag == "-" {
			continue
		}

		rules := strings.Split(tag, r.sep)
		if len(rules) == 0 {
			continue
		}

		if structField.Type.Kind() == reflect.Slice {
			items := rv.Field(i)
			for k := 0; k < items.Len(); k++ {
				path := structField.Name + "." + strconv.Itoa(k)
				results = append(results, r.checkRules(rules, items.Index(k), path)...)
			}
		} else {
			results = append(results, r.checkRules(rules, rv.Field(i), structField.Name)...)
		}
	}

	if len(results) == 0 {
		return nil
	}

	return results
}
