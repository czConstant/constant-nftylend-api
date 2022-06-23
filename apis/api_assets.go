package apis

import (
	"net/http"

	"github.com/czConstant/constant-nftylend-api/errs"
	"github.com/czConstant/constant-nftylend-api/serializers"
	"github.com/gin-gonic/gin"
)

func (s *Server) GetAssetDetail(c *gin.Context) {
	ctx := s.requestContext(c)
	m, err := s.nls.GetAssetDetail(ctx, s.stringFromContextParam(c, "seo_url"))
	if err != nil {
		ctxAbortWithStatusJSON(c, http.StatusBadRequest, &serializers.Resp{Error: errs.NewError(err)})
		return
	}
	resp := serializers.NewAssetResp(m)
	if m != nil {
		stats, err := s.nls.GetAssetStatsInfo(ctx, m.ID)
		if err != nil {
			ctxAbortWithStatusJSON(c, http.StatusBadRequest, &serializers.Resp{Error: errs.NewError(err)})
			return
		}
		resp.Stats = stats
	}
	ctxJSON(c, http.StatusOK, &serializers.Resp{Result: resp})
}

func (s *Server) GetAssetDetailInfo(c *gin.Context) {
	ctx := s.requestContext(c)
	m, err := s.nls.GetAssetDetailInfo(ctx, s.stringFromContextQuery(c, "contract_address"), s.stringFromContextQuery(c, "token_id"))
	if err != nil {
		ctxAbortWithStatusJSON(c, http.StatusBadRequest, &serializers.Resp{Error: errs.NewError(err)})
		return
	}
	resp := serializers.NewAssetResp(m)
	if m != nil {
		stats, err := s.nls.GetAssetStatsInfo(ctx, m.ID)
		if err != nil {
			ctxAbortWithStatusJSON(c, http.StatusBadRequest, &serializers.Resp{Error: errs.NewError(err)})
			return
		}
		resp.Stats = stats
	}
	ctxJSON(c, http.StatusOK, &serializers.Resp{Result: resp})
}

func (s *Server) GetAseetTransactions(c *gin.Context) {
	ctx := s.requestContext(c)
	page, limit := s.pagingFromContext(c)
	assetId, err := s.uintFromContextQuery(c, "asset_id")
	if err != nil {
		ctxAbortWithStatusJSON(c, http.StatusBadRequest, &serializers.Resp{Error: errs.NewError(err)})
		return
	}
	tnxs, count, err := s.nls.GetAseetTransactions(
		ctx,
		assetId,
		page,
		limit,
	)
	if err != nil {
		ctxAbortWithStatusJSON(c, http.StatusBadRequest, &serializers.Resp{Error: errs.NewError(err)})
		return
	}
	ctxJSON(c, http.StatusOK, &serializers.Resp{Result: serializers.NewAssetTransactionRespArr(tnxs), Count: &count})
}
