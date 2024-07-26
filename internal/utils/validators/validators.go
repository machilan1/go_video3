package validators

import (
	"regexp"
	"slices"
	"strings"
	"unicode/utf8"
)

type Validator struct {
	FieldErrors map[string]string
}

func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0
}

func (v *Validator) CheckField(ok bool, key, message string) {
	if !ok {
		v.addFieldError(key, message)
	}
}

func (v *Validator) addFieldError(key, message string) {
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}

	if _, exist := v.FieldErrors[key]; !exist {
		v.FieldErrors[key] = message
	}
}

func NotBlank(v any) bool {
	switch v := v.(type) {
	case string:
		return strings.TrimSpace(v) != ""
	case int:
		return v != 0
	default:
		return v != nil
	}
}

func MaxChars(v string, l int) bool {
	return utf8.RuneCountInString(v) <= l
}

func MinChars(v string, l int) bool {
	return utf8.RuneCountInString(v) >= l
}

func PermittedValue[T comparable](v T, permittedValues ...T) bool {
	return slices.Contains(permittedValues, v)
}

func IsEmail(v string) bool {
	r, _ := regexp.Compile("^[\\w-\\.]+@([\\w-]+\\.)+[\\w-]{2,4}$")
	return r.Match([]byte(v))
}

func CharAndNumOnly(v string) bool {
	r, _ := regexp.Compile("^[a-zA-Z0-9]+$")
	return r.Match([]byte(v))
}
