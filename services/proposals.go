package services

import (
	"context"
	"encoding/json"
	"math/big"
	"strings"
	"time"

	"github.com/czConstant/constant-nftylend-api/daos"
	"github.com/czConstant/constant-nftylend-api/errs"
	"github.com/czConstant/constant-nftylend-api/helpers"
	"github.com/czConstant/constant-nftylend-api/models"
	"github.com/czConstant/constant-nftylend-api/serializers"
	"github.com/czConstant/constant-nftylend-api/types/numeric"
	"github.com/jinzhu/gorm"
)

func (s *NftLend) GetProposals(ctx context.Context, statuses []string, page int, limit int) ([]*models.Proposal, uint, error) {
	filters := map[string][]interface{}{}
	if len(statuses) > 0 {
		filters["status in (?)"] = []interface{}{statuses}
	}
	proposals, count, err := s.pd.Find4Page(
		daos.GetDBMainCtx(ctx),
		filters,
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

func (s *NftLend) GetProposalVotes(ctx context.Context, proposalID uint, address string, statuses []string, page int, limit int) ([]*models.ProposalVote, uint, error) {
	filters := map[string][]interface{}{
		"proposal_id = ?": []interface{}{proposalID},
	}
	if address != "" {
		filters["address = ?"] = []interface{}{address}
	}
	if len(statuses) > 0 {
		filters["status in (?)"] = []interface{}{statuses}
	}
	proposalVotes, count, err := s.pvd.Find4Page(
		daos.GetDBMainCtx(ctx),
		filters,
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
			if msg.Payload.End < time.Now().Unix() {
				return errs.NewError(errs.ErrBadRequest)
			}
			if msg.Payload.Start >= msg.Payload.End {
				return errs.NewError(errs.ErrBadRequest)
			}
			choiceType := models.ProposalChoiceType(msg.Payload.Type)
			switch choiceType {
			case models.ProposalChoiceTypeSingleChoice:
				{
				}
			default:
				{
					return errs.NewError(errs.ErrBadRequest)
				}
			}
			pwpToken, err := s.getLendCurrencyBySymbol(
				tx,
				models.SymbolPWPToken,
				req.Network,
			)
			if err != nil {
				return errs.NewError(err)
			}
			pwpBalance, err := s.getEvmClientByNetwork(req.Network).Erc20Balance(
				pwpToken.ContractAddress,
				req.Address,
			)
			if err != nil {
				return errs.NewError(err)
			}
			if pwpBalance.Cmp(big.NewInt(0)) <= 0 {
				return errs.NewError(errs.ErrBadRequest)
			}
			powerVote := models.ConvertWeiToBigFloat(pwpBalance, pwpToken.Decimals)
			if powerVote.Cmp(&pwpToken.ProposalThreshold.Float) < 0 {
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
				Network:           req.Network,
				Address:           req.Address,
				Type:              msg.Type,
				Timestamp:         helpers.TimeFromUnix(msg.Timestamp),
				ChoiceType:        choiceType,
				Msg:               req.Msg,
				Sig:               req.Sig,
				Start:             helpers.TimeFromUnix(msg.Payload.Start),
				End:               helpers.TimeFromUnix(msg.Payload.End),
				Snapshot:          msg.Payload.Snapshot,
				Name:              msg.Payload.Name,
				Body:              msg.Payload.Body,
				IpfsHash:          ipfsHash,
				ProposalThreshold: pwpToken.ProposalThreshold,
				Status:            models.ProposalStatusPending,
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
					Choice   []int  `json:"choice"`
				} `json:"payload"`
			}
			err = json.Unmarshal([]byte(req.Msg), &msg)
			if err != nil {
				return errs.NewError(err)
			}
			if msg.Type == "" ||
				msg.Timestamp <= 0 ||
				msg.Payload.Proposal == "" ||
				len(msg.Payload.Choice) <= 0 {
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
			if proposal.Status != models.ProposalStatusCreated {
				return errs.NewError(errs.ErrBadRequest)
			}
			if proposal.End.Before(time.Now()) {
				return errs.NewError(errs.ErrBadRequest)
			}
			switch proposal.ChoiceType {
			case models.ProposalChoiceTypeSingleChoice:
				{
					if len(msg.Payload.Choice) != 1 {
						return errs.NewError(errs.ErrBadRequest)
					}
				}
			case models.ProposalChoiceTypeMultipleChoice:
				{
				}
			default:
				{
					return errs.NewError(errs.ErrBadRequest)
				}
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
			proposalChoices, err := s.pcd.Find(
				tx,
				map[string][]interface{}{
					"proposal_id = ?": []interface{}{proposal.ID},
					"choice in (?)":   []interface{}{msg.Payload.Choice},
				},
				map[string][]interface{}{},
				[]string{},
				0,
				999999,
			)
			if err != nil {
				return errs.NewError(err)
			}
			if len(proposalChoices) == len(msg.Payload.Choice) {
				return errs.NewError(errs.ErrBadRequest)
			}
			// get power vote
			pwpToken, err := s.getLendCurrencyBySymbol(
				tx,
				models.SymbolPWPToken,
				req.Network,
			)
			if err != nil {
				return errs.NewError(err)
			}
			pwpBalance, err := s.getEvmClientByNetwork(req.Network).Erc20Balance(
				pwpToken.ContractAddress,
				req.Address,
			)
			if err != nil {
				return errs.NewError(err)
			}
			if pwpBalance.Cmp(big.NewInt(0)) <= 0 {
				return errs.NewError(errs.ErrBadRequest)
			}
			powerVote := models.ConvertWeiToBigFloat(pwpBalance, pwpToken.Decimals)
			// end
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
			for _, proposalChoice := range proposalChoices {
				proposalVote = &models.ProposalVote{
					Network:          req.Network,
					ProposalID:       proposal.ID,
					ProposalChoiceID: proposalChoice.ID,
					Address:          req.Address,
					Type:             msg.Type,
					Timestamp:        helpers.TimeFromUnix(msg.Timestamp),
					PowerVote:        numeric.BigFloat{*powerVote},
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
			}
			return nil
		},
	)
	if err != nil {
		return nil, errs.NewError(err)
	}
	return proposalVote, nil
}

func (s *NftLend) ProposalUnVote(ctx context.Context, address string, txHash string, blockNumber uint64) error {
	err := daos.WithTransaction(
		daos.GetDBMainCtx(ctx),
		func(tx *gorm.DB) error {
			proposalVotes, err := s.pvd.Find(
				tx,
				map[string][]interface{}{
					"address = ?": []interface{}{address},
					"status = ?":  []interface{}{models.ProposalVoteStatusCreated},
					`exists(
						select 1
						from proposals
						where proposal_votes.proposal_id = proposals.id
						  and proposals.status = ?
					)`: []interface{}{models.ProposalStatusCreated},
				},
				map[string][]interface{}{},
				[]string{},
				0,
				999999,
			)
			if err != nil {
				return errs.NewError(err)
			}
			for _, proposalVote := range proposalVotes {
				proposalVote, err = s.pvd.FirstByID(
					tx,
					proposalVote.ID,
					map[string][]interface{}{},
					true,
				)
				if err != nil {
					return errs.NewError(err)
				}
				if proposalVote.Status == models.ProposalVoteStatusCreated {
					proposal, err := s.pd.FirstByID(
						tx,
						proposalVote.ProposalID,
						map[string][]interface{}{},
						true,
					)
					if err != nil {
						return errs.NewError(err)
					}
					if proposal.Status == models.ProposalStatusPending {
						proposalVote.CancelledHash = txHash
						proposalVote.Status = models.ProposalVoteStatusCancelled
						err = s.pvd.Save(
							tx,
							proposalVote,
						)
						if err != nil {
							return errs.NewError(err)
						}
						proposalChoice, err := s.pcd.FirstByID(
							tx,
							proposalVote.ProposalChoiceID,
							map[string][]interface{}{},
							true,
						)
						if err != nil {
							return errs.NewError(err)
						}
						proposalChoice.PowerVote = numeric.BigFloat{*models.SubBigFloats(&proposalChoice.PowerVote.Float, &proposalVote.PowerVote.Float)}
						err = s.pcd.Save(
							tx,
							proposalChoice,
						)
						if err != nil {
							return errs.NewError(err)
						}
						proposal.TotalVote = numeric.BigFloat{*models.SubBigFloats(&proposal.TotalVote.Float, &proposalVote.PowerVote.Float)}
						err = s.pd.Save(
							tx,
							proposal,
						)
						if err != nil {
							return errs.NewError(err)
						}
					}
				}
			}
			return nil
		},
	)
	if err != nil {
		return errs.NewError(err)
	}
	return nil
}

func (s *NftLend) JobProposalUnVote(ctx context.Context) error {
	var retErr error
	pwpToken, err := s.GetCurrencyBySymbol(ctx, models.SymbolPWPToken, models.NetworkAURORA)
	if err != nil {
		return errs.NewError(err)
	}
	transferLogs, err := s.bcs.Aurora.Erc20TransferLogs(
		[]string{pwpToken.ContractAddress},
		300,
	)
	if err != nil {
		return errs.NewError(err)
	}
	for _, transferLog := range transferLogs {
		if strings.EqualFold(transferLog.Address, pwpToken.ContractAddress) {
			err = s.ProposalUnVote(
				ctx,
				transferLog.From,
				transferLog.Hash,
				transferLog.BlockNumber,
			)
			if err != nil {
				retErr = errs.MergeError(retErr, err)
			}
		}
	}
	return retErr
}

func (s *NftLend) JobProposalStatus(ctx context.Context) error {
	var retErr error
	proposals, err := s.pd.Find(
		daos.GetDBMainCtx(ctx),
		map[string][]interface{}{
			"status = ?": []interface{}{models.ProposalStatusPending},
			"start <= ?": []interface{}{time.Now()},
		},
		map[string][]interface{}{},
		[]string{},
		0,
		999999,
	)
	if err != nil {
		return errs.NewError(err)
	}
	for _, proposal := range proposals {
		err = s.ProposalStatusCreated(ctx, proposal.ID)
		if err != nil {
			retErr = errs.MergeError(retErr, err)
		}
	}
	proposals, err = s.pd.Find(
		daos.GetDBMainCtx(ctx),
		map[string][]interface{}{
			"status = ?":                       []interface{}{models.ProposalStatusCreated},
			"end <= ?":                         []interface{}{time.Now()},
			"total_vote >= proposal_threshold": []interface{}{},
		},
		map[string][]interface{}{},
		[]string{},
		0,
		999999,
	)
	if err != nil {
		return errs.NewError(err)
	}
	for _, proposal := range proposals {
		err = s.ProposalStatusSucceeded(ctx, proposal.ID)
		if err != nil {
			retErr = errs.MergeError(retErr, err)
		}
	}
	proposals, err = s.pd.Find(
		daos.GetDBMainCtx(ctx),
		map[string][]interface{}{
			"status = ?":                      []interface{}{models.ProposalStatusCreated},
			"end <= ?":                        []interface{}{time.Now()},
			"total_vote < proposal_threshold": []interface{}{},
		},
		map[string][]interface{}{},
		[]string{},
		0,
		999999,
	)
	if err != nil {
		return errs.NewError(err)
	}
	for _, proposal := range proposals {
		err = s.ProposalStatusDefeated(ctx, proposal.ID)
		if err != nil {
			retErr = errs.MergeError(retErr, err)
		}
	}
	proposals, err = s.pd.Find(
		daos.GetDBMainCtx(ctx),
		map[string][]interface{}{
			"status = ?": []interface{}{models.ProposalStatusSucceeded},
			"end <= ?":   []interface{}{time.Now().Add(-2 * 24 * time.Hour)},
		},
		map[string][]interface{}{},
		[]string{},
		0,
		999999,
	)
	if err != nil {
		return errs.NewError(err)
	}
	for _, proposal := range proposals {
		err = s.ProposalStatusQueued(ctx, proposal.ID)
		if err != nil {
			retErr = errs.MergeError(retErr, err)
		}
	}
	return retErr
}

func (s *NftLend) ProposalStatusCreated(ctx context.Context, proposalID uint) error {
	err := daos.WithTransaction(
		daos.GetDBMainCtx(ctx),
		func(tx *gorm.DB) error {
			proposal, err := s.pd.FirstByID(
				tx,
				proposalID,
				map[string][]interface{}{},
				true,
			)
			if err != nil {
				return errs.NewError(err)
			}
			if proposal.Status != models.ProposalStatusPending {
				return errs.NewError(errs.ErrBadRequest)
			}
			if proposal.Start.After(time.Now()) {
				return errs.NewError(errs.ErrBadRequest)
			}
			proposal.Status = models.ProposalStatusCreated
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
		return errs.NewError(err)
	}
	return nil
}

func (s *NftLend) ProposalStatusSucceeded(ctx context.Context, proposalID uint) error {
	err := daos.WithTransaction(
		daos.GetDBMainCtx(ctx),
		func(tx *gorm.DB) error {
			proposal, err := s.pd.FirstByID(
				tx,
				proposalID,
				map[string][]interface{}{},
				true,
			)
			if err != nil {
				return errs.NewError(err)
			}
			if proposal.Status != models.ProposalStatusCreated {
				return errs.NewError(errs.ErrBadRequest)
			}
			if proposal.End.After(time.Now()) {
				return errs.NewError(errs.ErrBadRequest)
			}
			if proposal.TotalVote.Float.Cmp(&proposal.ProposalThreshold.Float) < 0 {
				return errs.NewError(errs.ErrBadRequest)
			}
			proposal.Status = models.ProposalStatusSucceeded
			err = s.pd.Save(
				tx,
				proposal,
			)
			if err != nil {
				return errs.NewError(err)
			}
			proposalChoices, err := s.pcd.Find(
				tx,
				map[string][]interface{}{
					"proposal_id = ?": []interface{}{proposal.ID},
				},
				map[string][]interface{}{},
				[]string{},
				0,
				999999,
			)
			if err != nil {
				return errs.NewError(err)
			}
			var proposalChoiceSucceeded *models.ProposalChoice
			for _, proposalChoice := range proposalChoices {
				proposalChoice, err = s.pcd.FirstByID(
					tx,
					proposalChoice.ID,
					map[string][]interface{}{},
					true,
				)
				if err != nil {
					return errs.NewError(err)
				}
				if proposalChoiceSucceeded == nil {
					proposalChoiceSucceeded = proposalChoice
				} else {
					if proposalChoiceSucceeded.PowerVote.Float.Cmp(&proposalChoice.PowerVote.Float) < 0 {
						proposalChoiceSucceeded = proposalChoice
					}
				}
			}
			for _, proposalChoice := range proposalChoices {
				proposalChoice, err = s.pcd.FirstByID(
					tx,
					proposalChoice.ID,
					map[string][]interface{}{},
					true,
				)
				if err != nil {
					return errs.NewError(err)
				}
				if proposalChoice.ID == proposalChoiceSucceeded.ID {
					proposalChoice.Status = models.ProposalChoiceStatusSucceeded
				} else {
					proposalChoice.Status = models.ProposalChoiceStatusDefeated
				}
				err = s.pcd.Save(
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
		return errs.NewError(err)
	}
	return nil
}

func (s *NftLend) ProposalStatusDefeated(ctx context.Context, proposalID uint) error {
	err := daos.WithTransaction(
		daos.GetDBMainCtx(ctx),
		func(tx *gorm.DB) error {
			proposal, err := s.pd.FirstByID(
				tx,
				proposalID,
				map[string][]interface{}{},
				true,
			)
			if err != nil {
				return errs.NewError(err)
			}
			if proposal.Status != models.ProposalStatusCreated {
				return errs.NewError(errs.ErrBadRequest)
			}
			if proposal.End.After(time.Now()) {
				return errs.NewError(errs.ErrBadRequest)
			}
			if proposal.TotalVote.Float.Cmp(&proposal.ProposalThreshold.Float) >= 0 {
				return errs.NewError(errs.ErrBadRequest)
			}
			proposal.Status = models.ProposalStatusDefeated
			err = s.pd.Save(
				tx,
				proposal,
			)
			if err != nil {
				return errs.NewError(err)
			}
			proposalChoices, err := s.pcd.Find(
				tx,
				map[string][]interface{}{
					"proposal_id = ?": []interface{}{proposal.ID},
				},
				map[string][]interface{}{},
				[]string{},
				0,
				999999,
			)
			if err != nil {
				return errs.NewError(err)
			}
			for _, proposalChoice := range proposalChoices {
				proposalChoice, err = s.pcd.FirstByID(
					tx,
					proposalChoice.ID,
					map[string][]interface{}{},
					true,
				)
				if err != nil {
					return errs.NewError(err)
				}
				proposalChoice.Status = models.ProposalChoiceStatusDefeated
				err = s.pcd.Save(
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
		return errs.NewError(err)
	}
	return nil
}

func (s *NftLend) ProposalStatusQueued(ctx context.Context, proposalID uint) error {
	err := daos.WithTransaction(
		daos.GetDBMainCtx(ctx),
		func(tx *gorm.DB) error {
			proposal, err := s.pd.FirstByID(
				tx,
				proposalID,
				map[string][]interface{}{},
				true,
			)
			if err != nil {
				return errs.NewError(err)
			}
			if proposal.Status != models.ProposalStatusSucceeded {
				return errs.NewError(errs.ErrBadRequest)
			}
			proposal.Status = models.ProposalStatusQueued
			err = s.pd.Save(
				tx,
				proposal,
			)
			if err != nil {
				return errs.NewError(err)
			}
			proposalChoice, err := s.pcd.First(
				tx,
				map[string][]interface{}{
					"proposal_id = ?": []interface{}{proposal.ID},
					"status = ?":      []interface{}{models.ProposalChoiceStatusSucceeded},
				},
				map[string][]interface{}{},
				[]string{},
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
			proposalChoice.Status = models.ProposalChoiceStatusQueued
			err = s.pcd.Save(
				tx,
				proposalChoice,
			)
			if err != nil {
				return errs.NewError(err)
			}
			return nil
		},
	)
	if err != nil {
		return errs.NewError(err)
	}
	return nil
}
