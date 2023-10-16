package sentry

import (
	"reflect"
	"time"

	"github.com/exactlylabs/go-errors/pkg/errors"
	"github.com/getsentry/sentry-go"
)

var notifiedErrors map[error]struct{}

func init() {
	notifiedErrors = make(map[error]struct{})
}

type Context = sentry.Context

// Typeable interface is used to define a type for the error, to show in Sentry instead of the usual fmt.Wrap or errors.withStack
type Typeable interface {
	Type() string
}

// NotifyIfPanic recovers from the panic and sends the error to Sentry, panicking shortly thereafter
func NotifyIfPanic() {
	// Clone the current hub so that modifications of the scope are visible only
	// within this function.
	hub := sentry.CurrentHub().Clone()
	var err error
	errors.RecoverPanic(recover(), &err)
	if err != nil {
		CaptureException(hub, err)

		hub.Flush(2 * time.Second)

		// Raise the panic back again
		panic(err)
	}
}

// NotifyError to sentry
func NotifyError(err error, contexts map[string]sentry.Context) {
	hub := sentry.CurrentHub().Clone()
	metadata := errors.GetMetadata(err)
	if metadata != nil {
		contexts["Error Metadata"] = *metadata
	}
	hub.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetContexts(contexts)
	})
	hub.CaptureException(err)
}

// NotifyErrorOnce will send an error if it has never occurred in this current process
func NotifyErrorOnce(err error, contexts map[string]sentry.Context) {
	if _, exists := notifiedErrors[err]; !exists {
		NotifyError(err, contexts)
		notifiedErrors[err] = struct{}{}
	}
}

func createEventFromError(hub *sentry.Hub, err error) *sentry.Event {
	// create a custom event
	// Based on https://github.com/getsentry/sentry-go/blob/85b380d192353dc9ca3df14fc4f8fa727a33cb2c/client.go
	// Additions: TypeableError interface
	topLevelErr := err
	event := sentry.NewEvent()
	event.Level = sentry.LevelFatal
	for i := 0; i < hub.Client().Options().MaxErrorDepth && err != nil; i++ {
		exc := sentry.Exception{
			Value:      err.Error(),
			Type:       reflect.TypeOf(err).String(),
			Stacktrace: sentry.ExtractStacktrace(err),
		}
		if tErr, ok := err.(Typeable); ok {
			exc.Type = tErr.Type()
		}
		event.Exception = append(event.Exception, exc)
		switch previous := err.(type) {
		case interface{ Unwrap() error }:
			err = previous.Unwrap()
		case interface{ Cause() error }:
			err = previous.Cause()
		default:
			err = nil
		}
	}
	// Add a trace of the current stack to the most recent error in a chain if
	// it doesn't have a stack trace yet.
	// We only add to the most recent error to avoid duplication and because the
	// current stack is most likely unrelated to errors deeper in the chain.
	if event.Exception[0].Stacktrace == nil {
		event.Exception[0].Stacktrace = sentry.NewStacktrace()
	}

	meta := errors.GetMetadata(topLevelErr)
	if meta != nil {
		if event.Contexts == nil {
			event.Contexts = make(map[string]sentry.Context)
		}
		event.Contexts["Error Metadata"] = *meta
	}

	// event.Exception should be sorted such that the most recent error is last.
	reverse(event.Exception)
	return event
}

func reverse(a []sentry.Exception) {
	// https://github.com/getsentry/sentry-go/blob/85b380d192353dc9ca3df14fc4f8fa727a33cb2c/client.go
	for i := len(a)/2 - 1; i >= 0; i-- {
		opp := len(a) - 1 - i
		a[i], a[opp] = a[opp], a[i]
	}
}

func CaptureException(hub *sentry.Hub, err error) {
	event := createEventFromError(hub, err)
	hub.CaptureEvent(event)
}

// Setup sentry global scope, and add some initial context tags: environment and application name
func Setup(dsn string, release, environment, name string, opts ...Option) {
	err := sentry.Init(sentry.ClientOptions{
		Dsn:              dsn,
		Release:          release,
		AttachStacktrace: true,
	})
	if err != nil {
		panic(err)
	}
	sentry.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetTag("environment", environment)
		scope.SetTag("application", name)
	})
	for _, opt := range opts {
		opt(sentry.CurrentHub())
	}
}
