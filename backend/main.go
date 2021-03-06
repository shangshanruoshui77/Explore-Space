package main

import (
    "sse.tongji.edu.cn/leon/handler"

    "github.com/labstack/echo"
    "github.com/labstack/echo/middleware"
    "github.com/labstack/gommon/log"
    "gopkg.in/mgo.v2"
)

func main() {
    e := echo.New()
    e.Logger.SetLevel(log.ERROR)
    e.Use(middleware.Logger())
    e.Use(middleware.JWTWithConfig(middleware.JWTConfig{
        SigningKey: []byte(handler.Key),
        Skipper: func(c echo.Context) bool {
            // Skip authentication for and signup login requests
            if c.Path() == "/login" || c.Path() == "/signup" || c.Path() == "/data" {
                return true
            }
            return false
        },
    }))

    e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
        AllowOrigins: []string{"*"},
        AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE, },
        // AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
    }))

    // Database connection
    db, err := mgo.Dial("127.0.0.1")
    if err != nil {
        e.Logger.Fatal(err)
    }

    // Create indices
    if err = db.Copy().DB("twitter").C("users").EnsureIndex(mgo.Index{
        Key:    []string{"email"},
        Unique: true,
    }); err != nil {
        log.Fatal(err)
    }

    // Initialize handler
    h := &handler.Handler{DB: db}

    // Routes
    e.POST("/signup", h.Signup)
    e.POST("/login", h.Login)
    e.POST("/follow/:id", h.Follow)
    e.POST("/posts", h.CreatePost)
    e.GET("/feed", h.FetchPost)
    // post csv data
    e.POST("/data", h.HandleData)

    // Start server
    e.Logger.Fatal(e.Start(":1323"))
}
