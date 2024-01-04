package fsl

import (
	"context"
	"testing"

	. "github.com/onsi/gomega"
)

func TestPathApply(t *testing.T) {
	type testCase struct {
		desp   string
		p      Path
		pc     pathContext
		expect []byte
		err    string
	}
	testCases := []testCase{
		{
			desp: "normal example",
			p:    Path(`$.numbers`),
			pc: pathContext{
				Input: []byte(`{"title": "Numbers to add", "numbers": {"val1": 3, "val2": 4}}`),
			},
			expect: []byte(`{"val1":3,"val2":4}`),
			err:    "",
		},
		{
			desp: "empty path",
			p:    Path(``),
			pc: pathContext{
				Input: []byte(`{"title": "Numbers to add", "numbers": {"val1": 3, "val2": 4}}`),
			},
			expect: []byte(`{}`),
			err:    "",
		},
		{
			desp: "root path",
			p:    Path(`$`),
			pc: pathContext{
				Input: []byte(`{"title": "Numbers to add", "numbers": {"val1": 3, "val2": 4}}`),
			},
			expect: []byte(`{"title": "Numbers to add", "numbers": {"val1": 3, "val2": 4}}`),
			err:    "",
		},
		{
			desp: "input unmarshal failed",
			p:    Path(`$.numbers`),
			pc: pathContext{
				Input: []byte(`{"title": "Numbers to add", "numbers": {"val1": 3, "val2": 4}`),
			},
			expect: []byte(`{"val1":3,"val2":4}`),
			err:    `unexpected end of JSON input`,
		},
		{
			desp: "normal context example",
			p:    Path(`$$.numbers`),
			pc: pathContext{
				ContextData: []byte(`{"title": "Numbers to add", "numbers": {"val1": 3, "val2": 4}}`),
			},
			expect: []byte(`{"val1":3,"val2":4}`),
			err:    "",
		},
		{
			desp: "context unmarshal failed",
			p:    Path(`$$.numbers`),
			pc: pathContext{
				ContextData: []byte(`{"title": "Numbers to add", "numbers": {"val1": 3, "val2": 4}`),
			},
			expect: []byte(`{"val1":3,"val2":4}`),
			err:    "unexpected end of JSON input",
		},
		{
			desp: "invalid expr",
			p:    Path(`$.`),
			pc: pathContext{
				Input: []byte(`{"title": "Numbers to add", "numbers": {"val1": 3, "val2": 4}}`),
			},
			expect: []byte(`{"val1":3,"val2":4}`),
			err:    "not terminated at 3 in",
		},
		{
			desp: "normal example without elements",
			p:    Path(`$.numbers1`),
			pc: pathContext{
				Input: []byte(`{"title": "Numbers to add", "numbers": {"val1": 3, "val2": 4}}`),
			},
			expect: []byte(`{"val1":3,"val2":4}`),
			err:    "must at least one element InputPath",
		},
		{
			desp: "normal example without multi elements",
			p:    Path(`$.numbers[*]`),
			pc: pathContext{
				Input: []byte(`{"title": "Numbers to add", "numbers": [3, 4]}`),
			},
			expect: []byte(`[3,4]`),
			err:    "",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			g := NewWithT(t)

			actual, err := tc.p.Apply(context.Background(), tc.pc)
			if tc.err != `` {
				g.Expect(err).To(HaveOccurred())
				g.Expect(err.Error()).To(MatchRegexp(tc.err))
				return
			}

			g.Expect(err).ToNot(HaveOccurred())
			g.Expect(actual).To(Equal(tc.expect))
		})
	}
}
