package command

import (
	"context"

	grpc "google.golang.org/grpc"

	core "github.com/v2fly/v2ray-core/v5"
	"github.com/v2fly/v2ray-core/v5/common"
	"github.com/v2fly/v2ray-core/v5/common/serial"
	"github.com/v2fly/v2ray-core/v5/features/inbound"
	"github.com/v2fly/v2ray-core/v5/features/outbound"
	"github.com/v2fly/v2ray-core/v5/proxy"
)

// InboundOperation is the interface for operations that applies to inbound handlers.
type InboundOperation interface {
	// ApplyInbound applies this operation to the given inbound handler.
	ApplyInbound(context.Context, inbound.Handler) error
}

// OutboundOperation is the interface for operations that applies to outbound handlers.
type OutboundOperation interface {
	// ApplyOutbound applies this operation to the given outbound handler.
	ApplyOutbound(context.Context, outbound.Handler) error
}

func getInbound(handler inbound.Handler) (proxy.Inbound, error) {
	gi, ok := handler.(proxy.GetInbound)
	if !ok {
		return nil, newError("can't get inbound proxy from handler.")
	}
	return gi.GetInbound(), nil
}

// ApplyInbound implements InboundOperation.
func (op *AddUserOperation) ApplyInbound(ctx context.Context, handler inbound.Handler) error {
	p, err := getInbound(handler)
	if err != nil {
		return err
	}
	um, ok := p.(proxy.UserManager)
	if !ok {
		return newError("proxy is not a UserManager")
	}
	mUser, err := op.User.ToMemoryUser()
	if err != nil {
		return newError("failed to parse user").Base(err)
	}
	return um.AddUser(ctx, mUser)
}

// ApplyInbound implements InboundOperation.
func (op *RemoveUserOperation) ApplyInbound(ctx context.Context, handler inbound.Handler) error {
	p, err := getInbound(handler)
	if err != nil {
		return err
	}
	um, ok := p.(proxy.UserManager)
	if !ok {
		return newError("proxy is not a UserManager")
	}
	return um.RemoveUser(ctx, op.Email)
}

type handlerServer struct {
	s   *core.Instance
	ihm inbound.Manager
	ohm outbound.Manager
}

func (s *handlerServer) AddInbound(ctx context.Context, request *AddInboundRequest) (*AddInboundResponse, error) {
	if err := core.AddInboundHandler(s.s, request.Inbound); err != nil {
		return nil, err
	}

	return &AddInboundResponse{}, nil
}

func (s *handlerServer) RemoveInbound(ctx context.Context, request *RemoveInboundRequest) (*RemoveInboundResponse, error) {
	return &RemoveInboundResponse{}, s.ihm.RemoveHandler(ctx, request.Tag)
}

func (s *handlerServer) ListUsers(ctx context.Context, request *ListUsersRequest) (*ListUsersResponse, error) {
	handler, err := s.ihm.GetHandler(ctx, request.Tag)
	if err != nil {
		return &ListUsersResponse{}, newError("failed to get handler: ", request.Tag).Base(err)
	}
	p, err := getInbound(handler)
	if err != nil {
		return &ListUsersResponse{}, newError("failed to get handler: ", request.Tag).Base(err)
	}
	if v,ok := p.(proxy.ListUser);ok {
		newError("List users success ", request.Tag).AtDebug().WriteToLog()
		ret := v.ListUser(ctx)
		return &ListUsersResponse{Users:ret},nil
		
	}
	return &ListUsersResponse{}, newError("Do not support user list: ", request.Tag).Base(err)
}

