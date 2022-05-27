package services

import (
	"context"
	"encoding/json"
	"math/big"

	"github.com/czConstant/constant-nftylend-api/daos"
	"github.com/czConstant/constant-nftylend-api/errs"
	"github.com/czConstant/constant-nftylend-api/helpers"
	"github.com/czConstant/constant-nftylend-api/models"
	"github.com/czConstant/constant-nftylend-api/serializers"
	"github.com/czConstant/constant-nftylend-api/types/numeric"
	"github.com/jinzhu/gorm"
)

func (s *NftLend) GetProposals(ctx context.Context, page int, limit int) ([]*models.Proposal, uint, error) {
	proposals, count, err := s.pd.Find4Page(
		daos.GetDBMainCtx(ctx),
		map[string][]interface{}{},
		map[string][]interface{}{
			"Choices": []interface{}{},
		},
		[]string{"id desc"},
		page,
		limit,
	)
	if err != nil {
		return nil, 0, errs.NewError(err)
	}
	return proposals, count, nil
}

func (s *NftLend) CreateProposal(ctx context.Context, req *serializers.CreateProposalReq) (*models.Proposal, error) {
	switch req.Network {
	case models.NetworkAURORA:
		{
		}
	default:
		{
			return nil, errs.NewError(errs.ErrBadRequest)
		}
	}
	var err error
	var proposal *models.Proposal
	err = daos.WithTransaction(
		daos.GetDBMainCtx(ctx),
		func(tx *gorm.DB) error {
			var msg struct {
				Timestamp int64  `json:"timestamp"`
				Type      string `json:"type"`
				Payload   struct {
					Name     string   `json:"name"`
					Body     string   `json:"body"`
					Start    int64    `json:"start"`
					End      int64    `json:"end"`
					Snapshot int64    `json:"snapshot"`
					Type     string   `json:"type"`
					Choices  []string `json:"choices"`
				} `json:"payload"`
			}
			err = json.Unmarshal([]byte(req.Msg), &msg)
			if err != nil {
				return errs.NewError(err)
			}
			if msg.Type == "" ||
				msg.Timestamp <= 0 ||
				msg.Payload.Snapshot <= 0 ||
				msg.Payload.Start <= 0 ||
				msg.Payload.End <= 0 ||
				msg.Payload.Name == "" ||
				msg.Payload.Body == "" ||
				len(msg.Payload.Choices) <= 0 {
				return errs.NewError(errs.ErrBadRequest)
			}
			var proposal = &models.Proposal{
				Network:    req.Network,
				Address:    req.Address,
				Type:       msg.Type,
				Timestamp:  helpers.TimeFromUnix(msg.Timestamp),
				ChoiceType: msg.Payload.Type,
				Msg:        req.Msg,
				Sig:        req.Sig,
				Start:      helpers.TimeFromUnix(msg.Payload.Start),
				End:        helpers.TimeFromUnix(msg.Payload.End),
				Snapshot:   msg.Payload.Snapshot,
				Name:       msg.Payload.Name,
				Body:       msg.Payload.Body,
				Status:     models.ProposalStatusCreated,
			}
			err = s.pd.Create(
				tx,
				proposal,
			)
			if err != nil {
				return errs.NewError(err)
			}
			for _, choice := range msg.Payload.Choices {
				proposalChoice := &models.ProposalChoice{
					Network:    proposal.Network,
					ProposalID: proposal.ID,
					Name:       choice,
					PowerVote:  numeric.BigFloat{*big.NewFloat(0)},
				}
				err = s.pcd.Create(
					tx,
					proposalChoice,
				)
				if err != nil {
					return errs.NewError(err)
				}
			}
			return nil
		},
	)
	if err != nil {
		return nil, errs.NewError(err)
	}
	return proposal, nil
}
