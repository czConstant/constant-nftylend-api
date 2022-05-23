package errs

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/czConstant/constant-nftylend-api/helpers"
	"github.com/czConstant/constant-nftylend-api/logger"
	"github.com/getsentry/raven-go"
	"go.uber.org/zap"
)

var (
	ErrSystemError              = &Error{Code: -1001, Message: "system error"}
	ErrInvalidEmail             = &Error{Code: -1002, Message: "invalid email"}
	ErrInvalidPassword          = &Error{Code: -1003, Message: "invalid password"}
	ErrEmailNotExists           = &Error{Code: -1004, Message: "email doesn't exist"}
	ErrEmailAlreadyExists       = &Error{Code: -1005, Message: "email already exists"}
	ErrInvalidCredentials       = &Error{Code: -1006, Message: "invalid credentials"}
	ErrBadRequest               = &Error{Code: -1007, Message: "bad request"}
	ErrBadPermission            = &Error{Code: -1008, Message: "bad permission"}
	ErrBadBodyRequest           = &Error{Code: -1009, Message: "bad body request"}
	ErrVerificationTokenExpired = &Error{Code: -1010, Message: "verification token expired"}
	ErrOTPIsInvalid             = &Error{Code: -1045, Message: "OTP not matched or invalidated!"}
	ErrUserNotFound             = &Error{Code: -1046, Message: "User not found"}
	ErrInvalidRecaptcha         = &Error{Code: -1076, Message: "invalid recaptcha"}

	ErrAddressInvalid          = &Error{Code: -333001, Message: "address invalid"}
	ErrTokenInvalid            = &Error{Code: -333002, Message: "token invalid"}
	ErrNetworkInvalid          = &Error{Code: -333003, Message: "Network invalid"}
	ErrStrategyNotFound        = &Error{Code: -333004, Message: "Strategy not found"}
	ErrTeamNotFound            = &Error{Code: -333005, Message: "Team not found"}
	ErrCurrencyNotFound        = &Error{Code: -333006, Message: "Currency not found"}
	ErrOrderNotFound           = &Error{Code: -333007, Message: "Order not found"}
	ErrContractNotFound        = &Error{Code: -333008, Message: "Contract not found"}
	ErrContractAlreadyReviewed = &Error{Code: -333009, Message: "Contract has reviewed"}
	ErrQueueNotFound           = &Error{Code: -333008, Message: "Contract queue not found"}
	ErrRatingInvalid           = &Error{Code: -333010, Message: "Invalid rating"}
	ErrPlayerNotFound          = &Error{Code: -333011, Message: "Player not found"}

	ErrPriceOutOfDate = &Error{Code: -9036, Message: "price is out of date"}
)

type Error struct {
	Code       int    `json:"code"`
	Message    string `json:"message"`
	stacktrace string
	extra      []interface{}
}

func (e *Error) SetStacktrace(stacktrace string) {
	e.stacktrace = stacktrace
}

func (e *Error) Stacktrace() string {
	return e.stacktrace
}

func (e *Error) Error() string {
	return e.Message
}

func (e *Error) SetExtra(extra []interface{}) {
	e.extra = extra
}

func (e *Error) Extra() []interface{} {
	return e.extra
}

func (e *Error) ExtraJson() string {
	return helpers.ConvertJsonString(e.extra)
}

func NewErrorWithId(err error, id interface{}) error {
	if err != nil {
		msg := err.Error()
		err = NewError(err)
		err.(*Error).Message = fmt.Sprintf("%v : %s", id, msg)
	}
	return err
}

func NewError(err error, extras ...interface{}) error {
	if err == nil {
		return nil
	}
	_, ok := err.(*Error)
	if ok {
		sterr := err.(*Error).Stacktrace()
		retErr := &Error{
			Code:    err.(*Error).Code,
			Message: err.(*Error).Message,
		}
		if sterr == "" {
			retErr.SetStacktrace(fmt.Sprintf("%s\n\n%s", err.Error(), NewStacktraceString(extras...)))
			err.(*Error).SetExtra(extras)
		} else {
			retErr.SetStacktrace(sterr)
		}
		return retErr
	}
	retErr := &Error{
		Code:    ErrSystemError.Code,
		Message: err.Error(),
	}
	retErr.SetStacktrace(fmt.Sprintf("%s\n\n%s", err.Error(), NewStacktraceString(extras...)))
	return retErr
}

