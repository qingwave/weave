package repository

import (
	"github.com/qingwave/weave/pkg/database"

	"gorm.io/gorm"
)

func NewRepository(db *gorm.DB, rdb *database.RedisDB) Repository {
	r := &repository{
		user:  newUserRepository(db, rdb),
		group: newGroupRepository(db, rdb),
		post:  newPostRepository(db, rdb),
	}

	r.migrants = getMigrants(
		r.user,
		r.group,
		r.post,
	)

	return r
}

func getMigrants(objs ...interface{}) []Migrant {
	var migrants []Migrant
	for _, obj := range objs {
		if m, ok := obj.(Migrant); ok {
			migrants = append(migrants, m)
		}
	}
	return migrants
}

type repository struct {
	user     UserRepository
	group    GroupRepository
	post     PostRepository
	migrants []Migrant
}

func (r *repository) User() UserRepository {
	return r.user
}

func (r *repository) Group() GroupRepository {
	return r.group
}

func (r *repository) Post() PostRepository {
	return r.post
}

func (r *repository) Migrate() error {
	for _, m := range r.migrants {
		if err := m.Migrate(); err != nil {
			return err
		}
	}
	return nil
}
