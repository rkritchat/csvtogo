package csvtogo

import (
	"fmt"
	"reflect"
	"strconv"
)

const (
	tagMin = "min"
	tagMax = "max"
)

func validateStruct[T any](f T, row int) error {
	v := reflect.ValueOf(&f).Elem()
	for i := 0; i < v.NumField(); i++ {
		//check minimum value length
		err := checkMin[T](f, i, row, v)
		if err != nil {
			return err
		}

		//check maximum value length
		err = checkMax[T](f, i, row, v)
		if err != nil {
			return err
		}
	}
	return nil
}

func checkMin[T any](field T, sequence, row int, v reflect.Value) error {
	minimum, err := isTagFound[T](field, sequence, tagMin, v)
	if err != nil {
		return err
	}
	if minimum < 0 {
		//no tag found, then skip validate
		return nil
	}

	value := fmt.Sprintf("%v", v.Field(sequence).Interface())
	if len(value) < minimum {
		return fmt.Errorf("value of %v at row %v is invalid, value length must more than or equal %v, but got: %v",
			v.Type().Field(sequence).Name,
			row,
			minimum,
			len(value),
		)
	}
	return nil
}

func checkMax[T any](field T, sequence, row int, v reflect.Value) error {
	maximum, err := isTagFound[T](field, sequence, tagMax, v)
	if err != nil {
		return err
	}
	if maximum < 0 {
		//no tag found, then skip validate
		return nil
	}

	value := fmt.Sprintf("%v", v.Field(sequence).Interface())
	if len(value) > maximum {
		return fmt.Errorf("value of %v at row %v is invalid, value length must less than or equal %v, but got: %v",
			v.Type().Field(sequence).Name,
			row,
			maximum,
			len(value),
		)
	}
	return nil
}

func isTagFound[T any](field T, sequence int, tag string, v reflect.Value) (int, error) {
	tmp := reflect.TypeOf(&field).Elem().Field(sequence).Tag.Get(tag)
	if len(tmp) > 0 {
		val, err := strconv.Atoi(tmp)
		if err != nil {
			return -1, fmt.Errorf("tag %v of field %v must be integer, got: %v", tag, v.Type().Field(sequence).Name, tmp)
		}
		if val < 0 {
			return -1, fmt.Errorf("tag %v of field %v must more than zero, got: %v", tag, v.Type().Field(sequence).Name, tmp)
		}
		return val, nil
	}
	return -1, nil
}
