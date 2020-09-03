package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"gtihub.com/hariolate/tonneau/shared"
	"gtihub.com/hariolate/tonneau/shared/models"
)

func main() {
	res := shared.SuccessDataResponse(gin.H{
		"cards": models.GenCards(),
	})
	data, _ := json.MarshalIndent(res, "", "	")

	fmt.Println(string(data))
}
