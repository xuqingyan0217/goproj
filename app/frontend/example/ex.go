package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "github.com/cloudwego/hertz/pkg/app"
    "github.com/cloudwego/hertz/pkg/app/server"
    "github.com/cloudwego/hertz/pkg/common/utils"
    "github.com/hertz-contrib/jwt"
)

type login struct {
    Username string `form:"username,required" json:"username,required"`
    Password string `form:"password,required" json:"password,required"`
}

var identityKey = "id"

func PingHandler(ctx context.Context, c *app.RequestContext) {
    user, _ := c.Get(identityKey)
    c.JSON(200, utils.H{
        "message": fmt.Sprintf("username:%v", user.(*User).UserName),
    })
}

// User demo
type User struct {
    UserName  string
    FirstName string
    LastName  string
}

func main() {
    h := server.Default()

    // the jwt middleware
    authMiddleware, err := jwt.New(&jwt.HertzJWTMiddleware{
        Realm:       "test zone",
        Key:         []byte("secret key"),
        Timeout:     time.Hour,
        MaxRefresh:  time.Hour,
        IdentityKey: identityKey,
        PayloadFunc: func(data interface{}) jwt.MapClaims {
            if v, ok := data.(*User); ok {
                return jwt.MapClaims{
                    identityKey: v.UserName,
                }
            }
            return jwt.MapClaims{}
        },
        IdentityHandler: func(ctx context.Context, c *app.RequestContext) interface{} {
            claims := jwt.ExtractClaims(ctx, c)
            return &User{
                UserName: claims[identityKey].(string),
            }
        },
        Authenticator: func(ctx context.Context, c *app.RequestContext) (interface{}, error) {
            var loginVals login
            if err := c.BindAndValidate(&loginVals); err != nil {
                return "", jwt.ErrMissingLoginValues
            }
            userID := loginVals.Username
            password := loginVals.Password

            if (userID == "admin" && password == "admin") || (userID == "test" && password == "test") {
                return &User{
                    UserName:  userID,
                    LastName:  "Hertz",
                    FirstName: "CloudWeGo",
                }, nil
            }

            return nil, jwt.ErrFailedAuthentication
        },
        Authorizator: func(data interface{}, ctx context.Context, c *app.RequestContext) bool {
            if v, ok := data.(*User); ok && v.UserName == "admin" {
                return true
            }

            return false
        },
        Unauthorized: func(ctx context.Context, c *app.RequestContext, code int, message string) {
            c.JSON(code, map[string]interface{}{
                "code":    code,
                "message": message,
            })
        },
    })
    if err != nil {
        log.Fatal("JWT Error:" + err.Error())
    }

    // When you use jwt.New(), the function is already automatically called for checking,
    // which means you don't need to call it again.
    errInit := authMiddleware.MiddlewareInit()

    if errInit != nil {
        log.Fatal("authMiddleware.MiddlewareInit() Error:" + errInit.Error())
    }

    h.POST("/login", authMiddleware.LoginHandler)

    h.NoRoute(authMiddleware.MiddlewareFunc(), func(ctx context.Context, c *app.RequestContext) {
        claims := jwt.ExtractClaims(ctx, c)
        log.Printf("NoRoute claims: %#v\n", claims)
        c.JSON(404, map[string]string{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
    })

    auth := h.Group("/auth")
    // Refresh time can be longer than token timeout
    auth.GET("/refresh_token", authMiddleware.RefreshHandler)
    auth.Use(authMiddleware.MiddlewareFunc())
    {
        auth.GET("/ping", PingHandler)
    }

    h.Spin()
}
