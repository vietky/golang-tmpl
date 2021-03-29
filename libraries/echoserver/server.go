package echoserver

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"time"

	"github.com/vietky/golang-tmpl/libraries/echopprof"
	"github.com/vietky/golang-tmpl/libraries/echoprometheus"
	"github.com/vietky/golang-tmpl/libraries/logger"
	"go.uber.org/zap"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// NewTokenString stands for "new_token"
const NewTokenString = "new_token"

// AuthPayload jwt payload
type AuthPayload struct {
	jwt.StandardClaims
	Typ string `json:"omitempty,typ"`
	// Vrf []string `json:"omitempty,vrf"`
}

//Counter type, support atomic inc/dec
type Counter int32

// Inc counter++
func (cnt *Counter) Inc() int32 {
	return atomic.AddInt32((*int32)(cnt), 1)
}

// Dec counter--
func (cnt *Counter) Dec() int32 {
	return atomic.AddInt32((*int32)(cnt), -1)
}

// Value current counter value
func (cnt *Counter) Value() int32 {
	return atomic.LoadInt32((*int32)(cnt))
}

// HookFunction HookFunction function
type HookFunction func(e *echo.Echo)

// Config config
type Config struct {
	MaxConnection    int32
	Port             int
	EnablePprof      bool
	EnablePrometheus bool
	JwtSecret        string
	HookFunction     HookFunction
}

// Server echo server
type Server struct {
	counter Counter
	config  Config
	echo    *echo.Echo
	log     *zap.SugaredLogger
}

// NewServer new server
func NewServer(config Config) (*Server, error) {
	logne := logger.GetLogger("echo-api-server")
	return &Server{
		counter: Counter(0),
		config:  config,
		echo:    echo.New(),
		log:     logne,
	}, nil
}

// Start server
func (server *Server) Start() { // declare all the method
	defer server.close()

	server.log.Info("Initializing handlers")

	server.log.Info("Bind custom http handler")
	server.echo.HTTPErrorHandler = server.errorsHandler
	// # Run before routing
	server.echo.Pre(server.checkMaxConn)
	// # End

	// # Run after routing
	server.echo.Use(server.writeStats)
	server.echo.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize: 1 << 10, // 1 KB
	}))

	if server.config.EnablePprof {
		echopprof.Wrapper(server.echo)
	}

	if server.config.EnablePrometheus {
		echoprometheus.NewPrometheus("").Use(server.echo)
	}

	jwtConfig := middleware.JWTConfig{
		Claims:     &AuthPayload{},
		SigningKey: []byte(server.config.JwtSecret),
	}

	private := server.echo.Group("/api/v1/private")

	private.Use(middleware.JWTWithConfig(jwtConfig))

	// # End
	// health, metrics handlers
	server.echo.GET("/health", server.healthCheck)
	server.config.HookFunction(server.echo)

	go server.start()

	server.waitForInterruptSignal()
}

// Start server
func (server *Server) start() error {
	server.log.Info("Starting at Port: ", server.config.Port)
	err := server.echo.Start(fmt.Sprintf(":%d", server.config.Port))
	if err != nil {
		server.log.Error(err)
		return err
	}
	return nil
}

// Close handler
func (server *Server) close() {
}

func (server *Server) waitForInterruptSignal() {
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	server.log.Debug("start shutting down")
	if err := server.echo.Shutdown(ctx); err != nil {
		server.log.Error(err)
	}
	server.log.Debug("finish shutting down")
}

func (server *Server) checkMaxConn(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if server.counter.Inc() > server.config.MaxConnection {
			server.log.Infof("Reach max connection limit: %v", server.config.MaxConnection)
			return echo.NewHTTPError(http.StatusServiceUnavailable, "REACH_LIMIT_CONNECTION")
		}
		return next(c)
	}
}

func (server *Server) writeStats(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		startTime := time.Now()
		defer func(t time.Time) {
			server.counter.Dec()
			server.log.Infof("%s runtime: %f", c.Path(), time.Since(startTime).Seconds())
		}(startTime)
		if err := next(c); err != nil {
			return err
		}
		return nil
	}
}

func (server *Server) healthCheck(c echo.Context) error {
	server.log.Info("Health check ok")
	return c.String(http.StatusOK, "OK")
}

// TODO this function should handle not only errors but also others
// Do the general error handling for db errors.
// Because we still use the blocket procedures, so we should implement this.
// It should be in db package, but I think we should put it in handlers package.
// Because we have to return the request status, for ex: 200 or 404,...
// the response should have the following format
// {
// 	"status": "SPINE_OK" or "SPINE_ERROR"
// 	"message": "error_msg"
// 	"data1":"data2",
// 	...
// }
func (server *Server) errorsHandler(err error, c echo.Context) {
	server.log.Info("Start handler error: ", err)
	var (
		code = http.StatusInternalServerError
		msg  interface{}
	)

	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		msg = he.Message
	} else {
		msg = err.Error() // http.StatusText(code)
	}

	if _, ok := msg.(string); ok {
		msg = echo.Map{"message": msg}
	} else {
		msg = echo.Map{"message": fmt.Sprintf("%v", msg)}
	}

	token := c.Get(NewTokenString)
	if tokenStr, ok := token.(string); ok && tokenStr != "" {
		msg.(echo.Map)["token"] = tokenStr
	}

	if !c.Response().Committed {
		if c.Request().Method == echo.HEAD { // Issue #608
			if err = c.NoContent(code); err != nil {
				goto ERROR
			}
		} else {
			if err = c.JSON(code, msg); err != nil {
				goto ERROR
			}
		}
	}
ERROR:
	if code < 500 {
		c.Logger().Debug(err)
	} else {
		c.Logger().Error(err)
	}
}
