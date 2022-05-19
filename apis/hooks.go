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

func (s *Server) BlockchainScanBlock(c *gin.Context) {
	ctx := s.requestContext(c)
	blockNumber, err := s.uint64FromContextParam(c, "block")
	if err != nil {
		ctxAbortWithStatusJSON(c, http.StatusBadRequest, &serializers.Resp{Error: errs.NewError(err)})
		return
	}
	network := models.Network(c.Param("network"))
	switch network {
	case models.NetworkSOL:
		{
			err = s.nls.LendNftLendUpdateBlock(ctx, blockNumber)
			if err != nil {
				ctxAbortWithStatusJSON(c, http.StatusBadRequest, &serializers.Resp{Error: errs.NewError(err)})
				return
			}
		}
	case models.NetworkMATIC,
		models.NetworkAVAX,
		models.NetworkBSC,
		models.NetworkBOBA:
		{
			err = s.nls.JobEvmNftypawnFilterLogs(ctx, network, blockNumber)
			if err != nil {
				ctxAbortWithStatusJSON(c, http.StatusBadRequest, &serializers.Resp{Error: errs.NewError(err)})
				return
			}
		}
	default:
		{
			ctxAbortWithStatusJSON(c, http.StatusBadRequest, &serializers.Resp{Error: errs.NewError(errs.ErrBadRequest)})
			return
		}
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
	err := s.nls.InternalHookSolanaInstruction(ctx, models.NetworkSOL, req.BlockNumber, req.BlockTime, req.TransactionHash, req.TransactionIndex, req.InstructionIndex, req.Program, req.Instruction, req.Data)
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

func (s *Server) NearSync(c *gin.Context) {
	ctx := s.requestContext(c)
	var req struct {
		Hash        string `json:"hash"`
		NftContract string `json:"nft_contract"`
		TokenID     string `json:"token_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		ctxJSON(c, http.StatusBadRequest, &serializers.Resp{Error: errs.NewError(err)})
		return
	}
	_, isUpdated, err := s.nls.NearUpdateLoan(
		ctx,
		&serializers.CreateLoanNearReq{
			ContractAddress: req.NftContract,
			TokenID:         req.TokenID,
		},
		"worker",
	)
	if err != nil {
		ctxAbortWithStatusJSON(c, http.StatusBadRequest, &serializers.Resp{Error: errs.NewError(err)})
		return
	}
	ctxJSON(c, http.StatusOK, &serializers.Resp{Result: isUpdated})
}
