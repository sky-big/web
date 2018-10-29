package http_helpers

import (
	. "common/clog"
	"math"
	"net/http"
	"reflect"
	"strconv"
)

const (
	DefaultPageSize    = 20
	DefaultCurrentPage = 1
)

//data struct for paging request
type Page struct {
	PageSize    int
	CurrentPage int
	TotalPage   int
	TotalNumber int
	Data        []interface{}
}

// parse paging params from query string
var ParsePagingParams = func(req *http.Request) (pageSize int, currentPage int) {
	q := req.URL.Query()
	pageSize, err := strconv.Atoi(q.Get("page_size"))
	if err != nil || pageSize <= 0 {
		pageSize = DefaultPageSize
	}
	currentPage, err = strconv.Atoi(q.Get("current_page"))
	if err != nil || currentPage <= 0 {
		currentPage = DefaultCurrentPage
	}
	return pageSize, currentPage
}

//require slice not nil, pageSize > 0, currentPage > 0
var PagingSlice = func(slice interface{}, pageSize int, currentPage int) *Page {
	kind := reflect.TypeOf(slice).Kind()
	if kind != reflect.Slice {
		Blog.Errorf("Expect Slice but got %s", kind.String())
		return nil
	}
	s := reflect.ValueOf(slice)
	totalNumber := s.Len()

	page := &Page{
		PageSize:    pageSize,
		CurrentPage: currentPage,
		TotalNumber: totalNumber,
		TotalPage:   int(math.Ceil(float64(totalNumber) / float64(pageSize))),
		Data:        []interface{}{},
	}

	var rangeBegin int
	var rangeEnd int
	if currentPage > page.TotalPage || page.TotalPage == 0 {
		rangeBegin, rangeEnd = 0, 0
	} else if currentPage == page.TotalPage {
		rangeBegin, rangeEnd = (currentPage-1)*pageSize, totalNumber
	} else {
		rangeBegin, rangeEnd = (currentPage-1)*pageSize, currentPage*pageSize
	}

	for i := rangeBegin; i < rangeEnd; i++ {
		page.Data = append(page.Data, s.Index(i).Interface())
	}

	return page
}
