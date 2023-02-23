package db

import (
	"context"
	"fmt"

	"reflect"
	"strings"

	"github.com/powerpuffpenguin/easy-grpc/core/db/data"

	"github.com/powerpuffpenguin/easy-grpc/core/logger"

	"github.com/powerpuffpenguin/easy-grpc/core/cnf"

	"github.com/go-redis/redis/v8"
	"github.com/powerpuffpenguin/xormcache"
	store_redis "github.com/powerpuffpenguin/xormcache/redis"
	"go.uber.org/zap"
	"xorm.io/xorm"
)

var defaultCoder = xormcache.NewBinaryCoder()

func DefaultCoder() *xormcache.BinaryCoder {
	return defaultCoder
}

type storeManager struct {
	client map[string]store_redis.Redis
	store  map[string]xormcache.Store
}

func newStoreManager() *storeManager {
	return &storeManager{
		client: make(map[string]store_redis.Redis),
		store:  make(map[string]xormcache.Store),
	}
}
func (s *storeManager) getRedis(url string) (client store_redis.Redis, e error) {
	client, ok := s.client[url]
	if !ok {
		opt, e := redis.ParseURL(url)
		if e != nil {
			return nil, e
		}
		c := redis.NewClient(opt)
		client = store_redis.NewMerge(context.Background(), c, 100)
		s.client[url] = client
	}
	return client, nil
}
func (s *storeManager) Create(cnf *cnf.Redis) error {
	client, e := s.getRedis(cnf.Write)
	if e != nil {
		return e
	}
	key := cnf.Key()
	if _, ok := s.store[key]; ok {
		return nil
	}

	opts := []store_redis.Option{
		store_redis.WithTimeout(cnf.Timeout),
	}
	if strings.HasPrefix(cnf.Read, `redis://`) {
		r, e := s.getRedis(cnf.Read)
		if e != nil {
			return e
		}
		opts = append(opts, store_redis.WithRead(r))
	}
	store, e := store_redis.New(client, opts...)
	if e != nil {
		return e
	}
	s.store[key] = store
	return nil
}
func (s *storeManager) Store(key string) xormcache.Store {
	result, ok := s.store[key]
	if !ok {
		panic(fmt.Sprint("store key not exists: ", key))
	}
	return result
}

