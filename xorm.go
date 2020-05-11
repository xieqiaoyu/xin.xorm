package sql

import (
	"github.com/xieqiaoyu/xin"
	"github.com/xieqiaoyu/xin/db/sql"
	"xorm.io/xorm"
)

//XormConfig config support xorm setting
type XormConfig interface {
	sql.Config
	EnableDbLog() bool
}

func newXormEngineHandler(config XormConfig) sql.GenEngineFunc {
	return func(driverName, dataSourceName string) (engine interface{}, err error) {
		e, err := xorm.NewEngine(driverName, dataSourceName)
		if err != nil {
			return nil, xin.NewWrapEf("Fail to new xorm database engine driver [%s] source [%s], Err:%w", driverName, dataSourceName, err)

		}

		logEnable := config.EnableDbLog()
		if logEnable {
			e.ShowSQL(true)
			//engine.Logger().SetLevel(core.LOG_DEBUG)
		}
		return e, nil
	}
}

func closeXormEngine(engine interface{}) error {
	e, ok := engine.(*xorm.Engine)
	if !ok {
		return xin.NewWrapEf("engine is not a *xorm.Engine")
	}
	return e.Close()
}

//XormService xorm engine service
type XormService struct {
	*sql.Service
	config XormConfig
}

//NewXormService NewXormService
func NewXormService(config XormConfig) *XormService {
	return &XormService{
		config:  config,
		Service: sql.NewService(config, newXormEngineHandler(config), closeXormEngine),
	}
}

//Engine load an xorm engine by id
func (s *XormService) Engine(id string) (engine *xorm.Engine, err error) {
	e, err := s.Get(id)
	if err != nil {
		return nil, err
	}
	engine, ok := e.(*xorm.Engine)
	if !ok {
		return nil, xin.NewWrapEf("db id %s is not a *xorm.Engine", id)
	}
	return engine, nil
}

//Session  load session by id if giving inf is nil ,if isNew is true caller  should close session after everything is done
func (s *XormService) Session(id string) (session *xorm.Session, err error) {
	engine, err := s.Engine(id)
	if err != nil {
		return nil, err
	}
	return engine.NewSession(), nil

}
