package api

import (
	"encoding/json"
	"net/http"

	"log/slog"

	"github.com/birdie-ai/parser/gpt"
	"github.com/go-chi/chi"
)

// RegisterHeaderParserHandler registers the parser handler
func RegisterHeaderParserHandler(r chi.Router, client *gpt.Client) {
	const parserHeaderPath = "/parser/header"
	const parserItemPath = "/parser/item"
	r.Post(parserHeaderPath, func(w http.ResponseWriter, r *http.Request) {
		header(w, r, client)
	})

	r.Post(parserItemPath, func(w http.ResponseWriter, r *http.Request) {
		item(w, r, client)
	})
}

// ParserReqBody is the request body for the parser endpoint
type ParserReqBody struct {
	Text string `json:"text"`
}

func header(w http.ResponseWriter, r *http.Request, client *gpt.Client) {
	ctx := r.Context()
	var ParserReqBody ParserReqBody
	err := json.NewDecoder(r.Body).Decode(&ParserReqBody)
	if err != nil {
		sendErr(ctx, w, http.StatusBadRequest, err)
	}
	header, err := client.HeaderParser(ParserReqBody.Text)
	if err != nil {
		sendErr(ctx, w, http.StatusInternalServerError, err)
	}
	slog.Info("Header Parsed", "header", header)
	// TODO: save header in the database
	sendJSON(ctx, w, http.StatusOK, header)
}

func item(w http.ResponseWriter, r *http.Request, client *gpt.Client) {
	ctx := r.Context()
	var ParserReqBody ParserReqBody
	err := json.NewDecoder(r.Body).Decode(&ParserReqBody)
	if err != nil {
		sendErr(ctx, w, http.StatusBadRequest, err)
	}
	itens, err := client.ItemParser(ParserReqBody.Text)
	if err != nil {
		sendErr(ctx, w, http.StatusInternalServerError, err)
	}
	slog.Info("Itens Parsed", "itens", itens)
	// TODO: save item in the database
	sendJSON(ctx, w, http.StatusOK, itens)
}
