package fsl

import (
	"encoding/json"
	"testing"

	"github.com/lsytj0413/ena"
	. "github.com/onsi/gomega"
)

func TestPassStateUnmarshalJSON(t *testing.T) {
	type testCase struct {
		desp string
		b    string

		err    string
		expect *PassState
	}
	testCases := []testCase{
		{
			desp: "normal",
			b: `
			{"Type": "Pass", 
			"Result": {"x-datum": 0.381018, "y-datum": 622.2269926397355}, "ResultPath": "$.coords", "End": true}`,
			err: ``,
			expect: &PassState{
				Type:       StateTypePass,
				Result:     (*PayloadTemplate)(ena.PointerTo(`{"x-datum": 0.381018, "y-datum": 622.2269926397355}`)),
				InputPath:  "$",
				OutputPath: "$",
				ResultPath: "$.coords",
				End:        true,
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			g := NewWithT(t)

			actual := &PassState{}
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
