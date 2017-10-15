package handlers

import (
	"../database/models"
	"../daemon"
	"../utils"

	"encoding/json"
	"github.com/valyala/fasthttp"
	"github.com/qiangxue/fasthttp-routing"
	"log"
)

func CreateUser(c *routing.Context) error {
	nickname := c.Param("nickname")
	user := new(models.Users)
	if err := json.Unmarshal(c.PostBody(), user); err != nil {
		log.Fatal(err)
		return err
	}
	user.Nickname = nickname

	if err := user.CreateUser(daemon.DB.Pool); err != nil {
		if err == utils.UniqueError {
			users, err := user.GetUserByLoginAndEmail(daemon.DB.Pool)

			if err != nil {
				log.Fatal(err)
				daemon.Render.JSON(c.RequestCtx, fasthttp.StatusBadRequest, nil)
				return err
			}

			daemon.Render.JSON(c.RequestCtx, fasthttp.StatusConflict, users)
			return nil
		}
		log.Fatal(err)
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

	resultUser, err := user.GetUserByLogin(daemon.DB.Pool);
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
		log.Fatal(err)
		return err
	}
	user.Nickname = nickname

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
