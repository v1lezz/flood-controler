package http

import (
	"context"
	"net/http"
	"task/internal/app"
)

func AppRouter(a *app.App, ctx context.Context) {
	http.HandleFunc("/check", CheckHandler(a, ctx))
}
