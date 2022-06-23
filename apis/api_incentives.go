package apis

import (
	"net/http"

	"github.com/czConstant/constant-nftylend-api/errs"
	"github.com/czConstant/constant-nftylend-api/serializers"
	"github.com/gin-gonic/gin"
)

func (s *Server) JobIncentiveForUnlock(c *gin.Context) {
	ctx := s.requestContext(c)
	var retErr error
	err := s.nls.JobIncentiveForUnlock(ctx)
	if err != nil {
		retErr = errs.MergeError(retErr, err)
	}
	if retErr != nil {
		ctxAbortWithStatusJSON(c, http.StatusBadRequest, &serializers.Resp{Error: errs.NewError(retErr)})
		return
	}
	ctxJSON(c, http.StatusOK, &serializers.Resp{Result: true})
}
