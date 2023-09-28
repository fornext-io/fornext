package asl

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestApplyPathOnObject(t *testing.T) {
	type testCase struct {
		desp  string
		exp   string
		input interface{}

		expect interface{}
		err    string
	}
	testCases := []testCase{
		{
			desp: "with result sub-array",
			exp:  `$.a[0,1]`,
			input: map[string]interface{}{
				"a": []int{1, 2, 3, 4},
			},
			expect: []interface{}{1, 2},
			err:    ``,
		},
		{
			desp: "with result object",
			exp:  `$.numbers`,
			input: map[string]interface{}{
				"title": "Numbers to add",
				"numbers": map[string]interface{}{
					"val1": 3,
					"val2": 4,
				},
			},
			expect: map[string]interface{}{
				"val1": 3,
				"val2": 4,
			},
			err: ``,
		},
		{
			desp: "with non-exists field",
			exp:  `$.numbers1`,
			input: map[string]interface{}{
				"title": "Numbers to add",
				"numbers": map[string]interface{}{
					"val1": 3,
					"val2": 4,
				},
			},
			expect: map[string]interface{}{},
			err:    ``,
		},
		{
			desp: "with non exists sub-array",
			exp:  `$.a[100,101]`,
			input: map[string]interface{}{
				"a": []int{1, 2, 3, 4},
			},
			expect: map[string]interface{}{},
			err:    ``,
		},
		{
			desp: "with null field",
			exp:  `$.numbers`,
			input: map[string]interface{}{
				"title":   "Numbers to add",
				"numbers": nil,
			},
			expect: nil,
			err:    ``,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			g := NewWithT(t)

			actual, err := ApplyPathOnObject(tc.exp, tc.input)
			if tc.err != "" {
				g.Expect(err).To(HaveOccurred())
				g.Expect(err.Error()).To(MatchRegexp(tc.err))
				return
			}

			g.Expect(err).ToNot(HaveOccurred())
			if tc.expect == nil {
				g.Expect(actual).To(BeNil())
				return
			}

			g.Expect(actual).To(Equal(tc.expect))
		})
	}
}
