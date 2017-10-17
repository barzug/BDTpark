package main

import (
	"./daemon"

	"log"
	h "./handlers"

	"./router"

	"github.com/valyala/fasthttp"
)

const port = ":8000"


func addRoutes(r *router.Routing) {
	r.AddRoute(&router.Route{Method: "POST", Path: "/api/forum/create", Function: h.CreateForum})
	r.AddRoute(&router.Route{Method: "POST", Path: "/api/forum/<slug>/create", Function: h.CreateThread})
	r.AddRoute(&router.Route{Method: "GET", Path: "/api/forum/<slug>/details", Function: h.GetForumDetails})
	r.AddRoute(&router.Route{Method: "GET", Path: "/api/forum/<slug>/threads", Function: h.GetThread})
	//r.AddRoute(&router.Route{Method: "GET", Path: "/api/forum/:slug/users", Function: foo})
	//
	//r.AddRoute(&router.Route{Method: "GET", Path: "/api/post/:id/details", Function: foo})
	//r.AddRoute(&router.Route{Method: "POST", Path: "/api/post/:id/details", Function: foo})
	//
	//r.AddRoute(&router.Route{Method: "POST", Path: "/api/service/clear", Function: foo})
	//r.AddRoute(&router.Route{Method: "GET", Path: "/api/service/status", Function: foo})
	//
	r.AddRoute(&router.Route{Method: "POST", Path: "/api/thread/<slug_or_id>/create", Function: h.CreatePosts})
	//r.AddRoute(&router.Route{Method: "GET", Path: "/api/thread/:slug_or_id/details", Function: foo})
	//r.AddRoute(&router.Route{Method: "POST", Path: "/api/thread/:slug_or_id/details", Function: foo})
	//r.AddRoute(&router.Route{Method: "GET", Path: "/api/thread/:slug_or_id/posts", Function: foo})
	r.AddRoute(&router.Route{Method: "POST", Path: "/api/thread/<slug_or_id>/vote", Function: h.VoteForThread})
	//
	r.AddRoute(&router.Route{Method: "POST", Path: "/api/user/<nickname>/create", Function: h.CreateUser})
	r.AddRoute(&router.Route{Method: "GET", Path: "/api/user/<nickname>/profile", Function: h.GetUser})
	r.AddRoute(&router.Route{Method: "POST", Path: "/api/user/<nickname>/profile", Function: h.UpdateUser})
}


func main() {

	log.Printf("Server started")

	err := daemon.Init("localhost", "postgres", "docker", "docker", 100)
	if err != nil {
		log.Fatal(err)
	}

	r := new(router.Routing)

	addRoutes(r)

	err = r.Init()
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(fasthttp.ListenAndServe(port, r.Router.HandleRequest))
}
