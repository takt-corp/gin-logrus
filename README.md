# Logrus Logger for Gin

Gin middleware/handler to logger url path using [sirupsen/logrus](https://github.com/sirupsen/logrus).

```go

    r := gin.New()

    r.Use(ginlogrus.LoggerMiddleware(ginlogrus.LoggerMiddlewareParams{
        SkipRoutes: []string{"/ping", "/health"}
    }))

```