func initCache(engine *xorm.EngineGroup, cnf *cnf.DBCache) {
	// redis
	m := newStoreManager()
	if cnf.Redis.IsValid() {
		e := m.Create(&cnf.Redis)
		if e != nil {
			logger.Logger.Fatal(`redis store of db cache error`,
				zap.Error(e),
				zap.String(`write`, cnf.Redis.Write),
				zap.String(`read`, cnf.Redis.Read),
				zap.String(`timeout`, cnf.Redis.Timeout.String()),
			)
		}
	}
	if cnf.Modtime.IsValid() {
		e := m.Create(&cnf.Modtime)
		if e != nil {
			logger.Logger.Fatal(`redis store of db cache error`,
				zap.Error(e),
				zap.String(`table`, data.TableModtimeName),
				zap.String(`write`, cnf.Redis.Write),
				zap.String(`read`, cnf.Redis.Read),
				zap.String(`timeout`, cnf.Redis.Timeout.String()),
			)
		}
	}
	for _, special := range cnf.Special {
		if special.Name == `` || !special.Redis.IsValid() {
			continue
		}
		e := m.Create(&special.Redis)
		if e != nil {
			logger.Logger.Fatal(`redis store of db cache error`,
				zap.Error(e),
				zap.String(`table`, special.Name),
				zap.String(`write`, special.Redis.Write),
				zap.String(`read`, special.Redis.Read),
				zap.String(`timeout`, special.Redis.Timeout.String()),
			)
		}
	}
	// cache
	opt := xormcache.WithCoder(defaultCoder)
	if cnf.Redis.IsValid() {
		cacher, e := xormcache.New(m.Store(cnf.Redis.Key()), opt)
		if e != nil {
			logger.Logger.Fatal(`db default cacher error`,
				zap.Error(e),
				zap.String(`write`, cnf.Redis.Write),
				zap.String(`read`, cnf.Redis.Read),
				zap.String(`timeout`, cnf.Redis.Timeout.String()),
			)
		}
		engine.SetDefaultCacher(cacher)
		if ce := logger.Logger.Check(zap.InfoLevel, `db default cacher`); ce != nil {
			ce.Write(
				zap.String(`write`, cnf.Redis.Write),
				zap.String(`read`, cnf.Redis.Read),
				zap.String(`timeout`, cnf.Redis.Timeout.String()),
			)
		}
		for _, name := range cnf.Direct {
			if name == `` {
				continue
			}
			engine.SetCacher(name, nil)
			if ce := logger.Logger.Check(zap.InfoLevel, `db disable cacher`); ce != nil {
				ce.Write(
					zap.String(`table`, name),
				)
			}
		}
	}
	if cnf.Modtime.IsValid() {
		store := m.Store(cnf.Modtime.Key())
		defaultModtimeProxy.store = store
		if ce := logger.Logger.Check(zap.InfoLevel, `db set cacher`); ce != nil {
			ce.Write(
				zap.String(`table`, data.TableModtimeName),
				zap.String(`write`, cnf.Modtime.Write),
				zap.String(`read`, cnf.Modtime.Read),
				zap.String(`timeout`, cnf.Modtime.Timeout.String()),
			)
		}
	}
	for _, special := range cnf.Special {
		if special.Name == `` || !special.Redis.IsValid() {
			continue
		}
		cacher, e := xormcache.New(m.Store(special.Redis.Key()), opt)
		if e != nil {
			logger.Logger.Fatal(`db set cacher error`,
				zap.Error(e),
				zap.String(`table`, special.Name),
				zap.String(`write`, special.Redis.Write),
				zap.String(`read`, special.Redis.Read),
				zap.String(`timeout`, special.Redis.Timeout.String()),
			)
		}
		engine.SetDefaultCacher(cacher)
		if ce := logger.Logger.Check(zap.InfoLevel, `db set cacher`); ce != nil {
			ce.Write(
				zap.String(`table`, special.Name),
				zap.String(`write`, special.Redis.Write),
				zap.String(`read`, special.Redis.Read),
				zap.String(`timeout`, special.Redis.Timeout.String()),
			)
		}
	}
	engine.SetCacher(data.TableModtimeName, nil)
}
func DelCacheBean(tableName string, v interface{}) {
	delCacheBean(false, tableName, v)
}
func AsyncDelCacheBean(tableName string, v interface{}) {
	delCacheBean(true, tableName, v)
}
func getBeanID(tableName string, bean interface{}) (id string, e error) {
	engine := Engine()
	t, e := engine.TableInfo(bean)
	if e != nil {
		if ce := logger.Logger.Check(zap.WarnLevel, `get pk error`); ce != nil {
			ce.Write(
				zap.Error(e),
				zap.String(`table`, tableName),
				zap.String(`value`, fmt.Sprint(bean)),
			)
		}
		return
	}
	pk, e := t.IDOfV(reflect.ValueOf(bean))
	if e != nil {
		if ce := logger.Logger.Check(zap.WarnLevel, `get pk error`); ce != nil {
			ce.Write(
				zap.Error(e),
				zap.String(`table`, tableName),
				zap.String(`value`, fmt.Sprint(bean)),
			)
		}
		return
	}
	id, e = pk.ToString()
	if e != nil {
		if ce := logger.Logger.Check(zap.WarnLevel, `get pk string error`); ce != nil {
			ce.Write(
				zap.Error(e),
				zap.String(`table`, tableName),
				zap.String(`value`, fmt.Sprint(bean)),
			)
		}
		return
	}
	return
}
func delCacheBean(async bool, tableName string, bean interface{}) {
	engine := Engine()
	cacher := engine.GetCacher(tableName)
	if cacher == nil {
		return
	}
	id, e := getBeanID(tableName, bean)
	if e == nil {
		if async {
			go cacher.DelBean(tableName, id)
		} else {
			cacher.DelBean(tableName, id)
		}
	} else {
		if async {
			go cacher.ClearBeans(tableName)
		} else {
			cacher.ClearBeans(tableName)
		}
	}
}
