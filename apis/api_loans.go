package apis

import (
	"net/http"

	"github.com/czConstant/constant-nftylend-api/errs"
	"github.com/czConstant/constant-nftylend-api/models"
	"github.com/czConstant/constant-nftylend-api/serializers"
	"github.com/gin-gonic/gin"
)

func (s *Server) GetListingLoans(c *gin.Context) {
	ctx := s.requestContext(c)
	page, limit := s.pagingFromContext(c)
	network := s.stringFromContextQuery(c, "network")
	collectionId, _ := s.uintFromContextQuery(c, "collection_id")
	minPrice, _ := s.float64FromContextQuery(c, "min_price")
	maxPrice, _ := s.float64FromContextQuery(c, "max_price")
	minDuration, _ := s.uintFromContextQuery(c, "min_duration")
	maxDuration, _ := s.uintFromContextQuery(c, "max_duration")
	minInterestRate, _ := s.float64FromContextQuery(c, "min_interest_rate")
	maxInterestRate, _ := s.float64FromContextQuery(c, "max_interest_rate")
	excludeIds, _ := s.uintArrayFromContextQuery(c, "exclude_ids")
	var sort []string
	switch s.stringFromContextQuery(c, "sort") {
	case "created_at":
		{
			sort = []string{"created_at asc"}
		}
	case "-created_at":
		{
			sort = []string{"created_at desc"}
		}
	case "principal_amount":
		{
			sort = []string{"principal_amount asc"}
		}
	case "-principal_amount":
		{
			sort = []string{"principal_amount desc"}
		}
	}
	loans, count, err := s.nls.GetListingLoans(
		ctx,
		models.Network(network),
		collectionId,
		minPrice,
		maxPrice,
		minDuration,
		maxDuration,
		minInterestRate,
		maxInterestRate,
		excludeIds,
		sort,
		page,
		limit,
	)
	if err != nil {
		ctxAbortWithStatusJSON(c, http.StatusBadRequest, &serializers.Resp{Error: errs.NewError(err)})
		return
	}
	ctxJSON(c, http.StatusOK, &serializers.Resp{Result: serializers.NewLoanRespArr(loans), Count: &count})
}

func (s *Server) GetLoans(c *gin.Context) {
	ctx := s.requestContext(c)
	page, limit := s.pagingFromContext(c)
	assetId, _ := s.uintFromContextQuery(c, "asset_id")
	loans, count, err := s.nls.GetLoans(
		ctx,
		models.Network(s.stringFromContextQuery(c, "network")),
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

func (s *Server) GetLoanOffers(c *gin.Context) {
	ctx := s.requestContext(c)
	page, limit := s.pagingFromContext(c)
	offers, count, err := s.nls.GetLoanOffers(
		ctx,
		models.Network(s.stringFromContextQuery(c, "network")),
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

func (s *Server) GetLoanTransactions(c *gin.Context) {
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

func (s *Server) CreateLoan(c *gin.Context) {
	ctx := s.requestContext(c)
	var req serializers.CreateLoanReq
	if err := c.ShouldBindJSON(&req); err != nil {
		ctxJSON(c, http.StatusBadRequest, &serializers.Resp{Error: errs.NewError(err)})
		return
	}
	loan, err := s.nls.CreateLoan(ctx, &req)
	if err != nil {
		ctxJSON(c, http.StatusBadRequest, &serializers.Resp{Error: errs.NewError(err)})
		return
	}
	ctxJSON(c, http.StatusOK, &serializers.Resp{Result: serializers.NewLoanResp(loan)})
}

func (s *Server) CreateLoanOffer(c *gin.Context) {
	ctx := s.requestContext(c)
	loanID, err := s.uintFromContextParam(c, "loan_id")
	if err != nil {
		ctxJSON(c, http.StatusBadRequest, &serializers.Resp{Error: errs.NewError(err)})
		return
	}
	var req serializers.CreateLoanOfferReq
	if err := c.ShouldBindJSON(&req); err != nil {
		ctxJSON(c, http.StatusBadRequest, &serializers.Resp{Error: errs.NewError(err)})
		return
	}
	loanOffer, err := s.nls.CreateLoanOffer(ctx, loanID, &req)
	if err != nil {
		ctxJSON(c, http.StatusBadRequest, &serializers.Resp{Error: errs.NewError(err)})
		return
	}
	ctxJSON(c, http.StatusOK, &serializers.Resp{Result: serializers.NewLoanOfferResp(loanOffer)})
}

func (s *Server) NearUpdateLoan(c *gin.Context) {
	ctx := s.requestContext(c)
	var req serializers.CreateLoanNearReq
	if err := c.ShouldBindJSON(&req); err != nil {
		ctxJSON(c, http.StatusBadRequest, &serializers.Resp{Error: errs.NewError(err)})
		return
	}
	loan, _, err := s.nls.NearUpdateLoan(ctx, &req, "client")
	if err != nil {
		ctxJSON(c, http.StatusBadRequest, &serializers.Resp{Error: errs.NewError(err)})
		return
	}
	ctxJSON(c, http.StatusOK, &serializers.Resp{Result: serializers.NewLoanResp(loan)})
}
