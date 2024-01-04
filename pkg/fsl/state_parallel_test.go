package fsl

import (
	"encoding/json"
	"testing"

	"github.com/lsytj0413/ena"
	. "github.com/onsi/gomega"
)

func TestParallelStateUnmarshalJSON(t *testing.T) {
	type testCase struct {
		desp string
		b    string

		err    string
		expect *ParallelState
	}
	testCases := []testCase{
		{
			desp: "normal example",
			b: `
{
	"Type": "Parallel", 
	"End": true, 
	"Branches": [
		{
			"StartAt": "LookupAddress", 
			"States": {
				"LookupAddress": {
					"Type": "Task",
					"Resource": "arn:aws:lambda:us-east-1:123456789012:function:AddressFinder",
					"End": true
				}
			}
		},
		{
			"StartAt": "LookupPhone",
			"States": {
				"LookupPhone": {
					"Type": "Task",
					"Resource": "arn:aws:lambda:us-east-1:123456789012:function:PhoneFinder",
					"End": true
				}
			}
		}
	]
}`,
			err: ``,
			expect: &ParallelState{
				Type:       StateTypeParallel,
				InputPath:  "$",
				OutputPath: "$",
				ResultPath: "$",
				End:        true,
				Branches: []BranchProcessor{
					{
						StartAt: "LookupAddress",
						States: map[string]State{
							"LookupAddress": &TaskState{
								Type:           StateTypeTask,
								Resource:       "arn:aws:lambda:us-east-1:123456789012:function:AddressFinder",
								End:            true,
								InputPath:      "$",
								OutputPath:     "$",
								ResultPath:     "$",
								TimeoutSeconds: ena.PointerTo(60),
							},
						},
					},
					{
						StartAt: "LookupPhone",
						States: map[string]State{
							"LookupPhone": &TaskState{
								Type:           StateTypeTask,
								Resource:       "arn:aws:lambda:us-east-1:123456789012:function:PhoneFinder",
								End:            true,
								InputPath:      "$",
								OutputPath:     "$",
								ResultPath:     "$",
								TimeoutSeconds: ena.PointerTo(60),
							},
						},
					},
				},
			},
		},
		{
			desp: "normal example input & output process",
			b: `
{
	"Type": "Parallel", 
	"End": true, 
	"Branches": [
		{
			"StartAt": "Add", 
			"States": {
				"Add": {
					"Type": "Task",
					"Resource": "arn:aws:states:us-east-1:123456789012:activity:Add",
					"End": true
				}
			}
		},
		{
			"StartAt": "Subtract",
			"States": {
				"Subtract": {
					"Type": "Task",
					"Resource": "arn:aws:states:us-east-1:123456789012:activity:Subtract",
					"End": true
				}
			}
		}
	]
}`,
			err: ``,
			expect: &ParallelState{
				Type:       StateTypeParallel,
				End:        true,
				InputPath:  "$",
				OutputPath: "$",
				ResultPath: "$",
				Branches: []BranchProcessor{
					{
						StartAt: "Add",
						States: map[string]State{
							"Add": &TaskState{
								Type:           StateTypeTask,
								Resource:       "arn:aws:states:us-east-1:123456789012:activity:Add",
								End:            true,
								InputPath:      "$",
								OutputPath:     "$",
								ResultPath:     "$",
								TimeoutSeconds: ena.PointerTo(60),
							},
						},
					},
					{
						StartAt: "Subtract",
						States: map[string]State{
							"Subtract": &TaskState{
								Type:           StateTypeTask,
								Resource:       "arn:aws:states:us-east-1:123456789012:activity:Subtract",
								End:            true,
								InputPath:      "$",
								OutputPath:     "$",
								ResultPath:     "$",
								TimeoutSeconds: ena.PointerTo(60),
							},
						},
					},
				},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			g := NewWithT(t)

			actual := &ParallelState{}
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
