package handlers

import (
	"../database/models"
	"../daemon"
	"../utils"

	"encoding/json"
	"github.com/qiangxue/fasthttp-routing"
	"log"
	"github.com/valyala/fasthttp"

	"strconv"
)

func CreateThread(c *routing.Context) error {
	slug := c.Param("slug")
	thread := new(models.Threads)
	if err := json.Unmarshal(c.PostBody(), thread); err != nil {
		log.Print(err)
		return err
	}

	author := models.Users{Nickname: thread.Author}
	forum := models.Forums{Slug: slug}

	threadAuthor, errAuthor := author.GetUserByLogin(daemon.DB.Pool)
	threadForum, errForum := forum.GetForumBySlug(daemon.DB.Pool)

	if errAuthor != nil || errForum != nil {
		daemon.Render.JSON(c.RequestCtx, fasthttp.StatusNotFound, nil)
		return nil
	}

	thread.Author = threadAuthor.Nickname
	thread.Forum = threadForum.Slug

	if err := thread.CreateThread(daemon.DB.Pool); err != nil {
		if err == utils.UniqueError {
			prevForum, err := thread.GetThreadBySlug(daemon.DB.Pool)

			if err != nil {
				daemon.Render.JSON(c.RequestCtx, fasthttp.StatusBadRequest, nil)
				return err
			}

			daemon.Render.JSON(c.RequestCtx, fasthttp.StatusConflict, prevForum)
			return nil
		}
		daemon.Render.JSON(c.RequestCtx, fasthttp.StatusBadRequest, thread)
		return nil
	}

	daemon.Render.JSON(c.RequestCtx, fasthttp.StatusCreated, thread)
	return nil
}

func GetThread(c *routing.Context) error {
	slug := c.Param("slug")

	limit := string(c.QueryArgs().Peek("limit"))
	since := string(c.QueryArgs().Peek("since"))
	desc := string(c.QueryArgs().Peek("desc"))

	forum := new(models.Forums)
	forum.Slug = slug
	_, err := forum.GetForumBySlug(daemon.DB.Pool);
	if err != nil {
		daemon.Render.JSON(c.RequestCtx, fasthttp.StatusNotFound, nil)
		return nil
	}
	threads, err := forum.GetAllThreads(daemon.DB.Pool, limit, since, desc);
	if err != nil {
		log.Print(err)
		daemon.Render.JSON(c.RequestCtx, fasthttp.StatusBadRequest, nil)
	}
	daemon.Render.JSON(c.RequestCtx, fasthttp.StatusOK, threads)
	return nil
}

func GetThreadInfo(c *routing.Context) error {
	slugOrId := c.Param("slug_or_id")
	thread := new(models.Threads)

	resultTread := models.Threads{}
	var err error
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

	daemon.Render.JSON(c.RequestCtx, fasthttp.StatusOK, resultTread)
	return nil
}

func GetThreadPosts(c *routing.Context) error {
	slugOrId := c.Param("slug_or_id")
	thread := new(models.Threads)

	resultTread := models.Threads{}
	var err error
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

	limit := string(c.QueryArgs().Peek("limit"))
	sort := string(c.QueryArgs().Peek("sort"))
	desc := string(c.QueryArgs().Peek("desc"))
	since := string(c.QueryArgs().Peek("since"))


	var posts []models.Posts
	switch sort {
	case "tree":


	case "parent_tree":


	case "flat":
		fallthrough
	default:
		posts, err = resultTread.GetPostsWithFlatSort(daemon.DB.Pool, limit, since, desc);
	}


	if err != nil {
		log.Print(err)
		daemon.Render.JSON(c.RequestCtx, fasthttp.StatusBadRequest, nil)
	}
	daemon.Render.JSON(c.RequestCtx, fasthttp.StatusOK, posts)
	return nil
}
