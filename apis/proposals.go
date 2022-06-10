package apis

import (
	"bytes"
	"net/http"

	"github.com/czConstant/constant-nftylend-api/errs"
	"github.com/czConstant/constant-nftylend-api/serializers"
	"github.com/gin-gonic/gin"
)

func (s *Server) GetIpfsInfo(c *gin.Context) {
	hash := c.Param("hash")
	fileData, err := s.nls.GetIpfsInfo(hash)
	if err != nil {
		ctxAbortWithStatusJSON(c, http.StatusBadRequest, &serializers.Resp{Error: errs.NewError(err)})
		return
	}
	r := bytes.NewReader(fileData)
	extraHeaders := map[string]string{}
	c.DataFromReader(http.StatusOK, r.Size(), http.DetectContentType(fileData), r, extraHeaders)
}

func (s *Server) GetProposals(c *gin.Context) {
	ctx := s.requestContext(c)
	page, limit := s.pagingFromContext(c)
	statuses := s.stringArrayFromContextQuery(c, "status")
	proposals, count, err := s.nls.GetProposals(ctx, statuses, page, limit)
	if err != nil {
		ctxAbortWithStatusJSON(c, http.StatusBadRequest, &serializers.Resp{Error: errs.NewError(err)})
		return
	}
	ctxJSON(c, http.StatusOK, &serializers.Resp{Result: serializers.NewProposalRespArr(proposals), Count: &count})
}

func (s *Server) GetProposalDetail(c *gin.Context) {
	ctx := s.requestContext(c)
	proposalID, err := s.uintFromContextParam(c, "proposal_id")
	if err != nil {
		ctxAbortWithStatusJSON(c, http.StatusBadRequest, &serializers.Resp{Error: errs.NewError(err)})
		return
	}
	proposal, err := s.nls.GetProposalDetail(ctx, proposalID)
	if err != nil {
		ctxAbortWithStatusJSON(c, http.StatusBadRequest, &serializers.Resp{Error: errs.NewError(err)})
		return
	}
	ctxJSON(c, http.StatusOK, &serializers.Resp{Result: serializers.NewProposalResp(proposal)})
}

func (s *Server) GetProposalVotes(c *gin.Context) {
	ctx := s.requestContext(c)
	page, limit := s.pagingFromContext(c)
	proposalID, err := s.uintFromContextParam(c, "proposal_id")
	if err != nil {
		ctxAbortWithStatusJSON(c, http.StatusBadRequest, &serializers.Resp{Error: errs.NewError(err)})
		return
	}
	statuses := s.stringArrayFromContextQuery(c, "status")
	proposalVotes, count, err := s.nls.GetProposalVotes(ctx, proposalID, statuses, page, limit)
	if err != nil {
		ctxAbortWithStatusJSON(c, http.StatusBadRequest, &serializers.Resp{Error: errs.NewError(err)})
		return
	}
	ctxJSON(c, http.StatusOK, &serializers.Resp{Result: serializers.NewProposalVoteRespArr(proposalVotes), Count: &count})
}

func (s *Server) CreateProposal(c *gin.Context) {
	ctx := s.requestContext(c)
	var req serializers.CreateProposalReq
	if err := c.ShouldBindJSON(&req); err != nil {
		ctxJSON(c, http.StatusBadRequest, &serializers.Resp{Error: errs.NewError(err)})
		return
	}
	proposal, err := s.nls.CreateProposal(
		ctx,
		&req,
	)
	if err != nil {
		ctxAbortWithStatusJSON(c, http.StatusBadRequest, &serializers.Resp{Error: errs.NewError(err)})
		return
	}
	ctxJSON(c, http.StatusOK, &serializers.Resp{Result: serializers.NewProposalResp(proposal)})
}

func (s *Server) CreateProposalVote(c *gin.Context) {
	ctx := s.requestContext(c)
	var req serializers.CreateProposalVoteReq
	if err := c.ShouldBindJSON(&req); err != nil {
		ctxJSON(c, http.StatusBadRequest, &serializers.Resp{Error: errs.NewError(err)})
		return
	}
	proposalVote, err := s.nls.CreateProposalVote(
		ctx,
		&req,
	)
	if err != nil {
		ctxAbortWithStatusJSON(c, http.StatusBadRequest, &serializers.Resp{Error: errs.NewError(err)})
		return
	}
	ctxJSON(c, http.StatusOK, &serializers.Resp{Result: serializers.NewProposalVoteResp(proposalVote)})
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
