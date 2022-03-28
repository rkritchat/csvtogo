package csvtogo

import "testing"

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
	type Student struct {
		Firstname string
	}
	tt := []struct {
		name string
		ref  Student
		data []string
		row  int
	}{
		{
			name: "should ",
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			//c := Executor[Student]
		})
	}
}
