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
	if err := p.db.Omit("content").Preload("Creator").Preload("Tags").Preload("Categories").Order(clause.OrderByColumn{Column: clause.Column{Name: "created_at"}, Desc: true}).Find(&posts).Error; err != nil {
		return nil, err
	}

	ids := make([]uint, len(posts))
	for i, p := range posts {
		ids[i] = p.ID
	}

	type result struct {
		ID    uint
		Likes uint
	}

	results := []result{}
	if err := p.db.Model(&model.Like{}).Select("post_id as id, count(likes.post_id) as likes").Where("post_id in ?", ids).Group("post_id").Scan(&results).Error; err != nil {
		return nil, err
	}

	resMap := make(map[uint]uint, len(results))
	for _, r := range results {
		resMap[r.ID] = r.Likes
	}

	for i := range posts {
		posts[i].Likes = resMap[posts[i].ID]
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
	if err := p.db.Preload("Creator").Preload("Tags").Preload("Categories").Preload("Comments.User").Preload("Comments").First(post, id).Error; err != nil {
		return nil, err
	}

	likes, err := p.CountLike(id)
	if err != nil {
		return nil, err
	}

	post.Likes = uint(likes)

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
	err := p.db.Model(post).Omit("views", "creator_id").Updates(post).Error
	return post, err
}

func (p *postRepository) Delete(id uint) error {
	return p.db.Delete(&model.Post{}, id).Error
}

func (p *postRepository) IncView(id uint) error {
	post := &model.Post{ID: id}
	return p.db.Model(post).UpdateColumn("views", gorm.Expr("views + 1")).Error
}

func (p *postRepository) AddLike(pid, uid uint) error {
	like := &model.Like{PostID: pid, UserID: uid}
	return p.db.Create(like).Error
}

func (p *postRepository) DelLike(pid, uid uint) error {
	return p.db.Where("post_id = ? and user_id = ?", pid, uid).Delete(&model.Like{}).Error
}

func (p *postRepository) CountLike(id uint) (int64, error) {
	var count int64
	err := p.db.Model(&model.Like{}).Where("post_id = ?", id).Count(&count).Error
	return count, err
}

func (p *postRepository) GetLike(pid, uid uint) (bool, error) {
	var count int64
	err := p.db.Model(&model.Like{}).Where("post_id = ? and user_id = ?", pid, uid).Count(&count).Error
	return count > 0, err
}

func (p *postRepository) GetLikeByUser(uid uint) ([]model.Like, error) {
	likes := make([]model.Like, 0)
	err := p.db.Model(&model.Like{}).Where("user_id = ?", uid).Find(&likes).Error
	return likes, err
}

func (p *postRepository) AddComment(comment *model.Comment) (*model.Comment, error) {
	err := p.db.Create(comment).Error
	return comment, err
}

func (p *postRepository) DelComment(id string) error {
	return p.db.Delete(&model.Comment{}, id).Error
}

func (p *postRepository) ListComment(pid string) ([]model.Comment, error) {
	comments := make([]model.Comment, 0)
	err := p.db.Where("post_id = ?", pid).Find(comments).Error
	return comments, err
}

func (p *postRepository) Migrate() error {
	return p.db.AutoMigrate(&model.Post{}, &model.Like{}, &model.Tag{}, &model.Category{}, &model.Comment{})
}
