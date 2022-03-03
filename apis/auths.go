package apis

import (
	"net/http"

	"github.com/czConstant/constant-nftlend-api/serializers"
	"github.com/gin-gonic/gin"
)

func (s *Server) UserMe(c *gin.Context) {
	ctxJSON(c, http.StatusOK, &serializers.Resp{Result: true})
}
