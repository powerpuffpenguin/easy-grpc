package db

import (
	"context"
	"fmt"
	"time"

	grpc_utils "gateway/protocol/utils"

	"github.com/powerpuffpenguin/easy-grpc/core/db/data"
	"github.com/powerpuffpenguin/easy-grpc/core/logger"

	"github.com/powerpuffpenguin/xormcache"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

var defaultModtimeProxy = newModtimeProxy()

func DefaultModtimeProxy() *ModtimeProxy {
	return defaultModtimeProxy
}

type ModtimeProxy struct {
	store xormcache.Store
}

func newModtimeProxy() *ModtimeProxy {
	return &ModtimeProxy{}
}

func (m *ModtimeProxy) getLastModified(ctx context.Context, id int32) time.Time {
	if m.store != nil {
		modtime, e := m.getByStore(ctx, m.store, id)
		if e != nil {
			if ce := logger.Logger.Check(zap.WarnLevel, `get modtime from store error`); ce != nil {
				ce.Write(
					zap.Error(e),
					zap.Int32(`id`, id),
				)
			}
		} else if modtime != nil {
			return time.Unix(modtime.Unix, 0)
		}
	}
	modtime, _ := m.getByDB(ctx, id)
	if modtime != nil {
		return time.Unix(modtime.Unix, 0)
	}
	return time.Time{}
}
func (m *ModtimeProxy) key(id int32) string {
	return fmt.Sprint("cache_modtime_of_", id)
}
func (m *ModtimeProxy) getByStore(ctx context.Context, store xormcache.Store, id int32) (*grpc_utils.Modtime, error) {
	key := m.key(id)
	b, e := store.Get(key)
	if e != nil {
		return nil, e
	} else if len(b) == 0 {
		return nil, nil
	}
	var result grpc_utils.Modtime
	e = proto.Unmarshal(b, &result)
	if e != nil {
		return nil, e
	}
	return &result, nil
}
func (m *ModtimeProxy) getByDB(ctx context.Context, id int32) (*data.Modtime, error) {
	var result data.Modtime
	exists, e := Engine().NoCache().ID(id).Get(&result)
	if e != nil {
		return nil, e
	} else if !exists {
		return nil, nil
	}
	if m.store != nil {
		e := m.putStore(id, &result)
		if e != nil {
			if ce := logger.Logger.Check(zap.WarnLevel, `put modtime to store error`); ce != nil {
				ce.Write(
					zap.Error(e),
					zap.Int32(`id`, id),
				)
			}
		}
	}
	return &result, nil
}
func (m *ModtimeProxy) putStore(id int32, modtime *data.Modtime) (e error) {
	node := &grpc_utils.Modtime{
		ID:   modtime.ID,
		Unix: modtime.Unix,
		ETag: modtime.ETag,
	}
	b, e := proto.Marshal(node)
	if e != nil {
		return
	}
	key := m.key(id)
	e = m.store.Put(key, b)
	return
}

func (m *ModtimeProxy) LastModified(ctx context.Context, id int32) (modtime time.Time) {
	return m.getLastModified(ctx, id)
}

func (m *ModtimeProxy) SetLastModified(id int32, modtime time.Time) (int64, error) {
	unix := modtime.Unix()
	rows, e := Engine().NoCache().ID(id).Where(
		data.ColModtimeUnix+` < ?`, unix,
	).Cols(data.ColModtimeUnix).Update(&data.Modtime{
		Unix: unix,
	})
	if rows > 0 && m.store != nil {
		err := m.store.Del(m.key(id))
		if err != nil {
			if ce := logger.Logger.Check(zap.WarnLevel, `delete modtime from store error`); ce != nil {
				ce.Write(
					zap.Error(err),
					zap.Int32(`id`, id),
				)
			}
		}
	}
	return rows, e
}
func (m *ModtimeProxy) ResetCache(id int32) {
	if m.store != nil {
		err := m.store.Del(m.key(id))
		if err != nil {
			if ce := logger.Logger.Check(zap.WarnLevel, `delete modtime from store error`); ce != nil {
				ce.Write(
					zap.Error(err),
					zap.Int32(`id`, id),
				)
			}
		}
	}
}
