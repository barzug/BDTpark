package handlers

import (
	"../database/models"
	"../daemon"
	"../utils"

	"encoding/json"
	"github.com/valyala/fasthttp"
	"github.com/qiangxue/fasthttp-routing"
)

func CreateForum(c *routing.Context) error {

	forum := new(models.Forums)
	if err := json.Unmarshal(c.PostBody(), forum); err != nil {
		return err
	}

	author := models.Users{Nickname: forum.Author}

	forumAuthor, err := author.GetUserByLogin(daemon.DB.Pool)
	if err != nil {
		daemon.Render.JSON(c.RequestCtx, fasthttp.StatusNotFound, nil)
		return nil
	}

	forum.Author = forumAuthor.Nickname

	if err := forum.CreateForum(daemon.DB.Pool); err != nil {
		if err == utils.UniqueError {
			err := forum.GetForumBySlug(daemon.DB.Pool)

			if err != nil {
				daemon.Render.JSON(c.RequestCtx, fasthttp.StatusBadRequest, nil)
				return err
			}

			daemon.Render.JSON(c.RequestCtx, fasthttp.StatusConflict, forum)
			return nil
		}
		daemon.Render.JSON(c.RequestCtx, fasthttp.StatusBadRequest, nil)
		return nil
	}

	daemon.Render.JSON(c.RequestCtx, fasthttp.StatusCreated, forum)
	return nil
}

func GetForumDetails(c *routing.Context) error {
	slug := c.Param("slug")
	forum := new(models.Forums)
	forum.Slug = slug

	err := forum.GetForumBySlug(daemon.DB.Pool);
	if err != nil {
		daemon.Render.JSON(c.RequestCtx, fasthttp.StatusNotFound, nil)
		return nil
	}

	daemon.Render.JSON(c.RequestCtx, fasthttp.StatusOK, forum)
	return nil
}
