package fsl

import (
	"encoding/json"
	"testing"

	"github.com/lsytj0413/ena"
	. "github.com/onsi/gomega"
)

func TestTaskStateUnmarshalJSON(t *testing.T) {
	type testCase struct {
		desp string
		b    string

		err    string
		expect *TaskState
	}
	testCases := []testCase{
		{
			desp: "normal with timeout and heartbeat",
			b: `{
				"Type": "Task", 
				"Resource": "arn:aws:states:us-east-1:123456789012:activity:HelloWorld", 
				"TimeoutSeconds": 300, 
				"HeartbeatSeconds": 60 , 
				"Next": "NextState"}`,
			err: ``,
			expect: &TaskState{
				Type:             StateTypeTask,
				Resource:         "arn:aws:states:us-east-1:123456789012:activity:HelloWorld",
				TimeoutSeconds:   ena.PointerTo(300),
				HeartbeatSeconds: ena.PointerTo(60),
				Next:             "NextState",
				InputPath:        "$",
				OutputPath:       "$",
				ResultPath:       "$",
			},
		},
		{
			desp: "normal with dynamic timeout and heartbeat",
			b: `{
				"Type": "Task", 
				"Resource": "arn:aws:states:::glue:startJobRun.sync", 
				"Parameters":{"JobName": "myGlueJob"},  
				"TimeoutSecondsPath": "$.params.maxTime", 
				"HeartbeatSecondsPath": "$.params.heartbeat", 
				"Next": "NextState"}`,
			err: ``,
			expect: &TaskState{
				Type:                 StateTypeTask,
				Resource:             "arn:aws:states:::glue:startJobRun.sync",
				Parameters:           (*PayloadTemplate)(ena.PointerTo(`{"JobName": "myGlueJob"}`)),
				TimeoutSecondsPath:   (*ReferencePath)(ena.PointerTo("$.params.maxTime")),
				HeartbeatSecondsPath: (*ReferencePath)(ena.PointerTo("$.params.heartbeat")),
				Next:                 "NextState",
				InputPath:            "$",
				OutputPath:           "$",
				ResultPath:           "$",
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			g := NewWithT(t)

			actual := &TaskState{}
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

// func TestTaskStateGetJobInput(t *testing.T) {
// 	type testCase struct {
// 		desp          string
// 		s             *TaskState
// 		obj           interface{}
// 		contextObject interface{}

// 		expect interface{}
// 		err    string
// 	}
// 	testCases := []testCase{
// 		{
// 			desp: "normal test",
// 			s: &TaskState{
// 				InputPath: "$.key1",
// 				Parameters: ena.PointerTo(PayloadTemplate(`{
// 					"p1": "pv1",
// 					"pk11.$": "$.k11"
// 				}`)),
// 				ResultSelector: ena.PointerTo(PayloadTemplate(`{
// 					"r1": "v1",
// 					"rk11.$": "$"
// 				}`)),
// 				ResultPath: "$.rp",
// 				OutputPath: "$.rp.r1",
// 			},
// 			obj: map[string]interface{}{
// 				"key1": map[string]interface{}{
// 					"k11": "v11",
// 				},
// 			},
// 			contextObject: nil,
// 			expect:        "v1",
// 			err:           "",
// 		},
// 	}
// 	for _, tc := range testCases {
// 		t.Run(tc.desp, func(t *testing.T) {
// 			g := NewWithT(t)

// 			actual, err := tc.s.GetJobInput(tc.obj, tc.contextObject)
// 			if tc.err != "" {
// 				g.Expect(err).To(HaveOccurred())
// 				g.Expect(err.Error()).To(MatchRegexp(tc.err))
// 				return
// 			}

// 			g.Expect(err).ToNot(HaveOccurred())
// 			g.Expect(actual).To(Equal(tc.expect))
// 		})
// 	}
// }
