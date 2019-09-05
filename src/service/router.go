package service

import "src/action"

func setRouters() {
	{
		//user := app.Group("/user", midware)
		//user.GET("/list", handler1,handler2,handler3)

		app.POST("/post", action.PostHandler)
		app.GET("/get/:id", action.GetHandler)
		app.GET("/panic", action.PanicTest)
		app.GET("/err", action.ErrTest)
		app.GET("/broken", action.BrokenPipeTest)
	}
}
