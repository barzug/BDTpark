package handlers

import (
	"../database/models"
	"../daemon"
	"../utils"

	"encoding/json"
	"github.com/qiangxue/fasthttp-routing"
	"log"
	"github.com/valyala/fasthttp"
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

	log.Print(thread)

	if err := thread.CreateThread(daemon.DB.Pool); err != nil {
		if err == utils.UniqueError {
			prevForum, err := thread.GetThreadBySlug(daemon.DB.Pool)

			if err != nil {
				//log.Fatal(err)
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
