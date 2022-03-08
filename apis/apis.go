package apis

import (
	"net/http"

	"github.com/czConstant/constant-nftylend-api/errs"
	"github.com/czConstant/constant-nftylend-api/serializers"
	"github.com/gin-gonic/gin"
)

func (s *Server) AppConfigs(c *gin.Context) {
	ctxJSON(c, http.StatusOK, &serializers.Resp{Result: gin.H{
		"program_id": s.conf.Contract.ProgramID,
	}})
}

func (s *Server) NftLendUpdateBlock(c *gin.Context) {
	ctx := s.requestContext(c)
	blockNumber, err := s.uint64FromContextParam(c, "block")
	if err != nil {
		ctxAbortWithStatusJSON(c, http.StatusBadRequest, &serializers.Resp{Error: errs.NewError(err)})
		return
	}
	err = s.nls.LendNftLendUpdateBlock(ctx, blockNumber)
	if err != nil {
		ctxAbortWithStatusJSON(c, http.StatusBadRequest, &serializers.Resp{Error: errs.NewError(err)})
		return
	}
	ctxJSON(c, http.StatusOK, &serializers.Resp{Result: true})
}

func (s *Server) LenInternalHookSolanaInstruction(c *gin.Context) {
	ctx := s.requestContext(c)
	var req struct {
		BlockNumber      uint64      `json:"block_number"`
		BlockTime        uint64      `json:"block_time"`
		TransactionHash  string      `json:"transaction_hash"`
		TransactionIndex uint        `json:"transaction_index"`
		InstructionIndex uint        `json:"instruction_index"`
		Program          string      `json:"program"`
		Instruction      string      `json:"instruction"`
		Data             interface{} `json:"data"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		ctxJSON(c, http.StatusBadRequest, &serializers.Resp{Error: errs.NewError(err)})
		return
	}
	err := s.nls.InternalHookSolanaInstruction(ctx, req.BlockNumber, req.BlockTime, req.TransactionHash, req.TransactionIndex, req.InstructionIndex, req.Program, req.Instruction, req.Data)
	if err != nil {
		ctxAbortWithStatusJSON(c, http.StatusBadRequest, &serializers.Resp{Error: errs.NewError(err)})
		return
	}
	ctxJSON(c, http.StatusOK, &serializers.Resp{Result: true})
}

func (s *Server) GetAssetDetail(c *gin.Context) {
	ctx := s.requestContext(c)
	m, err := s.nls.GetAssetDetail(ctx, s.stringFromContextParam(c, "seo_url"))
	if err != nil {
		ctxAbortWithStatusJSON(c, http.StatusBadRequest, &serializers.Resp{Error: errs.NewError(err)})
		return
	}
	ctxJSON(c, http.StatusOK, &serializers.Resp{Result: serializers.NewAssetResp(m)})
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
	currencies, err := s.nls.GetCurrencies(ctx)
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
