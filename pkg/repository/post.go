package repository

import (
	"github.com/qingwave/weave/pkg/database"
	"github.com/qingwave/weave/pkg/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type postRepository struct {
	db  *gorm.DB
	rdb *database.RedisDB
}

func newPostRepository(db *gorm.DB, rdb *database.RedisDB) PostRepository {
	return &postRepository{
		db:  db,
		rdb: rdb,
	}
}

func (p *postRepository) List() ([]model.Post, error) {
	posts := make([]model.Post, 0)
	if err := p.db.Omit("content").Order(clause.OrderByColumn{Column: clause.Column{Name: "created_at"}, Desc: true}).Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}

func (p *postRepository) Create(user *model.User, post *model.Post) (*model.Post, error) {
	post.CreatorID = user.ID
	post.Creator = *user
	err := p.db.Create(post).Error
	return post, err
}

func (p *postRepository) GetTags(post *model.Post) ([]model.Tag, error) {
	tags := make([]model.Tag, 0)
	err := p.db.Model(post).Association(model.TagAssociation).Find(&tags)
	return tags, err
}

func (p *postRepository) GetCategories(post *model.Post) ([]model.Category, error) {
	categories := make([]model.Category, 0)
	err := p.db.Model(post).Association(model.CategoriesAssociation).Find(&categories)
	return categories, err
}

func (p *postRepository) GetPostByID(id uint) (*model.Post, error) {
	post := new(model.Post)
	if err := p.db.Joins("Creator").First(post, id).Error; err != nil {
		return nil, err
	}

	return post, nil
}

func (p *postRepository) GetPostByName(name string) (*model.Post, error) {
	post := new(model.Post)
	if err := p.db.Where("name = ?", name).First(post).Error; err != nil {
		return nil, err
	}

	return post, nil
}

func (p *postRepository) Update(post *model.Post) (*model.Post, error) {
	err := p.db.Model(post).Updates(post).Error
	return post, err
}

func (p *postRepository) Delete(id uint) error {
	return p.db.Delete(&model.Post{}, id).Error
}

func (p *postRepository) Migrate() error {
	return p.db.AutoMigrate(&model.Post{})
}
