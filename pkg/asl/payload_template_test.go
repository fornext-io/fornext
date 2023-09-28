package asl

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestApplyPayloadTemplateOnObject(t *testing.T) {
	type testCase struct {
		desp          string
		payloadObject interface{}
		obj           interface{}
		contextObj    interface{}

		expect interface{}
		err    string
	}
	testCases := []testCase{
		{
			desp: "normal 1",
			payloadObject: map[string]interface{}{
				"flagged": true,
				"parts": map[string]interface{}{
					"first.$": "$.vals[0]",
					"last3.$": "$.vals[-3:]",
				},
				"weekday.$": "$$.DayOfWeek",
				// "formattedOutput.$": "States.Format('Today is {}', $$.DayOfWeek)"
			},
			obj: map[string]interface{}{
				"flagged": 7,
				"vals":    []int{0, 10, 20, 30, 40, 50},
			},
			contextObj: map[string]interface{}{
				"DayOfWeek": "TUESDAY",
			},
			expect: map[string]interface{}{
				"flagged": true,
				"parts": map[string]interface{}{
					"first": 0,
					// NOTE: this is an bug, the jp return it with Reverse order
					"last3": []interface{}{50, 40, 30},
				},
				"weekday": "TUESDAY",
				// "formattedOutput": "Today is TUESDAY",
			},
			err: ``,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			g := NewWithT(t)

			actual, err := ApplyPayloadTemplateOnObject(tc.payloadObject, tc.obj, tc.contextObj)
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
