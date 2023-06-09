package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/pprof"

	"github.com/fatih/color"
	"github.com/go-logr/zerologr"
	grpczerolog "github.com/grpc-ecosystem/go-grpc-middleware/providers/zerolog/v2"
	grpcauth "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	grpclog "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	grpcprom "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/jzelinskie/cobrautil/v2"
	"github.com/jzelinskie/cobrautil/v2/cobraotel"
	"github.com/jzelinskie/cobrautil/v2/cobrazerolog"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"github.com/authzed/spicedb/internal/dispatch"
	"github.com/authzed/spicedb/internal/logging"
	consistencymw "github.com/authzed/spicedb/internal/middleware/consistency"
	datastoremw "github.com/authzed/spicedb/internal/middleware/datastore"
	dispatchmw "github.com/authzed/spicedb/internal/middleware/dispatcher"
	"github.com/authzed/spicedb/internal/middleware/servicespecific"
	"github.com/authzed/spicedb/pkg/datastore"
	logmw "github.com/authzed/spicedb/pkg/middleware/logging"
	"github.com/authzed/spicedb/pkg/middleware/requestid"
	"github.com/authzed/spicedb/pkg/middleware/serverversion"
	"github.com/authzed/spicedb/pkg/releases"
)

var DisableTelemetryHandler *prometheus.Registry

// ServeExample creates an example usage string with the provided program name.
func ServeExample(programName string) string {
	return fmt.Sprintf(`	%[1]s:
		%[3]s serve --grpc-preshared-key "somerandomkeyhere"

	%[2]s:
		%[3]s serve --grpc-preshared-key "realkeyhere" --grpc-tls-cert-path path/to/tls/cert --grpc-tls-key-path path/to/tls/key \
			--http-tls-cert-path path/to/tls/cert --http-tls-key-path path/to/tls/key \
			--datastore-engine postgres --datastore-conn-uri "postgres-connection-string-here"
`,
		color.YellowString("No TLS and in-memory"),
		color.GreenString("TLS and a real datastore"),
		programName,
	)
}

// DefaultPreRunE sets up viper, zerolog, and OpenTelemetry flag handling for a
// command.
func DefaultPreRunE(programName string) cobrautil.CobraRunFunc {
	return cobrautil.CommandStack(
		cobrautil.SyncViperPreRunE(programName),
		cobrazerolog.New(
			cobrazerolog.WithTarget(func(logger zerolog.Logger) {
				logging.SetGlobalLogger(logger)
			}),
		).RunE(),
		cobraotel.New("spicedb",
			cobraotel.WithLogger(zerologr.New(&logging.Logger)),
		).RunE(),
		releases.CheckAndLogRunE(),
	)
}

// MetricsHandler sets up an HTTP server that handles serving Prometheus
// metrics and pprof endpoints.
func MetricsHandler(telemetryRegistry *prometheus.Registry, c *Config) http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/metrics", promhttp.Handler())
	if telemetryRegistry != nil {
		mux.Handle("/telemetry", promhttp.HandlerFor(telemetryRegistry, promhttp.HandlerOpts{}))
	}

	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	mux.HandleFunc("/debug/pprof/cmdline", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "This profile type has been disabled to avoid leaking private command-line arguments")
	})
	mux.HandleFunc("/debug/config", func(w http.ResponseWriter, r *http.Request) {
		if c == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		json, err := json.MarshalIndent(c.DebugMap(), "", "  ")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "%s", string(json))
	})

	return mux
}

var defaultGRPCLogOptions = []grpclog.Option{
	// the server has a deadline set, so we consider it a normal condition
	// this makes sure we don't log them as errors
	grpclog.WithLevels(func(code codes.Code) grpclog.Level {
		if code == codes.DeadlineExceeded {
			return grpclog.INFO
		}
		return grpclog.DefaultServerCodeToLevel(code)
	}),
	// changes default logging behaviour to only log finish call message
	grpclog.WithDecider(func(_ string, _ error) grpclog.Decision {
		return grpclog.LogFinishCall
	}),
}

const (
	DefaultMiddlewareRequestID     = "requestid"
	DefaultMiddlewareLog           = "log"
	DefaultMiddlewareGRPCLog       = "grpclog"
	DefaultMiddlewareOTelGRPC      = "otelgrpc"
	DefaultMiddlewareGRPCAuth      = "grpcauth"
	DefaultMiddlewareGRPCProm      = "grpcprom"
	DefaultMiddlewareServerVersion = "serverversion"

	DefaultInternalMiddlewareDispatch       = "dispatch"
	DefaultInternalMiddlewareDatastore      = "datastore"
	DefaultInternalMiddlewareConsistency    = "consistency"
	DefaultInternalMiddlewareServerSpecific = "servicespecific"
)

