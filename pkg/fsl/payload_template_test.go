package fsl

import (
	"context"
	"testing"

	. "github.com/onsi/gomega"

	"github.com/fornext-io/fornext/pkg/utils"
)

func TestPayloadTemplateApply(t *testing.T) {
	type testCase struct {
		desp   string
		p      PayloadTemplate
		pc     payloadTemplateContext
		expect []byte
		err    string
	}
	testCases := []testCase{
		{
			desp: "normal example",
			p: PayloadTemplate(`{
				  "flagged": true,
				  "parts": {
					"first.$": "$.vals[0]",
					"last3.$": "$.vals[-3:]"
				  },
				  "weekday.$": "$$.DayOfWeek"
				}`),
			pc: payloadTemplateContext{
				Input: []byte(`{
					"flagged": 7,
					"vals": [0, 10, 20, 30, 40, 50]
				  }`),
				ContextData: []byte(`{
					"DayOfWeek": "TUESDAY"
				  }`),
			},
			expect: []byte(`{"flagged":true,"parts":{"first":0,"last3":[30,40,50]},"weekday":"TUESDAY"}`),
			err:    ``,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			g := NewWithT(t)

			actual, err := tc.p.Apply(context.Background(), tc.pc)
			if tc.err != `` {
				g.Expect(err).To(HaveOccurred())
				g.Expect(err.Error()).To(MatchRegexp(tc.err))
				return
			}

			g.Expect(err).ToNot(HaveOccurred())
			g.Expect(utils.MustUnmarshalJSON(actual)).To(Equal(utils.MustUnmarshalJSON(tc.expect)))
		})
	}
}
