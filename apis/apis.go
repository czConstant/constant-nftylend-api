package apis

import (
	"net/http"

	"github.com/czConstant/constant-nftlend-api/errs"
	"github.com/czConstant/constant-nftlend-api/serializers"
	"github.com/gin-gonic/gin"
)

func (s *Server) LendNftLendUpdateBlock(c *gin.Context) {
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

func (s *Server) LendGetAssetDetail(c *gin.Context) {
	ctx := s.requestContext(c)
	m, err := s.nls.GetAssetDetail(ctx, s.stringFromContextParam(c, "seo_url"))
	if err != nil {
		ctxAbortWithStatusJSON(c, http.StatusBadRequest, &serializers.Resp{Error: errs.NewError(err)})
		return
	}
	ctxJSON(c, http.StatusOK, &serializers.Resp{Result: serializers.NewAssetResp(m)})
}

func (s *Server) LendGetCollections(c *gin.Context) {
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

func (s *Server) LendGetCollectionDetail(c *gin.Context) {
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

func (s *Server) LendGetCurrencies(c *gin.Context) {
	ctx := s.requestContext(c)
	currencies, err := s.nls.GetCurrencies(ctx)
	if err != nil {
		ctxAbortWithStatusJSON(c, http.StatusBadRequest, &serializers.Resp{Error: errs.NewError(err)})
		return
	}
	ctxJSON(c, http.StatusOK, &serializers.Resp{Result: serializers.NewCurrencyRespArr(currencies)})
}