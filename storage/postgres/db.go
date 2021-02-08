package postgres

import (
	"context"
	"database/sql"
	"errors"
	"sync"
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/libreria/models"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	Host         string        `mapstructure:"host"          default:"localhost"`
	Port         string        `mapstructure:"port"          default:"5432"`
	Name         string        `mapstructure:"name"          default:"libreria"`
	User         string        `mapstructure:"user"          default:"postgres"`
	Password     string        `mapstructure:"password"      default:"12345"`
	MaxRetries   int           `mapstructure:"max_retries"   default:"5"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"  default:"10s"`
	WriteTimeout time.Duration `mapstructure:"write_timeout" default:"10s"`
}

type Storage struct {
	db *pg.DB
}

func New(globalCtx context.Context, wg *sync.WaitGroup, cfg Config) (*Storage, error) {
	db := pg.Connect(&pg.Options{
		Addr:         cfg.Host + ":" + cfg.Port,
		User:         cfg.User,
		Password:     cfg.Password,
		Database:     cfg.Name,
		WriteTimeout: cfg.WriteTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		MaxRetries:   cfg.MaxRetries,
	})
	// Check connection to a database
	err := db.Ping(globalCtx)
	if err != nil {
		return nil, err
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-globalCtx.Done()
		err := db.Close()
		if err != nil {
			log.Errorf("db connection close error: %v", err)
			return
		}
		log.Info("db connection is closed")
	}()
	return &Storage{db: db}, nil
}

func toServiceError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pg.ErrNoRows) {
		return models.ErrNotFound{}
	}
	return err
}
