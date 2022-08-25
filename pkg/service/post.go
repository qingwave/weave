package service

import (
	"regexp"
	"strconv"

	"github.com/qingwave/weave/pkg/model"
	"github.com/qingwave/weave/pkg/repository"
)

type postService struct {
	postRepository repository.PostRepository
}

func NewPostService(postRepository repository.PostRepository) PostService {
	return &postService{
		postRepository: postRepository,
	}
}

func (p *postService) List() ([]model.Post, error) {
	return p.postRepository.List()
}

func (p *postService) Create(user *model.User, post *model.Post) (*model.Post, error) {
	if len(post.Summary) == 0 {
		post.Summary = getSummary(post.Content)
	}
	return p.postRepository.Create(user, post)
}

func (p *postService) Get(id string) (*model.Post, error) {
	pid, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	return p.postRepository.GetPostByID(uint(pid))
}

func (p *postService) Update(id string, post *model.Post) (*model.Post, error) {
	pid, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	post.ID = uint(pid)
	return p.postRepository.Update(post)
}

func (p *postService) Delete(id string) error {
	pid, err := strconv.Atoi(id)
	if err != nil {
		return err
	}

	return p.postRepository.Delete(uint(pid))
}

func (p *postService) GetTags(id string) ([]model.Tag, error) {
	pid, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	return p.postRepository.GetTags(&model.Post{ID: uint(pid)})
}

func (p *postService) GetCategories(id string) ([]model.Category, error) {
	pid, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	return p.postRepository.GetCategories(&model.Post{ID: uint(pid)})
}

const summaryLen = 128

var summaryRe = regexp.MustCompile(`(#.+\n)|(\n)|(<.+>)|(\(.+\))|(\[.+\])`)

func getSummary(content string) string {
	l := summaryLen
	if len(content) < l {
		l = len(content)
	}

	sum := string([]rune(content)[:l])
	return summaryRe.ReplaceAllString(sum, "")
}
