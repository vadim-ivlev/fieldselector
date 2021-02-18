package fieldselector

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var query = `
  query CaseFileCounts($filters: FiltersArgument) {
	viewer {
	  repository {
		cases {
		  hits(first: 1, filters: $filters) {
			edges {
			  node {
				case_id
				files {
				  hits(first: 0) {
					total
				  }
				}
				summary {
				  experimental_strategies {
					experimental_strategy
					file_count
				  }
				  data_categories {
					data_category
					file_count
				  }
				}
			  }
			}
		  }
		}
	  }
	}
  }
`

// go test ./pkg/fieldselector
func Test_findBrackets(t *testing.T) {
	var testString = "aaa(bbb){cccddd}eee{fff}ggg"

	i0, i1 := findBrackets("{}", testString)
	assert.Equal(t, 8, i0)
	assert.Equal(t, 15, i1)
	i0, i1 = findBrackets("()", testString)
	assert.Equal(t, 3, i0)
	assert.Equal(t, 7, i1)

	assert.Equal(t, "", insideBrackets("()", ""))
	assert.Equal(t, "123", insideBrackets("()", "(123)4567"))
	assert.Equal(t, "345", insideBrackets("()", "12(345)67"))
	assert.Equal(t, "7", insideBrackets("()", "123456(7)"))
	assert.Equal(t, "", insideBrackets("()", "123456(7"))
	assert.Equal(t, "", insideBrackets("()", "1234567"))

	assert.Equal(t, "", outsideBrackets("()", ""))
	assert.Equal(t, "4567", outsideBrackets("()", "(123)4567"))
	assert.Equal(t, "1267", outsideBrackets("()", "12(345)67"))
	assert.Equal(t, "123456", outsideBrackets("()", "123456(7)"))
	assert.Equal(t, "123456(7", outsideBrackets("()", "123456(7"))
	assert.Equal(t, "1234567", outsideBrackets("()", "1234567"))

	assert.Equal(t, "", selectionAfter("aaa", ""))
	assert.Equal(t, "", selectionAfter("bbb", "zz aaa{12345}6"))
	assert.Equal(t, "12345", selectionAfter("aaa", "zz aaa{12345}6"))
	assert.Equal(t, "12345", selectionAfter("aaa", "zz aaa{12345}"))
	assert.Equal(t, "", selectionAfter("aaa", "zz aaa{1234"))
	assert.Equal(t, "1234", selectionAfter("aaa", "zz aaa(zz) {1234}"))

	assert.Equal(t, []string{"f1", "f2", "f3"}, GetSelectedFields([]string{"aaa"}, "zz aaa(ab){ f1{cyz} \nf2(aa)   f3\n\t  (bb){ff1 ff2(cc) ff3}}"))
	assert.Equal(t, []string{"edges"}, GetSelectedFields([]string{"hits"}, query))
	assert.Equal(t, []string{"case_id", "files", "summary"}, GetSelectedFields([]string{"node"}, query))
	assert.Equal(t, []string{"case_id", "files", "summary"}, GetSelectedFields([]string{"repository", "node"}, query))

	assert.Equal(t, []string{"edges"}, GetSelectedFields([]string{"hits"}, query))
	assert.Equal(t, []string{"edges"}, GetSelectedFields([]string{"cases", "hits"}, query))
	assert.Equal(t, []string{"total"}, GetSelectedFields([]string{"cases", "files", "hits"}, query))
	assert.Equal(t, []string{"total"}, GetSelectedFields([]string{"files", "hits"}, query))

	assert.Equal(t, []string{"case_id", "files", "summary"}, GetSelectedFields([]string{"node"}, query))
	assert.Equal(t, []string{"experimental_strategies", "data_categories"}, GetSelectedFields([]string{"summary"}, query))
	assert.Equal(t, []string{"data_category", "file_count"}, GetSelectedFields([]string{"data_categories"}, query))

}

//go test -benchmem -run=XXXXX -bench=.  ./pkg/fieldselector
func Benchmark_FieldSelector(b *testing.B) {

	for n := 0; n < b.N; n++ {
		GetSelectedFields([]string{"node"}, query)
	}

}
