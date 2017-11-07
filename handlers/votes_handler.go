package handlers

import (
	"../database/models"
	"../daemon"

	"github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
	"strconv"
	"sync"
	"log"

	"encoding/json"
)

func VoteForThread(c *routing.Context) error {
	slugOrId := c.Param("slug_or_id")
	thread := models.Threads{}

	waitData := &sync.WaitGroup{}

	var errGetThread error

	waitData.Add(1)
	go func(waitData *sync.WaitGroup) {
		defer waitData.Done()

		if id, parseErr := strconv.ParseInt(slugOrId, 10, 64); parseErr == nil {
			thread.TID = id
			errGetThread = thread.GetThreadById(daemon.DB.Pool);
		} else {
			thread.Slug = slugOrId
			errGetThread = thread.GetThreadBySlug(daemon.DB.Pool);
		}

	}(waitData)


	var err, errGetUser error
	var userNickname string


	vote := models.Votes{}
	waitData.Add(1)
	go func(waitData *sync.WaitGroup) {
		defer waitData.Done()
		err = json.Unmarshal(c.PostBody(), &vote);
		log.Print(err)
		log.Print(vote)


		user := models.Users{}
		user.Nickname = vote.User
		user, errGetUser = user.GetUserByLogin(daemon.DB.Pool);
		userNickname = user.Nickname
	}(waitData)

	waitData.Wait()

	if errGetThread != nil || errGetUser != nil {
		daemon.Render.JSON(c.RequestCtx, fasthttp.StatusNotFound, nil)
		return nil
	}

	if err != nil {
		daemon.Render.JSON(c.RequestCtx, fasthttp.StatusBadRequest, nil)
		return nil
	}

	vote.Thread = thread.TID
	vote.User = userNickname

	thread.Votes, err = vote.VoteForThreadAndReturningVotes(daemon.DB.Pool);
	if err != nil {
		daemon.Render.JSON(c.RequestCtx, fasthttp.StatusBadRequest, nil)
		return nil
	}

	daemon.Render.JSON(c.RequestCtx, fasthttp.StatusOK, thread)
	return nil
}
