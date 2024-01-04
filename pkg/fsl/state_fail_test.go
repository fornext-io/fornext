package fsl

import (
	"encoding/json"
	"testing"

	. "github.com/onsi/gomega"
)

func TestFailStateUnmarshalJSON(t *testing.T) {
	type testCase struct {
		desp string
		b    string

		err    string
		expect *FailState
	}
	testCases := []testCase{
		{
			desp: "normal",
			b:    `{"Type": "Fail", "Error": "ErrorA", "Cause": "Kaiju attack"}`,
			err:  ``,
			expect: &FailState{
				Type:  StateTypeFail,
				Error: "ErrorA",
				Cause: "Kaiju attack",
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			g := NewWithT(t)

			actual := &FailState{}
			err := json.Unmarshal([]byte(tc.b), actual)
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
