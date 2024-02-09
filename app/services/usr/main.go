package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"go.uber.org/zap"
)

const (
	AppName = "usr"
)

type UsrApp struct {
	Conf   *Config
	Logger *zap.SugaredLogger
	Pg     *sql.DB
}

type Config struct {
	Debug    bool `required:"true"`
	Logger   *LogConfig
	DbConfig *PgConfig
}

type db_user struct {
	id   uint64
	name string
}

func main() {
	var cfg Config
	err := ParseConfig(AppName, &cfg)
	if err != nil {
		panic(err)
	}
	log, err := InitLogger(AppName, cfg.Logger)
	defer log.Sync() // nolint: errcheck
	if err != nil {
		fmt.Println(err)
		return
	}
	pg, err := InitDbConnection(cfg.DbConfig)
	if err != nil {
		fmt.Println(err)
		return
	}

	app := UsrApp{Conf: &cfg, Logger: log, Pg: pg}
	app.run()
}

func (a *UsrApp) run() {

	a.Logger.Infow("startup", "STATUS", "OK")
	a.Logger.Infow("usr config", "config debug value:", a.Conf.Debug)

	ctx, stop := context.WithCancel(context.Background())
	go waitSignalExit(stop)
	defer stop()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				a.Logger.Infow("startup", "STATUS", "Break process")
				return
			case <-time.After(1 * time.Second):
				a.process()
			}
		}
	}()
	wg.Wait()
	a.stop()
}

func (a *UsrApp) stop() {
	a.Pg.Close()
	a.Logger.Infow("startup", "STATUS", "Graceful shutdown.")
}
func waitSignalExit(cancel func()) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-ch
	cancel()
}

func (a *UsrApp) process() {
	for i := 0; i < 9; i++ {
		fmt.Printf("step: %d\n", i)
		rows, err := a.Pg.Query("select * from users")
		if err != nil {
			a.Logger.Errorw("user", "Get user from db", zap.Error(err))
			return
		}
		defer rows.Close()
		users := []db_user{}

		for rows.Next() {
			p := db_user{}
			err := rows.Scan(&p.id, &p.name)
			if err != nil {
				a.Logger.Errorw("user", "Scan user from db", zap.Error(err))
				continue
			}
			users = append(users, p)
		}
		fmt.Println(users)
		time.Sleep(time.Second)
	}
}
