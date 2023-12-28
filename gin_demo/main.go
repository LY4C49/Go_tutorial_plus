package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func func1(c *gin.Context) {
	fmt.Println("func1 front")

	c.Next()

	fmt.Println("func1 after")
}

func func2(c *gin.Context) {
	fmt.Println("this is func2")
}

func func3(c *gin.Context) {
	fmt.Println("this func3")
	c.Set("the_key", "the_value")
}

func func4(c *gin.Context) {
	fmt.Println("this func4")
	v, ok := c.Get("the_key")
	if ok {
		fmt.Println("We get: ", v.(string))
	}
	c.Abort()
}

func func5(c *gin.Context) {
	fmt.Println("this func5")
}

func main() {
	r := gin.Default()

	// r.Use() 手动添加中间件
	//或者, 对某一个路由组生效的中间件
	// func1(c *gin.Context)
	//shopGroup := r.Group("/shop", func1, func2)
	//shopGroup.Use(func3)
	//{
	//	shopGroup.GET("/index", func4, func5)
	//}

	shopGroup := r.Group("/shop", func1, func2)
	shopGroup.Use(func3)
	//{
	shopGroup.GET("/index", func4, func5)
	//}
	/*
		c.Next() 调用HandlerChain的下一个
		c.Abort() 终止（通过将Index设成最大值）
		c.Set() 存值 eg:func1存一个值，func2可以取
		c.Get() 取值
	*/

	//r.GET("/", func(c *gin.Context) {
	//	c.String(http.StatusOK, "ok")
	//})

	r.Run()
}
