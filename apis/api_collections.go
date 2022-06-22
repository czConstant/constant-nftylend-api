package apis

import (
	"net/http"
	"strings"

	"github.com/czConstant/constant-nftylend-api/errs"
	"github.com/czConstant/constant-nftylend-api/serializers"
	"github.com/gin-gonic/gin"
)

func (s *Server) CreateCollectionSubmitted(c *gin.Context) {
	ctx := s.requestContext(c)
	var req serializers.CollectionSubmittedReq
	if err := c.ShouldBindJSON(&req); err != nil {
		ctxJSON(c, http.StatusBadRequest, &serializers.Resp{Error: errs.NewError(err)})
		return
	}
	err := s.nls.CreateCollectionSubmitted(ctx, &req)
	if err != nil {
		ctxJSON(c, http.StatusBadRequest, &serializers.Resp{Error: errs.NewError(err)})
		return
	}
	ctxJSON(c, http.StatusOK, &serializers.Resp{Result: true})
}

func (s *Server) GetNearApprovedCreators(c *gin.Context) {
	ctx := s.requestContext(c)
	creators, err := s.nls.GetNearApprovedCreators(ctx)
	if err != nil {
		ctxJSON(c, http.StatusBadRequest, &serializers.Resp{Error: errs.NewError(err)})
		return
	}
	ctxSTRING(c, http.StatusOK, strings.Join(creators, "\n"))
}

func (s *Server) GetNearApprovedCollections(c *gin.Context) {
	ctx := s.requestContext(c)
	ms, err := s.nls.GetNearApprovedCollections(ctx)
	if err != nil {
		ctxJSON(c, http.StatusBadRequest, &serializers.Resp{Error: errs.NewError(err)})
		return
	}
	ctxJSON(c, http.StatusOK, &serializers.Resp{Result: serializers.NewCollectionSubmittedRespArr(ms)})
}
