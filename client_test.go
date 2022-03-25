package csvtogo

import "testing"

func Test_validateOps(t *testing.T) {
	tt := []struct {
		name      string
		ops       *Options
		expectedR *Options
		expectedE error
	}{
		{
			name:      "should do nothing when ops is nil",
			ops:       nil,
			expectedR: nil,
			expectedE: nil,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			r, e := validateOps(tc.ops)
			if tc.expectedE != e {
				t.Errorf("want: %v, but got: %v", tc.expectedR, e)
			}
			if tc.expectedR != r {
				t.Errorf("want: %v, but got: %v", tc.expectedR, r)
			}
		})
	}
}
