package app

import (
	"context"
	"golang.org/x/sync/errgroup"
	"net/http"
	"time"
)

type FloodControl interface {
	Check(ctx context.Context, userID int64) (bool, error)
}

type App struct {
	HttpSRV        *http.Server
	FloodControler FloodControl
}

func NewApp(httpSRV *http.Server, control FloodControl) (*App, error) {
	return &App{
		HttpSRV:        httpSRV,
		FloodControler: control,
	}, nil
}

func (a *App) Run(eg *errgroup.Group, gCtx context.Context) error {
	eg.Go(func() error {
		if err := a.HttpSRV.ListenAndServe(); err != nil {
			return err
		}
		return nil
	})
	<-gCtx.Done()
	ctxDown, cancelDown := context.WithTimeout(context.Background(), time.Second*10)
	defer cancelDown()
	if err := a.HttpSRV.Shutdown(ctxDown); err != nil {
		return err
	}
	if err := eg.Wait(); err != nil {
		return err
	}
	return nil
}
