package presenters

import (
	"fmt"
	"net/http"

	"github.com/rs/zerolog"
	"gitlab.com/g6834/team17/task-service/internal/domain/models"
	"gitlab.com/g6834/team17/task-service/internal/utils"
)

type Presenters struct {
	logger *zerolog.Logger
}

func New(logger *zerolog.Logger) *Presenters {
	return &Presenters{logger: logger}
}

func (p *Presenters) JSON(w http.ResponseWriter, r *http.Request, v interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := utils.WriteJson(w, v)
	if err != nil {
		p.Error(w, r, err)
	}
}

func (p *Presenters) Error(w http.ResponseWriter, r *http.Request, err error) {
	switch e := err.(type) {
	case models.StatusError:
		p.logger.Error().
			Err(err).
			Str("caller", err.(models.StatusError).Caller).
			//Str("request-id", mw.GetReqID(r.Context())).
			//Str("trace.id", span.SpanContext().TraceID().String()).
			Msg("error.go")

		http.Error(w, fmt.Sprintf("{\"error\": \"%s\"}", err), e.Code)
		return
	default:
		p.logger.Error().
			Err(err).
			//Str("request-id", mw.GetReqID(r.Context())).
			//Str("trace.id", span.SpanContext().TraceID().String()).
			Msg("unhandled error.go")

		http.Error(w, fmt.Sprintf("{\"error\": \"%s\"}", err), http.StatusInternalServerError)
		return
	}
}
