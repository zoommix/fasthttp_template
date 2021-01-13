package utils

import (
	"encoding/json"

	http "github.com/valyala/fasthttp"
)

type errorResp map[string]string

var (
	strContentType     = []byte("Content-Type")
	strApplicationJSON = []byte("application/json")
)

// WriteJSON ...
func WriteJSON(ctx *http.RequestCtx, obj interface{}) {
	ctx.Response.Header.SetCanonical(strContentType, strApplicationJSON)

	if err := json.NewEncoder(ctx).Encode(obj); err != nil {
		WriteError(ctx, http.StatusInternalServerError, newRespError(err.Error()))
	}
}

// SetStatus ...
func SetStatus(ctx *http.RequestCtx, code int) {
	ctx.Response.SetStatusCode(code)
}

// WriteError ...
func WriteError(ctx *http.RequestCtx, code int, obj interface{}) {
	ctx.Response.SetStatusCode(code)

	if err := json.NewEncoder(ctx).Encode(obj); err != nil {
		LogError(err.Error())
		panic(err)
	}
}

// RenderNotFoundError ...
func RenderNotFoundError(ctx *http.RequestCtx, msg string) {
	WriteError(ctx, http.StatusNotFound, newRespError(msg))
}

// RenderUnauthorized ...
func RenderUnauthorized(ctx *http.RequestCtx, msg string) {
	WriteError(ctx, http.StatusUnauthorized, newRespError(msg))
}

// RenderInternalError ...
func RenderInternalError(ctx *http.RequestCtx, msg string) {
	WriteError(ctx, http.StatusInternalServerError, newRespError(msg))
}

func newRespError(msg string) errorResp {
	return errorResp{"error": msg}
}

// RenderValidationErrors ...
func RenderValidationErrors(ctx *http.RequestCtx, verrors map[string][]string) {
	verrorsJSON, err := json.MarshalIndent(verrors, "", "  ")
	LogInfo("Validation errors:\n", string(verrorsJSON))

	if err != nil {
		LogError("Unable to marshal validation errors: ", err)
		RenderInternalError(ctx, err.Error())
		return
	}

	WriteError(ctx, http.StatusUnprocessableEntity, verrors)
}
