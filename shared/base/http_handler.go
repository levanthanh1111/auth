package base

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/schema"
	"github.com/tpp/msf/external-adapter/db"
	"github.com/tpp/msf/shared/context"
	"github.com/tpp/msf/shared/log"
	"github.com/tpp/msf/shared/validator"
	"gorm.io/gorm"
)

type ParseType uint8

const (
	ParseTypeParam         ParseType = 0
	ParseTypePostForm      ParseType = 1
	ParseTypeJSON          ParseType = 2
	ParseTypeMultipartForm ParseType = 3
	ParseTypeNone          ParseType = 255 // temporary util multiple parser implemented
)

var decoder = schema.NewDecoder()

// responseJSON function response as json with ResponseWriter
func responseJSON(w http.ResponseWriter, code int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if data != nil {
		b, _ := json.Marshal(data)
		w.Write(b)
	}
}

type HTTPHandler interface {
	Logger
	// additional method helper for http handler
	Parse(r *http.Request, out any, parseType ParseType) (context.Context, error)
	ResponseSuccess(w http.ResponseWriter, data any, code ...int)
	ResponseBadRequest(w http.ResponseWriter, data ...any)
	ResponseNotFound(w http.ResponseWriter)
	ResponseInternalServerError(w http.ResponseWriter)
	ResponseUnauthorized(w http.ResponseWriter)
	Validate(Validatee) (isValidationError bool, err error)
	Start(ctx context.Context)
	Commit(ctx context.Context)
	Rollback(ctx context.Context)
	QueryOnly(ctx context.Context)
}

type httpHandler struct {
	Logger
	// additional helper for http handler
	db *gorm.DB
	// readonly db
	// dbro *gorm.DB
}

// Parse request parsing with form type
func (h *httpHandler) Parse(r *http.Request, out any, parseType ParseType) (context.Context, error) {
	var err error
	var ctx = context.FromBaseContext(r.Context())

	// TODO: future's multiple parser: for each parse type use bit and operation to specify multiple parse type
	// for eg. pareType = 3, ((paseType >> ParseTypeParam) & 1) == 1, ((partType >> ParseTypePostForm) & 1) == 1
	// then both ParseParam and ParsePostForm action should be executed
	switch parseType {
	case ParseTypePostForm:
		if err = r.ParseForm(); err != nil {
			h.Error(ctx).Err(err).Msg("could not parse form")
			break
		}

		if err = decoder.Decode(out, r.PostForm); err != nil {
			h.Error(ctx).Err(err).Msg("could not decode form data")
		}
	case ParseTypeParam:
		if err = r.ParseForm(); err != nil {
			h.Error(ctx).Err(err).Msg("could not parse form")
			break
		}

		if err = decoder.Decode(out, r.Form); err != nil {
			h.Error(ctx).Err(err).Msg("could not decode form data")
		}
	case ParseTypeJSON:
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			h.Error(ctx).Err(err).Msg("could not read JSON request body")
			break
		}
		defer r.Body.Close()

		if err = json.Unmarshal(body, out); err != nil {
			h.Error(ctx).Err(err).Msg("could not decode JSON request body")
		}
	case ParseTypeMultipartForm:
		// TODO: unsupported yet, implement logics here
	}

	return ctx, err
}

// ResponseSuccess responses status code 200 and json.
func (h *httpHandler) ResponseSuccess(w http.ResponseWriter, data any, code ...int) {
	if len(code) == 0 {
		// status code 200
		responseJSON(w, http.StatusOK, data)
	} else {
		responseJSON(w, code[0], data)
	}
}

//ResponseBadRequest responses status code 400 and json.
func (h *httpHandler) ResponseBadRequest(w http.ResponseWriter, data ...any) {
	// status code 400
	if len(data) == 0 {
		responseJSON(w, http.StatusBadRequest, nil)
		return
	}
	responseJSON(w, http.StatusBadRequest, data[0])
}

//ResponseNotFound responses status code 404.
func (h *httpHandler) ResponseNotFound(w http.ResponseWriter) {
	// status code 404
	w.WriteHeader(http.StatusNotFound)
}

//ResponseInternalServerError responses 500.
func (h *httpHandler) ResponseInternalServerError(w http.ResponseWriter) {
	// status code 500
	w.WriteHeader(http.StatusInternalServerError)
}

// ResponseUStatusUnauthorized response 401.
func (h *httpHandler) ResponseUnauthorized(w http.ResponseWriter) {
	// status code 401
	w.WriteHeader(http.StatusUnauthorized)
}

// Validate validate and pre-process validation error message
func (h *httpHandler) Validate(v Validatee) (isValidationErrs bool, err error) {
	if err = v.IsValid(); err == nil {
		return false, nil
	}

	var ve validator.ValidationErrors
	if !errors.As(err, &ve) {
		return false, err
	}

	var apiErrors = APIErrors{}
	for _, fe := range ve {
		apiErrors = append(apiErrors, &APIError{
			Field:   fe.Field(),
			Message: mapMsgValidationTag[fe.Tag()](fe.Value(), fe.Param()),
		})
	}
	return true, apiErrors
}

// Start tx
func (h *httpHandler) Start(ctx context.Context) {
	tx := h.db.Begin()
	ctx.WithDBTx(tx)
}

// Commit tx
func (h *httpHandler) Commit(ctx context.Context) {
	if tx := context.DBTxFromContext(ctx); tx != nil {
		tx.Commit()
	}
}

// Rollback tx
func (h *httpHandler) Rollback(ctx context.Context) {
	if tx := context.DBTxFromContext(ctx); tx != nil {
		tx.Rollback()
	}
}

// OnlyQuery not modify data, no need to create transation
func (h *httpHandler) QueryOnly(ctx context.Context) {
	// ctx.WithDBTx(h.dbro)
	ctx.WithDBTx(h.db)
}

// NewBaseHTTPHandler
func NewBaseHTTPHandler(handlerName string) HTTPHandler {
	return &httpHandler{
		Logger: newBaseLogger(log.Logger.With().Str("layer", fmt.Sprintf("handler:%s", handlerName)).Logger()),
		db:     db.GetDBInstance(),
	}
}
