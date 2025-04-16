package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"sca-integrator/app/exception"
)

func ErrorHandler() gin.HandlerFunc {
	return func(context *gin.Context) {
		defer func() {
			if err := recover(); err != nil {

				var errType interface{}
				var code int

				switch err.(type) {
				case exception.InternalServerError:
					errType = err.(exception.InternalServerError)
					code = http.StatusInternalServerError
					fmt.Printf("%s %+v \n", code, errType)
				case exception.NotFoundError:
					errType = err.(exception.NotFoundError)
					code = http.StatusNotFound
				case exception.ValidationError:
					errType = err.(exception.ValidationError)
					code = http.StatusBadRequest
				case exception.NotImplementedError:
					errType = err.(exception.NotImplementedError)
					code = http.StatusNotImplemented
				}

				context.JSON(code, errType)
			}
		}()
		context.Next()
	}
}
