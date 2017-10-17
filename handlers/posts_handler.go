package handlers

import (
	"../database/models"
	"../daemon"
	//"../utils"

	"encoding/json"
	"github.com/valyala/fasthttp"
	"github.com/qiangxue/fasthttp-routing"
	"time"
	"strconv"
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
		//if err == utils.UniqueError {
		//	prevForum, err := thread.GetThreadBySlug(daemon.DB.Pool)
		//
		//	if err != nil {
		//		//log.Fatal(err)
		//		daemon.Render.JSON(c.RequestCtx, fasthttp.StatusBadRequest, nil)
		//		return err
		//	}
		//
		//	daemon.Render.JSON(c.RequestCtx, fasthttp.StatusConflict, prevForum)
		//	return nil
		//}
		daemon.Render.JSON(c.RequestCtx, fasthttp.StatusConflict, nil)
		return nil
	}

	daemon.Render.JSON(c.RequestCtx, fasthttp.StatusCreated, posts)
	return nil
}