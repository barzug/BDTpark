package handlers

import (
	"github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
	"../database/models"
	"../daemon"
	"log"
)

func GetStatus(c *routing.Context) error {
	type StatusResponse struct {
		Post   int32 `json:"post"`
		Thread int32 `json:"thread"`
		Forum  int32 `json:"forum"`
		User   int32 `json:"user"`
	}

	statusResponse := new(StatusResponse)
	var err error

	statusResponse.Forum, err = models.ForumsCount(daemon.DB.Pool)
	if err != nil {
		log.Print(err)
		daemon.Render.JSON(c.RequestCtx, fasthttp.StatusBadRequest, nil)
		return nil
	}

	statusResponse.User, err = models.UsersCount(daemon.DB.Pool)
	if err != nil {
		log.Print(err)
		daemon.Render.JSON(c.RequestCtx, fasthttp.StatusBadRequest, nil)
		return nil
	}

	statusResponse.Thread, err = models.ThreadsCount(daemon.DB.Pool)
	if err != nil {
		log.Print(err)
		daemon.Render.JSON(c.RequestCtx, fasthttp.StatusBadRequest, nil)
		return nil
	}

	statusResponse.Post, err = models.PostsCount(daemon.DB.Pool)
	if err != nil {
		log.Print(err)
		daemon.Render.JSON(c.RequestCtx, fasthttp.StatusBadRequest, nil)
		return nil
	}




	daemon.Render.JSON(c.RequestCtx, fasthttp.StatusOK, statusResponse)
	return nil
}

func ClearDB(c *routing.Context) error {
	return daemon.DB.InitSchema()
}