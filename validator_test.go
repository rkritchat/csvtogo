package csvtogo

import (
	"errors"
	"fmt"
	"testing"
)

func Test_validateStruct(t *testing.T) {
	type Student struct {
		Firstname string `json:"firstname" max:"2" min:"1"`
		Lastname  string `json:"lastname" max:"10" min:"1"`
		Age       int    `json:"age" min:"1"`
	}
	tt := []struct {
		name        string
		t           Student
		expectedErr error
	}{
		{
			name: "should return err when Firstname is more than max",
			t: Student{
				Firstname: "This is my first name",
				Lastname:  "doe",
				Age:       28,
			},
			expectedErr: errors.New("value of Firstname at row 1 is invalid, value length must less than or equal 2, but got: 21"),
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			err := validateStruct[Student](tc.t, 1)
			if fmt.Sprintf("%v", tc.expectedErr) != fmt.Sprintf("%v", err) {
				t.Errorf("must:%v, but got: %v", tc.expectedErr, err)
			}
		})
	}

}
