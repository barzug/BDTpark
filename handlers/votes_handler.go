package handlers

import (
	"../database/models"
	"../daemon"

	"encoding/json"
	"github.com/qiangxue/fasthttp-routing"
	"log"
	"github.com/valyala/fasthttp"
	"strconv"
)

func VoteForThread(c *routing.Context) error {
	slugOrId := c.Param("slug_or_id")
	thread := new(models.Threads)
	resultTread := models.Threads{}

	var err error
	if id, err := strconv.ParseInt(slugOrId, 10, 64); err == nil {
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

	vote := new(models.Votes)
	err = json.Unmarshal(c.PostBody(), vote);
	if err != nil {
		log.Print(err)
		return err
	}

	user := new(models.Users)
	user.Nickname = vote.User
	log.Print(user.Nickname)
	resultUser, err := user.GetUserByLogin(daemon.DB.Pool);
	if err != nil {
		daemon.Render.JSON(c.RequestCtx, fasthttp.StatusNotFound, nil)
		return nil
	}

	vote.Thread = resultTread.TID
	vote.User = resultUser.Nickname
	log.Print(vote)

	resultTread.Votes, err = vote.VoteForThreadAndReturningVotes(daemon.DB.Pool);
	if err != nil {
		daemon.Render.JSON(c.RequestCtx, fasthttp.StatusBadRequest, nil)
		return nil
	}
	log.Print(resultTread.Votes)

	daemon.Render.JSON(c.RequestCtx, fasthttp.StatusOK, resultTread)
	return nil
}
