package sentry

import "github.com/getsentry/sentry-go"

type Option func(*sentry.Hub)

func WithTags(tags map[string]string) Option {
	return func(hub *sentry.Hub) {
		hub.ConfigureScope(func(scope *sentry.Scope) {
			for key, value := range tags {
				scope.SetTag(key, value)
			}
		})
	}
}

func WithTag(key, value string) Option {
	return func(hub *sentry.Hub) {
		hub.ConfigureScope(func(scope *sentry.Scope) {
			scope.SetTag(key, value)
		})
	}
}
