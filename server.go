package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime/debug"

	"github.com/getsentry/raven-go"
	"github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"

	"github.com/czConstant/blockchain-api/bcclient"
	"github.com/czConstant/constant-nftylend-api/apis"
	"github.com/czConstant/constant-nftylend-api/configs"
	"github.com/czConstant/constant-nftylend-api/daos"
	"github.com/czConstant/constant-nftylend-api/databases"
	"github.com/czConstant/constant-nftylend-api/logger"
	"github.com/czConstant/constant-nftylend-api/services"
	"github.com/czConstant/constant-nftylend-api/services/3rd/ipfs"
	"github.com/czConstant/constant-nftylend-api/services/3rd/kitwallet"
	"github.com/czConstant/constant-nftylend-api/services/3rd/mailer"
	"github.com/czConstant/constant-nftylend-api/services/3rd/moralis"
	"github.com/czConstant/constant-nftylend-api/services/3rd/saletrack"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	sqltrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/database/sql"
	gintrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/gin-gonic/gin"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	"gopkg.in/DataDog/dd-trace-go.v1/profiler"
)

func init() {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
}

func main() {
	conf := configs.GetConfig()
	logger.NewLogger("nftpawn-api", conf.Env, conf.LogPath, true)
	defer logger.Sync()
	raven.SetDSN(conf.RavenDNS)
	raven.SetEnvironment(conf.RavenENV)
	defer func() {
		if err := recover(); err != nil {
			panicErr := errors.Wrap(errors.New("panic start server"), string(debug.Stack()))
			logger.Info(
				logger.LOGGER_API_APP_PANIC,
				"panic start server",
				zap.Error(panicErr),
			)
			fmt.Println(err)
			fmt.Println(panicErr)
			return
		}
	}()
	// for dd-trace
	tracer.Start(
		tracer.WithEnv(conf.Datadog.Env),
		tracer.WithService(conf.Datadog.Service),
		tracer.WithServiceVersion(conf.Datadog.Version),
	)
	err := profiler.Start(
		profiler.WithEnv(conf.Datadog.Env),
		profiler.WithService(conf.Datadog.Service),
		profiler.WithVersion(conf.Datadog.Version),
	)
	if err != nil {
		log.Fatal(err)
	}

	sqltrace.Register("mysql", &mysql.MySQLDriver{}, sqltrace.WithServiceName(fmt.Sprintf("%s-gorm", conf.Datadog.Service)))

	defer profiler.Stop()
	defer tracer.Stop()

	var migrateDBMainFunc func(db *gorm.DB) error
	if os.Getenv("DEV") != "true" {
		migrateDBMainFunc = databases.MigrateDBMain
	}
	dbMain, err := databases.Init(
		conf.DbURL,
		migrateDBMainFunc,
		10,
		20,
		conf.Debug,
	)
	if err != nil {
		panic(err)
	}
	daos.InitDBConn(
		dbMain,
	)
	mailer.SetURL(conf.Mailer.URL)
	var (
		bcs = bcclient.NewBlockchainClient(
			conf.Blockchain,
		)
		ud   = &daos.User{}
		cd   = &daos.Currency{}
		cld  = &daos.Collection{}
		clsd = &daos.CollectionSubmission{}
		ad   = &daos.Asset{}
		atd  = &daos.AssetTransaction{}
		ld   = &daos.Loan{}
		lod  = &daos.LoanOffer{}
		ltd  = &daos.LoanTransaction{}
		id   = &daos.Instruction{}
		pd   = &daos.Proposal{}
		pcd  = &daos.ProposalChoice{}
		pvd  = &daos.ProposalVote{}
		ntd  = &daos.NotificationTemplate{}
		nd   = &daos.Notification{}

		ubd  = &daos.UserBalance{}
		ubtd = &daos.UserBalanceTransaction{}
		ubhd = &daos.UserBalanceHistory{}
		ipd  = &daos.IncentiveProgram{}
		ipdd = &daos.IncentiveProgramDetail{}
		itd  = &daos.IncentiveTransaction{}

		vd = &daos.Verification{}

		asd = &daos.AffiliateSubmission{}

		lbd = &daos.Leaderboard{}

		stc = &saletrack.Client{
			NftbankKey: conf.SaleTrack.NftbankKey,
		}
		mc = &moralis.Client{
			APIKey: conf.Moralis.APIKey,
		}

		ifc = &ipfs.Client{
			URL:       conf.Ipfs.URL,
			BasicAuth: conf.Ipfs.BasicAuth,
		}

		kwc = &kitwallet.Client{
			URL:     conf.Kitwallet.URL,
			Origin:  conf.Kitwallet.Origin,
			Referer: conf.Kitwallet.Referer,
		}

		s = services.NewNftLend(
			conf,
			bcs,
			stc,
			ifc,
			mc,
			kwc,
			ud,
			cd,
			cld,
			clsd,
			ad,
			atd,
			ld,
			lod,
			ltd,
			id,
			pd,
			pcd,
			pvd,
			ntd,
			nd,

			ubd,
			ubtd,
			ubhd,
			ipd,
			ipdd,
			itd,

			vd,
			asd,

			lbd,
		)
	)

	r := gin.Default()
	r.Use(gintrace.Middleware(fmt.Sprintf("%s-gin", conf.Datadog.Service), gintrace.WithAnalytics(true)))
	srv := apis.NewServer(
		r,
		conf,
		s,
	)
	srv.Routers()
	if conf.Port == 0 {
		conf.Port = 8080
	}
	if err := r.Run(fmt.Sprintf(":%d", conf.Port)); err != nil {
		logger.WrapError(
			logger.LOGGER_API_APP_ERROR,
			err,
		)
	}
}
