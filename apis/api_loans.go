package apis

import (
	"net/http"

	"github.com/czConstant/constant-nftlend-api/errs"
	"github.com/czConstant/constant-nftlend-api/serializers"
	"github.com/gin-gonic/gin"
)

func (s *Server) LendGetListingLoans(c *gin.Context) {
	ctx := s.requestContext(c)
	page, limit := s.pagingFromContext(c)
	collectionId, _ := s.uintFromContextQuery(c, "collection_id")
	minPrice, _ := s.float64FromContextQuery(c, "min_price")
	maxPrice, _ := s.float64FromContextQuery(c, "max_price")
	excludeIds, _ := s.uintArrayFromContextQuery(c, "exclude_ids")
	loans, count, err := s.nls.GetListingLoans(ctx, collectionId, minPrice, maxPrice, excludeIds, page, limit)
	if err != nil {
		ctxAbortWithStatusJSON(c, http.StatusBadRequest, &serializers.Resp{Error: errs.NewError(err)})
		return
	}
	ctxJSON(c, http.StatusOK, &serializers.Resp{Result: serializers.NewLoanRespArr(loans), Count: &count})
}

func (s *Server) LendGetLoans(c *gin.Context) {
	ctx := s.requestContext(c)
	page, limit := s.pagingFromContext(c)
	assetId, _ := s.uintFromContextQuery(c, "asset_id")
	loans, count, err := s.nls.GetLoans(
		ctx,
		s.stringFromContextQuery(c, "owner"),
		s.stringFromContextQuery(c, "lender"),
		assetId,
		s.stringArrayFromContextQuery(c, "status"),
		page,
		limit,
	)
	if err != nil {
		ctxAbortWithStatusJSON(c, http.StatusBadRequest, &serializers.Resp{Error: errs.NewError(err)})
		return
	}
	ctxJSON(c, http.StatusOK, &serializers.Resp{Result: serializers.NewLoanRespArr(loans), Count: &count})
}

func (s *Server) LendGetLoanOffers(c *gin.Context) {
	ctx := s.requestContext(c)
	page, limit := s.pagingFromContext(c)
	offers, count, err := s.nls.GetLoanOffers(
		ctx,
		s.stringFromContextQuery(c, "borrower"),
		s.stringFromContextQuery(c, "lender"),
		s.stringArrayFromContextQuery(c, "status"),
		page,
		limit,
	)
	if err != nil {
		ctxAbortWithStatusJSON(c, http.StatusBadRequest, &serializers.Resp{Error: errs.NewError(err)})
		return
	}
	ctxJSON(c, http.StatusOK, &serializers.Resp{Result: serializers.NewLoanOfferRespArr(offers), Count: &count})
}

func (s *Server) LendGetLoanTransactions(c *gin.Context) {
	ctx := s.requestContext(c)
	page, limit := s.pagingFromContext(c)
	assetId, err := s.uintFromContextQuery(c, "asset_id")
	if err != nil {
		ctxAbortWithStatusJSON(c, http.StatusBadRequest, &serializers.Resp{Error: errs.NewError(err)})
		return
	}
	tnxs, count, err := s.nls.GetLoanTransactions(
		ctx,
		assetId,
		page,
		limit,
	)
	if err != nil {
		ctxAbortWithStatusJSON(c, http.StatusBadRequest, &serializers.Resp{Error: errs.NewError(err)})
		return
	}
	ctxJSON(c, http.StatusOK, &serializers.Resp{Result: serializers.NewLoanTransactionRespArr(tnxs), Count: &count})
}
