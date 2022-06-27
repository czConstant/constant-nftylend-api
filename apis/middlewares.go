package apis

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/czConstant/constant-nftylend-api/configs"
	"github.com/czConstant/constant-nftylend-api/errs"
	"github.com/czConstant/constant-nftylend-api/helpers"
	"github.com/czConstant/constant-nftylend-api/logger"
	"github.com/czConstant/constant-nftylend-api/models"
	"github.com/czConstant/constant-nftylend-api/serializers"
	"github.com/getsentry/raven-go"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

const (
	CONTEXT_USER_DATA       = "context_user_data"
	CONTEXT_ERROR_DATA      = "context_error_data"
	CONTEXT_STACKTRACE_DATA = "context_stacktrace_data"
)

func ctxJSON(c *gin.Context, respCode int, resp interface{}) {
	if respCode != http.StatusOK {
		WrapRespError(c, resp)
	}
	c.JSON(respCode, resp)
}

func ctxSTRING(c *gin.Context, respCode int, resp string) {
	if respCode != http.StatusOK {
		WrapRespError(c, resp)
	}
	c.String(respCode, resp)
}

func ctxData(c *gin.Context, respCode int, contentType string, resp []byte) {
	if respCode != http.StatusOK {
		WrapRespError(c, resp)
	}
	c.Data(respCode, contentType, resp)
}

func ctxAbortWithStatusJSON(c *gin.Context, respCode int, resp interface{}) {
	if respCode != http.StatusOK {
		WrapRespError(c, resp)
	}
	c.AbortWithStatusJSON(respCode, resp)
}

func WrapRespError(c *gin.Context, resp interface{}) {
	var retErr *errs.Error
	switch resp.(type) {
	case *serializers.Resp:
		{
			retResp, ok := resp.(*serializers.Resp)
			if ok &&
				retResp.Error != nil {
				retResp.Error = errs.NewError(retResp.Error)
				retErr, _ = retResp.Error.(*errs.Error)
			}
		}
	}
	if retErr != nil {
		c.Set(CONTEXT_ERROR_DATA, retErr.Error())
		c.Set(CONTEXT_STACKTRACE_DATA, retErr.Stacktrace())
	}
}

func (s *Server) requestContext(c *gin.Context) context.Context {
	return c.Request.Context()
}

func (s *Server) GetUserToken(c *gin.Context) (string, error) {
	auth := c.GetHeader("Authorization")
	if auth == "" {
		auth = c.Query("auth_token")
		if auth == "" {
			return "", errs.NewError(errs.ErrTokenInvalid)
		}
		return auth, nil
	}
	auths := strings.Split(auth, " ")
	if len(auths) < 2 {
		return "", errs.NewError(errs.ErrTokenInvalid)
	}
	return auths[1], nil
}

func (s *Server) authorizeJobMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if configs.GetConfig().JobToken != "" &&
			c.GetHeader("Authorization") != configs.GetConfig().JobToken {
			ctxAbortWithStatusJSON(c, http.StatusBadRequest, &serializers.Resp{Error: errs.NewError(errs.ErrBadRequest)})
			return
		}
		c.Next()
	}
}

func (s *Server) recaptchaV3Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if configs.GetConfig().RecaptchaV3Serect != "" {
			recaptcha := c.GetHeader("recaptcha")
			if recaptcha == "" {
				recaptcha = c.GetHeader("x-recaptcha")
				if recaptcha == "" {
					ctxAbortWithStatusJSON(c, http.StatusBadRequest, &serializers.Resp{Error: errs.NewError(errs.ErrInvalidRecaptcha)})
					return
				}
			}
			ok, err := helpers.ValidateRecaptcha(configs.GetConfig().RecaptchaV3Serect, recaptcha)
			if err != nil {
				ctxAbortWithStatusJSON(c, http.StatusBadRequest, &serializers.Resp{Error: errs.NewError(errs.ErrInvalidRecaptcha)})
				return
			}
			if !ok {
				ctxAbortWithStatusJSON(c, http.StatusBadRequest, &serializers.Resp{Error: errs.NewError(errs.ErrInvalidRecaptcha)})
				return
			}
		}
		c.Next()
	}
}

func (s *Server) otpFromContext(c *gin.Context) string {
	myOtp := c.GetHeader("OTP")
	return myOtp
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}
func (w bodyLogWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

func (s *Server) loggerDisabledBodyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("log_body", false)
		c.Next()
	}
}

