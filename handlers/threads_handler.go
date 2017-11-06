package handlers

import (
	"../database/models"
	"../daemon"
	"../utils"

	"encoding/json"
	"github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"

	"strconv"
	"sync"
)

func CreateThread(c *routing.Context) error {
	slug := c.Param("slug")
	thread := new(models.Threads)
	if err := json.Unmarshal(c.PostBody(), thread); err != nil {
		return err
	}

	waitData := &sync.WaitGroup{}

	var authorNickname string
	var errAuthor error
	waitData.Add(1)

	go func(waitData *sync.WaitGroup) {
		defer waitData.Done()
		author := models.Users{Nickname: thread.Author}
		author, errAuthor = author.GetUserByLogin(daemon.DB.Pool)
		authorNickname = author.Nickname
	}(waitData)

	var forumSlug string
	var errForum error
	waitData.Add(1)

	go func(waitData *sync.WaitGroup) {
		defer waitData.Done()
		forum := models.Forums{Slug: slug}
		forum, errForum = forum.GetForumBySlug(daemon.DB.Pool)
		forumSlug = forum.Slug
	}(waitData)

	waitData.Wait()

	if errAuthor != nil || errForum != nil {
		daemon.Render.JSON(c.RequestCtx, fasthttp.StatusNotFound, nil)
		return nil
	}

	thread.Author = authorNickname
	thread.Forum = forumSlug

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
		daemon.Render.JSON(c.RequestCtx, fasthttp.StatusBadRequest, nil)
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
		posts, err = resultTread.GetPostsWithTreeSort(daemon.DB.Pool, limit, since, desc);

	case "parent_tree":
		posts, err = resultTread.GetPostsWithParentTreeSort(daemon.DB.Pool, limit, since, desc);

	case "flat":
		fallthrough
	default:
		posts, err = resultTread.GetPostsWithFlatSort(daemon.DB.Pool, limit, since, desc);
	}

	if err != nil {
		daemon.Render.JSON(c.RequestCtx, fasthttp.StatusBadRequest, nil)
	}
	daemon.Render.JSON(c.RequestCtx, fasthttp.StatusOK, posts)
	return nil
}

func UpdateThread(c *routing.Context) error {
	slugOrId := c.Param("slug_or_id")

	thread := new(models.Threads)
	if err := json.Unmarshal(c.PostBody(), thread); err != nil {
		return err
	}

	prevThread := models.Threads{}
	var err error
	if id, parseErr := strconv.ParseInt(slugOrId, 10, 64); parseErr == nil {
		thread.TID = id
		prevThread, err = thread.GetThreadById(daemon.DB.Pool);
	} else {
		thread.Slug = slugOrId
		prevThread, err = thread.GetThreadBySlug(daemon.DB.Pool);
	}
	if err != nil {
		daemon.Render.JSON(c.RequestCtx, fasthttp.StatusNotFound, nil)
		return nil
	}

	if utils.CheckEmpty(thread) {
		utils.AdditionObject(thread, &prevThread)

	}

	if err := thread.UpdateThread(daemon.DB.Pool); err != nil {
		daemon.Render.JSON(c.RequestCtx, fasthttp.StatusBadRequest, nil)
		return nil
	}

	daemon.Render.JSON(c.RequestCtx, fasthttp.StatusOK, thread)
	return nil
}
