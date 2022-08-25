package controller

import (
	"fmt"
	"net/http"

	"github.com/qingwave/weave/pkg/common"
	"github.com/qingwave/weave/pkg/model"
	"github.com/qingwave/weave/pkg/service"
	"github.com/qingwave/weave/pkg/utils/trace"

	"github.com/gin-gonic/gin"
)

type PostController struct {
	postService service.PostService
}

func NewPostController(postService service.PostService) Controller {
	return &PostController{
		postService: postService,
	}
}

// @Summary List post
// @Description List post
// @Produce json
// @Tags post
// @Security JWT
// @Success 200 {object} common.Response{data=[]model.Post}
// @Router /api/v1/posts [get]
func (p *PostController) List(c *gin.Context) {
	common.TraceStep(c, "start list post")
	posts, err := p.postService.List()
	if err != nil {
		common.ResponseFailed(c, http.StatusBadRequest, err)
		return
	}
	common.TraceStep(c, "list post done")
	common.ResponseSuccess(c, posts)
}

// @Summary Get post
// @Description Get post
// @Produce json
// @Tags post
// @Security JWT
// @Param id path int true "post id"
// @Success 200 {object} common.Response{data=model.Post}
// @Router /api/v1/posts/{id} [get]
func (p *PostController) Get(c *gin.Context) {
	post, err := p.postService.Get(c.Param("id"))
	if err != nil {
		common.ResponseFailed(c, http.StatusBadRequest, err)
		return
	}
	common.ResponseSuccess(c, post)
}

// @Summary Create post
// @Description Create post and storage
// @Accept json
// @Produce json
// @Tags post
// @Security JWT
// @Param post body model.Post true "post info"
// @Success 200 {object} common.Response{data=model.Post}
// @Router /api/v1/posts [post]
func (p *PostController) Create(c *gin.Context) {
	user := common.GetUser(c)
	if user == nil {
		common.ResponseFailed(c, http.StatusBadRequest, fmt.Errorf("failed to get user"))
		return
	}

	post := new(model.Post)
	if err := c.BindJSON(post); err != nil {
		common.ResponseFailed(c, http.StatusBadRequest, err)
		return
	}

	common.TraceStep(c, "start create post", trace.Field{"post", post.Name})
	defer common.TraceStep(c, "create post done", trace.Field{"post", post.Name})

	post, err := p.postService.Create(user, post)
	if err != nil {
		common.ResponseFailed(c, http.StatusInternalServerError, err)
		return
	}

	common.ResponseSuccess(c, post)
}

// @Summary Update post
// @Description Update post and storage
// @Accept json
// @Produce json
// @Tags post
// @Security JWT
// @Param post body model.UpdatedUser true "post info"
// @Param id   path      int  true  "post id"
// @Success 200 {object} common.Response{data=model.Post}
// @Router /api/v1/posts/{id} [put]
func (p *PostController) Update(c *gin.Context) {
	user := common.GetUser(c)
	if user == nil {
		common.ResponseFailed(c, http.StatusBadRequest, fmt.Errorf("failed to get user"))
		return
	}

	id := c.Param("id")

	new := new(model.Post)
	if err := c.BindJSON(new); err != nil {
		common.ResponseFailed(c, http.StatusBadRequest, err)
		return
	}

	common.TraceStep(c, "start update post", trace.Field{"post", new.Name})
	defer common.TraceStep(c, "update post done", trace.Field{"post", new.Name})

	post, err := p.postService.Update(id, new)
	if err != nil {
		common.ResponseFailed(c, http.StatusInternalServerError, err)
		return
	}

	common.ResponseSuccess(c, post)
}

// @Summary Delete post
// @Description Delete post
// @Produce json
// @Tags post
// @Security JWT
// @Param id path int true "post id"
// @Success 200 {object} common.Response
// @Router /api/v1/posts/{id} [delete]
func (p *PostController) Delete(c *gin.Context) {
	user := common.GetUser(c)
	if user == nil {
		common.ResponseFailed(c, http.StatusBadRequest, fmt.Errorf("failed to get user"))
		return
	}

	if err := p.postService.Delete(c.Param("id")); err != nil {
		common.ResponseFailed(c, http.StatusBadRequest, err)
		return
	}

	common.ResponseSuccess(c, nil)
}

func (p *PostController) RegisterRoute(api *gin.RouterGroup) {
	api.GET("/posts", p.List)
	api.POST("/posts", p.Create)
	api.GET("/posts/:id", p.Get)
	api.PUT("/posts/:id", p.Update)
	api.DELETE("/posts/:id", p.Delete)
}

func (p *PostController) Name() string {
	return "Post"
}
