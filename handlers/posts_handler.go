package handlers

import (
	"../database/models"
	"../daemon"
	"../utils"

	"encoding/json"
	"github.com/valyala/fasthttp"
	"github.com/qiangxue/fasthttp-routing"
	"time"
	"strconv"
	"strings"
)

func CreatePosts(c *routing.Context) error {
	slugOrId := c.Param("slug_or_id")
	posts := []models.Posts{}
	err := json.Unmarshal(c.PostBody(), &posts);
	if err != nil {
		return err
	}

	created := time.Now()

	thread := new(models.Threads)

	resultTread := models.Threads{}
	if id, parseErr := strconv.ParseInt(slugOrId, 10, 64); parseErr == nil {
		thread.TID = id
		resultTread, err = thread.GetThreadById(daemon.DB.Pool);
	} else {
		thread.Slug = slugOrId
		resultTread, err = thread.GetThreadBySlug(daemon.DB.Pool);
	}
	if err != nil {
		daemon.Render.JSON(c.RequestCtx, fasthttp.StatusNotFound, nil)
		return nil
	}

	if err := models.CreatePostsBySlice(daemon.DB.Pool, posts, resultTread.TID, created, resultTread.Forum); err != nil {
		if err == utils.NotFoundError {
			daemon.Render.JSON(c.RequestCtx, fasthttp.StatusNotFound, nil)
			return nil
		}
		daemon.Render.JSON(c.RequestCtx, fasthttp.StatusConflict, nil)
		return nil
	}
	daemon.Render.JSON(c.RequestCtx, fasthttp.StatusCreated, posts)
	return nil
}

func GetPost(c *routing.Context) error {
	type Response struct {
		Post   *models.Posts   `json:"post"`
		Forum  *models.Forums  `json:"forum,omitempty"`
		Author *models.Users   `json:"author,omitempty"`
		Thread *models.Threads `json:"thread,omitempty"`
	}

	response := new(Response)

	stringId := c.Param("id")
	post := new(models.Posts)

	var parseErr error
	post.PID, parseErr = strconv.ParseInt(stringId, 10, 64)
	if parseErr != nil {
		daemon.Render.JSON(c.RequestCtx, fasthttp.StatusBadRequest, nil)
		return nil
	}

	resultPost, err := post.GetPostById(daemon.DB.Pool)
	if err != nil {
		daemon.Render.JSON(c.RequestCtx, fasthttp.StatusNotFound, nil)
		return nil
	}

	response.Post = &resultPost

	related := string(c.QueryArgs().Peek("related"))
	if related != "" {
		splitRelated := strings.Split(related, ",")
		for _, entity := range splitRelated {
			var err error
			switch entity {
			case "forum":
				forum := new(models.Forums)
				forum.Slug = resultPost.Forum

				var resultForum models.Forums
				resultForum, err = forum.GetForumBySlug(daemon.DB.Pool)
				response.Forum = &resultForum
			case "user":
				user := new(models.Users)
				user.Nickname = resultPost.Author

				var resultUser models.Users
				resultUser, err = user.GetUserByLogin(daemon.DB.Pool)
				response.Author = &resultUser
			case "thread":
				thread := new(models.Threads)
				thread.TID = resultPost.Thread

				var resultThread models.Threads
				resultThread, err = thread.GetThreadById(daemon.DB.Pool)
				response.Thread = &resultThread
			}
			if err != nil {
				daemon.Render.JSON(c.RequestCtx, fasthttp.StatusBadRequest, nil)
				return nil
			}
		}
	}
	daemon.Render.JSON(c.RequestCtx, fasthttp.StatusOK, response)
	return nil
}

func UpdatePost(c *routing.Context) error {
	stringId := c.Param("id")

	post := new(models.Posts)
	if err := json.Unmarshal(c.PostBody(), post); err != nil {
		return err
	}

	var parseErr error
	post.PID, parseErr = strconv.ParseInt(stringId, 10, 64)
	if parseErr != nil {
		daemon.Render.JSON(c.RequestCtx, fasthttp.StatusBadRequest, nil)
		return nil
	}

	prevPost, err := post.GetPostById(daemon.DB.Pool)
	if err != nil {
		daemon.Render.JSON(c.RequestCtx, fasthttp.StatusNotFound, nil)
		return nil
	}
	if post.Message == "" || strings.Compare(prevPost.Message, post.Message) == 0 {
		daemon.Render.JSON(c.RequestCtx, fasthttp.StatusOK, prevPost)
		return nil
	}

	if err = post.UpdatePost(daemon.DB.Pool); err != nil {
		if err == utils.UniqueError {
			daemon.Render.JSON(c.RequestCtx, fasthttp.StatusConflict, nil)
			return nil
		}
		daemon.Render.JSON(c.RequestCtx, fasthttp.StatusNotFound, nil)
		return nil
	}

	daemon.Render.JSON(c.RequestCtx, fasthttp.StatusOK, post)
	return nil
}