func NewStacktraceString(extras ...interface{}) string {
	var rets []string
	if len(extras) > 0 {
		rets = append(rets, fmt.Sprintf("Extras -> %s", helpers.ConvertJsonString(extras)))
	}
	st := raven.NewStacktrace(1, 3, nil)
	for i := len(st.Frames) - 1; i >= 0; i-- {
		frame := st.Frames[i]
		if strings.TrimSpace(frame.Filename) != "" {
			rets = append(rets, fmt.Sprintf("%s\t%s\t%d", frame.Filename, frame.Function, frame.Lineno))
			rets = append(rets, fmt.Sprintf("\t%s", strings.Join(frame.PreContext, "\n\t")))
			rets = append(rets, fmt.Sprintf("%d.\t%s", frame.Lineno, frame.ContextLine))
			rets = append(rets, fmt.Sprintf("\t%s", strings.Join(frame.PostContext, "\n\t")))
		}
	}
	return strings.Join(rets, "\n")
}

func MergeError(err1 error, errss ...error) error {
	var msgs, sterrs []string
	if err1 != nil {
		err1 = NewError(err1)
		_, ok := err1.(*Error)
		if ok {
			msgs = append(msgs, strings.TrimSpace(err1.Error()))
			sterrs = append(sterrs,
				err1.(*Error).Stacktrace(),
			)
		}
	}
	for _, err := range errss {
		if err != nil {
			err = NewError(err)
			_, ok := err.(*Error)
			if ok {
				msgs = append(msgs, strings.TrimSpace(err.Error()))
				sterrs = append(sterrs,
					fmt.Sprintf(
						"------------------------------------------------------------------------------------------------------------------------------------\n\n%s\n\n%s",
						strings.TrimSpace(err.Error()),
						err.(*Error).Stacktrace()),
				)
			}
		}
	}
	if len(msgs) <= 0 {
		return nil
	}
	err := &Error{
		Code:    ErrSystemError.Code,
		Message: strings.Join(msgs, "\n"),
	}
	err.SetStacktrace(
		strings.Join(
			sterrs,
			"\n\n",
		),
	)
	return err
}

func LoggerFunc(fn func() error, path string, userID uint, email string, extras ...interface{}) {
	var err error
	start := time.Now()
	defer func() {
		end := time.Now()
		latency := end.Sub(start).Seconds()
		if rval := recover(); rval != nil {
			if rval := recover(); rval != nil {
				err = NewError(errors.New(fmt.Sprint(rval)))
			}
		}
		if path == "" {
			path = "default"
		}
		path = fmt.Sprintf("nft-marketet-api-fun-%s", path)
		var stacktrace, errText string
		errCode := 200
		if err != nil {
			errCode = 400
			err = NewError(err)
			errText = err.Error()
			retErr, ok := err.(*Error)
			if ok {
				stacktrace = retErr.Stacktrace()
			}
		}
		logger.Info(
			"logger_func_error",
			"msg info",
			zap.Any("referer", ""),
			zap.Any("ip", ""),
			zap.Any("method", "FUN"),
			zap.Any("path", path),
			zap.Any("raw_query", ""),
			zap.Any("latency", latency),
			zap.Any("status", errCode),
			zap.Any("user_agent", ""),
			zap.Any("platform", ""),
			zap.Any("os", ""),
			zap.Any("country", ""),
			zap.Any("email", email),
			zap.Any("user_id", userID),
			zap.Any("error_text", errText),
			zap.Any("stacktrace", stacktrace),
			zap.Any("body_request", helpers.ConvertJsonString(extras)),
			zap.Any("body_response", ""),
		)
		if os.Getenv("DEV") == "true" {
			if stacktrace != "" {
				fmt.Println(stacktrace)
			}
		}
	}()
	err = fn()
}
