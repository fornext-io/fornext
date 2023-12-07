package jsonpath

import (
	"testing"

	"github.com/ohler55/ojg/jp"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
)

// Examples from
// https://ietf-wg-jsonpath.github.io/draft-ietf-jsonpath-base/draft-ietf-jsonpath-base.html

func TestJSONPathQueryExample(t *testing.T) {
	bookstore := map[string]interface{}{
		"store": map[string]interface{}{
			"book": []map[string]interface{}{
				{
					"category": "reference",
					"author":   "Nigel Rees",
					"title":    "Sayings of the Century",
					"price":    8.95,
				},
				{
					"category": "fiction",
					"author":   "Evelyn Waugh",
					"title":    "Sword of Honour",
					"price":    12.99,
				},
				{
					"category": "fiction",
					"author":   "Herman Melville",
					"title":    "Moby Dick",
					"isbn":     "0-553-21311-3",
					"price":    8.99,
				},
				{
					"category": "fiction",
					"author":   "J. R. R. Tolkien",
					"title":    "The Lord of the Rings",
					"isbn":     "0-395-19395-8",
					"price":    22.99,
				},
			},
			"bicycle": map[string]interface{}{
				"color": "red",
				"price": 399,
			},
		},
	}
	type testCase struct {
		skip      string
		desp      string
		obj       interface{}
		expr      string
		expect    []interface{}
		matcherFn func(expect interface{}) types.GomegaMatcher
		err       string
	}
	testCases := []testCase{
		{
			desp: "multi bracket notation",
			obj:  bookstore,
			expr: "$['store']['book'][0]['title']",
			expect: []interface{}{
				"Sayings of the Century",
			},
			err: ``,
		},
		{
			desp: "compact dot notation",
			obj:  bookstore,
			expr: "$.store.book[0].title",
			expect: []interface{}{
				"Sayings of the Century",
			},
			err: ``,
		},
		{
			desp: "logical expr filter",
			obj:  bookstore,
			expr: "$.store.book[?@.price < 10].title",
			expect: []interface{}{
				"Sayings of the Century",
				"Moby Dick",
			},
			err: ``,
		},
		{
			desp: "the authors of all books",
			obj:  bookstore,
			expr: "$.store.book[*].author",
			expect: []interface{}{
				"Nigel Rees",
				"Evelyn Waugh",
				"Herman Melville",
				"J. R. R. Tolkien",
			},
			err: ``,
		},
		{
			desp: "all authors",
			obj:  bookstore,
			expr: "$..author",
			expect: []interface{}{
				"Nigel Rees",
				"Evelyn Waugh",
				"Herman Melville",
				"J. R. R. Tolkien",
			},
			err: ``,
		},
		{
			desp: "all things in store",
			obj:  bookstore,
			expr: "$.store.*",
			expect: []interface{}{
				[]interface{}{
					bookstore["store"].(map[string]interface{})["book"],
					bookstore["store"].(map[string]interface{})["bicycle"],
				},
				[]interface{}{
					bookstore["store"].(map[string]interface{})["bicycle"],
					bookstore["store"].(map[string]interface{})["book"],
				},
			},
			matcherFn: func(expect interface{}) types.GomegaMatcher {
				return BeElementOf(expect)
			},
			err: ``,
		},
		{
			desp: "prices of everything in the store",
			obj:  bookstore,
			expr: "$.store..price",
			expect: []interface{}{
				[]interface{}{
					8.95,
					12.99,
					8.99,
					22.99,
					399,
				},
				[]interface{}{
					399,
					8.95,
					12.99,
					8.99,
					22.99,
				},
			},
			matcherFn: func(expect interface{}) types.GomegaMatcher {
				return BeElementOf(expect)
			},
			err: ``,
		},
		{
			desp: "the third book",
			obj:  bookstore,
			expr: "$..book[2]",
			expect: []interface{}{
				map[string]interface{}{
					"category": "fiction",
					"author":   "Herman Melville",
					"title":    "Moby Dick",
					"isbn":     "0-553-21311-3",
					"price":    8.99,
				},
			},
			err: ``,
		},
		{
			desp: "the third book's author",
			obj:  bookstore,
			expr: "$..book[2].author",
			expect: []interface{}{
				"Herman Melville",
			},
			err: ``,
		},
		{
			desp:   "the third book's publisher",
			obj:    bookstore,
			expr:   "$..book[2].publisher",
			expect: nil,
			err:    ``,
		},
		{
			desp: "the last book in order",
			obj:  bookstore,
			expr: "$..book[-1]",
			expect: []interface{}{
				map[string]interface{}{
					"category": "fiction",
					"author":   "J. R. R. Tolkien",
					"title":    "The Lord of the Rings",
					"isbn":     "0-395-19395-8",
					"price":    22.99,
				},
			},
			err: ``,
		},
		{
			desp: "the first two books",
			obj:  bookstore,
			expr: "$..book[0,1]",
			expect: []interface{}{
				map[string]interface{}{
					"category": "reference",
					"author":   "Nigel Rees",
					"title":    "Sayings of the Century",
					"price":    8.95,
				},
				map[string]interface{}{
					"category": "fiction",
					"author":   "Evelyn Waugh",
					"title":    "Sword of Honour",
					"price":    12.99,
				},
			},
			err: ``,
		},
		{
			desp: "the first two books #2",
			obj:  bookstore,
			expr: "$..book[:2]",
			expect: []interface{}{
				map[string]interface{}{
					"category": "reference",
					"author":   "Nigel Rees",
					"title":    "Sayings of the Century",
					"price":    8.95,
				},
				map[string]interface{}{
					"category": "fiction",
					"author":   "Evelyn Waugh",
					"title":    "Sword of Honour",
					"price":    12.99,
				},
			},
			err: ``,
		},
		{
			skip: `select books without isbn attribute`,
			desp: "all books with isbn",
			obj:  bookstore,
			expr: "$..book[?@.isbn!=null]",
			expect: []interface{}{
				map[string]interface{}{
					"category": "fiction",
					"author":   "Herman Melville",
					"title":    "Moby Dick",
					"isbn":     "0-553-21311-3",
					"price":    8.99,
				},
				map[string]interface{}{
					"category": "fiction",
					"author":   "J. R. R. Tolkien",
					"title":    "The Lord of the Rings",
					"isbn":     "0-395-19395-8",
					"price":    22.99,
				},
			},
			err: ``,
		},
		{
			desp: "all books with isbn #2",
			obj:  bookstore,
			expr: "$..book[?@.isbn]",
			expect: []interface{}{
				map[string]interface{}{
					"category": "fiction",
					"author":   "Herman Melville",
					"title":    "Moby Dick",
					"isbn":     "0-553-21311-3",
					"price":    8.99,
				},
				map[string]interface{}{
					"category": "fiction",
					"author":   "J. R. R. Tolkien",
					"title":    "The Lord of the Rings",
					"isbn":     "0-395-19395-8",
					"price":    22.99,
				},
			},
			err: ``,
		},
		{
			desp: "all books cheaper than 10",
			obj:  bookstore,
			expr: "$..book[?@.price<10]",
			expect: []interface{}{
				map[string]interface{}{
					"category": "reference",
					"author":   "Nigel Rees",
					"title":    "Sayings of the Century",
					"price":    8.95,
				},
				map[string]interface{}{
					"category": "fiction",
					"author":   "Herman Melville",
					"title":    "Moby Dick",
					"isbn":     "0-553-21311-3",
					"price":    8.99,
				},
			},
			err: ``,
		},
		{
			desp: "all member values and array elements",
			obj: map[string]interface{}{
				"k1": "v1",
				"k2": []interface{}{
					"v2",
					"v3",
				},
			},
			expr: "$..*",
			expect: []interface{}{
				[]interface{}{
					"v2",
					"v3",
					"v1",
					[]interface{}{
						"v2",
						"v3",
					},
				},
				[]interface{}{
					"v2",
					"v3",
					[]interface{}{
						"v2",
						"v3",
					},
					"v1",
				},
			},
			matcherFn: func(expect interface{}) types.GomegaMatcher {
				return BeElementOf(expect)
			},
			err: ``,
		},
		{
			desp: "wildcard select",
			obj: map[string]interface{}{
				"a": []interface{}{
					map[string]interface{}{
						"b": 0,
					},
					map[string]interface{}{
						"b": 1,
					},
					map[string]interface{}{
						"c": 2,
					},
				},
			},
			expr: "$.a[*].b",
			expect: []interface{}{
				0,
				1,
			},
			err: ``,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			if tc.skip != `` {
				t.Skipf(tc.skip)
			}

			g := NewWithT(t)
			expr, err := jp.Parse([]byte(tc.expr))
			if tc.err != "" {
				g.Expect(err).To(HaveOccurred())
				g.Expect(err.Error()).To(MatchRegexp(tc.err))
				return
			}
			g.Expect(err).ToNot(HaveOccurred())

			actual := expr.Get(tc.obj)
			matcherFn := Equal
			if tc.matcherFn != nil {
				matcherFn = tc.matcherFn
			}
			g.Expect(actual).To(matcherFn(tc.expect))
		})
	}
}

