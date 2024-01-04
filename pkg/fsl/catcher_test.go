package fsl

import (
	"context"
	"encoding/json"
	"testing"

	. "github.com/onsi/gomega"
)

func TestCatcherUnmarshalJSON(t *testing.T) {
	type testCase struct {
		desp string
		data string

		expect *Catcher
		err    string
	}
	testCases := []testCase{
		{
			desp: "normal test",
			data: `{"ErrorEquals": ["1"], "Next": "2", "ResultPath": "3"}`,
			expect: &Catcher{
				ErrorEquals: []string{"1"},
				Next:        "2",
				ResultPath:  "3",
			},
			err: ``,
		},
		{
			desp: "normal test without ResultPath",
			data: `{"ErrorEquals": ["1"], "Next": "2"}`,
			expect: &Catcher{
				ErrorEquals: []string{"1"},
				Next:        "2",
				ResultPath:  "$",
			},
			err: ``,
		},
		{
			desp: "unmarshal failed",
			data: `{"ErrorEquals": ["1"], "Next": 2}`,
			expect: &Catcher{
				ErrorEquals: []string{"1"},
				Next:        "2",
				ResultPath:  "$",
			},
			err: `cannot unmarshal number into Go struct field CatcherForUnmarshal.Next of type string`,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			g := NewWithT(t)

			actual := &Catcher{}
			err := json.Unmarshal([]byte(tc.data), actual)
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
				ResultPath:  "r",
			},
			err: ``,
		},
		{
			desp: "empty ErrorEquals",
			c: &Catcher{
				ErrorEquals: []string{},
				Next:        "n",
				ResultPath:  "r",
			},
			err: `ErrorEquals cannot be empty`,
		},
		{
			desp: "empty next",
			c: &Catcher{
				ErrorEquals: []string{ErrorCodeStatesBranchFailed, ErrorCodeStatesHeartbeatTimeout, "tt"},
				Next:        "Next cannot be empty",
				ResultPath:  "r",
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
