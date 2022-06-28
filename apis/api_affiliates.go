package apis

import (
	"net/http"

	"github.com/czConstant/constant-nftylend-api/errs"
	"github.com/czConstant/constant-nftylend-api/models"
	"github.com/czConstant/constant-nftylend-api/serializers"
	"github.com/gin-gonic/gin"
)

func (s *Server) CreateAffiliateSubmitted(c *gin.Context) {
	ctx := s.requestContext(c)
	var req serializers.AffiliateSubmittedReq
	if err := c.ShouldBindJSON(&req); err != nil {
		ctxJSON(c, http.StatusBadRequest, &serializers.Resp{Error: errs.NewError(err)})
		return
	}
	err := s.validateTimestampWithSignature(
		ctx,
		req.Network,
		req.Address,
		req.Signature,
		req.Timestamp,
	)
	if err != nil {
		ctxJSON(c, http.StatusBadRequest, &serializers.Resp{Error: errs.NewError(err)})
		return
	}
	err = s.nls.EmailForAffiliateSubmitted(ctx, &req)
	if err != nil {
		ctxJSON(c, http.StatusBadRequest, &serializers.Resp{Error: errs.NewError(err)})
		return
	}
	ctxJSON(c, http.StatusOK, &serializers.Resp{Result: true})
}

func (s *Server) GetAffiliateStats(c *gin.Context) {
	ctx := s.requestContext(c)
	network, address, err := s.getNetworkAddress(c)
	if err != nil {
		ctxJSON(c, http.StatusBadRequest, &serializers.Resp{Error: errs.NewError(err)})
		return
	}
	affiliateStats, err := s.nls.GetAffiliateStats(
		ctx,
		network,
		address,
	)
	if err != nil {
		ctxJSON(c, http.StatusBadRequest, &serializers.Resp{Error: errs.NewError(err)})
		return
	}
	affiliateShareRate, err := s.nls.GetAffiliateShareRate(
		ctx,
		network,
		address,
	)
	if err != nil {
		ctxJSON(c, http.StatusBadRequest, &serializers.Resp{Error: errs.NewError(err)})
		return
	}
	ctxJSON(c, http.StatusOK, &serializers.Resp{Result: serializers.NewAffiliateStatsRespResp(affiliateStats, affiliateShareRate)})
}

func (s *Server) GetAffiliateVolumes(c *gin.Context) {
	ctx := s.requestContext(c)
	network, address, err := s.getNetworkAddress(c)
	if err != nil {
		ctxJSON(c, http.StatusBadRequest, &serializers.Resp{Error: errs.NewError(err)})
		return
	}
	_, limit := s.pagingFromContext(c)
	rpts, err := s.nls.GetAffiliateVolumes(
		ctx,
		network,
		address,
		s.stringFromContextQuery(c, "rpt_by"),
		uint(limit),
	)
	if err != nil {
		ctxAbortWithStatusJSON(c, http.StatusBadRequest, &serializers.Resp{Error: errs.NewError(err)})
		return
	}

	ctxJSON(c, http.StatusOK, &serializers.Resp{Result: serializers.NewAffiliateVolumesRespArr(rpts)})
}

func (s *Server) GetAffiliateTransactions(c *gin.Context) {
	ctx := s.requestContext(c)
	network, address, err := s.getNetworkAddress(c)
	if err != nil {
		ctxJSON(c, http.StatusBadRequest, &serializers.Resp{Error: errs.NewError(err)})
		return
	}
	page, limit := s.pagingFromContext(c)
	statuses := s.stringArrayFromContextQuery(c, "status")
	txns, count, err := s.nls.GetIncentiveTransactions(ctx,
		network,
		address,
		[]string{
			string(models.IncentiveTransactionTypeAffiliateBorrowerLoanDone),
			string(models.IncentiveTransactionTypeAffiliateLenderLoanDone),
		},
		statuses,
		page,
		limit,
	)
	if err != nil {
		ctxAbortWithStatusJSON(c, http.StatusBadRequest, &serializers.Resp{Error: errs.NewError(err)})
		return
	}
	ctxJSON(c, http.StatusOK, &serializers.Resp{Result: serializers.NewIncentiveTransactionRespArr(txns), Count: &count})
}
