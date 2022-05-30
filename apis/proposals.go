package apis

import (
	"net/http"

	"github.com/czConstant/constant-nftylend-api/errs"
	"github.com/czConstant/constant-nftylend-api/serializers"
	"github.com/gin-gonic/gin"
)

func (s *Server) GetProposals(c *gin.Context) {
	ctx := s.requestContext(c)
	page, limit := s.pagingFromContext(c)
	proposals, count, err := s.nls.GetProposals(ctx, page, limit)
	if err != nil {
		ctxAbortWithStatusJSON(c, http.StatusBadRequest, &serializers.Resp{Error: errs.NewError(err)})
		return
	}
	ctxJSON(c, http.StatusOK, &serializers.Resp{Result: serializers.NewProposalRespArr(proposals), Count: &count})
}

func (s *Server) GetProposalVotes(c *gin.Context) {
	ctx := s.requestContext(c)
	page, limit := s.pagingFromContext(c)
	proposalID, err := s.uintFromContextParam(c, "proposal_id")
	if err != nil {
		ctxAbortWithStatusJSON(c, http.StatusBadRequest, &serializers.Resp{Error: errs.NewError(err)})
		return
	}
	address := s.stringFromContextQuery(c, "address")
	statuses := s.stringArrayFromContextQuery(c, "status")
	proposalVotes, count, err := s.nls.GetProposalVotes(ctx, proposalID, address, statuses, page, limit)
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
