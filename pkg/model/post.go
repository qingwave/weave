package model

import "time"

const (
	PostAssociation       = "Posts"
	TagAssociation        = "Tags"
	CategoriesAssociation = "Categories"
)

type Post struct {
	ID         uint       `json:"id" gorm:"autoIncrement;primaryKey"`
	Name       string     `json:"name" gorm:"size:256;not null;unique"`
	Content    string     `json:"content" gorm:"type:text;not null"`
	Summary    string     `json:"summary" gorm:"size:512"`
	CreatorID  uint       `json:"creatorId"`
	Creator    User       `json:"creator" gorm:"foreignKey:CreatorID"`
	Tags       []Tag      `json:"tags"  gorm:"many2many:tag_posts"`
	Categories []Category `json:"categories" gorm:"many2many:category_posts"`
	Comments   []Comment  `json:"comments"`

	Views     uint `json:"views" gorm:"type:uint"`
	Likes     uint `json:"likes" gorm:"-"`
	UserLiked bool `json:"userLiked" gorm:"-"`

	BaseModel
}

type Tag struct {
	ID   uint   `json:"id" gorm:"autoIncrement;primaryKey"`
	Name string `json:"name" gorm:"size:256;not null;unique"`
}

type Category struct {
	ID   uint   `json:"id" gorm:"autoIncrement;primaryKey"`
	Name string `json:"name" gorm:"size:256;not null;unique"`
}

type Like struct {
	ID     uint `json:"id" gorm:"autoIncrement;primaryKey"`
	UserID uint `json:"userId" gorm:"uniqueIndex:user_post"`
	User   User `json:"-" gorm:"foreignKey:UserID"`
	PostID uint `json:"postId" gorm:"uniqueIndex:user_post"`
	Post   Post `json:"-" gorm:"foreignKey:PostID"`
}

type Comment struct {
	ID        uint      `json:"id" gorm:"autoIncrement;primaryKey"`
	ParentID  *uint     `json:"parentId"`
	Parent    *Comment  `json:"parent" gorm:"foreignKey:ParentID"`
	UserID    uint      `json:"userId"`
	User      User      `json:"user" gorm:"foreignKey:UserID"`
	PostID    uint      `json:"postId"`
	Post      Post      `json:"-" gorm:"foreignKey:PostID"`
	Content   string    `json:"content" gorm:"size:1024"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
