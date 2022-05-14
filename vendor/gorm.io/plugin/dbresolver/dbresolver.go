package dbresolver

import (
	"errors"
	"sync"

	"gorm.io/gorm"
)

const (
	Write Operation = "write"
	Read  Operation = "read"
)

type DBResolver struct {
	*gorm.DB
	configs          []Config
	resolvers        map[string]*resolver
	global           *resolver
	prepareStmtStore map[gorm.ConnPool]*gorm.PreparedStmtDB
	compileCallbacks []func(gorm.ConnPool) error
}

type Config struct {
	Sources  []gorm.Dialector
	Replicas []gorm.Dialector
	Policy   Policy
	datas    []interface{}
}

func Register(config Config, datas ...interface{}) *DBResolver {
	return (&DBResolver{}).Register(config, datas...)
}

func (dr *DBResolver) Register(config Config, datas ...interface{}) *DBResolver {
	if dr.prepareStmtStore == nil {
		dr.prepareStmtStore = map[gorm.ConnPool]*gorm.PreparedStmtDB{}
	}

	if dr.resolvers == nil {
		dr.resolvers = map[string]*resolver{}
	}

	if config.Policy == nil {
		config.Policy = RandomPolicy{}
	}

	config.datas = datas
	dr.configs = append(dr.configs, config)
	if dr.DB != nil {
		dr.compileConfig(config)
	}
	return dr
}

func (dr *DBResolver) Name() string {
	return "gorm:db_resolver"
}

func (dr *DBResolver) Initialize(db *gorm.DB) error {
	dr.DB = db
	dr.registerCallbacks(db)
	return dr.compile()
}

func (dr *DBResolver) compile() error {
	for _, config := range dr.configs {
		if err := dr.compileConfig(config); err != nil {
			return err
		}
	}
	return nil
}

func (dr *DBResolver) compileConfig(config Config) (err error) {
	var (
		connPool = dr.DB.Config.ConnPool
		r        = resolver{
			dbResolver: dr,
			policy:     config.Policy,
		}
	)

	if preparedStmtDB, ok := connPool.(*gorm.PreparedStmtDB); ok {
		connPool = preparedStmtDB.ConnPool
	}

	if len(config.Sources) == 0 {
		r.sources = []gorm.ConnPool{connPool}
	} else if r.sources, err = dr.convertToConnPool(config.Sources); err != nil {
		return err
	}

	if len(config.Replicas) == 0 {
		r.replicas = r.sources
	} else if r.replicas, err = dr.convertToConnPool(config.Replicas); err != nil {
		return err
	}

	if len(config.datas) > 0 {
		for _, data := range config.datas {
			if t, ok := data.(string); ok {
				dr.resolvers[t] = &r
			} else {
				stmt := &gorm.Statement{DB: dr.DB}
				if err := stmt.Parse(data); err == nil {
					dr.resolvers[stmt.Table] = &r
				} else {
					return err
				}
			}
		}
	} else if dr.global == nil {
		dr.global = &r
	} else {
		return errors.New("conflicted global resolver")
	}

	for _, fc := range dr.compileCallbacks {
		if err = r.call(fc); err != nil {
			return err
		}
	}

	return nil
}

func (dr *DBResolver) convertToConnPool(dialectors []gorm.Dialector) (connPools []gorm.ConnPool, err error) {
	config := *dr.DB.Config
	for _, dialector := range dialectors {
		if db, err := gorm.Open(dialector, &config); err == nil {
			connPool := db.Config.ConnPool
			if preparedStmtDB, ok := connPool.(*gorm.PreparedStmtDB); ok {
				connPool = preparedStmtDB.ConnPool
			}

			dr.prepareStmtStore[connPool] = &gorm.PreparedStmtDB{
				ConnPool:    db.Config.ConnPool,
				Stmts:       map[string]gorm.Stmt{},
				Mux:         &sync.RWMutex{},
				PreparedSQL: make([]string, 0, 100),
			}

			connPools = append(connPools, connPool)
		} else {
			return nil, err
		}
	}

	return connPools, err
}

func (dr *DBResolver) resolve(stmt *gorm.Statement, op Operation) gorm.ConnPool {
	if len(dr.resolvers) > 0 {
		if u, ok := stmt.Clauses[usingName].Expression.(using); ok && u.Use != "" {
			if r, ok := dr.resolvers[u.Use]; ok {
				return r.resolve(stmt, op)
			}
		}

		if stmt.Table != "" {
			if r, ok := dr.resolvers[stmt.Table]; ok {
				return r.resolve(stmt, op)
			}
		}

		if stmt.Schema != nil {
			if r, ok := dr.resolvers[stmt.Schema.Table]; ok {
				return r.resolve(stmt, op)
			}
		}

		if rawSQL := stmt.SQL.String(); rawSQL != "" {
			if r, ok := dr.resolvers[getTableFromRawSQL(rawSQL)]; ok {
				return r.resolve(stmt, op)
			}
		}
	}

	if dr.global != nil {
		return dr.global.resolve(stmt, op)
	}

	return stmt.ConnPool
}
