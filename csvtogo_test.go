package csvtogo

import (
	"errors"
	"fmt"
	"testing"
)

func Test_realNoOfCol(t *testing.T) {
	tt := []struct {
		name      string
		noOfCol   int
		skip      int
		expectedR int
	}{
		{
			name:      "should return valid result when no of column is less than skip",
			noOfCol:   5,
			skip:      6,
			expectedR: 5,
		},
		{
			name:      "should return valid result when no of column is more than skip",
			noOfCol:   10,
			skip:      6,
			expectedR: 4,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			r := realNoOfCol(tc.noOfCol, tc.skip)
			if tc.expectedR != r {
				t.Errorf("must:%v, but got: %v", tc.expectedR, r)
			}
		})
	}
}

func Test_valueSetter(t *testing.T) {
	skipper := make(map[int]int)
	skipper[1] = 1
	type Student struct {
		Age int `min:"2" max:"5"`
	}
	tt := []struct {
		name      string
		ops       Options
		ref       Student
		data      []string
		row       int
		expectedE error
	}{
		{
			name: "should return nil when value is valid to set",
			ops: Options{
				SkipHeader: false,
				skipper:    skipper,
			},
			ref:       Student{},
			data:      []string{"29", "ignore me"},
			row:       0,
			expectedE: nil,
		},
		{
			name: "should do nothing at row 0 when skip header is true",
			ops: Options{
				SkipHeader: true,
				skipper:    skipper,
			},
			ref:       Student{},
			data:      []string{"12", "this is first row"},
			row:       0,
			expectedE: nil,
		},
		{
			name: "should return err when value is not match with struct",
			ops: Options{
				SkipHeader: false,
				skipper:    skipper,
			},
			ref:       Student{},
			data:      []string{"e"},
			row:       0,
			expectedE: errors.New("invalid csv value at row: 0, the struct accept type int"),
		},
		{
			name: "should return err when value length is more than max",
			ops: Options{
				SkipHeader: false,
				skipper:    skipper,
			},
			ref:       Student{},
			data:      []string{"123456"},
			row:       0,
			expectedE: errors.New("value of Age at row 0 is invalid, value length must less than or equal 5, but got: 6"),
		},
		{
			name: "should return err when value length is less than min",
			ops: Options{
				SkipHeader: false,
				skipper:    skipper,
			},
			ref:       Student{},
			data:      []string{"1"},
			row:       0,
			expectedE: errors.New("value of Age at row 0 is invalid, value length must more than or equal 2, but got: 1"),
		},
		{
			name: "should return err when number of array string not match with struct",
			ops: Options{
				SkipHeader: false,
			},
			ref:       Student{},
			data:      []string{"12", "firstname"},
			row:       0,
			expectedE: errors.New("number of column is not match with struct at row: 0, expected: 1, got: 2"),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			c := Executor[Student]{
				ops:      tc.ops,
				outsChan: make(chan []Student, 1),
				outChan:  make(chan Student, 1),
				endChan:  make(chan bool, 1),
				nextChan: make(chan bool, 1),
				errChan:  make(chan error),
				run:      true,
			}
			//channel watching
			go func() {
				for c.run {
					c.Read()
					c.Next()
				}
			}()
			e := c.valueSetter(tc.ref, tc.data, 0)
			if fmt.Sprintf("%v", tc.expectedE) != fmt.Sprintf("%v", e) {
				t.Errorf("must:%v, but got: %v", tc.expectedE, e)
			}
			c.run = false
			c.Close()
		})
	}
}

//func Test_typeSafe(t *testing.T) {
//	tt := []struct {
//		name string
//	}{
//		{
//			name: "should return nil when value is valid",
//		},
//	}
//	for _, tc := range tt {
//		t.Run(tc.name, func(t *testing.T) {
//
//		})
//	}
//}
