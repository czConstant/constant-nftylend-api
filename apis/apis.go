package apis

import (
	"net/http"

	"github.com/czConstant/constant-nftylend-api/errs"
	"github.com/czConstant/constant-nftylend-api/models"
	"github.com/czConstant/constant-nftylend-api/serializers"
	"github.com/gin-gonic/gin"
)

func (s *Server) AppConfigs(c *gin.Context) {
	ctxJSON(c, http.StatusOK, &serializers.Resp{Result: gin.H{
		"program_id":               s.conf.Contract.ProgramID,
		"matic_nftypawn_address":   s.conf.Contract.MaticNftypawnAddress,
		"matic_nftypawn_admin_fee": 100,
		"avax_nftypawn_address":    s.conf.Contract.AvaxNftypawnAddress,
		"avax_nftypawn_admin_fee":  100,
		"bsc_nftypawn_address":     s.conf.Contract.BscNftypawnAddress,
		"bsc_nftypawn_admin_fee":   100,
		"boba_nftypawn_address":    s.conf.Contract.BobaNftypawnAddress,
		"boba_nftypawn_admin_fee":  100,
		"near_nftypawn_address":    s.conf.Contract.NearNftypawnAddress,
		"near_nftypawn_admin_fee":  100,
	}})
}

func (s *Server) MoralisGetNFTs(c *gin.Context) {
	limit, _ := s.uintFromContextQuery(c, "limit")
	rs, err := s.nls.MoralisGetNFTs(s.requestContext(c), s.stringFromContextQuery(c, "chain"), s.stringFromContextParam(c, "address"), s.stringFromContextQuery(c, "cursor"), int(limit))
	if err != nil {
		ctxAbortWithStatusJSON(c, http.StatusBadRequest, &serializers.Resp{Error: errs.NewError(err)})
		return
	}
	ctxJSON(c, http.StatusOK, rs)
}

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

func (s *Server) GetCollectionAssetVerified(c *gin.Context) {
	ctx := s.requestContext(c)
	m, err := s.nls.GetCollectionVerified(ctx, s.stringFromContextQuery(c, "mint"))
	if err != nil {
		ctxAbortWithStatusJSON(c, http.StatusBadRequest, &serializers.Resp{Error: errs.NewError(err)})
		return
	}
	ctxJSON(c, http.StatusOK, &serializers.Resp{Result: serializers.NewCollectionResp(m)})
}

func (s *Server) GetCollections(c *gin.Context) {
	ctx := s.requestContext(c)
	page, limit := s.pagingFromContext(c)
	collections, count, err := s.nls.GetCollections(ctx, page, limit)
	if err != nil {
		ctxAbortWithStatusJSON(c, http.StatusBadRequest, &serializers.Resp{Error: errs.NewError(err)})
		return
	}
	rptCollections, err := s.nls.GetRPTListingCollection(ctx)
	if err != nil {
		ctxAbortWithStatusJSON(c, http.StatusBadRequest, &serializers.Resp{Error: errs.NewError(err)})
		return
	}
	mapRPTCollection := map[uint]uint{}
	for _, rptCollection := range rptCollections {
		mapRPTCollection[rptCollection.CollectionID] = rptCollection.Total
	}
	resps := serializers.NewCollectionRespArr(collections)
	for _, resp := range resps {
		resp.ListingTotal = mapRPTCollection[resp.ID]
	}
	ctxJSON(c, http.StatusOK, &serializers.Resp{Result: resps, Count: &count})
}

func (s *Server) GetCollectionDetail(c *gin.Context) {
	ctx := s.requestContext(c)
	m, err := s.nls.GetCollectionDetail(ctx, s.stringFromContextParam(c, "seo_url"))
	if err != nil {
		ctxAbortWithStatusJSON(c, http.StatusBadRequest, &serializers.Resp{Error: errs.NewError(err)})
		return
	}
	if m == nil {
		ctxAbortWithStatusJSON(c, http.StatusBadRequest, &serializers.Resp{Error: errs.NewError(errs.ErrBadRequest)})
		return
	}
	rpt, err := s.nls.GetRPTCollectionLoan(ctx, m.ID)
	if err != nil {
		ctxAbortWithStatusJSON(c, http.StatusBadRequest, &serializers.Resp{Error: errs.NewError(err)})
		return
	}
	resp := serializers.NewCollectionResp(m)
	resp.TotalVolume = rpt.TotalVolume
	resp.Avg24hAmount = rpt.Avg24hAmount
	resp.TotalListed = rpt.TotalListed
	ctxJSON(c, http.StatusOK, &serializers.Resp{Result: resp})
}

func (s *Server) GetCurrencies(c *gin.Context) {
	ctx := s.requestContext(c)
	currencies, err := s.nls.GetCurrencies(ctx, models.Network(s.stringFromContextQuery(c, "network")))
	if err != nil {
		ctxAbortWithStatusJSON(c, http.StatusBadRequest, &serializers.Resp{Error: errs.NewError(err)})
		return
	}
	ctxJSON(c, http.StatusOK, &serializers.Resp{Result: serializers.NewCurrencyRespArr(currencies)})
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
