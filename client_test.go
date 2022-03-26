package csvtogo

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"
)

func Test_NewClient(t *testing.T) {
	type Student struct {
		Firstname string
	}
	tt := []struct {
		name      string
		ops       []*Options
		file      string
		genFile   func()
		expectedR bool
		expectedE error
	}{
		{
			name: "should return valid result when options is nil",
			ops:  nil,
			file: "./new_client_test.csv",
			genFile: func() {
				f, err := os.Create("./new_client_test.csv")
				if err != nil {
					panic(err)
				}
				w := csv.NewWriter(f)
				_ = w.Write([]string{"Firstname"})
				_ = w.Write([]string{strings.Repeat("e", 5)})
				w.Flush()
				_ = f.Close()
			},
			expectedR: true,
			expectedE: nil,
		},
		{
			name: "should return valid result when options is not nil",
			ops: []*Options{
				{
					SkipHeader: true,
					Comma:      'b',
				},
			},
			file: "./new_client_test.csv",
			genFile: func() {
				f, err := os.Create("./new_client_test.csv")
				if err != nil {
					panic(err)
				}
				w := csv.NewWriter(f)
				_ = w.Write([]string{"Firstname"})
				_ = w.Write([]string{strings.Repeat("e", 5)})
				w.Flush()
				_ = f.Close()
			},
			expectedR: true,
			expectedE: nil,
		},
		{
			name: "should return nil and err when file is not found",
			ops: []*Options{
				{
					SkipHeader: true,
					Comma:      'b',
				},
			},
			file:      "./new_client_test.csv",
			genFile:   nil,
			expectedR: false,
			expectedE: errors.New("open ./new_client_test.csv: no such file or directory"),
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if tc.genFile != nil {
				tc.genFile()
			}
			r, e := NewClient[Student](tc.file, tc.ops...)
			if fmt.Sprintf("%v", tc.expectedE) != fmt.Sprintf("%v", e) {
				t.Errorf("must:%v, but got: %v", tc.expectedE, e)
			}

			if tc.expectedR && r == nil {
				t.Errorf("must not nil, but got: nil")
			}

			if !tc.expectedR && r != nil {
				t.Errorf("must nil, but got: object")
			}

			if r != nil {
				r.Close()
			}
			_ = os.Remove("./new_client_test.csv")
		})
	}
}
func Test_initSkipper(t *testing.T) {
	m := make(map[int]int)
	m[0] = 0
	m[3] = 3
	tt := []struct {
		name      string
		skipCols  []int
		expectedR map[int]int
	}{
		{
			name:      "should return nil when skipCals is nil",
			skipCols:  nil,
			expectedR: nil,
		},
		{
			name:      "should return valid result when skipCols is not empty",
			skipCols:  []int{0, 3},
			expectedR: m,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			r := initSkipper(tc.skipCols)
			if fmt.Sprintf("%v", tc.expectedR) != fmt.Sprintf("%v", r) {
				t.Errorf("must:%v, but got: %v", tc.expectedR, r)
			}
		})
	}
}
