package csvtogo

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"reflect"
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

func Test_typeSafe(t *testing.T) {
	//string
	f := reflect.ValueOf(&struct {
		Firstname string
	}{}).Elem().Field(0)
	t.Run("should return nil when value is valid", func(t *testing.T) {
		e := typeSafe(f, "test", 0)
		if e != nil {
			t.Errorf("must:nil, but got: %v", e)
		}
	})

	//valid int
	t.Run("should return nil when value is valid", func(t *testing.T) {
		e := typeSafe(
			reflect.ValueOf(&struct {
				Age int
			}{}).Elem().Field(0),
			"1",
			0,
		)
		if e != nil {
			t.Errorf("must:nil, but got: %v", e)
		}
	})
	//invalid int
	t.Run("should return nil when value is valid", func(t *testing.T) {
		e := typeSafe(
			reflect.ValueOf(&struct {
				Age int
			}{}).Elem().Field(0),
			"x",
			0,
		)
		if e == nil {
			t.Errorf("must: err, but got: nil")
		}
	})

	//valid float
	t.Run("should return nil when value is valid", func(t *testing.T) {
		e := typeSafe(
			reflect.ValueOf(&struct {
				Salary float64
			}{}).Elem().Field(0),
			"1",
			0,
		)
		if e != nil {
			t.Errorf("must:nil, but got: %v", e)
		}
	})
	//invalid float
	t.Run("should return nil when value is valid", func(t *testing.T) {
		e := typeSafe(
			reflect.ValueOf(&struct {
				Salary float64
			}{}).Elem().Field(0),
			"x",
			0,
		)
		if e == nil {
			t.Errorf("must: err, but got: nil")
		}
	})

	//valid bool
	t.Run("should return nil when value is valid", func(t *testing.T) {
		e := typeSafe(
			reflect.ValueOf(&struct {
				Married bool
			}{}).Elem().Field(0),
			"true",
			0,
		)
		if e != nil {
			t.Errorf("must:nil, but got: %v", e)
		}
	})
	//invalid bool
	t.Run("should return nil when value is valid", func(t *testing.T) {
		e := typeSafe(
			reflect.ValueOf(&struct {
				Married bool
			}{}).Elem().Field(0),
			"yes",
			0,
		)
		if e == nil {
			t.Errorf("must: err, but got: nil")
		}
	})

	//un support type
	t.Run("should return nil when value is valid", func(t *testing.T) {
		e := typeSafe(
			reflect.ValueOf(&struct {
				Married []string
			}{}).Elem().Field(0),
			"something",
			0,
		)
		if e == nil {
			t.Errorf("must: err, but got: nil")
		}
	})
}

func Test_CsvToStruct(t *testing.T) {
	type Customer struct {
		Name    string
		Age     int
		Salary  float64
		Married bool
	}
	tt := []struct {
		name      string
		filename  string
		genFile   func()
		expectedR []Customer
		expectedE error
	}{
		{
			name:     "should return valid result",
			filename: "./for_test.csv",
			genFile: func() {
				f, err := os.Create("./for_test.csv")
				if err != nil {
					panic(err)
				}
				w := csv.NewWriter(f)
				_ = w.Write([]string{"NAME", "AGE", "SALARY", "MARRIED"})
				_ = w.Write([]string{"Sarah", "12", "200.00", "false"})
				_ = w.Write([]string{"John", "21", "10.59", "true"})
				w.Flush()
				_ = f.Close()
			},
			expectedR: []Customer{
				{
					Name:    "Sarah",
					Age:     12,
					Salary:  200,
					Married: false,
				},
				{
					Name:    "John",
					Age:     21,
					Salary:  10.59,
					Married: true,
				},
			},
			expectedE: nil,
		},
		{
			name:     "should return error when some value is not match with field type",
			filename: "./for_test.csv",
			genFile: func() {
				f, err := os.Create("./for_test.csv")
				if err != nil {
					panic(err)
				}
				w := csv.NewWriter(f)
				_ = w.Write([]string{"NAME", "AGE", "SALARY", "MARRIED"})
				_ = w.Write([]string{"Sarah", "12", "200.00", "false"})
				_ = w.Write([]string{"John", "21", "NOT FOUND", "true"})
				w.Flush()
				_ = f.Close()
			},
			expectedR: nil,
			expectedE: errors.New("invalid csv value at row: 2, the struct accept type float"),
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if tc.genFile != nil {
				tc.genFile()
			}
			c, _ := NewClient[Customer](tc.filename)
			r, e := c.CsvToStruct()
			if fmt.Sprintf("%v", tc.expectedE) != fmt.Sprintf("%v", e) {
				t.Errorf("must:%v, but got: %v", tc.expectedE, e)
			}
			//deep check
			err := deepEqual[Customer](tc.expectedR, r)
			if err != nil {
				t.Error(err)
			}
			_ = os.Remove(tc.filename)
		})
	}
}

func deepEqual[T any](expect []T, actual []*T) error {
	if len(expect) != len(actual) {
		return fmt.Errorf("array size is not match")
	}
	for k, v := range expect {
		if fmt.Sprintf("%v", v) != fmt.Sprintf("%v", *actual[k]) {
			return fmt.Errorf("must: %v, but got: %v", v, *actual[k])
		}
	}
	return nil
}
