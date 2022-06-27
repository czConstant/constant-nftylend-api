package services

import (
	"context"

	"github.com/czConstant/blockchain-api/bcclient"
	"github.com/czConstant/blockchain-api/bcclient/ethereum"
	"github.com/czConstant/constant-nftylend-api/configs"
	"github.com/czConstant/constant-nftylend-api/daos"
	"github.com/czConstant/constant-nftylend-api/errs"
	"github.com/czConstant/constant-nftylend-api/models"
	"github.com/czConstant/constant-nftylend-api/services/3rd/ipfs"
	"github.com/czConstant/constant-nftylend-api/services/3rd/moralis"
	"github.com/czConstant/constant-nftylend-api/services/3rd/saletrack"
	"github.com/jinzhu/gorm"
)

type NftLend struct {
	conf *configs.Config
	bcs  *bcclient.Client
	stc  *saletrack.Client
	ifc  *ipfs.Client
	mc   *moralis.Client
	ud   *daos.User
	cd   *daos.Currency
	cld  *daos.Collection
	clsd *daos.CollectionSubmitted
	ad   *daos.Asset
	atd  *daos.AssetTransaction
	ld   *daos.Loan
	lod  *daos.LoanOffer
	ltd  *daos.LoanTransaction
	id   *daos.Instruction
	pd   *daos.Proposal
	pcd  *daos.ProposalChoice
	pvd  *daos.ProposalVote
	ntd  *daos.NotificationTemplate
	nd   *daos.Notification

	// for incentive
	ubd  *daos.UserBalance
	ubtd *daos.UserBalanceTransaction
	ubhd *daos.UserBalanceHistory
	ipd  *daos.IncentiveProgram
	ipdd *daos.IncentiveProgramDetail
	itd  *daos.IncentiveTransaction

	vd *daos.Verification
}

func NewNftLend(
	conf *configs.Config,
	bcs *bcclient.Client,
	stc *saletrack.Client,
	ifc *ipfs.Client,
	mc *moralis.Client,
	ud *daos.User,
	cd *daos.Currency,
	cld *daos.Collection,
	clsd *daos.CollectionSubmitted,
	ad *daos.Asset,
	atd *daos.AssetTransaction,
	ld *daos.Loan,
	lod *daos.LoanOffer,
	ltd *daos.LoanTransaction,
	id *daos.Instruction,
	pd *daos.Proposal,
	pcd *daos.ProposalChoice,
	pvd *daos.ProposalVote,
	ntd *daos.NotificationTemplate,
	nd *daos.Notification,

	// for incentive
	ubd *daos.UserBalance,
	ubtd *daos.UserBalanceTransaction,
	ubhd *daos.UserBalanceHistory,
	ipd *daos.IncentiveProgram,
	ipdd *daos.IncentiveProgramDetail,
	itd *daos.IncentiveTransaction,

	vd *daos.Verification,

) *NftLend {
	s := &NftLend{
		conf: conf,
		bcs:  bcs,
		stc:  stc,
		ifc:  ifc,
		mc:   mc,
		ud:   ud,
		cd:   cd,
		cld:  cld,
		clsd: clsd,
		ad:   ad,
		atd:  atd,
		ld:   ld,
		lod:  lod,
		ltd:  ltd,
		id:   id,
		pd:   pd,
		pcd:  pcd,
		pvd:  pvd,
		ntd:  ntd,
		nd:   nd,

		// for incentive
		ubd:  ubd,
		ubtd: ubtd,
		ubhd: ubhd,
		ipd:  ipd,
		ipdd: ipdd,
		itd:  itd,
		vd:   vd,
	}
	if s.conf.Contract.ProgramID != "" {
		go stc.StartWssSolsea(s.solseaMsgReceived)
	}
	return s
}

func (s *NftLend) getEvmClientByNetwork(network models.Network) *ethereum.Client {
	switch network {
	case models.NetworkETH:
		{
			return s.bcs.Ethereum
		}
	case models.NetworkMATIC:
		{
			return s.bcs.Matic
		}
	case models.NetworkAVAX:
		{
			return s.bcs.Avax
		}
	case models.NetworkBSC:
		{
			return s.bcs.BSC
		}
	case models.NetworkBOBA:
		{
			return s.bcs.Boba
		}
	case models.NetworkAURORA:
		{
			return s.bcs.Aurora
		}
	}
	return nil
}

func (s *NftLend) getEvmAdminFee(network models.Network) int64 {
	switch network {
	case models.NetworkETH:
		{
			return 100
		}
	case models.NetworkMATIC:
		{
			return 100
		}
	case models.NetworkAVAX:
		{
			return 100
		}
	case models.NetworkBSC:
		{
			return 100
		}
	case models.NetworkBOBA:
		{
			return 100
		}
	}
	return 0
}

func (s *NftLend) getSupportedNetworks() []models.Network {
	ns := []models.Network{}
	if s.conf.Contract.ProgramID != "" {
		ns = append(ns, models.NetworkSOL)
	}
	if s.conf.Contract.BscNftypawnAddress != "" {
		ns = append(ns, models.NetworkBSC)
	}
	if s.conf.Contract.AvaxNftypawnAddress != "" {
		ns = append(ns, models.NetworkAVAX)
	}
	if s.conf.Contract.MaticNftypawnAddress != "" {
		ns = append(ns, models.NetworkMATIC)
	}
	if s.conf.Contract.BobaNftypawnAddress != "" {
		ns = append(ns, models.NetworkBOBA)
	}
	if s.conf.Contract.NearNftypawnAddress != "" {
		ns = append(ns, models.NetworkNEAR)
	}
	return ns
}

