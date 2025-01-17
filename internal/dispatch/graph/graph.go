package graph

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog/log"
	"github.com/shopspring/decimal"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/authzed/spicedb/internal/datastore"
	"github.com/authzed/spicedb/internal/dispatch"
	"github.com/authzed/spicedb/internal/graph"
	"github.com/authzed/spicedb/internal/namespace"
	graphwalk "github.com/authzed/spicedb/pkg/graph"
	core "github.com/authzed/spicedb/pkg/proto/core/v1"
	v1 "github.com/authzed/spicedb/pkg/proto/dispatch/v1"
	"github.com/authzed/spicedb/pkg/tuple"
)

const errDispatch = "error dispatching request: %w"

var tracer = otel.Tracer("spicedb/internal/dispatch/local")

var slowLookupCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
	Name: "spicedb_dispatch_slow_lookup_total",
	Help: "count of how many times the slow lookup path is taken",
}, []string{"prefix"})

// NewLocalOnlyDispatcher creates a dispatcher that consults with the graph to formulate a response.
func NewLocalOnlyDispatcher(
	nsm namespace.Manager,
) dispatch.Dispatcher {
	d := &localDispatcher{nsm: nsm}

	d.checker = graph.NewConcurrentChecker(d, nsm)
	d.expander = graph.NewConcurrentExpander(d, nsm)
	d.lookupHandler = graph.NewConcurrentLookup(d, d, nsm)

	return d
}

// NewDispatcher creates a dispatcher that consults with the graph and redispatches subproblems to
// the provided redispatcher.
func NewDispatcher(
	redispatcher dispatch.Dispatcher,
	nsm namespace.Manager,
) dispatch.Dispatcher {
	checker := graph.NewConcurrentChecker(redispatcher, nsm)
	expander := graph.NewConcurrentExpander(redispatcher, nsm)
	lookupHandler := graph.NewConcurrentLookup(redispatcher, redispatcher, nsm)

	return &localDispatcher{checker, expander, lookupHandler, nsm}
}

type localDispatcher struct {
	checker       *graph.ConcurrentChecker
	expander      *graph.ConcurrentExpander
	lookupHandler *graph.ConcurrentLookup

	nsm namespace.Manager
}

func (ld *localDispatcher) loadRelation(ctx context.Context, nsName, relationName string, revision decimal.Decimal) (*core.Relation, error) {
	// Load namespace and relation from the datastore
	ns, err := ld.nsm.ReadNamespace(ctx, nsName, revision)
	if err != nil {
		return nil, rewriteError(err)
	}

	var relation *core.Relation
	for _, candidate := range ns.Relation {
		if candidate.Name == relationName {
			relation = candidate
			break
		}
	}

	if relation == nil {
		return nil, NewRelationNotFoundErr(nsName, relationName)
	}

	return relation, nil
}

type stringableOnr struct {
	*core.ObjectAndRelation
}

func (onr stringableOnr) String() string {
	return tuple.StringONR(onr.ObjectAndRelation)
}

type stringableRelRef struct {
	*core.RelationReference
}

func (rr stringableRelRef) String() string {
	return fmt.Sprintf("%s::%s", rr.Namespace, rr.Relation)
}

// DispatchCheck implements dispatch.Check interface
func (ld *localDispatcher) DispatchCheck(ctx context.Context, req *v1.DispatchCheckRequest) (*v1.DispatchCheckResponse, error) {
	ctx, span := tracer.Start(ctx, "DispatchCheck", trace.WithAttributes(
		attribute.Stringer("start", stringableOnr{req.ObjectAndRelation}),
		attribute.Stringer("subject", stringableOnr{req.Subject}),
	))
	defer span.End()

	err := dispatch.CheckDepth(ctx, req)
	if err != nil {
		return &v1.DispatchCheckResponse{Metadata: emptyMetadata}, err
	}

	revision, err := decimal.NewFromString(req.Metadata.AtRevision)
	if err != nil {
		return &v1.DispatchCheckResponse{Metadata: emptyMetadata}, err
	}

	relation, err := ld.loadRelation(ctx, req.ObjectAndRelation.Namespace, req.ObjectAndRelation.Relation, revision)
	if err != nil {
		return &v1.DispatchCheckResponse{Metadata: emptyMetadata}, err
	}

	validatedReq := graph.ValidatedCheckRequest{
		DispatchCheckRequest: req,
		Revision:             revision,
	}

	return ld.checker.Check(ctx, validatedReq, relation)
}

