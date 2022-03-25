package csvtogo

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

const (
	min   = "min"
	max   = "max"
	empty = ""
)

func CheckMin[T any](f T) error {
	v := reflect.ValueOf(&f).Elem()
	for i := 0; i < v.NumField(); i++ {
		must, err := isTagFound[T](f, i, min)
		if err != nil {
			return err
		}
		if must < 0 {
			//skip validate
			continue
		}

		value := v.Field(i).Interface()
		if len(fmt.Sprintf("%s", value)) < must {
			fmt.Println(value)
			return fmt.Errorf("value of %v is invalid", reflect.ValueOf(&f).Elem().Type().Field(i).Name) //TODO change err to must more than
		}
	}
	return nil
}

func isTagFound[T any](field T, s int, tag string) (int, error) {
	tmp := reflect.TypeOf(&field).Elem().Field(s).Tag.Get(tag)
	if len(tmp) > 0 {
		v, err := strconv.Atoi(tmp)
		if err != nil {
			return -1, errors.New("value of tag is invalid")
		}
		return v, nil
	}
	return -1, nil
}
