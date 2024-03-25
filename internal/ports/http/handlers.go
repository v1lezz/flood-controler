package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"task/internal/app"
)

func CheckHandler(a *app.App, ctx context.Context) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// не ограничивал методы, поскольку в тз изначально не сказано об этом ничего
		defer r.Body.Close()
		queries := r.URL.Query()
		sID := queries.Get("userID")
		if sID == "" {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("query userID is empty"))
			return
		}
		ID, err := strconv.ParseInt(sID, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("query userID is not a number: %s", err.Error())))
			return
		}
		flag, err := a.FloodControler.Check(ctx, ID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("error checker: %s", err.Error())))
			return
		}
		if flag {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusTooManyRequests)
		}
		b, err := json.Marshal(map[string]interface{}{
			"flag": flag,
		})
		if err != nil {
			w.Write([]byte(fmt.Sprintf("error marshal ans to json: %s", err.Error())))
			return
		}
		w.Write(b)
	}
}
