package http

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/libreria/server/http/handlers"
	log "github.com/sirupsen/logrus"
)

const version1 = "/v1"

type Config struct {
	Port      int    `mapstructure:"PORT" default:"8080"`
	URLPrefix string `mapstructure:"URL_PREFIX" default:"/api"`
}

type Server struct {
	config Config
	server *http.Server
	oh     *handlers.Book
}

func New(cfg Config, oh *handlers.Book) *Server {
	s := &Server{
		config: cfg,
		oh:     oh,
	}
	// build http server
	httpSrv := &http.Server{Addr: fmt.Sprintf(":%d", cfg.Port)}
	httpSrv.Handler = s.BuildHandler()
	s.server = httpSrv
	return s
}

func (s *Server) Run(globalCtx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Debugf("http server started listening on addr %s", s.server.Addr)
		err := s.server.ListenAndServe()
		log.Infof("http server has stopped: %s", err)
	}()
	go func() {
		<-globalCtx.Done()
		sdCtx, _ := context.WithTimeout(context.Background(), 5*time.Second) // nolint
		err := s.server.Shutdown(sdCtx)
		if err != nil {
			log.Infof("http server shutdown error %s", err)
		}
	}()
}

func (s *Server) BuildHandler() http.Handler {
	var (
		router        = mux.NewRouter()
		serviceRouter = router.PathPrefix(s.config.URLPrefix).Subrouter()
		v1Router      = serviceRouter.PathPrefix(version1).Subrouter()
	)
	// routes
	v1Router.HandleFunc("/books", s.oh.AddBook).Methods(http.MethodPost)
	v1Router.HandleFunc("/books", s.oh.ListBooks).Methods(http.MethodGet)
	v1Router.HandleFunc("/books/{id}", s.oh.GetBook).Methods(http.MethodGet)
	v1Router.HandleFunc("/books/{id}", s.oh.UpdateBook).Methods(http.MethodPut)
	v1Router.HandleFunc("/books/{id}/in", s.oh.CheckinBook).Methods(http.MethodPatch)
	v1Router.HandleFunc("/books/{id}/out", s.oh.CheckoutBook).Methods(http.MethodPatch)
	v1Router.HandleFunc("/books/{id}/rate", s.oh.RateBook).Methods(http.MethodPatch)
	v1Router.HandleFunc("/books/{id}", s.oh.DeleteBook).Methods(http.MethodDelete)
	return router
}