func TestJSONPathQueryRootIdentifier(t *testing.T) {
	type testCase struct {
		desp      string
		obj       interface{}
		expr      string
		expect    []interface{}
		matcherFn func(expect interface{}) types.GomegaMatcher
		err       string
	}
	testCases := []testCase{
		{
			desp: "root select",
			obj: map[string]interface{}{
				"k": "v",
			},
			expr: "$",
			expect: []interface{}{
				map[string]interface{}{
					"k": "v",
				},
			},
			err: ``,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			g := NewWithT(t)
			expr, err := jp.Parse([]byte(tc.expr))
			if tc.err != "" {
				g.Expect(err).To(HaveOccurred())
				g.Expect(err.Error()).To(MatchRegexp(tc.err))
				return
			}
			g.Expect(err).ToNot(HaveOccurred())

			actual := expr.Get(tc.obj)
			matcherFn := Equal
			if tc.matcherFn != nil {
				matcherFn = tc.matcherFn
			}
			g.Expect(actual).To(matcherFn(tc.expect))
		})
	}
}

func TestJSONPathQueryNameSelector(t *testing.T) {
	nameSelector := map[string]interface{}{
		"o": map[string]interface{}{
			"j j": map[string]interface{}{
				"k.k": 3,
			},
		},
		"'": map[string]interface{}{
			"@": 2,
		},
	}
	type testCase struct {
		desp      string
		obj       interface{}
		expr      string
		expect    []interface{}
		matcherFn func(expect interface{}) types.GomegaMatcher
		err       string
	}
	testCases := []testCase{
		{
			desp: "named value in nested object",
			obj:  nameSelector,
			expr: "$.o['j j']",
			expect: []interface{}{
				map[string]interface{}{
					"k.k": 3,
				},
			},
			err: ``,
		},
		{
			desp: "nesting further down",
			obj:  nameSelector,
			expr: "$.o['j j']['k.k']",
			expect: []interface{}{
				3,
			},
			err: ``,
		},
		{
			desp: "nesting further down with different delimiter",
			obj:  nameSelector,
			expr: `$.o["j j"]["k.k"]`,
			expect: []interface{}{
				3,
			},
			err: ``,
		},
		{
			desp: "unusual member names",
			obj:  nameSelector,
			expr: `$["'"]["@"]`,
			expect: []interface{}{
				2,
			},
			err: ``,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			g := NewWithT(t)
			expr, err := jp.Parse([]byte(tc.expr))
			if tc.err != "" {
				g.Expect(err).To(HaveOccurred())
				g.Expect(err.Error()).To(MatchRegexp(tc.err))
				return
			}
			g.Expect(err).ToNot(HaveOccurred())

			actual := expr.Get(tc.obj)
			matcherFn := Equal
			if tc.matcherFn != nil {
				matcherFn = tc.matcherFn
			}
			g.Expect(actual).To(matcherFn(tc.expect))
		})
	}
}

