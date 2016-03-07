package main

import (
	"errors"
	"fmt"

	"net/http"
	"os"

	"golang.org/x/net/context"

	log "github.com/Sirupsen/logrus"
	"github.com/xtracdev/xavi/plugin"
	"github.com/xtracdev/xavi/runner"
)

type RecoveryContext struct {
	LogFn          func(interface{})
	ErrorMessageFn func(interface{}) string
}

var defaultGlobalRecoveryContext = &RecoveryContext{
	LogFn: func(r interface{}) {
		var err error
		switch t := r.(type) {
		case string:
			err = errors.New(t)
		case error:
			err = t
		default:
			err = errors.New("Unknown error")
		}
		fmt.Println("Handled panic: ", err.Error())
	},
	ErrorMessageFn: func(r interface{}) string {
		return "PANIC recovered"
	},
}

type RecoveryWrapper struct{}

func (rw RecoveryWrapper) Wrap(h plugin.ContextHandler) plugin.ContextHandler {
	return plugin.ContextHandlerFunc(func(c context.Context, w http.ResponseWriter, r *http.Request) {
		rc := defaultGlobalRecoveryContext
		if rc == nil {
			rc = defaultGlobalRecoveryContext
		}

		defer func() {
			r := recover()
			if r != nil {
				rc.LogFn(r)
				http.Error(w, rc.ErrorMessageFn(r), http.StatusInternalServerError)
			}
		}()
		h.ServeHTTPContext(c, w, r)
	},
	)
}

func NewRecoveryWrapper() plugin.Wrapper {
	return &RecoveryWrapper{}
}

type PanicWrapper struct{}

func NewPanicWrapper() plugin.Wrapper {
	return &PanicWrapper{}
}

func (pw PanicWrapper) Wrap(h plugin.ContextHandler) plugin.ContextHandler {
	return plugin.ContextHandlerFunc(func(c context.Context, w http.ResponseWriter, r *http.Request) {
		fmt.Println("###################")
		log.Debugf("###################")
		fmt.Fprintf(w, "{PANIC}")
		h.ServeHTTPContext(c, w, r)
		panic("noooooooo!")
	},
	)
}

func registerPlugins() {
	err := plugin.RegisterWrapperFactory("Recovery", NewRecoveryWrapper)
	if err != nil {
		log.Fatal("Error registering recovery plugin factory")
	}
	err = plugin.RegisterWrapperFactory("Panic", NewPanicWrapper)
	if err != nil {
		log.Fatal("Error registering panic plugin factory")
	}
}

func main() {
	runner.Run(os.Args[1:], registerPlugins)
}
