package main

import (
	"go-project/internal/store"
	"net/http"
)

// getUserFeedHandler godoc
//
//	@Summary		Fetches the user feed
//	@Description	Fetches the user feed
//	@Tags			feed
//	@Accept			json
//	@Produce		json
//	@Param			since	query		string	false	"Since"
//	@Param			until	query		string	false	"Until"
//	@Param			limit	query		int		false	"Limit"
//	@Param			offset	query		int		false	"Offset"
//	@Param			sort	query		string	false	"Sort"
//	@Param			tags	query		string	false	"Tags"
//	@Param			search	query		string	false	"Search"
//	@Success		200		{object}	[]store.PostswithMetadata
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/feed [get]
func (app *application) getUserFeedHandler(w http.ResponseWriter, r *http.Request){

	var feedQuery store.PaginatedFeed

	feedQuery.Limit = 20
	feedQuery.Offset = 0
	feedQuery.Sort = "desc"
	feedQuery.Search = ""
	feedQuery.Tags = []string{}
	feedQuery.Since = "2024-07-16T08:28:28Z"
	feedQuery.Until = "2024-11-26T08:28:28ZZ"

	fq, err := feedQuery.Parse(r)

	if err != nil{
		app.badRequest(w, r, err)
		return
	}

	if err := Validate.Struct(fq); err != nil{
		app.badRequest(w, r, err)
		return
	}

	ctx := r.Context()
	feed, err := app.store.Posts.GetUserFeed(ctx, int64(720), fq)

	if err != nil{
		app.internalServerError(w, r, err)
		return
	}


	if err := app.jsonResponse(w, http.StatusOK, feed); err !=nil{
		app.internalServerError(w, r, err)
	}
}