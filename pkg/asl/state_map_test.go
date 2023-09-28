package asl

import (
	"encoding/json"
	"testing"

	"github.com/lsytj0413/ena"
	. "github.com/onsi/gomega"
)

func TestMapStateUnmarshalJSON(t *testing.T) {
	type testCase struct {
		desp string
		b    string

		err    string
		expect *MapState
	}
	testCases := []testCase{
		{
			desp: "normal",
			b: `
{
	"Type": "Map",
	"InputPath": "$.detail",
	"ItemsPath": "$.shipped",
	"MaxConcurrency": 0,
	"ItemProcessor": {
		"StartAt": "Validate",
		"States": {
		  "Validate": {
			"Type": "Task",
			"Resource": "arn:aws:lambda:us-east-1:123456789012:function:ship-val",
			"End": true
		  }
		}
	  },
	  "ResultPath": "$.detail.shipped",
	  "End": true
}`,
			err: ``,
			expect: &MapState{
				Type:           StateTypeMap,
				InputPath:      (*Path)(ena.PointerTo("$.detail")),
				ItemsPath:      (*ReferencePath)(ena.PointerTo("$.shipped")),
				MaxConcurrency: ena.PointerTo(0),
				ItemProcessor: ItemProcessor{
					StartAt: "Validate",
					States: map[string]State{
						"Validate": &TaskState{
							Type:     StateTypeTask,
							Resource: "arn:aws:lambda:us-east-1:123456789012:function:ship-val",
							End:      true,
						},
					},
				},
				ResultPath: (*ReferencePath)(ena.PointerTo(`$.detail.shipped`)),
				End:        true,
			},
		},
		{
			desp: "normal with item selector",
			b: `
{
	"Type": "Map",
	"InputPath": "$.detail",
	"ItemsPath": "$.shipped",
	"MaxConcurrency": 0,
	"ItemSelector": {"parcel.$": "$$.Map.Item.Value","courier.$": "$.delivery-partner"},
	"ItemProcessor": {
		"StartAt": "Validate",
		"States": {
		  "Validate": {
			"Type": "Task",
			"Resource": "arn:aws:lambda:us-east-1:123456789012:function:ship-val",
			"End": true
		  }
		}
	  },
	  "ResultPath": "$.detail.shipped",
	  "End": true
}`,
			err: ``,
			expect: &MapState{
				Type:           StateTypeMap,
				InputPath:      (*Path)(ena.PointerTo("$.detail")),
				ItemsPath:      (*ReferencePath)(ena.PointerTo("$.shipped")),
				MaxConcurrency: ena.PointerTo(0),
				ItemSelector: (*PayloadTemplate)(ena.PointerTo(
					`{"parcel.$": "$$.Map.Item.Value","courier.$": "$.delivery-partner"}`)),
				ItemProcessor: ItemProcessor{
					StartAt: "Validate",
					States: map[string]State{
						"Validate": &TaskState{
							Type:     StateTypeTask,
							Resource: "arn:aws:lambda:us-east-1:123456789012:function:ship-val",
							End:      true,
						},
					},
				},
				ResultPath: (*ReferencePath)(ena.PointerTo(`$.detail.shipped`)),
				End:        true,
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			g := NewWithT(t)

			actual := &MapState{}
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
