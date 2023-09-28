package asl

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestApplyReferencePath(t *testing.T) {
	type testCase struct {
		desp   string
		exp    string
		object interface{}
		value  interface{}

		expect interface{}
		err    string
	}
	testCases := []testCase{
		{
			desp: "with array overwrite",
			exp:  `$.master.detail`,
			object: map[string]interface{}{
				"master": map[string]interface{}{
					"detail": []int{1, 2, 3},
				},
			},
			value: 6,
			expect: map[string]interface{}{
				"master": map[string]interface{}{
					"detail": 6,
				},
			},
			err: ``,
		},
		{
			desp: "with field insert",
			exp:  `$.master.result.sum`,
			object: map[string]interface{}{
				"master": map[string]interface{}{
					"detail": []int{1, 2, 3},
				},
			},
			value: 6,
			expect: map[string]interface{}{
				"master": map[string]interface{}{
					"detail": []int{1, 2, 3},
					"result": map[string]interface{}{
						"sum": 6,
					},
				},
			},
			err: ``,
		},
		{
			desp: "with sum field insert",
			exp:  `$.sum`,
			object: map[string]interface{}{
				"title": "Numbers to add",
				"numbers": map[string]interface{}{
					"val1": 3,
					"val2": 4,
				},
			},
			value: 7,
			expect: map[string]interface{}{
				"title": "Numbers to add",
				"numbers": map[string]interface{}{
					"val1": 3,
					"val2": 4,
				},
				"sum": 7,
			},
			err: ``,
		},
		{
			desp: `with object field insert`,
			exp:  `$.coords`,
			object: map[string]interface{}{
				"x-datum": 0.381018,
				"y-datum": 622.2269926397355,
			},
			value: map[string]interface{}{
				"georefOf": "Home",
			},
			expect: map[string]interface{}{
				"x-datum": 0.381018,
				"y-datum": 622.2269926397355,
				"coords": map[string]interface{}{
					"georefOf": "Home",
				},
			},
			err: ``,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			g := NewWithT(t)

			actual, err := ApplyReferencePath(&tc.exp, tc.object, tc.value)
			if tc.err != "" {
				g.Expect(err).To(HaveOccurred())
				g.Expect(err.Error()).To(MatchRegexp(tc.err))
				return
			}

			g.Expect(err).ToNot(HaveOccurred())
			g.Expect(actual).To(Equal(tc.expect))
		})
	}
}
