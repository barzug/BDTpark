package handlers

import (
	"sync"

	"../daemon"
	"../database/models"
	"../utils"

	"encoding/json"

	"github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"

)

func CreateUser(c *routing.Context) error {
	nickname := c.Param("nickname")
	user := new(models.Users)
	if err := json.Unmarshal(c.PostBody(), user); err != nil {
		return err
	}
	user.Nickname = nickname

	if err := user.CreateUser(daemon.DB.Pool); err != nil {
		if err == utils.UniqueError {
			users, err := user.GetUserByLoginAndEmail(daemon.DB.Pool)

			if err != nil {
				daemon.Render.JSON(c.RequestCtx, fasthttp.StatusBadRequest, nil)
				return err
			}

			daemon.Render.JSON(c.RequestCtx, fasthttp.StatusConflict, users)
			return nil
		}
		daemon.Render.JSON(c.RequestCtx, fasthttp.StatusBadRequest, nil)
		return nil
	}

	daemon.Render.JSON(c.RequestCtx, fasthttp.StatusCreated, user)
	return nil
}

func GetUser(c *routing.Context) error {
	nickname := c.Param("nickname")
	user := new(models.Users)
	user.Nickname = nickname

	resultUser, err := user.GetUserByLogin(daemon.DB.Pool)
	if err != nil {
		daemon.Render.JSON(c.RequestCtx, fasthttp.StatusNotFound, nil)
		return nil
	}

	daemon.Render.JSON(c.RequestCtx, fasthttp.StatusOK, resultUser)
	return nil
}

func UpdateUser(c *routing.Context) error {
	nickname := c.Param("nickname")
	user := new(models.Users)
	if err := json.Unmarshal(c.PostBody(), user); err != nil {
		return err
	}

	user.Nickname = nickname

	if utils.CheckEmpty(user) {
		prevUser, err := user.GetUserByLogin(daemon.DB.Pool)
		if err != nil {
			daemon.Render.JSON(c.RequestCtx, fasthttp.StatusNotFound, nil)
			return nil
		}
		utils.AdditionObject(user, &prevUser)

	}

	if err := user.UpdateUser(daemon.DB.Pool); err != nil {
		if err == utils.UniqueError {
			daemon.Render.JSON(c.RequestCtx, fasthttp.StatusConflict, nil)
			return nil
		}
		daemon.Render.JSON(c.RequestCtx, fasthttp.StatusNotFound, nil)
		return nil
	}

	daemon.Render.JSON(c.RequestCtx, fasthttp.StatusOK, user)
	return nil
}

func GetForumUsers(c *routing.Context) error {
 	slug := c.Param("slug")
 	forum := new(models.Forums)
	forum.Slug = slug
	 
	waitGetForum := &sync.WaitGroup{}

	var errInGetForum error

	waitGetForum.Add(1)
	go func(waitGetForum *sync.WaitGroup, err *error) {
		defer waitGetForum.Done()
 		*err = forum.GetForumBySlug(daemon.DB.Pool);
	 }(waitGetForum, &errInGetForum)


 	limit := string(c.QueryArgs().Peek("limit"))
 	since := string(c.QueryArgs().Peek("since"))
	desc := string(c.QueryArgs().Peek("desc"))
	 
	if errInGetForum != nil {
		daemon.Render.JSON(c.RequestCtx, fasthttp.StatusNotFound, nil)
		return nil
	}

 	users, err := forum.GetMembers(daemon.DB.Pool, limit, since, desc)
 	if err != nil {
 		daemon.Render.JSON(c.RequestCtx, fasthttp.StatusBadRequest, nil)
 		return nil
	 }
	 

	waitGetForum.Wait()
	if errInGetForum != nil {
		daemon.Render.JSON(c.RequestCtx, fasthttp.StatusNotFound, nil)
		return nil
	}

 	daemon.Render.JSON(c.RequestCtx, fasthttp.StatusOK, users)
 	return nil
 }
