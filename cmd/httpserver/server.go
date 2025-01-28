package httpserver

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"orderservice/state"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Serve(app *state.State) error {

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.Config.ApplicationPort),
		Handler:      routes(app),
		ErrorLog:     log.New(app.Logger, "", 0),
		IdleTimeout:  5 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 3 * time.Second,
	}

	shutdownError := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit
		app.Logger.PrintInfo("shutting down server", map[string]string{
			"signal": s.String(),
		})

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := srv.Shutdown(ctx)
		if err != nil {
			shutdownError <- err
		}

		app.Logger.PrintInfo("completing background tasks", map[string]string{
			"addr": srv.Addr,
		})
		app.Wg.Wait()
		app.Repository.Close()
		shutdownError <- nil

	}()

	app.Logger.PrintInfo("starting server", map[string]string{
		"addr": srv.Addr,
	})

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownError
	if err != nil {
		return err
	}

	app.Logger.PrintInfo("stopped server", map[string]string{
		"addr": srv.Addr,
	})

	return nil
}
