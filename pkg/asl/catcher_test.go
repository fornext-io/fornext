package asl

import (
	"context"
	"testing"

	. "github.com/onsi/gomega"

	"github.com/lsytj0413/ena"
)

func TestCatcherValidate(t *testing.T) {
	type testCase struct {
		desp string
		c    *Catcher
		err  string
	}
	testCases := []testCase{
		{
			desp: "normal multi",
			c: &Catcher{
				ErrorEquals: []string{ErrorCodeStatesBranchFailed, ErrorCodeStatesHeartbeatTimeout, "tt"},
				Next:        "n",
				ResultPath:  ena.PointerTo("r"),
			},
			err: ``,
		},
		{
			desp: "empty ErrorEquals",
			c: &Catcher{
				ErrorEquals: []string{},
				Next:        "n",
				ResultPath:  ena.PointerTo("r"),
			},
			err: `ErrorEquals cannot be empty`,
		},
		{
			desp: "empty next",
			c: &Catcher{
				ErrorEquals: []string{ErrorCodeStatesBranchFailed, ErrorCodeStatesHeartbeatTimeout, "tt"},
				Next:        "Next cannot be empty",
				ResultPath:  ena.PointerTo("r"),
			},
			err: ``,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			g := NewWithT(t)
			err := tc.c.Validate(context.Background())

			if tc.err != "" {
				g.Expect(err).To(HaveOccurred())
				g.Expect(err.Error()).To(MatchRegexp(tc.err))
				return
			}

			g.Expect(err).ToNot(HaveOccurred())
		})
	}
}

func TestCatcherIsAnyErrorWildcard(t *testing.T) {
	type testCase struct {
		desp   string
		c      *Catcher
		expect bool
	}
	testCases := []testCase{
		{
			desp: "is not wildcard",
			c: &Catcher{
				ErrorEquals: []string{ErrorCodeStatesBranchFailed, ErrorCodeStatesHeartbeatTimeout},
			},
			expect: false,
		},
		{
			desp: "is wildcard",
			c: &Catcher{
				ErrorEquals: []string{ErrorCodeStatesAll},
			},
			expect: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			g := NewWithT(t)

			g.Expect(tc.c.IsAnyErrorWildcard()).To(Equal(tc.expect))
		})
	}
}

func TestCatchersValidate(t *testing.T) {
	type testCase struct {
		desp string
		c    Catchers
		err  string
	}
	testCases := []testCase{
		{
			desp: "empty catchers",
			c:    Catchers([]Catcher{}),
			err:  ``,
		},
		{
			desp: "multi catchers",
			c: Catchers([]Catcher{
				{
					ErrorEquals: []string{ErrorCodeStatesBranchFailed, ErrorCodeStatesHeartbeatTimeout, "ERR"},
					Next:        "n",
				},
				{
					ErrorEquals: []string{ErrorCodeStatesAll},
					Next:        "n",
				},
			}),
			err: ``,
		},
		{
			desp: "catcher valid failed",
			c: Catchers([]Catcher{
				{
					ErrorEquals: []string{ErrorCodeStatesBranchFailed, ErrorCodeStatesHeartbeatTimeout, "ERR"},
					Next:        "",
				},
				{
					ErrorEquals: []string{ErrorCodeStatesAll},
					Next:        "n",
				},
			}),
			err: `Next cannot be empty`,
		},
		{
			desp: "ALL is not last",
			c: Catchers([]Catcher{
				{
					ErrorEquals: []string{ErrorCodeStatesAll},
					Next:        "n",
				},
				{
					ErrorEquals: []string{ErrorCodeStatesBranchFailed, ErrorCodeStatesHeartbeatTimeout, "ERR"},
					Next:        "n",
				},
			}),
			err: `States\.ALL must be the last catcher`,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			g := NewWithT(t)
			err := tc.c.Validate(context.Background())

			if tc.err != "" {
				g.Expect(err).To(HaveOccurred())
				g.Expect(err.Error()).To(MatchRegexp(tc.err))
				return
			}

			g.Expect(err).ToNot(HaveOccurred())
		})
	}
}