// DispatchExpand implements dispatch.Expand interface
func (ld *localDispatcher) DispatchExpand(ctx context.Context, req *v1.DispatchExpandRequest) (*v1.DispatchExpandResponse, error) {
	ctx, span := tracer.Start(ctx, "DispatchExpand", trace.WithAttributes(
		attribute.Stringer("start", stringableOnr{req.ObjectAndRelation}),
	))
	defer span.End()

	err := dispatch.CheckDepth(ctx, req)
	if err != nil {
		return &v1.DispatchExpandResponse{Metadata: emptyMetadata}, err
	}

	revision, err := decimal.NewFromString(req.Metadata.AtRevision)
	if err != nil {
		return &v1.DispatchExpandResponse{Metadata: emptyMetadata}, err
	}

	relation, err := ld.loadRelation(ctx, req.ObjectAndRelation.Namespace, req.ObjectAndRelation.Relation, revision)
	if err != nil {
		return &v1.DispatchExpandResponse{Metadata: emptyMetadata}, err
	}

	validatedReq := graph.ValidatedExpandRequest{
		DispatchExpandRequest: req,
		Revision:              revision,
	}

	return ld.expander.Expand(ctx, validatedReq, relation)
}

// DispatchLookup implements dispatch.Lookup interface
func (ld *localDispatcher) DispatchLookup(ctx context.Context, req *v1.DispatchLookupRequest) (*v1.DispatchLookupResponse, error) {
	ctx, span := tracer.Start(ctx, "DispatchLookup", trace.WithAttributes(
		attribute.Stringer("start", stringableRelRef{req.ObjectRelation}),
		attribute.Stringer("subject", stringableOnr{req.Subject}),
		attribute.Int64("limit", int64(req.Limit)),
	))
	defer span.End()

	err := dispatch.CheckDepth(ctx, req)
	if err != nil {
		return &v1.DispatchLookupResponse{Metadata: emptyMetadata}, err
	}

	revision, err := decimal.NewFromString(req.Metadata.AtRevision)
	if err != nil {
		return &v1.DispatchLookupResponse{Metadata: emptyMetadata}, err
	}

	if req.Limit <= 0 {
		return &v1.DispatchLookupResponse{Metadata: emptyMetadata, ResolvedOnrs: []*core.ObjectAndRelation{}}, nil
	}

	validatedReq := graph.ValidatedLookupRequest{
		DispatchLookupRequest: req,
		Revision:              revision,
	}
	relation, err := ld.loadRelation(ctx, req.ObjectRelation.Namespace, req.ObjectRelation.Relation, revision)
	if err != nil {
		return &v1.DispatchLookupResponse{Metadata: emptyMetadata, ResolvedOnrs: []*core.ObjectAndRelation{}}, nil
	}

	if ld.requiresLookupViaChecks(relation) {
		log.Warn().Msg("slow path dispatch lookup")
		slowLookupCounter.WithLabelValues(prefix(validatedReq.Subject.Namespace)).Inc()
		return ld.lookupHandler.LookupViaChecks(ctx, validatedReq)
	}

	return ld.lookupHandler.Lookup(ctx, validatedReq)
}

func prefix(namespace string) (prefix string) {
	parts := strings.SplitN(namespace, "/", 2)
	if len(parts) > 1 {
		prefix = parts[0]
	}
	return prefix
}

func (ld *localDispatcher) requiresLookupViaChecks(relation *core.Relation) bool {
	// TODO: refactor walker so that we don't have to make two separate checks
	// check top-level rewrite
	if rw := relation.GetUsersetRewrite(); rw != nil {
		switch rw.RewriteOperation.(type) {
		case *core.UsersetRewrite_Intersection:
			return true
		case *core.UsersetRewrite_Exclusion:
			return true
		}
	}

	// check child rewrites
	childIntersectionExclusion := graphwalk.WalkRewrite(relation.GetUsersetRewrite(), func(childOneof *core.SetOperation_Child) interface{} {
		switch child := childOneof.ChildType.(type) {
		case *core.SetOperation_Child_UsersetRewrite:
			switch child.UsersetRewrite.RewriteOperation.(type) {
			case *core.UsersetRewrite_Intersection:
				return true
			case *core.UsersetRewrite_Exclusion:
				return true
			}
		}
		return nil
	})
	return childIntersectionExclusion != nil
}

func (ld *localDispatcher) Close() error {
	return nil
}

func rewriteError(original error) error {
	nsNotFound := datastore.ErrNamespaceNotFound{}

	switch {
	case errors.As(original, &nsNotFound):
		return NewNamespaceNotFoundErr(nsNotFound.NotFoundNamespaceName())
	case errors.As(original, &ErrNamespaceNotFound{}):
		fallthrough
	case errors.As(original, &ErrRelationNotFound{}):
		return original
	default:
		return fmt.Errorf(errDispatch, original)
	}
}

var emptyMetadata *v1.ResponseMeta = &v1.ResponseMeta{
	DispatchCount: 0,
}
