package asl

import (
	"context"
	"testing"

	"github.com/lsytj0413/ena"
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
				IntervalSeconds: ena.PointerTo(1),
				MaxAttempts:     ena.PointerTo(3),
				BackoffRate:     ena.PointerTo(float32(1.0)),
			},
			err: ``,
		},
		{
			desp: "normal without interval seconds",
			r: &Retrier{
				ErrorEquals: []string{ErrorCodeStatesBranchFailed, ErrorCodeStatesHeartbeatTimeout, "tt"},
				MaxAttempts: ena.PointerTo(3),
				BackoffRate: ena.PointerTo(float32(1.0)),
			},
			err: ``,
		},
		{
			desp: "normal without max attempts",
			r: &Retrier{
				ErrorEquals:     []string{ErrorCodeStatesBranchFailed, ErrorCodeStatesHeartbeatTimeout, "tt"},
				IntervalSeconds: ena.PointerTo(1),
				BackoffRate:     ena.PointerTo(float32(1.0)),
			},
			err: ``,
		},
		{
			desp: "normal without backoff rate",
			r: &Retrier{
				ErrorEquals:     []string{ErrorCodeStatesBranchFailed, ErrorCodeStatesHeartbeatTimeout, "tt"},
				IntervalSeconds: ena.PointerTo(1),
				MaxAttempts:     ena.PointerTo(3),
			},
			err: ``,
		},
		{
			desp: "empty ErrorEquals",
			r: &Retrier{
				ErrorEquals:     []string{},
				IntervalSeconds: ena.PointerTo(1),
				MaxAttempts:     ena.PointerTo(3),
				BackoffRate:     ena.PointerTo(float32(1.0)),
			},
			err: `ErrorEquals cannot be empty`,
		},
		{
			desp: "invalid interval seconds",
			r: &Retrier{
				ErrorEquals:     []string{ErrorCodeStatesBranchFailed, ErrorCodeStatesHeartbeatTimeout, "tt"},
				IntervalSeconds: ena.PointerTo(0),
				MaxAttempts:     ena.PointerTo(3),
				BackoffRate:     ena.PointerTo(float32(1.0)),
			},
			err: `IntervalSeconds must be an positive integer`,
		},
		{
			desp: "invalid max attempts",
			r: &Retrier{
				ErrorEquals:     []string{ErrorCodeStatesBranchFailed, ErrorCodeStatesHeartbeatTimeout, "tt"},
				IntervalSeconds: ena.PointerTo(1),
				MaxAttempts:     ena.PointerTo(-1),
				BackoffRate:     ena.PointerTo(float32(1.0)),
			},
			err: `MaxAttempts must be a non-negative integer`,
		},
		{
			desp: "invalid backoff rate",
			r: &Retrier{
				ErrorEquals:     []string{ErrorCodeStatesBranchFailed, ErrorCodeStatesHeartbeatTimeout, "tt"},
				IntervalSeconds: ena.PointerTo(1),
				MaxAttempts:     ena.PointerTo(3),
				BackoffRate:     ena.PointerTo(float32(0.5)),
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
					IntervalSeconds: ena.PointerTo(1),
				},
				{
					ErrorEquals:     []string{ErrorCodeStatesAll},
					IntervalSeconds: ena.PointerTo(1),
				},
			}),
			err: ``,
		},
		{
			desp: "retrier valid failed",
			r: Retriers([]Retrier{
				{
					ErrorEquals:     []string{ErrorCodeStatesBranchFailed, ErrorCodeStatesHeartbeatTimeout, "ERR"},
					IntervalSeconds: ena.PointerTo(0),
				},
				{
					ErrorEquals:     []string{ErrorCodeStatesAll},
					IntervalSeconds: ena.PointerTo(1),
				},
			}),
			err: `IntervalSeconds must be an positive integer`,
		},
		{
			desp: "ALL is not last",
			r: Retriers([]Retrier{
				{
					ErrorEquals:     []string{ErrorCodeStatesAll},
					IntervalSeconds: ena.PointerTo(1),
				},
				{
					ErrorEquals:     []string{ErrorCodeStatesBranchFailed, ErrorCodeStatesHeartbeatTimeout, "ERR"},
					IntervalSeconds: ena.PointerTo(1),
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
