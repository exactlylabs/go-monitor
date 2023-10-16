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

func WithTransport(transport sentry.Transport) Option {
	return func(hub *sentry.Hub) {
		opts := hub.Client().Options()
		opts.Transport = transport
		c, err := sentry.NewClient(opts)
		if err != nil {
			panic(err)
		}
		hub.BindClient(c)
	}
}

func WithClientOptions(opts ClientOptions) Option {
	return func(hub *sentry.Hub) {
		opts.AttachStacktrace = true
		c, err := sentry.NewClient(opts)
		if err != nil {
			panic(err)
		}
		hub.BindClient(c)
	}
}
