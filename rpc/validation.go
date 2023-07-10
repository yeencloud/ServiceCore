package rpc

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

func (v *Validator) formatHostname(val string) bool {
	//regexp:
	r, _ := regexp.Compile("^(([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\\-]*[a-zA-Z0-9])\\.)*([A-Za-z0-9]|[A-Za-z0-9][A-Za-z0-9\\-]*[A-Za-z0-9])$")

	if !r.MatchString(val) {
		return false
	}

	return true
}

func (v *Validator) formatIpv4(val string) bool {
	return false
}

func (v *Validator) formatIpv6(val string) bool {
	return false
}

func (v *Validator) format(name string, val any, rules []string, errors *[]string) {
	str := val.(string)

	isValid := false

	for _, rule := range rules {
		switch rule {
		case "hostname":
			isValid = v.formatHostname(str)
		case "ipv4":
			isValid = v.formatIpv4(str)
		case "ipv6":
			isValid = v.formatIpv6(str)
		}
	}

	if !isValid {
		*errors = append(*(errors), fmt.Sprintf("%s format doesn't match any of the following %s", name, strings.Join(rules, ", ")))
		return
	}

	return

}

func (v *Validator) minimumLength(name string, val any, rules []string, errors *[]string) {
	defer func() {
		if r := recover(); r != nil {
			*errors = append(*(errors), fmt.Sprintf("%s: cannot use minimumLength on unsupported type: %s", name, reflect.TypeOf(val).Kind().String()))
			return
		}
	}()

	if len(rules) != 1 {
		*errors = append(*(errors), fmt.Sprintf("%s: minimumLength's rule must have one rule", name))
		return
	}

	minimumLength, err := strconv.Atoi(rules[0])

	if err != nil {
		*errors = append(*(errors), fmt.Sprintf("%s: minimumLength's rule cannot be cast to int", name))
		return
	}

	if minimumLength < 1 {
		*errors = append(*(errors), fmt.Sprintf("%s: minimumLength's rule should be at least 1", name))
		return
	}

	length := reflect.ValueOf(val).Len()

	if length < minimumLength {
		*errors = append(*(errors), fmt.Sprintf("%s should have a minimumLength of at least %d", name, minimumLength))
		return
	}
}

func (v *Validator) minimumValue(name string, val any, rules []string, errors *[]string) {
	defer func() {
		if r := recover(); r != nil {
			*errors = append(*(errors), fmt.Sprintf("%s: cannot use minimumValue on unsupported type: %s", name, reflect.TypeOf(val).Kind().String()))
			return
		}
	}()

	value := reflect.ValueOf(val).Int()

	if len(rules) != 1 {
		*errors = append(*(errors), fmt.Sprintf("%s: minimumValue's rule must have one rule", name))
		return
	}

	minimumValue, err := strconv.Atoi(rules[0])

	if err != nil {
		*errors = append(*(errors), fmt.Sprintf("%s: minimumValue's rule cannot be cast to int", name))
		return
	}

	if int(value) < minimumValue {
		*errors = append(*(errors), fmt.Sprintf("%s should be at least %d", name, minimumValue))
		return
	}
}

func (v *Validator) maximumValue(name string, val any, rules []string, errors *[]string) {
	defer func() {
		if r := recover(); r != nil {
			*errors = append(*(errors), fmt.Sprintf("%s: cannot use maximumValue on unsupported type: %s", name, reflect.TypeOf(val).Kind().String()))
			return
		}
	}()

	value := reflect.ValueOf(val).Int()

	if len(rules) != 1 {
		*errors = append(*(errors), fmt.Sprintf("%s: maximumValue's rule must have one rule", name))
		return
	}

	maximumValue, err := strconv.Atoi(rules[0])

	if err != nil {
		*errors = append(*(errors), fmt.Sprintf("%s: maximumValue's rule cannot be cast to int", name))
		return
	}

	if int(value) > maximumValue {
		*errors = append(*(errors), fmt.Sprintf("%s should be at most %d", name, maximumValue))
		return
	}
}