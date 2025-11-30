package db

import (
	"context"
	"runtime"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DBRequest struct {
	QueryFunc func(*gorm.DB) (any, error)

	Resp chan DBResponse

	Ctx context.Context
}

type DBResponse struct {
	Result any
	Error  error
}
type DBThread struct {
	reqChan  chan DBRequest
	stopOnce chan struct{}
}

func StartDBThread(dbPath string) (*DBThread, error) {
	reqCh := make(chan DBRequest, buffer)
	stopOnce := make(chan struct{})

	go func() {
		runtime.LockOSThread()
		defer runtime.UnlockOSThread()

		db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
		if err != nil {
			panic("failed to connect database")
		}

		for req := range reqCh {
			result, err := req.QueryFunc(db.WithContext(req.Ctx))
			resp := DBResponse{
				Result: result,
				Error:  err,
			}
			req.Resp <- resp
		}

		close(stopOnce)
	}()

	return &DBThread{
		reqChan:  reqCh,
		stopOnce: stopOnce,
	}, nil
}

var globalDBThread *DBThread

func DBThread()