// DefaultUnaryMiddleware generates the default middleware chain used for the public SpiceDB Unary gRPC methods
func DefaultUnaryMiddleware(logger zerolog.Logger, authFunc grpcauth.AuthFunc, enableVersionResponse bool, dispatcher dispatch.Dispatcher, ds datastore.Datastore) (*MiddlewareChain[grpc.UnaryServerInterceptor], error) {
	chain, err := NewMiddlewareChain([]ReferenceableMiddleware[grpc.UnaryServerInterceptor]{
		{
			Name:       DefaultMiddlewareRequestID,
			Middleware: requestid.UnaryServerInterceptor(requestid.GenerateIfMissing(true)),
		},
		{
			Name:       DefaultMiddlewareLog,
			Middleware: logmw.UnaryServerInterceptor(logmw.ExtractMetadataField("x-request-id", "requestID")),
		},
		{
			Name:       DefaultMiddlewareGRPCLog,
			Middleware: grpclog.UnaryServerInterceptor(grpczerolog.InterceptorLogger(logger), defaultGRPCLogOptions...),
		},
		{
			Name:       DefaultMiddlewareOTelGRPC,
			Middleware: otelgrpc.UnaryServerInterceptor(),
		},
		{
			Name:       DefaultMiddlewareGRPCProm,
			Middleware: grpcprom.UnaryServerInterceptor,
		},
		{
			Name:       DefaultMiddlewareGRPCAuth,
			Middleware: grpcauth.UnaryServerInterceptor(authFunc),
		},
		{
			Name:       DefaultMiddlewareServerVersion,
			Middleware: serverversion.UnaryServerInterceptor(enableVersionResponse),
		},
		{
			Name:       DefaultInternalMiddlewareDispatch,
			Internal:   true,
			Middleware: dispatchmw.UnaryServerInterceptor(dispatcher),
		},
		{
			Name:       DefaultInternalMiddlewareDatastore,
			Internal:   true,
			Middleware: datastoremw.UnaryServerInterceptor(ds),
		},
		{
			Name:       DefaultInternalMiddlewareConsistency,
			Internal:   true,
			Middleware: consistencymw.UnaryServerInterceptor(),
		},
		{
			Name:       DefaultInternalMiddlewareServerSpecific,
			Internal:   true,
			Middleware: servicespecific.UnaryServerInterceptor,
		},
	}...)
	return &chain, err
}

// DefaultStreamingMiddleware generates the default middleware chain used for the public SpiceDB Streaming gRPC methods
func DefaultStreamingMiddleware(logger zerolog.Logger, authFunc grpcauth.AuthFunc, enableVersionResponse bool, dispatcher dispatch.Dispatcher, ds datastore.Datastore) (*MiddlewareChain[grpc.StreamServerInterceptor], error) {
	chain, err := NewMiddlewareChain([]ReferenceableMiddleware[grpc.StreamServerInterceptor]{
		{
			Name:       DefaultMiddlewareRequestID,
			Middleware: requestid.StreamServerInterceptor(requestid.GenerateIfMissing(true)),
		},
		{
			Name:       DefaultMiddlewareLog,
			Middleware: logmw.StreamServerInterceptor(logmw.ExtractMetadataField("x-request-id", "requestID")),
		},
		{
			Name:       DefaultMiddlewareGRPCLog,
			Middleware: grpclog.StreamServerInterceptor(grpczerolog.InterceptorLogger(logger), defaultGRPCLogOptions...),
		},
		{
			Name:       DefaultMiddlewareOTelGRPC,
			Middleware: otelgrpc.StreamServerInterceptor(),
		},
		{
			Name:       DefaultMiddlewareGRPCProm,
			Middleware: grpcprom.StreamServerInterceptor,
		},
		{
			Name:       DefaultMiddlewareGRPCAuth,
			Middleware: grpcauth.StreamServerInterceptor(authFunc),
		},
		{
			Name:       DefaultMiddlewareServerVersion,
			Middleware: serverversion.StreamServerInterceptor(enableVersionResponse),
		},
		{
			Name:       DefaultInternalMiddlewareDispatch,
			Internal:   true,
			Middleware: dispatchmw.StreamServerInterceptor(dispatcher),
		},
		{
			Name:       DefaultInternalMiddlewareDatastore,
			Internal:   true,
			Middleware: datastoremw.StreamServerInterceptor(ds),
		},
		{
			Name:       DefaultInternalMiddlewareConsistency,
			Internal:   true,
			Middleware: consistencymw.StreamServerInterceptor(),
		},
		{
			Name:       DefaultInternalMiddlewareServerSpecific,
			Internal:   true,
			Middleware: servicespecific.StreamServerInterceptor,
		},
	}...)
	return &chain, err
}

// DefaultDispatchMiddleware generates the default middleware chain used for the internal dispatch SpiceDB gRPC API
func DefaultDispatchMiddleware(logger zerolog.Logger, authFunc grpcauth.AuthFunc, ds datastore.Datastore) ([]grpc.UnaryServerInterceptor, []grpc.StreamServerInterceptor) {
	return []grpc.UnaryServerInterceptor{
			requestid.UnaryServerInterceptor(requestid.GenerateIfMissing(true)),
			logmw.UnaryServerInterceptor(logmw.ExtractMetadataField("x-request-id", "requestID")),
			grpclog.UnaryServerInterceptor(grpczerolog.InterceptorLogger(logger), defaultGRPCLogOptions...),
			otelgrpc.UnaryServerInterceptor(),
			grpcprom.UnaryServerInterceptor,
			grpcauth.UnaryServerInterceptor(authFunc),
			datastoremw.UnaryServerInterceptor(ds),
			servicespecific.UnaryServerInterceptor,
		}, []grpc.StreamServerInterceptor{
			requestid.StreamServerInterceptor(requestid.GenerateIfMissing(true)),
			logmw.StreamServerInterceptor(logmw.ExtractMetadataField("x-request-id", "requestID")),
			grpclog.StreamServerInterceptor(grpczerolog.InterceptorLogger(logger), defaultGRPCLogOptions...),
			otelgrpc.StreamServerInterceptor(),
			grpcprom.StreamServerInterceptor,
			grpcauth.StreamServerInterceptor(authFunc),
			datastoremw.StreamServerInterceptor(ds),
			servicespecific.StreamServerInterceptor,
		}
}