func (s *handlerServer) AddUsers(ctx context.Context, request *AddUsersRequest) (*AddUsersResponse, error) {
	handler, err := s.ihm.GetHandler(ctx, request.Tag)
	if err != nil {
		return &AddUsersResponse{}, newError("failed to get handler: ", request.Tag).Base(err)
	}
	p, err := getInbound(handler)
	if err != nil {
		return &AddUsersResponse{}, newError("failed to get handler: ", request.Tag).Base(err)
	}
	um, ok := p.(proxy.UserManager)
	if !ok {
		return &AddUsersResponse{},newError("proxy is not a UserManager")
	}
	for _,u := range request.Users {
		mUser, err := u.ToMemoryUser()
		if err != nil {
			return &AddUsersResponse{},newError("failed to parse user", u.Email).Base(err)
		}
		if err = um.AddUser(ctx, mUser); err != nil {
			return &AddUsersResponse{},newError("failed to add user", u.Email).Base(err)
		}
	}
	return &AddUsersResponse{},nil
}

func (s *handlerServer) RmUsers(ctx context.Context, request *RmUsersRequest) (*RmUsersResponse, error) {
	handler, err := s.ihm.GetHandler(ctx, request.Tag)
	if err != nil {
		return &RmUsersResponse{}, newError("failed to get handler: ", request.Tag).Base(err)
	}
	p, err := getInbound(handler)
	if err != nil {
		return &RmUsersResponse{}, newError("failed to get handler: ", request.Tag).Base(err)
	}
	um, ok := p.(proxy.UserManager)
	if !ok {
		return &RmUsersResponse{},newError("proxy is not a UserManager")
	}
	for _,u := range request.Users {
		if err := um.RemoveUser(ctx, u.Email); err != nil {
			return &RmUsersResponse{},newError("failed to rm user", u.Email).Base(err)
		}
	}
	return &RmUsersResponse{},nil
}

func (s *handlerServer) AlterInbound(ctx context.Context, request *AlterInboundRequest) (*AlterInboundResponse, error) {
	rawOperation, err := serial.GetInstanceOf(request.Operation)
	if err != nil {
		return nil, newError("unknown operation").Base(err)
	}
	operation, ok := rawOperation.(InboundOperation)
	if !ok {
		return nil, newError("not an inbound operation")
	}

	handler, err := s.ihm.GetHandler(ctx, request.Tag)
	if err != nil {
		return nil, newError("failed to get handler: ", request.Tag).Base(err)
	}

	return &AlterInboundResponse{}, operation.ApplyInbound(ctx, handler)
}

func (s *handlerServer) AddOutbound(ctx context.Context, request *AddOutboundRequest) (*AddOutboundResponse, error) {
	if err := core.AddOutboundHandler(s.s, request.Outbound); err != nil {
		return nil, err
	}
	return &AddOutboundResponse{}, nil
}

func (s *handlerServer) RemoveOutbound(ctx context.Context, request *RemoveOutboundRequest) (*RemoveOutboundResponse, error) {
	return &RemoveOutboundResponse{}, s.ohm.RemoveHandler(ctx, request.Tag)
}

func (s *handlerServer) AlterOutbound(ctx context.Context, request *AlterOutboundRequest) (*AlterOutboundResponse, error) {
	rawOperation, err := serial.GetInstanceOf(request.Operation)
	if err != nil {
		return nil, newError("unknown operation").Base(err)
	}
	operation, ok := rawOperation.(OutboundOperation)
	if !ok {
		return nil, newError("not an outbound operation")
	}

	handler := s.ohm.GetHandler(request.Tag)
	return &AlterOutboundResponse{}, operation.ApplyOutbound(ctx, handler)
}

func (s *handlerServer) mustEmbedUnimplementedHandlerServiceServer() {}

type service struct {
	v *core.Instance
}

func (s *service) Register(server *grpc.Server) {
	hs := &handlerServer{
		s: s.v,
	}
	common.Must(s.v.RequireFeatures(func(im inbound.Manager, om outbound.Manager) {
		hs.ihm = im
		hs.ohm = om
	}))
	RegisterHandlerServiceServer(server, hs)
}

func init() {
	common.Must(common.RegisterConfig((*Config)(nil), func(ctx context.Context, cfg interface{}) (interface{}, error) {
		s := core.MustFromContext(ctx)
		return &service{v: s}, nil
	}))
}
