package db

import (
	"context"
	"fmt"

	"github.com/powerpuffpenguin/easy-grpc/core/cnf"
	"github.com/powerpuffpenguin/easy-grpc/core/logger"
	"go.uber.org/zap"
	"xorm.io/xorm"
	"xorm.io/xorm/schemas"
)

var _Engine *xorm.EngineGroup

// Init Initialize the database
func Init(cnf *cnf.DB) {
	engine, e := xorm.NewEngineGroup(cnf.Driver, cnf.Source)
	if e != nil {
		logger.Logger.Panic(`new db engine group error`,
			zap.Error(e),
			zap.String(`driver`, cnf.Driver),
			zap.Strings(`source`, cnf.Source),
		)
	}
	_Engine = engine
	logger.Logger.Info(`db group`,
		zap.String(`driver`, cnf.Driver),
		zap.Strings(`source`, cnf.Source),
	)
	engine.ShowSQL(cnf.ShowSQL)
	// connect pool
	if cnf.MaxOpen > 1 {
		engine.SetMaxOpenConns(cnf.MaxOpen)
	}
	if cnf.MaxIdle > 1 {
		engine.SetMaxIdleConns(cnf.MaxIdle)
	}
	// cache
	initCache(engine, &cnf.Cache)
}
func Engine() *xorm.EngineGroup {
	return _Engine
}
func Master() *xorm.Engine {
	return _Engine.Master()
}
func Slave() *xorm.Engine {
	return _Engine.Slave()
}
func Session(ctx ...context.Context) *xorm.Session {
	s := _Engine.NewSession()
	for _, c := range ctx {
		s = s.Context(c)
	}
	return s
}
func Begin(ctx ...context.Context) (*xorm.Session, error) {
	s := Session(ctx...)
	e := s.Begin()
	if e != nil {
		s.Close()
		return nil, e
	}
	return s, nil
}
func DeleteCache(tableName string, pk ...interface{}) error {
	if len(pk) == 0 {
		return nil
	}
	cacher := _Engine.GetCacher(tableName)
	if cacher == nil {
		return nil
	}
	id, e := schemas.NewPK(pk...).ToString()
	if e != nil {
		if ce := logger.Logger.Check(zap.WarnLevel, `DeleteCache`); ce != nil {
			ce.Write(
				zap.Error(e),
				zap.String(`table name`, tableName),
				zap.String(`pk`, fmt.Sprint(pk...)),
			)
		}
		return e
	}
	cacher.ClearIds(tableName)
	cacher.DelBean(tableName, id)
	return nil
}