func TestJSONPathQueryWildcardSelector(t *testing.T) {
	wildcardSelector := map[string]interface{}{
		"o": map[string]interface{}{
			"j": 1,
			"k": 2,
		},
		"a": []interface{}{
			5,
			3,
		},
	}
	type testCase struct {
		skip      string
		desp      string
		obj       interface{}
		expr      string
		expect    []interface{}
		matcherFn func(expect interface{}) types.GomegaMatcher
		err       string
	}
	testCases := []testCase{
		{
			desp: "object values",
			obj:  wildcardSelector,
			expr: `$[*]`,
			expect: []interface{}{
				[]interface{}{
					[]interface{}{
						5,
						3,
					},
					map[string]interface{}{
						"j": 1,
						"k": 2,
					},
				},
				[]interface{}{
					map[string]interface{}{
						"j": 1,
						"k": 2,
					},
					[]interface{}{
						5,
						3,
					},
				},
			},
			matcherFn: func(expect interface{}) types.GomegaMatcher {
				return BeElementOf(expect)
			},
			err: ``,
		},
		{
			desp: "object values #2",
			obj:  wildcardSelector,
			expr: `$.o[*]`,
			expect: []interface{}{
				[]interface{}{
					1,
					2,
				},
				[]interface{}{
					2,
					1,
				},
			},
			matcherFn: func(expect interface{}) types.GomegaMatcher {
				return BeElementOf(expect)
			},
			err: ``,
		},
		{
			skip: `not terminated at 7`,
			desp: "non-deterministic ordering",
			obj:  wildcardSelector,
			expr: `$.o[*,*]`,
			expect: []interface{}{
				[]interface{}{
					1,
					2,
				},
				[]interface{}{
					2,
					1,
				},
			},
			matcherFn: func(expect interface{}) types.GomegaMatcher {
				return BeElementOf(expect)
			},
			err: ``,
		},
		{
			desp: "array member",
			obj:  wildcardSelector,
			expr: `$.a[*]`,
			expect: []interface{}{
				5,
				3,
			},
			err: ``,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			if tc.skip != `` {
				t.Skipf(tc.skip)
			}

			g := NewWithT(t)
			expr, err := jp.Parse([]byte(tc.expr))
			if tc.err != "" {
				g.Expect(err).To(HaveOccurred())
				g.Expect(err.Error()).To(MatchRegexp(tc.err))
				return
			}
			g.Expect(err).ToNot(HaveOccurred())

			actual := expr.Get(tc.obj)
			matcherFn := Equal
			if tc.matcherFn != nil {
				matcherFn = tc.matcherFn
			}
			g.Expect(actual).To(matcherFn(tc.expect))
		})
	}
}

