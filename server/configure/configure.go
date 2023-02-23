package configure

import (
	"encoding/json"

	"easy-grpc/core/cnf"
	"server/logger"
)

var defaultConfigure Configure

// DefaultConfigure return default Configure
func DefaultConfigure() *Configure {
	return &defaultConfigure
}

type Configure struct {
	HTTP    HTTP
	Session Session
	DB      DB
	Logger  logger.Options
}

func (c *Configure) String() string {
	if c == nil {
		return "nil"
	}
	b, e := json.MarshalIndent(c, ``, `	`)
	if e != nil {
		return e.Error()
	}
	return string(b)
}

func (c *Configure) Load(filename string) (e error) {
	e = cnf.Load(filename, c)
	if e != nil {
		return
	}

	var formats = []format{
		&c.DB, &c.Session,
	}
	for _, format := range formats {
		e = format.format()
		if e != nil {
			return
		}
	}
	return
}

type format interface {
	format() (e error)
}
