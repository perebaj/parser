package api

import (
	"encoding/json"
	"net/http"

	"github.com/birdie-ai/parser/gpt"
	"github.com/go-chi/chi"
)

// RegisterHeaderParserHandler registers the parser handler
func RegisterHeaderParserHandler(r chi.Router, client *gpt.Client) {
	const parserHeaderPath = "/parsers/header/veneza"

	r.Post(parserHeaderPath, func(w http.ResponseWriter, r *http.Request) {
		parser(w, r, client)
	})

	// TODO: remove this endpoint
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		sendJSON(r.Context(), w, http.StatusOK, "ok")
	})
}

// ParserReqBody is the request body for the parser endpoint
type ParserReqBody struct {
	Text string `json:"text"`
}

func parser(w http.ResponseWriter, r *http.Request, client *gpt.Client) {
	ctx := r.Context()
	var ExtractorReqBody ParserReqBody
	err := json.NewDecoder(r.Body).Decode(&ExtractorReqBody)
	if err != nil {
		sendErr(ctx, w, http.StatusBadRequest, err)
	}
	header, err := client.HeaderParser(ExtractorReqBody.Text)
	if err != nil {
		sendErr(ctx, w, http.StatusInternalServerError, err)
	}

	sendJSON(ctx, w, http.StatusOK, header)
}