func TestJSONPathQueryIndexSelector(t *testing.T) {
	indexSelector := []interface{}{
		"a",
		"b",
	}
	type testCase struct {
		desp      string
		obj       interface{}
		expr      string
		expect    []interface{}
		matcherFn func(expect interface{}) types.GomegaMatcher
		err       string
	}
	testCases := []testCase{
		{
			desp: "element of array",
			obj:  indexSelector,
			expr: `$[1]`,
			expect: []interface{}{
				"b",
			},
			err: ``,
		},
		{
			desp: "element of array from the end",
			obj:  indexSelector,
			expr: `$[-2]`,
			expect: []interface{}{
				"a",
			},
			err: ``,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			g := NewWithT(t)
			expr, err := jp.Parse([]byte(tc.expr))
			if tc.err != "" {
				g.Expect(err).To(HaveOccurred())
				g.Expect(err.Error()).To(MatchRegexp(tc.err))
				return
			}
			g.Expect(err).ToNot(HaveOccurred())

			actual := expr.Get(tc.obj)
			matcherFn := Equal
			if tc.matcherFn != nil {
				matcherFn = tc.matcherFn
			}
			g.Expect(actual).To(matcherFn(tc.expect))
		})
	}
}

func TestJSONPathQueryArraySliceSelector(t *testing.T) {
	arraySliceSelector := []interface{}{
		"a",
		"b",
		"c",
		"d",
		"e",
		"f",
		"g",
	}
	type testCase struct {
		skip      string
		desp      string
		obj       interface{}
		expr      string
		expect    []interface{}
		matcherFn func(expect interface{}) types.GomegaMatcher
		err       string
	}
	testCases := []testCase{
		{
			desp: "with default step",
			obj:  arraySliceSelector,
			expr: `$[1:3]`,
			expect: []interface{}{
				"b",
				"c",
			},
			err: ``,
		},
		{
			desp: "with no end index",
			obj:  arraySliceSelector,
			expr: `$[5:]`,
			expect: []interface{}{
				"f",
				"g",
			},
			err: ``,
		},
		{
			desp: "with step 2",
			obj:  arraySliceSelector,
			expr: `$[1:5:2]`,
			expect: []interface{}{
				"b",
				"d",
			},
			err: ``,
		},
		{
			desp: "with negative step",
			obj:  arraySliceSelector,
			expr: `$[5:1:-2]`,
			expect: []interface{}{
				"f",
				"d",
			},
			err: ``,
		},
		{
			skip: `default slice start:end in reverse order`,
			desp: "in reverse order",
			obj:  arraySliceSelector,
			expr: `$[::-1]`,
			expect: []interface{}{
				"g",
				"f",
				"e",
				"d",
				"c",
				"b",
				"a",
			},
			err: ``,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			if tc.skip != `` {
				t.Skipf(tc.skip)
			}

			g := NewWithT(t)
			expr, err := jp.Parse([]byte(tc.expr))
			if tc.err != "" {
				g.Expect(err).To(HaveOccurred())
				g.Expect(err.Error()).To(MatchRegexp(tc.err))
				return
			}
			g.Expect(err).ToNot(HaveOccurred())

			actual := expr.Get(tc.obj)
			matcherFn := Equal
			if tc.matcherFn != nil {
				matcherFn = tc.matcherFn
			}
			g.Expect(actual).To(matcherFn(tc.expect))
		})
	}
}

