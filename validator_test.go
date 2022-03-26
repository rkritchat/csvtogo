package csvtogo

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func Test_validateStruct(t *testing.T) {
	type Student struct {
		Firstname string `json:"firstname" max:"10" min:"1"`
		Lastname  string `json:"lastname" max:"10" min:"1"`
		Age       int    `json:"age"`
	}
	tt := []struct {
		name        string
		t           Student
		expectedErr error
	}{
		{
			name: "should return err is nil when struct is valid",
			t: Student{
				Firstname: strings.Repeat("e", 9),
				Lastname:  strings.Repeat("e", 9),
				Age:       1,
			},
			expectedErr: nil,
		},
		{
			name: "should return err when Firstname is more than max",
			t: Student{
				Firstname: strings.Repeat("e", 11),
				Lastname:  strings.Repeat("e", 9),
				Age:       28,
			},
			expectedErr: errors.New("value of Firstname at row 1 is invalid, value length must less than or equal 10, but got: 11"),
		},
		{
			name: "should return err when Lastname is less than min",
			t: Student{
				Firstname: strings.Repeat("e", 9),
				Lastname:  "",
				Age:       28,
			},
			expectedErr: errors.New("value of Lastname at row 1 is invalid, value length must more than or equal 1, but got: 0"),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			e := validateStruct[Student](tc.t, 1)
			if fmt.Sprintf("%v", tc.expectedErr) != fmt.Sprintf("%v", e) {
				t.Errorf("must:%v, but got: %v", tc.expectedErr, e)
			}
		})
	}
}

func Test_checkMin(t *testing.T) {
	type Student struct {
		Firstname string `json:"firstname" max:"10" min:"1"`
		Lastname  string `json:"lastname" max:"10" min:"1"`
		Age       int    `json:"age"`
	}
	tt := []struct {
		name        string
		t           Student
		expectedErr error
	}{
		{
			name: "should return nil when struct is valid",
			t: Student{
				Firstname: strings.Repeat("e", 9),
				Lastname:  strings.Repeat("e", 9),
				Age:       1,
			},
			expectedErr: nil,
		},
		{
			name: "should return nil when struct is valid and some field doesn't have min tag",
			t: Student{
				Firstname: strings.Repeat("e", 9),
				Lastname:  strings.Repeat("e", 9),
				Age:       1,
			},
			expectedErr: nil,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			e := checkMin[Student](tc.t, 1, 1, reflect.ValueOf(&tc.t).Elem())
			if fmt.Sprintf("%v", tc.expectedErr) != fmt.Sprintf("%v", e) {
				t.Errorf("must:%v, but got: %v", tc.expectedErr, e)
			}
		})
	}

}

func Test_checkMin_invalidTag(t *testing.T) {
	type InvalidTag struct {
		Firstname string `json:"firstname" min:"x"`
	}
	tt := []struct {
		name        string
		t           InvalidTag
		expectedErr error
	}{
		{
			name: "should return err when tag value of min is not integer",
			t: InvalidTag{
				Firstname: strings.Repeat("e", 9),
			},
			expectedErr: errors.New("tag min of field Firstname must be integer, got: x"),
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			e := checkMin[InvalidTag](tc.t, 0, 0, reflect.ValueOf(&tc.t).Elem())
			if fmt.Sprintf("%v", tc.expectedErr) != fmt.Sprintf("%v", e) {
				t.Errorf("must:%v, but got: %v", tc.expectedErr, e)
			}
		})
	}
}
