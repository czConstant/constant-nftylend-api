package apis

import (
	"net/http"

	"github.com/czConstant/constant-nftylend-api/errs"
	"github.com/czConstant/constant-nftylend-api/models"
	"github.com/czConstant/constant-nftylend-api/serializers"
	"github.com/gin-gonic/gin"
)

func (s *Server) JobUpdateCurrencyPrice(c *gin.Context) {
	ctx := s.requestContext(c)
	err := s.nls.JobUpdateCurrencyPrice(ctx)
	if err != nil {
		ctxAbortWithStatusJSON(c, http.StatusBadRequest, &serializers.Resp{Error: errs.NewError(err)})
		return
	}
	ctxJSON(c, http.StatusOK, &serializers.Resp{Result: true})
}

func (s *Server) JobEvmNftypawnFilterLogs(c *gin.Context) {
	ctx := s.requestContext(c)
	var retErr error
	err := s.nls.JobEvmNftypawnFilterLogs(ctx, models.NetworkMATIC, 0)
	if err != nil {
		retErr = errs.MergeError(retErr, err)
	}
	err = s.nls.JobEvmNftypawnFilterLogs(ctx, models.NetworkAVAX, 0)
	if err != nil {
		retErr = errs.MergeError(retErr, err)
	}
	err = s.nls.JobEvmNftypawnFilterLogs(ctx, models.NetworkBSC, 0)
	if err != nil {
		retErr = errs.MergeError(retErr, err)
	}
	err = s.nls.JobEvmNftypawnFilterLogs(ctx, models.NetworkBOBA, 0)
	if err != nil {
		retErr = errs.MergeError(retErr, err)
	}
	if retErr != nil {
		ctxAbortWithStatusJSON(c, http.StatusBadRequest, &serializers.Resp{Error: errs.NewError(retErr)})
		return
	}
	ctxJSON(c, http.StatusOK, &serializers.Resp{Result: true})
}

func (s *Server) JobEmailSchedule(c *gin.Context) {
	ctx := s.requestContext(c)
	var retErr error
	err := s.nls.JobEmailSchedule(ctx)
	if err != nil {
		retErr = errs.MergeError(retErr, err)
	}
	if retErr != nil {
		ctxAbortWithStatusJSON(c, http.StatusBadRequest, &serializers.Resp{Error: errs.NewError(retErr)})
		return
	}
	ctxJSON(c, http.StatusOK, &serializers.Resp{Result: true})
}

func (s *Server) JobIncentiveForUnlock(c *gin.Context) {
	ctx := s.requestContext(c)
	var retErr error
	err := s.nls.JobIncentiveStatus(ctx)
	if err != nil {
		retErr = errs.MergeError(retErr, err)
	}
	if retErr != nil {
		ctxAbortWithStatusJSON(c, http.StatusBadRequest, &serializers.Resp{Error: errs.NewError(retErr)})
		return
	}
	ctxJSON(c, http.StatusOK, &serializers.Resp{Result: true})
}

func (s *Server) JobProposalStatus(c *gin.Context) {
	ctx := s.requestContext(c)
	var retErr error
	// err := s.nls.JobProposalUnVote(ctx)
	// if err != nil {
	// 	retErr = errs.MergeError(retErr, err)
	// }
	err := s.nls.JobProposalStatus(ctx)
	if err != nil {
		retErr = errs.MergeError(retErr, err)
	}
	if retErr != nil {
		ctxAbortWithStatusJSON(c, http.StatusBadRequest, &serializers.Resp{Error: errs.NewError(retErr)})
		return
	}
	ctxJSON(c, http.StatusOK, &serializers.Resp{Result: true})
}

func (s *Server) JobUpdateStats(c *gin.Context) {
	ctx := s.requestContext(c)
	var retErr error
	err := s.nls.JobUpdateStatsCollection(ctx)
	if err != nil {
		retErr = errs.MergeError(retErr, err)
	}
	if retErr != nil {
		ctxAbortWithStatusJSON(c, http.StatusBadRequest, &serializers.Resp{Error: errs.NewError(retErr)})
		return
	}
	ctxJSON(c, http.StatusOK, &serializers.Resp{Result: true})
}
