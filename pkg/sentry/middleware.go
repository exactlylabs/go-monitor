package sentry

import (
	"log"
	"net/http"
	"runtime/debug"

	"github.com/exactlylabs/go-rest/pkg/restapi/apierrors"
	"github.com/exactlylabs/go-rest/pkg/restapi/webcontext"
	"github.com/getsentry/sentry-go"
)

type sentryMiddleware struct {
	handler http.Handler
}

func (sm *sentryMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			ctx := webcontext.New()
			ctx = ctx.PrepareRequest(w, r)
			ctx.Reject(http.StatusInternalServerError, &apierrors.InternalAPIError)
			log.Println(err)
			stack := string(debug.Stack())
			log.Println(stack)
			hub := sentry.CurrentHub().Clone()
			hub.ConfigureScope(func(scope *sentry.Scope) {
				scope.SetRequest(ctx.Request)
			})
			hub.Recover(err)
		}
	}()
	sm.handler.ServeHTTP(w, r)
}

func SentryMiddleware(handler http.Handler) http.Handler {
	return &sentryMiddleware{handler}
}
