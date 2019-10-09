package service

import "src/action"

func setRouters() {
	{
		//user := App.Group("/user", middleware)
		//user.GET("/list", handler1,handler2,handler3) // 按从左到右顺序执行

		App.POST("/post", action.PostHandler)
		App.GET("/get/:id", action.GetHandler)
		App.GET("/panic", action.PanicTest)
		App.GET("/err", action.ErrTest)
	}
}
