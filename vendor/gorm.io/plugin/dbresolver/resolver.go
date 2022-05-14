package dbresolver

import "gorm.io/gorm"

type resolver struct {
	sources    []gorm.ConnPool
	replicas   []gorm.ConnPool
	policy     Policy
	dbResolver *DBResolver
}

func (r *resolver) resolve(stmt *gorm.Statement, op Operation) (connPool gorm.ConnPool) {
	if op == Read {
		if len(r.replicas) == 1 {
			connPool = r.replicas[0]
		} else {
			connPool = r.policy.Resolve(r.replicas)
		}
	} else if len(r.sources) == 1 {
		connPool = r.sources[0]
	} else {
		connPool = r.policy.Resolve(r.sources)
	}

	if stmt.DB.PrepareStmt {
		if preparedStmt, ok := r.dbResolver.prepareStmtStore[connPool]; ok {
			return &gorm.PreparedStmtDB{
				ConnPool: connPool,
				Mux:      preparedStmt.Mux,
				Stmts:    preparedStmt.Stmts,
			}
		}
	}

	return
}

func (r *resolver) call(fc func(connPool gorm.ConnPool) error) error {
	for _, s := range r.sources {
		if err := fc(s); err != nil {
			return err
		}
	}

	for _, r := range r.replicas {
		if err := fc(r); err != nil {
			return err
		}
	}
	return nil
}
