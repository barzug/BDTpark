package main

import (
	"./database"
	"./database/services"
	"./router"

	"os"
	"log"
	"encoding/json"

	"github.com/jackc/pgx"
	"github.com/valyala/fasthttp"
	"github.com/qiangxue/fasthttp-routing"
	"github.com/fasthttp-contrib/render"
)

var (
	r  = render.New()
	db database.DbFacade
)


/////////////////////////////////////////
// forumHandler

func createForum(c *routing.Context) error {
	err := services.CreateForum(db.Pool, c.PostBody())
	if err != nil {
		r.JSON(c.RequestCtx, fasthttp.StatusBadRequest, nil)
		return nil
	}
	r.JSON(c.RequestCtx, fasthttp.StatusCreated, make(map[string]interface{}))
	return nil
}

//func createThreadBySlug(c *routing.Context) error {
//	slug := c.Param("slug")
//	err := services.CreateForum(db.Pool, c.PostBody())
//	if err != nil {
//		r.JSON(c.RequestCtx, fasthttp.StatusBadRequest, nil)
//		return nil
//	}
//	r.JSON(c.RequestCtx, fasthttp.StatusCreated, make(map[string]interface{}))
//	return nil
//}

func getForumDetails(c *routing.Context) error {
	slug := c.Param("slug")
	var byteData, err = services.GetForumBySlug(db.Pool, slug)
	if err != nil {
		r.JSON(c.RequestCtx, fasthttp.StatusNotFound, nil)
		return nil
	}

	var data map[string]interface{}
	json.Unmarshal(byteData, &data)
	r.JSON(c.RequestCtx, fasthttp.StatusOK, data)
	return nil
}


//////////////////////
// userHandler



func addRoutes(r *router.Routing) {
	r.AddRoute(&router.Route{Method: "POST", Path: "forum/create", Function: createForum})
	//r.AddRoute(&router.Route{Method: "POST", Path: "forum/:slug/create", Function: foo})
	r.AddRoute(&router.Route{Method: "GET", Path: "forum/:slug/details", Function: getForumDetails})
	//r.AddRoute(&router.Route{Method: "GET", Path: "forum/:slug/threads", Function: foo})
	//r.AddRoute(&router.Route{Method: "GET", Path: "forum/:slug/users", Function: foo})
	//
	//r.AddRoute(&router.Route{Method: "GET", Path: "/post/:id/details", Function: foo})
	//r.AddRoute(&router.Route{Method: "POST", Path: "/post/:id/details", Function: foo})
	//
	//r.AddRoute(&router.Route{Method: "POST", Path: "/service/clear", Function: foo})
	//r.AddRoute(&router.Route{Method: "GET", Path: "/service/status", Function: foo})
	//
	//r.AddRoute(&router.Route{Method: "POST", Path: "/thread/:slug_or_id/create", Function: foo})
	//r.AddRoute(&router.Route{Method: "GET", Path: "/thread/:slug_or_id/details", Function: foo})
	//r.AddRoute(&router.Route{Method: "POST", Path: "/thread/:slug_or_id/details", Function: foo})
	//r.AddRoute(&router.Route{Method: "GET", Path: "/thread/:slug_or_id/posts", Function: foo})
	//r.AddRoute(&router.Route{Method: "POST", Path: "/thread/:slug_or_id/vote", Function: foo})
	//
	r.AddRoute(&router.Route{Method: "POST", Path: "user/:nickname/create", Function: foo})
	//r.AddRoute(&router.Route{Method: "GET", Path: "user/:nickname/profile", Function: foo})
	//r.AddRoute(&router.Route{Method: "POST", Path: "user/:nickname/profile", Function: foo})
}

func afterConnect(conn *pgx.Conn) error {
	return nil
}

func main() {
	err := db.InitDB(pgx.ConnConfig{
		Host:     "localhost",
		User:     "docker",
		Password: "docker",
		Database: "postgres",
	}, 100, afterConnect) //?
	if err != nil {
		log.Fatal("Unable to create connection pool", "error", err)
		os.Exit(1)
	}
	r := new(router.Routing)
	addRoutes(r)
	r.Init()
	if err := r.Init(); err != nil {
		log.Fatal("Unable to setup router", "error", err)
		os.Exit(1)
	}
	db.InitSchema()

	log.Fatal(fasthttp.ListenAndServe(":8000/api", r.Router.HandleRequest))
}
