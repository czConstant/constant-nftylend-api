package apis

import (
	"net/http"
	"strings"

	"github.com/czConstant/constant-nftylend-api/errs"
	"github.com/czConstant/constant-nftylend-api/models"
	"github.com/czConstant/constant-nftylend-api/serializers"
	"github.com/gin-gonic/gin"
)

func (s *Server) CreateCollectionSubmission(c *gin.Context) {
	ctx := s.requestContext(c)
	var req serializers.CollectionSubmissionReq
	if err := c.ShouldBindJSON(&req); err != nil {
		ctxJSON(c, http.StatusBadRequest, &serializers.Resp{Error: errs.NewError(err)})
		return
	}
	err := s.nls.CreateCollectionSubmission(ctx, &req)
	if err != nil {
		ctxJSON(c, http.StatusBadRequest, &serializers.Resp{Error: errs.NewError(err)})
		return
	}
	ctxJSON(c, http.StatusOK, &serializers.Resp{Result: true})
}

func (s *Server) GetNearApprovedCreators(c *gin.Context) {
	ctx := s.requestContext(c)
	creators, err := s.nls.GetNearApprovedCreators(ctx)
	if err != nil {
		ctxJSON(c, http.StatusBadRequest, &serializers.Resp{Error: errs.NewError(err)})
		return
	}
	ctxSTRING(c, http.StatusOK, strings.Join(creators, "\n"))
}

func (s *Server) GetNearApprovedCollections(c *gin.Context) {
	ctx := s.requestContext(c)
	ms, err := s.nls.GetNearApprovedCollections(ctx)
	if err != nil {
		ctxJSON(c, http.StatusBadRequest, &serializers.Resp{Error: errs.NewError(err)})
		return
	}
	ctxJSON(c, http.StatusOK, &serializers.Resp{Result: serializers.NewCollectionSubmissionRespArr(ms)})
}

func (s *Server) GetCollectionAssetVerified(c *gin.Context) {
	ctx := s.requestContext(c)
	m, err := s.nls.GetCollectionVerified(
		ctx,
		models.Network(s.stringFromContextQuery(c, "network")),
		s.stringFromContextQuery(c, "contract_address"),
		s.stringFromContextQuery(c, "token_id"),
	)
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
	rpt, err := s.nls.GetCollectionStats(ctx, m.ID)
	if err != nil {
		ctxAbortWithStatusJSON(c, http.StatusBadRequest, &serializers.Resp{Error: errs.NewError(err)})
		return
	}
	resp := serializers.NewCollectionResp(m)
	resp.TotalVolume = rpt.TotalVolume
	resp.Avg24hAmount = rpt.Avg24hAmount
	resp.MinAmount = rpt.MinAmount
	resp.TotalListed = rpt.TotalListed
	ctxJSON(c, http.StatusOK, &serializers.Resp{Result: resp})
}
