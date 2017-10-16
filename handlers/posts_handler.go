package handlers

import (
	"../database/models"
	"../daemon"
	//"../utils"

	"encoding/json"
	"github.com/valyala/fasthttp"
	"github.com/qiangxue/fasthttp-routing"
	"log"
	"time"
	"strconv"
)

func CreatePosts(c *routing.Context) error {
	slugOrId := c.Param("slug_or_id")
	posts := []models.Posts{}
	err := json.Unmarshal(c.PostBody(), &posts);
	if err != nil {
		log.Print(err)
		return err
	}
	created := time.Now()

	thread := new(models.Threads)
	log.Print(slugOrId)

	resultTread := models.Threads{}
	if id, err := strconv.ParseInt(slugOrId, 10, 64); err == nil {
		thread.TID = id
		log.Print(id)
		resultTread, err = thread.GetThreadById(daemon.DB.Pool);
	} else {
		thread.Slug = slugOrId
		resultTread, err = thread.GetThreadBySlug(daemon.DB.Pool);
	}
	if err != nil {
		daemon.Render.JSON(c.RequestCtx, fasthttp.StatusNotFound, nil)
		return nil
	}

	log.Print(resultTread.Forum)
	log.Print(resultTread.TID)

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
		log.Print(err)
		daemon.Render.JSON(c.RequestCtx, fasthttp.StatusBadRequest, nil)
		return nil
	}

	daemon.Render.JSON(c.RequestCtx, fasthttp.StatusCreated, posts)
	return nil
}