package apis

import (
	"net/http"

	"github.com/czConstant/constant-nftylend-api/errs"
	"github.com/czConstant/constant-nftylend-api/models"
	"github.com/czConstant/constant-nftylend-api/serializers"
	"github.com/gin-gonic/gin"
)

func (s *Server) GetAffiliateStats(c *gin.Context) {
	ctx := s.requestContext(c)
	affiliateStats, err := s.nls.GetAffiliateStats(
		ctx,
		models.Network(s.stringFromContextQuery(c, "network")),
		s.stringFromContextQuery(c, "address"),
	)
	if err != nil {
		ctxJSON(c, http.StatusBadRequest, &serializers.Resp{Error: errs.NewError(err)})
		return
	}
	affiliateShareRate, err := s.nls.GetAffiliateShareRate(
		ctx,
		models.Network(s.stringFromContextQuery(c, "network")),
		s.stringFromContextQuery(c, "address"),
	)
	if err != nil {
		ctxJSON(c, http.StatusBadRequest, &serializers.Resp{Error: errs.NewError(err)})
		return
	}
	ctxJSON(c, http.StatusOK, &serializers.Resp{Result: serializers.NewAffiliateStatsRespResp(affiliateStats, affiliateShareRate)})
}

func (s *Server) GetAffiliateVolumes(c *gin.Context) {
	ctx := s.requestContext(c)
	_, limit := s.pagingFromContext(c)
	rpts, err := s.nls.GetAffiliateVolumes(
		ctx,
		models.Network(s.stringFromContextQuery(c, "network")),
		s.stringFromContextQuery(c, "address"),
		s.stringFromContextQuery(c, "rpt_by"),
		uint(limit),
	)
	if err != nil {
		ctxAbortWithStatusJSON(c, http.StatusBadRequest, &serializers.Resp{Error: errs.NewError(err)})
		return
	}
	ctxJSON(c, http.StatusOK, &serializers.Resp{Result: serializers.NewAffiliateVolumesRespArr(rpts)})
}
