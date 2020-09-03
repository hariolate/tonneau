package shared

import "github.com/gin-gonic/gin"

func ErrorMessageResponse(message string) gin.H {
	return gin.H{
		"status": "error",
		"data": gin.H{
			"message": message,
		},
	}
}

func SuccessMessageResponse(message string) gin.H {
	return gin.H{
		"status": "success",
		"data": gin.H{
			"message": message,
		},
	}
}

func SuccessDataResponse(data gin.H) gin.H {
	return gin.H{
		"status": "success",
		"data":   data,
	}
}