func (s *Server) logApiMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("log", true)
		c.Set("log_body", true)
		start := time.Now()
		bodyLogWriter := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = bodyLogWriter
		var bodyRequest string
		if (c.Request.Method == http.MethodPost ||
			c.Request.Method == http.MethodPut) &&
			strings.LastIndex(strings.ToLower(c.GetHeader("content-type")), "application/json") >= 0 {
			buf, bodyErr := ioutil.ReadAll(c.Request.Body)
			if bodyErr == nil {
				rdr2 := ioutil.NopCloser(bytes.NewBuffer(buf))
				bodyRequest = string(buf)
				c.Request.Body = rdr2
			}
		}
		c.Next()
		if c.GetBool("log") {
			end := time.Now()
			latency := end.Sub(start).Seconds()
			ipStr := c.Request.Header.Get("ip")
			if ipStr == "" {
				ipStr = c.ClientIP()
			}
			var errText, stacktraceText, bodyResponse string
			v, ok := c.Get(CONTEXT_ERROR_DATA)
			if ok {
				errText = v.(string)
			}
			v, ok = c.Get(CONTEXT_STACKTRACE_DATA)
			if ok {
				stacktraceText = v.(string)
				fmt.Println(stacktraceText)
			}
			bodyResponse = bodyLogWriter.body.String()
			if !c.GetBool("log_body") {
				bodyRequest = ""
			}
			logger.Info(
				"api_response_time",
				"request info",
				zap.Any("referer", c.Request.Referer()),
				zap.Any("ip", ipStr),
				zap.Any("method", c.Request.Method),
				zap.Any("path", c.Request.URL.Path),
				zap.Any("raw_query", c.Request.URL.RawQuery),
				zap.Any("latency", latency),
				zap.Any("status", c.Writer.Status()),
				zap.Any("user_agent", c.Request.UserAgent()),
				zap.Any("platform", c.Request.Header.Get("platform")),
				zap.Any("os", c.Request.Header.Get("os")),
				zap.Any("country", c.Request.Header.Get("country")),
				zap.Any("error_text", errText),
				zap.Any("stacktrace", stacktraceText),
				zap.Any("body_request", helpers.SubStringBodyResponse(bodyRequest, 1000)),
				zap.Any("body_response", helpers.SubStringBodyResponse(bodyResponse, 1000)),
			)
			if os.Getenv("DEV") == "true" {
				fmt.Println(stacktraceText)
			}
		}
	}
}

func (s *Server) recoveryMiddleware(client *raven.Client, onlyCrashes bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			flags := map[string]string{
				"endpoint": c.Request.RequestURI,
			}
			if rval := recover(); rval != nil {
				rvalStr := fmt.Sprint(rval)
				client.CaptureMessage(rvalStr, flags, raven.NewException(errors.New(rvalStr), raven.NewStacktrace(2, 3, nil)),
					raven.NewHttp(c.Request))
				ctxAbortWithStatusJSON(c, http.StatusInternalServerError, &serializers.Resp{
					Result: nil,
					Error:  errs.NewError(errors.New(rvalStr)),
				})
			}
			if !onlyCrashes {
				for _, item := range c.Errors {
					client.CaptureMessage(item.Error(), flags, &raven.Message{
						Message: item.Error(),
						Params:  []interface{}{item.Meta},
					},
						raven.NewHttp(c.Request))
				}
			}
		}()
		c.Next()
	}
}

func (s *Server) pagingFromContext(c *gin.Context) (int, int) {
	var (
		pageS  = c.DefaultQuery("page", "1")
		limitS = c.DefaultQuery("limit", "30")
		page   int
		limit  int
		err    error
	)

	page, err = strconv.Atoi(pageS)
	if err != nil {
		page = 1
	}

	limit, err = strconv.Atoi(limitS)
	if err != nil {
		limit = 10
	}

	if limit > 500 {
		limit = 500
	}

	return page, limit
}

func (s *Server) ipfsProxyMiddleware(hostPath string) gin.HandlerFunc {
	return func(c *gin.Context) {
		hash := c.Param("hash")
		r := c.Request
		w := c.Writer
		director := func(req *http.Request) {
			hostURL, err := url.Parse(hostPath)
			if err != nil {
				ctxAbortWithStatusJSON(c, http.StatusBadRequest, &serializers.Resp{Error: errs.NewError(err)})
				return
			}
			req.URL.Scheme = hostURL.Scheme
			req.URL.Host = hostURL.Host
			req.Host = hostURL.Host
			req.URL.Path = hostPath + "/" + hash + ".mp4"
		}
		proxy := &httputil.ReverseProxy{
			Director: director,
		}
		proxy.ServeHTTP(w, r)
	}
}

func (s *Server) uintFromContextParam(c *gin.Context, param string) (uint, error) {
	val := strings.TrimSpace(c.Param(param))
	if val == "" {
		return uint(0), nil
	}
	num, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(num), nil
}

func (s *Server) uint64FromContextParam(c *gin.Context, param string) (uint64, error) {
	val := strings.TrimSpace(c.Param(param))
	if val == "" {
		return uint64(0), nil
	}
	num, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint64(num), nil
}

