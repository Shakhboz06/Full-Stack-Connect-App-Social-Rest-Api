package store

import (
	"net/http"
	"strconv"
	"strings"
	"time"
)

type PaginatedFeed struct {
	Limit  int      `json:"limit" validate:"gte=1,lte=20"`
	Offset int      `json:"offset" validate:"gte=0"`
	Sort   string   `json:"sort" validate:"oneof=asc desc"`
	Tags   []string `json:"tags" validate:"max=5"`
	Search string   `json:"search" validate:"max=100"`
	Since  string   `json:"since"`
	Until  string   `json:"until"`
}

func (fq PaginatedFeed) Parse(r *http.Request) (PaginatedFeed, error) {
	queryString := r.URL.Query()

	limit := queryString.Get("limit")

	if limit != "" {
		l, err := strconv.Atoi(limit)
		if err != nil {
			return fq, nil
		}

		fq.Limit = l
	}

	offset := queryString.Get("offset")

	if offset != "" {
		off, err := strconv.Atoi(offset)
		if err != nil {
			return fq, nil
		}

		fq.Offset = off
	}

	sort := queryString.Get("sort")
	if sort != "" {
		fq.Sort = sort
	}
	
	tags := queryString.Get("tags")
	if tags != "" {
		fq.Tags = strings.Split(tags, ",")
	}
	
	search := queryString.Get("search")
	if search != "" {
		fq.Search = search
	}
	
	since := queryString.Get("since")
	if since != "" {
		fq.Since = parseTime(since)
	}

	until := queryString.Get("until")
	if until != "" {
		fq.Until = parseTime(until)
	}
	

	return fq, nil
}


func parseTime(s string)string{
	t, err := time.Parse(time.DateTime, s)
	if err != nil{
		return "Error parsing time"
	}

	return t.Format(time.DateTime)
}