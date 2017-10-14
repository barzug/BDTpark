package handlers
//
//import (
//	"github.com/valyala/fasthttp"
//	"github.com/qiangxue/fasthttp-routing"
//	"encoding/json"
//
//	"../database/services"
//)
//
//
//func createForum(c *routing.Context) error {
//	err := services.CreateForum(db.Pool, c.PostBody())
//	if err != nil {
//		r.JSON(c.RequestCtx, fasthttp.StatusBadRequest, nil)
//		return nil
//	}
//	r.JSON(c.RequestCtx, fasthttp.StatusCreated, make(map[string]interface{}))
//	return nil
//}
//
////func createThreadBySlug(c *routing.Context) error {
////	slug := c.Param("slug")
////	err := services.CreateForum(db.Pool, c.PostBody())
////	if err != nil {
////		r.JSON(c.RequestCtx, fasthttp.StatusBadRequest, nil)
////		return nil
////	}
////	r.JSON(c.RequestCtx, fasthttp.StatusCreated, make(map[string]interface{}))
////	return nil
////}
//
//func getForumDetails(c *routing.Context) error {
//	slug := c.Param("slug")
//	var byteData, err = services.GetForumBySlug(db.Pool, slug)
//	if err != nil {
//		r.JSON(c.RequestCtx, fasthttp.StatusNotFound, nil)
//		return nil
//	}
//
//	var data map[string]interface{}
//	json.Unmarshal(byteData, &data)
//	r.JSON(c.RequestCtx, fasthttp.StatusOK, data)
//	return nil
//}