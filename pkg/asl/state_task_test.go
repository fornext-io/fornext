package asl

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
