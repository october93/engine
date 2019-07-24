package gql

import (
	"context"
	"net/http"
	"time"

	"github.com/october93/engine/dataloader"
)

type ctxKeyType struct{ name string }

var ctxKey = ctxKeyType{"context"}

func (g *GraphQL) LoaderMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ldrs := dataloader.NewLoaders(g.Store, 350*time.Microsecond)
		dlCtx := context.WithValue(r.Context(), ctxKey, ldrs)
		next.ServeHTTP(w, r.WithContext(dlCtx))
	})
}

func ctxLoaders(ctx context.Context) dataloader.Loaders {
	return ctx.Value(ctxKey).(dataloader.Loaders)
}
