package services

import (
	"context"
	"encoding/json"
	"math/big"
	"time"

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

func (s *NftLend) GetProposalVotes(ctx context.Context, proposalID uint, page int, limit int) ([]*models.ProposalVote, uint, error) {
	proposalVotes, count, err := s.pvd.Find4Page(
		daos.GetDBMainCtx(ctx),
		map[string][]interface{}{
			"proposal_id = ?": []interface{}{proposalID},
		},
		map[string][]interface{}{
			"Proposal":       []interface{}{},
			"ProposalChoice": []interface{}{},
		},
		[]string{"id desc"},
		page,
		limit,
	)
	if err != nil {
		return nil, 0, errs.NewError(err)
	}
	return proposalVotes, count, nil
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
	err := s.bcs.Aurora.ValidateMessageSignature(
		req.Msg,
		req.Sig,
		req.Address,
	)
	if err != nil {
		return nil, errs.NewError(err)
	}
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
				len(msg.Payload.Choices) <= 1 {
				return errs.NewError(errs.ErrBadRequest)
			}
			if msg.Timestamp < time.Now().Add(-60*time.Second).Unix() ||
				msg.Timestamp > time.Now().Add(60*time.Second).Unix() {
				return errs.NewError(errs.ErrBadRequest)
			}
			if msg.Payload.Start < time.Now().Add(-60*time.Second).Unix() {
				return errs.NewError(errs.ErrBadRequest)
			}
			if msg.Payload.End < time.Now().Add(-60*time.Second).Unix() {
				return errs.NewError(errs.ErrBadRequest)
			}
			if msg.Payload.Start >= msg.Payload.End {
				return errs.NewError(errs.ErrBadRequest)
			}
			ipfsData, err := json.Marshal(&req)
			if err != nil {
				return errs.NewError(err)
			}
			ipfsHash, err := s.ifc.UploadString(string(ipfsData))
			if err != nil {
				return errs.NewError(err)
			}
			proposal, err = s.pd.First(
				tx,
				map[string][]interface{}{
					"ipfs_hash = ?": []interface{}{ipfsHash},
				},
				map[string][]interface{}{},
				[]string{},
			)
			if err != nil {
				return errs.NewError(err)
			}
			if proposal != nil {
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
				IpfsHash:   ipfsHash,
				Status:     models.ProposalStatusCreated,
			}
			err = s.pd.Create(
				tx,
				proposal,
			)
			if err != nil {
				return errs.NewError(err)
			}
			for idx, choice := range msg.Payload.Choices {
				proposalChoice := &models.ProposalChoice{
					Network:    proposal.Network,
					ProposalID: proposal.ID,
					Choice:     (idx + 1),
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

func (s *NftLend) CreateProposalVote(ctx context.Context, req *serializers.CreateProposalVoteReq) (*models.ProposalVote, error) {
	switch req.Network {
	case models.NetworkAURORA:
		{
		}
	default:
		{
			return nil, errs.NewError(errs.ErrBadRequest)
		}
	}
	err := s.bcs.Aurora.ValidateMessageSignature(
		req.Msg,
		req.Sig,
		req.Address,
	)
	if err != nil {
		return nil, errs.NewError(err)
	}
	var proposalVote *models.ProposalVote
	err = daos.WithTransaction(
		daos.GetDBMainCtx(ctx),
		func(tx *gorm.DB) error {
			var msg struct {
				Timestamp int64  `json:"timestamp"`
				Type      string `json:"type"`
				Payload   struct {
					Proposal string `json:"proposal"`
					Choice   int    `json:"choice"`
				} `json:"payload"`
			}
			err = json.Unmarshal([]byte(req.Msg), &msg)
			if err != nil {
				return errs.NewError(err)
			}
			if msg.Type == "" ||
				msg.Timestamp <= 0 ||
				msg.Payload.Proposal == "" ||
				msg.Payload.Choice < 0 {
				return errs.NewError(errs.ErrBadRequest)
			}
			if msg.Timestamp < time.Now().Add(-60*time.Second).Unix() ||
				msg.Timestamp > time.Now().Add(60*time.Second).Unix() {
				return errs.NewError(errs.ErrBadRequest)
			}
			proposal, err := s.pd.First(
				tx,
				map[string][]interface{}{
					"ipfs_hash = ?": []interface{}{msg.Payload.Proposal},
				},
				map[string][]interface{}{},
				[]string{},
			)
			if err != nil {
				return errs.NewError(err)
			}
			if proposal == nil {
				return errs.NewError(errs.ErrBadRequest)
			}
			proposalVote, err = s.pvd.First(
				tx,
				map[string][]interface{}{
					"proposal_id = ?": []interface{}{proposal.ID},
					"address = ?":     []interface{}{req.Address},
					"status = ?":      []interface{}{models.ProposalVoteStatusCreated},
				},
				map[string][]interface{}{},
				[]string{},
			)
			if err != nil {
				return errs.NewError(err)
			}
			if proposalVote != nil {
				return errs.NewError(errs.ErrBadRequest)
			}
			proposalChoice, err := s.pcd.First(
				tx,
				map[string][]interface{}{
					"proposal_id = ?": []interface{}{proposal.ID},
					"choice = ?":      []interface{}{msg.Payload.Choice},
				},
				map[string][]interface{}{},
				[]string{},
			)
			if err != nil {
				return errs.NewError(err)
			}
			if proposalChoice == nil {
				return errs.NewError(errs.ErrBadRequest)
			}
			ipfsData, err := json.Marshal(&req)
			if err != nil {
				return errs.NewError(err)
			}
			ipfsHash, err := s.ifc.UploadString(string(ipfsData))
			if err != nil {
				return errs.NewError(err)
			}
			proposalVote, err = s.pvd.First(
				tx,
				map[string][]interface{}{
					"ipfs_hash = ?": []interface{}{ipfsHash},
				},
				map[string][]interface{}{},
				[]string{},
			)
			if err != nil {
				return errs.NewError(err)
			}
			if proposalVote != nil {
				return errs.NewError(errs.ErrBadRequest)
			}
			proposalVote = &models.ProposalVote{
				Network:          req.Network,
				ProposalID:       proposal.ID,
				ProposalChoiceID: proposalChoice.ID,
				Address:          req.Address,
				Type:             msg.Type,
				Timestamp:        helpers.TimeFromUnix(msg.Timestamp),
				PowerVote:        numeric.BigFloat{*big.NewFloat(0)},
				IpfsHash:         ipfsHash,
				Status:           models.ProposalVoteStatusCreated,
			}
			err = s.pvd.Create(
				tx,
				proposalVote,
			)
			if err != nil {
				return errs.NewError(err)
			}
			proposalChoice, err = s.pcd.FirstByID(
				tx,
				proposalChoice.ID,
				map[string][]interface{}{},
				true,
			)
			if err != nil {
				return errs.NewError(err)
			}
			proposalChoice.PowerVote = numeric.BigFloat{*models.AddBigFloats(&proposalChoice.PowerVote.Float, &proposalVote.PowerVote.Float)}
			err = s.pcd.Save(
				tx,
				proposalChoice,
			)
			if err != nil {
				return errs.NewError(err)
			}
			proposal, err = s.pd.FirstByID(
				tx,
				proposal.ID,
				map[string][]interface{}{},
				true,
			)
			if err != nil {
				return errs.NewError(err)
			}
			proposal.TotalVote = numeric.BigFloat{*models.AddBigFloats(&proposal.TotalVote.Float, &proposalVote.PowerVote.Float)}
			err = s.pd.Save(
				tx,
				proposal,
			)
			if err != nil {
				return errs.NewError(err)
			}
			return nil
		},
	)
	if err != nil {
		return nil, errs.NewError(err)
	}
	return proposalVote, nil
}
