# context

installation
```go
go get -u git.sportsdots.ru/go-util/sportctx.git
```

example
```go
package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"git.sportsdots.ru/go-util/sportctx.git"
	"git.sportsdots.ru/api-service/internal/form"
	"git.sportsdots.ru/api-service/internal/sporterror"
)

func (hr *Handler) error(c *gin.Context, err error) {
	ctx := c.Request.Context()

	var requestID string
	if sc, ok := ctx.(*sportctx.Context); ok {
		requestID = sc.GetRequestID()
	}

	status := http.StatusInternalServerError

	var e sportctx.Error
	if errors.As(err, &e) {
		switch e.Type {
		case sportctx.BadRequestErrorType:
			status = http.StatusBadRequest
		case sportctx.NotFoundErrorType:
			status = http.StatusNotFound
        case sportctx.ErrorTypeForbidden:
			status = http.StatusForbidden
		case sportctx.ErrorTypeUnauthorized:
			status = http.StatusUnauthorized
		}
	}

	c.AbortWithStatusJSON(status, form.Error{
		Error:     err.Error(),
		RequestID: requestID,
	})
}
```

- содержит RequestID (генерирует сам или берет из http.Request), который в дальнейшем может использоваться для логирования / отправки этого идентификатора в другие сервисы (для трассировки)
- может содержать снапшот конфигурации, чтобы в рамках обработки запроса один и тот же параметр был идемпотентен