func (s *Server) uintFromContextQuery(c *gin.Context, query string) (uint, error) {
	val := strings.TrimSpace(c.Query(query))
	if val == "" {
		return uint(0), nil
	}
	num, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(num), nil
}

func (s *Server) float64FromContextQuery(c *gin.Context, query string) (float64, error) {
	val := strings.TrimSpace(c.Query(query))
	if val == "" {
		return 0, nil
	}
	num, err := strconv.ParseFloat(val, 10)
	if err != nil {
		return 0, err
	}
	return num, nil
}

func (s *Server) uint64FromContextQuery(c *gin.Context, query string) (uint64, error) {
	val := strings.TrimSpace(c.Query(query))
	if val == "" {
		return uint64(0), nil
	}
	num, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint64(num), nil
}

func (s *Server) stringFromContextQuery(c *gin.Context, query string) string {
	return strings.TrimSpace(c.Query(query))
}

func (s *Server) stringFromContextParam(c *gin.Context, query string) string {
	return strings.TrimSpace(c.Param(query))
}

func (s *Server) stringArrayFromContextQuery(c *gin.Context, query string) []string {
	val := strings.TrimSpace(c.Query(query))
	if val == "" {
		return []string{}
	}
	return strings.Split(val, ",")
}

func (s *Server) uintArrayFromContextQuery(c *gin.Context, query string) ([]uint, error) {
	val := strings.TrimSpace(c.Query(query))
	if val == "" {
		return []uint{}, nil
	}
	vals := strings.Split(val, ",")
	rets := []uint{}
	for _, val := range vals {
		num, err := strconv.ParseUint(val, 10, 64)
		if err != nil {
			return []uint{}, err
		}
		rets = append(rets, uint(num))
	}
	return rets, nil
}

func (s *Server) dateFromContextQuery(c *gin.Context, query string) (*time.Time, error) {
	val := strings.TrimSpace(c.Query(query))
	if val == "" {
		return nil, nil
	}
	t, err := time.Parse("2006-01-02", val)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (s *Server) timeFromContextQuery(c *gin.Context, query string) (*time.Time, error) {
	val := strings.TrimSpace(c.Query(query))
	if val == "" {
		return nil, nil
	}
	t, err := time.Parse(time.RFC3339, val)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (s *Server) boolFromContextQuery(c *gin.Context, query string) (*bool, error) {
	val := strings.TrimSpace(c.Query(query))
	if val == "" {
		return nil, nil
	}
	ret, err := strconv.ParseBool(val)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

func (s *Server) signatureBindJSON(c *gin.Context, resp interface{}) (*serializers.SignatureReq, error) {
	var req serializers.SignatureReq
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, errs.NewError(err)
	}
	err := json.Unmarshal([]byte(req.Message), &req)
	if err != nil {
		return nil, errs.NewError(err)
	}
	err = s.nls.VerifyAddressSignature(
		s.requestContext(c),
		req.Network,
		req.Address,
		req.Message,
		req.Signature,
	)
	if err != nil {
		return nil, errs.NewError(err)
	}
	return nil, nil
}

func (s *Server) validateTimestampWithSignature(ctx context.Context, network models.Network, address string, signature string, timestamp int64) error {
	if time.Unix(timestamp, 0).Before(time.Now().Add(-30 * time.Second)) {
		return errs.NewError(errs.ErrBadRequest)
	}
	if time.Unix(timestamp, 0).After(time.Now().Add(30 * time.Second)) {
		return errs.NewError(errs.ErrBadRequest)
	}
	err := s.nls.VerifyAddressSignature(
		ctx,
		network,
		address,
		fmt.Sprintf("%d", timestamp),
		signature,
	)
	if err != nil {
		return errs.NewError(err)
	}
	err = s.nls.VerifyUserTimestamp(
		ctx,
		network,
		address,
		timestamp,
	)
	if err != nil {
		return errs.NewError(err)
	}
	return nil
}

func (s *Server) getNetworkAddress(c *gin.Context) (models.Network, string, error) {
	network := s.stringFromContextQuery(c, "network")
	address := s.stringFromContextQuery(c, "address")
	if network == "" ||
		address == "" {
		return "", "", errs.NewError(errs.ErrBadRequest)
	}
	switch models.Network(network) {
	case models.NetworkSOL,
		models.NetworkAVAX,
		models.NetworkBOBA,
		models.NetworkBSC,
		models.NetworkETH,
		models.NetworkMATIC,
		models.NetworkNEAR:
		{
		}
	default:
		{
			return models.Network(""), "", errs.NewError(errs.ErrBadRequest)
		}
	}
	return models.Network(network), address, nil
}