func TestJSONPathQueryFilterSelector(t *testing.T) {
	filterSelector := map[string]interface{}{
		"a": []interface{}{
			3,
			5,
			1,
			2,
			4,
			6,
			map[string]interface{}{
				"b": "j",
			},
			map[string]interface{}{
				"b": "k",
			},
			map[string]interface{}{
				"b": map[string]interface{}{},
			},
			map[string]interface{}{
				"b": "kilo",
			},
		},
		"o": map[string]interface{}{
			"p": 1,
			"q": 2,
			"r": 3,
			"s": 5,
			"t": map[string]interface{}{
				"u": 6,
			},
		},
		"e": "f",
	}
	type testCase struct {
		skip      string
		desp      string
		obj       interface{}
		expr      string
		expect    []interface{}
		matcherFn func(expect interface{}) types.GomegaMatcher
		err       string
	}
	testCases := []testCase{
		{
			desp: "member value comparison",
			obj:  filterSelector,
			expr: `$.a[?@.b == 'kilo']`,
			expect: []interface{}{
				map[string]interface{}{
					"b": "kilo",
				},
			},
			err: ``,
		},
		{
			desp: "member value comparison with enclosing parentheses",
			obj:  filterSelector,
			expr: `$.a[?(@.b == 'kilo')]`,
			expect: []interface{}{
				map[string]interface{}{
					"b": "kilo",
				},
			},
			err: ``,
		},
		{
			desp: "array value comparison",
			obj:  filterSelector,
			expr: `$.a[?@ > 3.5]`,
			expect: []interface{}{
				5,
				4,
				6,
			},
			err: ``,
		},
		{
			desp: "array value existence",
			obj:  filterSelector,
			expr: `$.a[?@.b]`,
			expect: []interface{}{
				map[string]interface{}{
					"b": "j",
				},
				map[string]interface{}{
					"b": "k",
				},
				map[string]interface{}{
					"b": map[string]interface{}{},
				},
				map[string]interface{}{
					"b": "kilo",
				},
			},
			err: ``,
		},
		{
			desp: "existence of non-singular queries",
			obj:  filterSelector,
			expr: `$[?@.*]`,
			expect: []interface{}{
				[]interface{}{
					[]interface{}{
						3,
						5,
						1,
						2,
						4,
						6,
						map[string]interface{}{
							"b": "j",
						},
						map[string]interface{}{
							"b": "k",
						},
						map[string]interface{}{
							"b": map[string]interface{}{},
						},
						map[string]interface{}{
							"b": "kilo",
						},
					},
					map[string]interface{}{
						"p": 1,
						"q": 2,
						"r": 3,
						"s": 5,
						"t": map[string]interface{}{
							"u": 6,
						},
					},
				},
				[]interface{}{
					map[string]interface{}{
						"p": 1,
						"q": 2,
						"r": 3,
						"s": 5,
						"t": map[string]interface{}{
							"u": 6,
						},
					},
					[]interface{}{
						3,
						5,
						1,
						2,
						4,
						6,
						map[string]interface{}{
							"b": "j",
						},
						map[string]interface{}{
							"b": "k",
						},
						map[string]interface{}{
							"b": map[string]interface{}{},
						},
						map[string]interface{}{
							"b": "kilo",
						},
					},
				},
			},
			matcherFn: func(expect interface{}) types.GomegaMatcher {
				return BeElementOf(expect)
			},
			err: ``,
		},
		{
			desp: "nested filters",
			obj:  filterSelector,
			expr: `$[?@[?@.b]]`,
			expect: []interface{}{
				[]interface{}{
					3,
					5,
					1,
					2,
					4,
					6,
					map[string]interface{}{
						"b": "j",
					},
					map[string]interface{}{
						"b": "k",
					},
					map[string]interface{}{
						"b": map[string]interface{}{},
					},
					map[string]interface{}{
						"b": "kilo",
					},
				},
			},
			err: ``,
		},
		{
			skip: `not a valid operation at 9`,
			desp: "non-deterministic ordering",
			obj:  filterSelector,
			expr: `$.o[?@<3, ?@<3]`,
			expect: []interface{}{
				5,
				4,
				6,
			},
			err: ``,
		},
		{
			desp: "array avlue logical or",
			obj:  filterSelector,
			expr: `$.a[?@<2 || @.b == "k"]`,
			expect: []interface{}{
				1,
				map[string]interface{}{
					"b": "k",
				},
			},
			err: ``,
		},
		{
			desp: "array avlue regular expression match",
			obj:  filterSelector,
			expr: `$.a[?match(@.b, "[jk]")]`,
			expect: []interface{}{
				map[string]interface{}{
					"b": "j",
				},
				map[string]interface{}{
					"b": "k",
				},
			},
			err: ``,
		},
		{
			desp: "array avlue regular expression search",
			obj:  filterSelector,
			expr: `$.a[?search(@.b, "[jk]")]`,
			expect: []interface{}{
				map[string]interface{}{
					"b": "j",
				},
				map[string]interface{}{
					"b": "k",
				},
				map[string]interface{}{
					"b": "kilo",
				},
			},
			err: ``,
		},
		{
			desp: "object value logical and",
			obj:  filterSelector,
			expr: `$.o[?@>1 && @<4]`,
			expect: []interface{}{
				[]interface{}{
					2,
					3,
				},
				[]interface{}{
					3,
					2,
				},
			},
			matcherFn: func(expect interface{}) types.GomegaMatcher {
				return BeElementOf(expect)
			},
			err: ``,
		},
		{
			skip: `wrong result (empty)`,
			desp: "object value logical or",
			obj:  filterSelector,
			expr: `$.o[?@.u || @.x]`,
			expect: []interface{}{
				map[string]interface{}{
					"u": 6,
				},
			},
			err: ``,
		},
		{
			desp: "comparison of queries with no values",
			obj:  filterSelector,
			expr: `$.a[?@.b == $.x]`,
			expect: []interface{}{
				3,
				5,
				1,
				2,
				4,
				6,
			},
			err: ``,
		},
		{
			skip: `panic comparing uncomparable type map[string]interface {}`,
			desp: "comparison of primitive and of structured values",
			obj:  filterSelector,
			expr: `$.a[?@ == @]`,
			expect: []interface{}{
				3,
				5,
				1,
				2,
				4,
				6,
				map[string]interface{}{
					"b": "j",
				},
				map[string]interface{}{
					"b": "k",
				},
				map[string]interface{}{
					"b": map[string]interface{}{},
				},
				map[string]interface{}{
					"b": "kilo",
				},
			},
			err: ``,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			if tc.skip != `` {
				t.Skipf(tc.skip)
			}

			g := NewWithT(t)
			expr, err := jp.Parse([]byte(tc.expr))
			if tc.err != "" {
				g.Expect(err).To(HaveOccurred())
				g.Expect(err.Error()).To(MatchRegexp(tc.err))
				return
			}
			g.Expect(err).ToNot(HaveOccurred())

			actual := expr.Get(tc.obj)
			matcherFn := Equal
			if tc.matcherFn != nil {
				matcherFn = tc.matcherFn
			}
			g.Expect(actual).To(matcherFn(tc.expect))
		})
	}
}

