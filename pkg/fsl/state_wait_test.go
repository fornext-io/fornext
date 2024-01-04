package fsl

import (
	"encoding/json"
	"testing"

	"github.com/lsytj0413/ena"
	. "github.com/onsi/gomega"
)

func TestWaitStateUnmarshalJSON(t *testing.T) {
	type testCase struct {
		desp string
		b    string

		err    string
		expect *WaitState
	}
	testCases := []testCase{
		{
			desp: "normal with seconds",
			b:    `{"Type": "Wait", "Seconds": 10, "Next": "NextState"}`,
			err:  ``,
			expect: &WaitState{
				Type:       StateTypeWait,
				Seconds:    ena.PointerTo(10),
				Next:       "NextState",
				InputPath:  "$",
				OutputPath: "$",
			},
		},
		{
			desp: "normal with timestamp",
			b:    `{"Type": "Wait", "Timestamp": "2016-03-14T01:59:00Z", "Next": "NextState"}`,
			err:  ``,
			expect: &WaitState{
				Type:       StateTypeWait,
				Timestamp:  MustParseTime("2016-03-14T01:59:00Z"),
				Next:       "NextState",
				InputPath:  "$",
				OutputPath: "$",
			},
		},
		{
			desp: "normal with timestamppath",
			b:    `{"Type": "Wait", "TimestampPath": "$.expirydate", "Next": "NextState"}`,
			err:  ``,
			expect: &WaitState{
				Type:          StateTypeWait,
				TimestampPath: (*ReferencePath)(ena.PointerTo("$.expirydate")),
				Next:          "NextState",
				InputPath:     "$",
				OutputPath:    "$",
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			g := NewWithT(t)

			actual := &WaitState{}
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
