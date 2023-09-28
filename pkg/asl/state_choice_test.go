package asl

import (
	"encoding/json"
	"testing"

	"github.com/lsytj0413/ena"
	. "github.com/onsi/gomega"
)

func TestChoiceStateUnmarshalJSON(t *testing.T) {
	type testCase struct {
		desp string
		b    string

		err    string
		expect *ChoiceState
	}
	testCases := []testCase{
		{
			desp: "normal with seconds",
			b: `
{
	"Type": "Choice",
	"Choices": [
		{
			"Not": {
				"Variable": "$.type",
				"StringEquals": "Private"
			},
			"Next": "Public"
		},
		{
			"Variable": "$.value",
			"NumericEquals": 0,
			"Next": "ValueIsZero"
		},
		{
			"And": [
				{
					"Variable": "$.value",
					"NumericGreaterThanEquals": 20
				},
				{
					"Variable": "$.value",
					"NumericLessThan": 30
				}
			],
			"Next": "ValueInTwenties"
		}
	],
	"Default": "DefaultState"
}`,
			err: ``,
			expect: &ChoiceState{
				Type: StateTypeChoice,
				Choices: []Choice{
					{
						Next: "Public",
						ChoiceRule: ChoiceRule{
							Not: &ChoiceRule{
								Variable:     (*Path)(ena.PointerTo("$.type")),
								StringEquals: ena.PointerTo("Private"),
							},
						},
					},
					{
						Next: "ValueIsZero",
						ChoiceRule: ChoiceRule{
							Variable:      (*Path)(ena.PointerTo("$.value")),
							NumericEquals: ena.PointerTo(float64(0)),
						},
					},
					{
						Next: "ValueInTwenties",
						ChoiceRule: ChoiceRule{
							And: []ChoiceRule{
								{
									Variable:                 (*Path)(ena.PointerTo("$.value")),
									NumericGreaterThanEquals: ena.PointerTo(float64(20)),
								},
								{
									Variable:        (*Path)(ena.PointerTo("$.value")),
									NumericLessThan: ena.PointerTo(float64(30)),
								},
							},
						},
					},
				},
				Default: ena.PointerTo("DefaultState"),
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			g := NewWithT(t)

			actual := &ChoiceState{}
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