func TestJSONPathQueryChildSegment(t *testing.T) {
	childSegment := []interface{}{
		"a",
		"b",
		"c",
		"d",
		"e",
		"f",
		"g",
	}
	type testCase struct {
		skip      string
		desp      string
		obj       interface{}
		expr      string
		expect    []interface{}
		matcherFn func(expect interface{}) types.GomegaMatcher
		err       string
	}
	testCases := []testCase{
		{
			desp: "indices",
			obj:  childSegment,
			expr: `$[0, 3]`,
			expect: []interface{}{
				"a",
				"d",
			},
			err: ``,
		},
		{
			skip: `invalid slice syntax`,
			desp: "slice and indices",
			obj:  childSegment,
			expr: `$[0:2, 5]`,
			expect: []interface{}{
				"a",
				"b",
				"f",
			},
			err: ``,
		},
		{
			desp: "indices",
			obj:  childSegment,
			expr: `$[0, 0]`,
			expect: []interface{}{
				"a",
				"a",
			},
			err: ``,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			if tc.skip != `` {
				t.Skipf(tc.skip)
			}

			g := NewWithT(t)
			expr, err := jp.Parse([]byte(tc.expr))
			if tc.err != "" {
				g.Expect(err).To(HaveOccurred())
				g.Expect(err.Error()).To(MatchRegexp(tc.err))
				return
			}
			g.Expect(err).ToNot(HaveOccurred())

			actual := expr.Get(tc.obj)
			matcherFn := Equal
			if tc.matcherFn != nil {
				matcherFn = tc.matcherFn
			}
			g.Expect(actual).To(matcherFn(tc.expect))
		})
	}
}

