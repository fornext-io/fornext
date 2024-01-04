package fsl

import (
	"context"
	"testing"

	. "github.com/onsi/gomega"

	"github.com/fornext-io/fornext/pkg/utils"
)

func TestReferencePathApply(t *testing.T) {
	type testCase struct {
		desp string
		p    ReferencePath
		pc   referencePathContext

		expect []byte
		err    string
	}
	testCases := []testCase{
		{
			desp: "with array overwrite",
			p:    ReferencePath(`$.master.detail`),
			pc: referencePathContext{
				Input:  []byte(`{"master":{"detail":[1,2,3]}}`),
				Output: []byte(`6`),
			},
			expect: []byte(`{"master":{"detail":6}}`),
			err:    ``,
		},
		{
			desp: "with field insert",
			p:    ReferencePath(`$.master.result.sum`),
			pc: referencePathContext{
				Input:  []byte(`{"master":{"detail":[1,2,3]}}`),
				Output: []byte(`6`),
			},
			expect: []byte(`{"master":{"detail":[1,2,3],"result":{"sum":6}}}`),
			err:    ``,
		},
		{
			desp: "with sum fiedl insert",
			p:    ReferencePath(`$.sum`),
			pc: referencePathContext{
				Input:  []byte(`{"title":"Numbers to add","numbers":{"val1":3,"val2":4}}`),
				Output: []byte(`7`),
			},
			expect: []byte(`{"title":"Numbers to add","numbers":{"val1":3,"val2":4},"sum":7}`),
			err:    ``,
		},
		{
			desp: "with object field insert",
			p:    ReferencePath(`$.coords`),
			pc: referencePathContext{
				Input:  []byte(`{"x-datum":0.381018,"y-datum":622.2269926397355}`),
				Output: []byte(`{"georefOf":"Home"}`),
			},
			expect: []byte(`{"x-datum":0.381018,"y-datum":622.2269926397355,"coords":{"georefOf":"Home"}}`),
			err:    ``,
		},
		{
			desp: "empty reference path",
			p:    ReferencePath(``),
			pc: referencePathContext{
				Input:  []byte(`{"master":{"detail":[1,2,3]}}`),
				Output: []byte(`6`),
			},
			expect: []byte(`{"master":{"detail":[1,2,3]}}`),
			err:    ``,
		},
		{
			desp: "root reference path",
			p:    ReferencePath(`$`),
			pc: referencePathContext{
				Input:  []byte(`{"master":{"detail":[1,2,3]}}`),
				Output: []byte(`6`),
			},
			expect: []byte(`6`),
			err:    ``,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			g := NewWithT(t)

			actual, err := tc.p.Apply(context.Background(), tc.pc)
			if tc.err != "" {
				g.Expect(err).To(HaveOccurred())
				g.Expect(err.Error()).To(MatchRegexp(tc.err))
				return
			}

			g.Expect(err).ToNot(HaveOccurred())
			g.Expect(utils.MustUnmarshalJSON(actual)).To(Equal(utils.MustUnmarshalJSON(tc.expect)))
		})
	}
}
