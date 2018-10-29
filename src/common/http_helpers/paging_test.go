package http_helpers

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func doTestParsePagingParams(t *testing.T, path string, expectPageSize int, expectCurrentPage int) {
	req, err := http.NewRequest("GET", path, nil)
	assert.Nil(t, err)

	pageSize, currentPage := ParsePagingParams(req)

	assert.Equal(t, expectPageSize, pageSize)
	assert.Equal(t, expectCurrentPage, currentPage)
}

func TestParsePagingParams(t *testing.T) {
	doTestParsePagingParams(t, "/dummy", DefaultPageSize, DefaultCurrentPage)
	doTestParsePagingParams(t, "/dummy?page_size=asb&current_page=hhdh", DefaultPageSize, DefaultCurrentPage)
	doTestParsePagingParams(t, "/dummy?page_size=0&current_page=-1", DefaultPageSize, DefaultCurrentPage)
	doTestParsePagingParams(t, "/dummy?page_size=12&current_page=32", 12, 32)
}

func TestPagingSlice(t *testing.T) {
	arr := []string{"a", "b", "c", "d", "e"}

	page := PagingSlice(arr, 2, 4)
	assert.EqualValues(t, &Page{
		PageSize:    2,
		CurrentPage: 4,
		TotalNumber: 5,
		TotalPage:   3,
		Data:        []interface{}{},
	}, page)

	page = PagingSlice(arr, 2, 3)
	assert.EqualValues(t, &Page{
		PageSize:    2,
		CurrentPage: 3,
		TotalNumber: 5,
		TotalPage:   3,
		Data:        []interface{}{"e"},
	}, page)

	page = PagingSlice(arr, 2, 2)
	assert.EqualValues(t, &Page{
		PageSize:    2,
		CurrentPage: 2,
		TotalNumber: 5,
		TotalPage:   3,
		Data:        []interface{}{"c", "d"},
	}, page)

}