func TestJSONPathQueryDescendantSegment(t *testing.T) {
	descendantSegment := map[string]interface{}{
		"o": map[string]interface{}{
			"j": 1,
			"k": 2,
		},
		"a": []interface{}{
			5,
			3,
			[]interface{}{
				map[string]interface{}{
					"j": 4,
				},
				map[string]interface{}{
					"k": 6,
				},
			},
		},
	}
	type testCase struct {
		skip      string
		desp      string
		obj       interface{}
		expr      string
		expect    []interface{}
		matcherFn func(expect interface{}) types.GomegaMatcher
		err       string
	}
	testCases := []testCase{
		{
			desp: "object values",
			obj:  descendantSegment,
			expr: `$..j`,
			expect: []interface{}{
				[]interface{}{
					1,
					4,
				},
				[]interface{}{
					4,
					1,
				},
			},
			matcherFn: func(expect interface{}) types.GomegaMatcher {
				return BeElementOf(expect)
			},
			err: ``,
		},
		{
			desp: "array values",
			obj:  descendantSegment,
			expr: `$..[0]`,
			expect: []interface{}{
				[]interface{}{
					5,
					map[string]interface{}{
						"j": 4,
					},
				},
				[]interface{}{
					map[string]interface{}{
						"j": 4,
					},
					5,
				},
			},
			matcherFn: func(expect interface{}) types.GomegaMatcher {
				return BeElementOf(expect)
			},
			err: ``,
		},
		{
			skip: `too many results for validate`,
			desp: "all values",
			obj:  descendantSegment,
			expr: `$..[*]`,
			expect: []interface{}{
				[]interface{}{
					5,
					map[string]interface{}{
						"j": 4,
					},
				},
				[]interface{}{
					map[string]interface{}{
						"j": 4,
					},
					5,
				},
			},
			matcherFn: func(expect interface{}) types.GomegaMatcher {
				return BeElementOf(expect)
			},
			err: ``,
		},
		{
			desp: "input value is visited",
			obj:  descendantSegment,
			expr: `$..o`,
			expect: []interface{}{
				map[string]interface{}{
					"j": 1,
					"k": 2,
				},
			},
			err: ``,
		},
		{
			skip: `not terminated at 9`,
			desp: "non-deterministic ordering",
			obj:  descendantSegment,
			expr: `$.o..[*, *]`,
			expect: []interface{}{
				[]interface{}{
					1,
					2,
				},
				[]interface{}{
					2,
					1,
				},
			},
			err: ``,
		},
		{
			desp: "multiple segments",
			obj:  descendantSegment,
			expr: `$.a..[0,1]`,
			expect: []interface{}{
				[]interface{}{
					5,
					3,
					map[string]interface{}{
						"j": 4,
					},
					map[string]interface{}{
						"k": 6,
					},
				},
				[]interface{}{
					map[string]interface{}{
						"j": 4,
					},
					map[string]interface{}{
						"k": 6,
					},
					5,
					3,
				},
			},
			matcherFn: func(expect interface{}) types.GomegaMatcher {
				return BeElementOf(expect)
			},
			err: ``,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			if tc.skip != `` {
				t.Skipf(tc.skip)
			}

			g := NewWithT(t)
			expr, err := jp.Parse([]byte(tc.expr))
			if tc.err != "" {
				g.Expect(err).To(HaveOccurred())
				g.Expect(err.Error()).To(MatchRegexp(tc.err))
				return
			}
			g.Expect(err).ToNot(HaveOccurred())

			actual := expr.Get(tc.obj)
			matcherFn := Equal
			if tc.matcherFn != nil {
				matcherFn = tc.matcherFn
			}
			g.Expect(actual).To(matcherFn(tc.expect))
		})
	}
}

