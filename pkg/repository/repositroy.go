package repository

import (
	"context"

	"github.com/qingwave/weave/pkg/database"
	"github.com/qingwave/weave/pkg/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func NewRepository(db *gorm.DB, rdb *database.RedisDB) Repository {
	r := &repository{
		db:    db,
		rdb:   rdb,
		user:  newUserRepository(db, rdb),
		group: newGroupRepository(db, rdb),
		post:  newPostRepository(db, rdb),
		rbac:  newRBACRepository(db, rdb),
	}

	r.migrants = getMigrants(
		r.user,
		r.group,
		r.post,
		r.rbac,
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
	rbac     RBACRepository
	db       *gorm.DB
	rdb      *database.RedisDB
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

func (r *repository) RBAC() RBACRepository {
	return r.rbac
}

func (r *repository) Close() error {
	db, _ := r.db.DB()
	if db != nil {
		if err := db.Close(); err != nil {
			return err
		}
	}

	if r.rdb != nil {
		if err := r.rdb.Close(); err != nil {
			return err
		}
	}

	return nil
}

func (r *repository) Ping(ctx context.Context) error {
	db, err := r.db.DB()
	if err != nil {
		return err
	}
	if err = db.PingContext(ctx); err != nil {
		return err
	}

	if r.rdb == nil {
		return nil
	}
	if _, err := r.rdb.Ping(ctx).Result(); err != nil {
		return err
	}

	return nil
}

func (r *repository) Migrate() error {
	for _, m := range r.migrants {
		if err := m.Migrate(); err != nil {
			return err
		}
	}
	return nil
}

func (r *repository) Init() error {
	resources := []model.Resource{
		{
			Name:  model.ContainerResource,
			Scope: model.ClusterScope,
		},
		{
			Name:  model.ContainerResource + "/log",
			Scope: model.ClusterScope,
		},
		{
			Name:  model.ContainerResource + "/exec",
			Scope: model.ClusterScope,
		},
		{
			Name:  model.ContainerResource + "/proxy",
			Scope: model.ClusterScope,
		},
		{
			Name:  model.PostResource,
			Scope: model.ClusterScope,
		},
		{
			Name:  model.GroupResource,
			Scope: model.ClusterScope,
		},
		{
			Name:  model.UserResource,
			Scope: model.ClusterScope,
		},
		{
			Name:  model.RoleResource,
			Scope: model.ClusterScope,
		},
		{
			Name:  model.AuthResource,
			Scope: model.ClusterScope,
		},
		{
			Name:  model.NamespaceResource,
			Scope: model.ClusterScope,
		},
		{
			Name:  model.KubeDeployment,
			Scope: model.NamespaceScope,
		},
		{
			Name:  model.KubeStatefulset,
			Scope: model.NamespaceScope,
		},
		{
			Name:  model.KubeDaemonset,
			Scope: model.NamespaceScope,
		},
		{
			Name:  model.KubePod,
			Scope: model.NamespaceScope,
		},
		{
			Name:  model.KubeService,
			Scope: model.NamespaceScope,
		},
		{
			Name:  model.KubeIngress,
			Scope: model.NamespaceScope,
		},
	}

	if err := r.RBAC().CreateResources(resources, clause.OnConflict{DoNothing: true}); err != nil {
		return err
	}

	// create default group
	groups := []model.Group{
		{
			Name:     model.RootGroup,
			Kind:     model.SystemGroup,
			Describe: "system root group",
		},
		{
			Name:     model.AuthenticatedGroup,
			Kind:     model.SystemGroup,
			Describe: "system group contains all authenticated user",
		},
		{
			Name:     model.UnAuthenticatedGroup,
			Kind:     model.SystemGroup,
			Describe: "system group contains all unauthenticated user",
		},
	}
	if err := r.Group().CreateGroups(groups, clause.OnConflict{DoNothing: true}); err != nil {
		return err
	}

	return nil
}
