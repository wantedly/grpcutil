package statsutil

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"google.golang.org/grpc/stats"

	grpcutiltesting "github.com/wantedly/grpcutil/testing"
)

type handlerCall struct {
	id       int
	funcname string
	typename string
}

type handlerCallRecorder struct {
	calls []handlerCall
}

func (r *handlerCallRecorder) Record(id int, funcname string, obj interface{}) {
	r.calls = append(r.calls, handlerCall{id: id, funcname: funcname, typename: reflect.ValueOf(obj).Elem().Type().Name()})
}

func (r *handlerCallRecorder) Print() {
	rows := make([]string, 0, len(r.calls))
	for _, c := range r.calls {
		rows = append(rows, fmt.Sprintf("%d %s %s", c.id, c.funcname, c.typename))
	}
	fmt.Println(strings.Join(rows, "\n"))
}

type fakeHandler struct {
	id       int
	recorder *handlerCallRecorder
}

func (h *fakeHandler) TagRPC(ctx context.Context, info *stats.RPCTagInfo) context.Context {
	h.recorder.Record(h.id, "TagRPC", info)
	return ctx
}

func (h *fakeHandler) HandleRPC(ctx context.Context, s stats.RPCStats) {
	h.recorder.Record(h.id, "HandleRPC", s)
}

func (h *fakeHandler) TagConn(ctx context.Context, info *stats.ConnTagInfo) context.Context {
	h.recorder.Record(h.id, "TagConn", info)
	return ctx
}

func (h *fakeHandler) HandleConn(ctx context.Context, s stats.ConnStats) {
	h.recorder.Record(h.id, "HandleConn", s)
}

func ExampleChainHandlers() {
	recorder := &handlerCallRecorder{}

	tc := grpcutiltesting.CreateTestContext(nil)
	tc.SetServerStatsHandler(ChainHandlers(
		&fakeHandler{id: 1, recorder: recorder},
		&fakeHandler{id: 2, recorder: recorder},
		&fakeHandler{id: 3, recorder: recorder},
	))

	tc.Setup()
	defer tc.Teardown()

	_, _ = tc.Client.Echo(context.TODO(), &grpcutiltesting.EchoRequest{Message: "ping"})

	recorder.Print()
	// Output:
	// 1 TagConn ConnTagInfo
	// 2 TagConn ConnTagInfo
	// 3 TagConn ConnTagInfo
	// 1 HandleConn ConnBegin
	// 2 HandleConn ConnBegin
	// 3 HandleConn ConnBegin
	// 1 TagRPC RPCTagInfo
	// 2 TagRPC RPCTagInfo
	// 3 TagRPC RPCTagInfo
	// 1 HandleRPC InHeader
	// 2 HandleRPC InHeader
	// 3 HandleRPC InHeader
	// 1 HandleRPC Begin
	// 2 HandleRPC Begin
	// 3 HandleRPC Begin
	// 1 HandleRPC InPayload
	// 2 HandleRPC InPayload
	// 3 HandleRPC InPayload
	// 1 HandleRPC OutHeader
	// 2 HandleRPC OutHeader
	// 3 HandleRPC OutHeader
	// 1 HandleRPC OutPayload
	// 2 HandleRPC OutPayload
	// 3 HandleRPC OutPayload
	// 1 HandleRPC OutTrailer
	// 2 HandleRPC OutTrailer
	// 3 HandleRPC OutTrailer
	// 1 HandleRPC End
	// 2 HandleRPC End
	// 3 HandleRPC End
}