func (s *NftLend) getEvmContractAddress(network models.Network) string {
	switch network {
	case models.NetworkETH:
		{
			return ""
		}
	case models.NetworkMATIC:
		{
			return s.conf.Contract.MaticNftypawnAddress
		}
	case models.NetworkAVAX:
		{
			return s.conf.Contract.AvaxNftypawnAddress
		}
	case models.NetworkBSC:
		{
			return s.conf.Contract.BscNftypawnAddress
		}
	case models.NetworkBOBA:
		{
			return s.conf.Contract.BobaNftypawnAddress
		}
	case models.NetworkNEAR:
		{
			return s.conf.Contract.NearNftypawnAddress
		}
	}
	return ""
}

func (s *NftLend) getLendCurrency(tx *gorm.DB, address string) (*models.Currency, error) {
	c, err := s.cd.First(
		tx,
		map[string][]interface{}{
			"contract_address = ?": []interface{}{address},
		},
		map[string][]interface{}{},
		[]string{},
	)
	if err != nil {
		return nil, errs.NewError(err)
	}
	if c == nil {
		return nil, errs.NewError(errs.ErrBadRequest)
	}
	return c, nil
}

func (s *NftLend) GetCurrencyByID(tx *gorm.DB, id uint, chain models.Network) (*models.Currency, error) {
	c, err := s.cd.First(
		tx,
		map[string][]interface{}{
			"id = ?":      []interface{}{id},
			"network = ?": []interface{}{chain},
		},
		map[string][]interface{}{},
		[]string{},
	)
	if err != nil {
		return nil, errs.NewError(err)
	}
	if c == nil {
		return nil, errs.NewError(errs.ErrBadRequest)
	}
	return c, nil
}

func (s *NftLend) getCurrencyByNetworkSymbol(tx *gorm.DB, network models.Network, symbol string) (*models.Currency, error) {
	c, err := s.cd.First(
		tx,
		map[string][]interface{}{
			"network = ?": []interface{}{network},
			"symbol = ?":  []interface{}{symbol},
		},
		map[string][]interface{}{},
		[]string{},
	)
	if err != nil {
		return nil, errs.NewError(err)
	}
	if c == nil {
		return nil, errs.NewError(errs.ErrBadRequest)
	}
	return c, nil
}

func (s *NftLend) GetAssetDetailInfo(ctx context.Context, contractAddress string, tokenID string) (*models.Asset, error) {
	m, err := s.ad.First(
		daos.GetDBMainCtx(ctx),
		map[string][]interface{}{
			"network in (?)":       []interface{}{s.getSupportedNetworks()},
			"contract_address = ?": []interface{}{contractAddress},
			"token_id = ?":         []interface{}{tokenID},
		},
		map[string][]interface{}{
			"Collection": []interface{}{},
			"NewLoan": []interface{}{
				"status in (?)",
				[]models.LoanStatus{
					models.LoanStatusNew,
					models.LoanStatusCreated,
				},
			},
			"NewLoan.Currency": []interface{}{},
			"NewLoan.Offers": []interface{}{
				func(db *gorm.DB) *gorm.DB {
					return db.Order("loan_offers.id DESC")
				},
			},
			"NewLoan.ApprovedOffer": []interface{}{
				"status = ?",
				models.LoanOfferStatusApproved,
			},
		},
		[]string{"id desc"},
	)
	if err != nil {
		return nil, errs.NewError(err)
	}
	return m, nil
}

func (s *NftLend) GetCurrencyBySymbol(ctx context.Context, symbol string, network models.Network) (*models.Currency, error) {
	m, err := s.cd.First(
		daos.GetDBMainCtx(ctx),
		map[string][]interface{}{
			"symbol = ?":  []interface{}{symbol},
			"network = ?": []interface{}{network},
		},
		map[string][]interface{}{},
		[]string{},
	)
	if err != nil {
		return nil, errs.NewError(err)
	}
	return m, nil
}

func (s *NftLend) GetCurrencies(ctx context.Context, network models.Network) ([]*models.Currency, error) {
	currencies, err := s.cd.Find(
		daos.GetDBMainCtx(ctx),
		map[string][]interface{}{
			"network = ?": []interface{}{network},
			"enabled = ?": []interface{}{true},
		},
		map[string][]interface{}{},
		[]string{"id desc"},
		0,
		99999999,
	)
	if err != nil {
		return nil, errs.NewError(err)
	}
	return currencies, nil
}

func (s *NftLend) VerifyAddressSignature(ctx context.Context, network models.Network, address string, message string, signature string) error {
	switch network {
	case models.NetworkNEAR:
		{
			err := s.bcs.Near.ValidateMessageSignature(
				s.getEvmContractAddress(network),
				message,
				signature,
				address,
			)
			if err != nil {
				return errs.NewError(err)
			}
		}
	default:
		{
			return errs.NewError(errs.ErrBadRequest)
		}
	}
	err := daos.WithTransaction(
		daos.GetDBMainCtx(ctx),
		func(tx *gorm.DB) error {
			user, err := s.getUser(
				tx,
				network,
				address,
				true,
			)
			if err != nil {
				return errs.NewError(err)
			}
			err = s.ud.Save(
				tx,
				user,
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

func (s *NftLend) VerifyUserTimestamp(ctx context.Context, network models.Network, address string, timestamp int64) error {
	err := daos.WithTransaction(
		daos.GetDBMainCtx(ctx),
		func(tx *gorm.DB) error {
			user, err := s.getUser(
				tx,
				network,
				address,
				true,
			)
			if err != nil {
				return errs.NewError(err)
			}
			if user.TimestampReq >= timestamp {
				return errs.NewError(errs.ErrBadRequest)
			}
			user.TimestampReq = timestamp
			err = s.ud.Save(
				tx,
				user,
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
