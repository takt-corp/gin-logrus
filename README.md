# Logrus Logger for Gin

Gin middleware/handler to logger url path using [sirupsen/logrus](https://github.com/sirupsen/logrus).

```go

    r := gin.New()

    r.Use(ginlogrus.LoggerMiddleware(ginlogrus.LoggerMiddlewareParams{
        SkipPaths: []string{"/ping", "/health", "/:some_id/ping"}
    }))

```
