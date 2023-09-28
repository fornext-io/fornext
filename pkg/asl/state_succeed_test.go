package asl

import (
	"encoding/json"
	"testing"

	. "github.com/onsi/gomega"
)

func TestSucceedStateUnmarshalJSON(t *testing.T) {
	type testCase struct {
		desp string
		b    string

		err    string
		expect *SucceedState
	}
	testCases := []testCase{
		{
			desp: "normal",
			b:    `{"Type": "Succeed"}`,
			err:  ``,
			expect: &SucceedState{
				Type: StateTypeSucceed,
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			g := NewWithT(t)

			actual := &SucceedState{}
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