func TestJSONPathQuerySemanticsOfNull(t *testing.T) {
	nullSelector := map[string]interface{}{
		"a": nil,
		"b": []interface{}{
			nil,
		},
		"c": []interface{}{
			map[string]interface{}{},
		},
		"null": 1,
	}
	type testCase struct {
		desp   string
		obj    interface{}
		expr   string
		expect []interface{}
		err    string
	}
	testCases := []testCase{
		{
			desp: "object value",
			obj:  nullSelector,
			expr: `$.a`,
			expect: []interface{}{
				nil,
			},
			err: ``,
		},
		{
			desp:   "used as array",
			obj:    nullSelector,
			expr:   `$.a[0]`,
			expect: nil,
			err:    ``,
		},
		{
			desp:   "used as object",
			obj:    nullSelector,
			expr:   `$.a.b`,
			expect: nil,
			err:    ``,
		},
		{
			desp: "array value",
			obj:  nullSelector,
			expr: `$.b[0]`,
			expect: []interface{}{
				nil,
			},
			err: ``,
		},
		{
			desp: "array value #2",
			obj:  nullSelector,
			expr: `$.b[*]`,
			expect: []interface{}{
				nil,
			},
			err: ``,
		},
		{
			desp: "existence",
			obj:  nullSelector,
			expr: `$.b[?@]`,
			expect: []interface{}{
				nil,
			},
			err: ``,
		},
		{
			desp: "comparison",
			obj:  nullSelector,
			expr: `$.b[?@==null]`,
			expect: []interface{}{
				nil,
			},
			err: ``,
		},
		{
			desp:   "comparison of missing value",
			obj:    nullSelector,
			expr:   `$.c[?@.d==null]`,
			expect: nil,
			err:    ``,
		},
		{
			desp: "member name string of null",
			obj:  nullSelector,
			expr: `$.null`,
			expect: []interface{}{
				1,
			},
			err: ``,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			g := NewWithT(t)
			expr, err := jp.Parse([]byte(tc.expr))
			if tc.err != "" {
				g.Expect(err).To(HaveOccurred())
				g.Expect(err.Error()).To(MatchRegexp(tc.err))
				return
			}
			g.Expect(err).ToNot(HaveOccurred())

			actual := expr.Get(tc.obj)
			matcherFn := Equal
			g.Expect(actual).To(matcherFn(tc.expect))
		})
	}
}
