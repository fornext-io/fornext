package fsl

import (
	"context"
	"testing"

	. "github.com/onsi/gomega"
)

func TestRetrierValidate(t *testing.T) {
	type testCase struct {
		desp string
		r    *Retrier
		err  string
	}
	testCases := []testCase{
		{
			desp: "normal",
			r: &Retrier{
				ErrorEquals:     []string{ErrorCodeStatesBranchFailed, ErrorCodeStatesHeartbeatTimeout, "tt"},
				IntervalSeconds: 1,
				MaxAttempts:     3,
				BackoffRate:     1.0,
			},
			err: ``,
		},
		{
			desp: "empty ErrorEquals",
			r: &Retrier{
				ErrorEquals:     []string{},
				IntervalSeconds: 1,
				MaxAttempts:     3,
				BackoffRate:     1.0,
			},
			err: `ErrorEquals cannot be empty`,
		},
		{
			desp: "invalid interval seconds",
			r: &Retrier{
				ErrorEquals:     []string{ErrorCodeStatesBranchFailed, ErrorCodeStatesHeartbeatTimeout, "tt"},
				IntervalSeconds: 0,
				MaxAttempts:     3,
				BackoffRate:     1.0,
			},
			err: `IntervalSeconds must be an positive integer`,
		},
		{
			desp: "invalid max attempts",
			r: &Retrier{
				ErrorEquals:     []string{ErrorCodeStatesBranchFailed, ErrorCodeStatesHeartbeatTimeout, "tt"},
				IntervalSeconds: 1,
				MaxAttempts:     -1,
				BackoffRate:     1.0,
			},
			err: `MaxAttempts must be a non-negative integer`,
		},
		{
			desp: "invalid backoff rate",
			r: &Retrier{
				ErrorEquals:     []string{ErrorCodeStatesBranchFailed, ErrorCodeStatesHeartbeatTimeout, "tt"},
				IntervalSeconds: 1,
				MaxAttempts:     3,
				BackoffRate:     0.5,
			},
			err: `BackoffRate must be greater than or equal to 1.0`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			g := NewWithT(t)
			err := tc.r.Validate(context.Background())

			if tc.err != "" {
				g.Expect(err).To(HaveOccurred())
				g.Expect(err.Error()).To(MatchRegexp(tc.err))
				return
			}

			g.Expect(err).ToNot(HaveOccurred())
		})
	}
}

func TestRetrierIsAnyErrorWildcard(t *testing.T) {
	type testCase struct {
		desp   string
		r      *Retrier
		expect bool
	}
	testCases := []testCase{
		{
			desp: "is not wildcard",
			r: &Retrier{
				ErrorEquals: []string{ErrorCodeStatesBranchFailed, ErrorCodeStatesHeartbeatTimeout},
			},
			expect: false,
		},
		{
			desp: "is wildcard",
			r: &Retrier{
				ErrorEquals: []string{ErrorCodeStatesAll},
			},
			expect: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			g := NewWithT(t)

			g.Expect(tc.r.IsAnyErrorWildcard()).To(Equal(tc.expect))
		})
	}
}

func TestRetriersValidate(t *testing.T) {
	type testCase struct {
		desp string
		r    Retriers
		err  string
	}
	testCases := []testCase{
		{
			desp: "empty retriers",
			r:    Retriers([]Retrier{}),
			err:  ``,
		},
		{
			desp: "multi retriers",
			r: Retriers([]Retrier{
				{
					ErrorEquals:     []string{ErrorCodeStatesBranchFailed, ErrorCodeStatesHeartbeatTimeout, "ERR"},
					IntervalSeconds: 1,
					MaxAttempts:     3,
					BackoffRate:     2.0,
				},
				{
					ErrorEquals:     []string{ErrorCodeStatesAll},
					IntervalSeconds: 1,
					MaxAttempts:     3,
					BackoffRate:     2.0,
				},
			}),
			err: ``,
		},
		{
			desp: "retrier valid failed",
			r: Retriers([]Retrier{
				{
					ErrorEquals:     []string{ErrorCodeStatesBranchFailed, ErrorCodeStatesHeartbeatTimeout, "ERR"},
					IntervalSeconds: 0,
					MaxAttempts:     3,
					BackoffRate:     2.0,
				},
				{
					ErrorEquals:     []string{ErrorCodeStatesAll},
					IntervalSeconds: 1,
					MaxAttempts:     3,
					BackoffRate:     2.0,
				},
			}),
			err: `IntervalSeconds must be an positive integer`,
		},
		{
			desp: "ALL is not last",
			r: Retriers([]Retrier{
				{
					ErrorEquals:     []string{ErrorCodeStatesAll},
					IntervalSeconds: 1,
					MaxAttempts:     3,
					BackoffRate:     2.0,
				},
				{
					ErrorEquals:     []string{ErrorCodeStatesBranchFailed, ErrorCodeStatesHeartbeatTimeout, "ERR"},
					IntervalSeconds: 1,
					MaxAttempts:     3,
					BackoffRate:     2.0,
				},
			}),
			err: `States\.ALL must be the last retrier`,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			g := NewWithT(t)
			err := tc.r.Validate(context.Background())

			if tc.err != "" {
				g.Expect(err).To(HaveOccurred())
				g.Expect(err.Error()).To(MatchRegexp(tc.err))
				return
			}

			g.Expect(err).ToNot(HaveOccurred())
		})
	}
}
