package fsl

import (
	"testing"

	"github.com/lsytj0413/ena"
	. "github.com/onsi/gomega"
)

func TestIsAnyErrorWildcard(t *testing.T) {
	type testCase struct {
		desp        string
		errorEquals []string
		expect      bool
	}
	testCases := []testCase{
		{
			desp:        "not only 1 ErrorEquals",
			errorEquals: []string{ErrorCodeStatesBranchFailed, ErrorCodeStatesAll},
			expect:      false,
		},
		{
			desp:        "only 1 but not ALL",
			errorEquals: []string{ErrorCodeStatesBranchFailed},
			expect:      false,
		},
		{
			desp:        "only 1 and is ALL",
			errorEquals: []string{ErrorCodeStatesAll},
			expect:      true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			g := NewWithT(t)

			g.Expect(IsAnyErrorWildcard(tc.errorEquals)).To(Equal(tc.expect))
		})
	}
}

func TestValidateErrorEquals(t *testing.T) {
	type testCase struct {
		desp        string
		errorEquals []string
		err         string
	}
	testCases := []testCase{
		{
			desp:        "normal one",
			errorEquals: []string{ErrorCodeStatesBranchFailed},
			err:         ``,
		},
		{
			desp:        "normal wildcard",
			errorEquals: []string{ErrorCodeStatesAll},
			err:         ``,
		},
		{
			desp:        "normal multi",
			errorEquals: []string{ErrorCodeStatesBranchFailed, ErrorCodeStatesHeartbeatTimeout, "tt"},
			err:         ``,
		},
		{
			desp:        "empty ErrorEquals",
			errorEquals: []string{},
			err:         `ErrorEquals cannot be empty`,
		},
		{
			desp:        "empty ErrorEqual item",
			errorEquals: []string{ErrorCodeStatesBranchFailed, ErrorCodeStatesHeartbeatTimeout, ""},
			err:         `ErrorEqual cannot be empty`,
		},
		{
			desp:        "not pre defined",
			errorEquals: []string{ErrorCodeStatesBranchFailed + "1", ErrorCodeStatesHeartbeatTimeout, "tt"},
			err:         `ErrorEqual is not pre defined`,
		},
		{
			desp:        "multi with ALL",
			errorEquals: []string{ErrorCodeStatesBranchFailed, ErrorCodeStatesAll, "tt"},
			err:         `States\.ALL must be be only element in ErrorEquals`,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			g := NewWithT(t)
			err := ValidateErrorEquals(tc.errorEquals)

			if tc.err != "" {
				g.Expect(err).To(HaveOccurred())
				g.Expect(err.Error()).To(MatchRegexp(tc.err))
				return
			}

			g.Expect(err).ToNot(HaveOccurred())
		})
	}
}

func TestValidateTerminalState(t *testing.T) {
	type testCase struct {
		desp string
		next *string
		end  bool
		err  string
	}
	testCases := []testCase{
		{
			desp: "only next",
			next: ena.PointerTo("n"),
			end:  false,
			err:  ``,
		},
		{
			desp: "only end",
			next: nil,
			end:  true,
			err:  ``,
		},
		{
			desp: "nil next & end is false",
			next: nil,
			end:  false,
			err:  `Next is nil and End is false`,
		},
		{
			desp: "non-nil next but empty",
			next: ena.PointerTo(""),
			end:  false,
			err:  `Next is configured but is empty`,
		},
		{
			desp: "both next and end configured",
			next: ena.PointerTo("n"),
			end:  true,
			err:  `Next is configured and End is true`,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			g := NewWithT(t)
			err := ValidateTerminalState(tc.next, tc.end)

			if tc.err != "" {
				g.Expect(err).To(HaveOccurred())
				g.Expect(err.Error()).To(MatchRegexp(tc.err))
				return
			}

			g.Expect(err).ToNot(HaveOccurred())
		})
	}
}
