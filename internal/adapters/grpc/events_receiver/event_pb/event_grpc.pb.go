// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.12.4
// source: event.proto

package event_pb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// EventWriterClient is the client API for EventWriter service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type EventWriterClient interface {
	WriteEvent(ctx context.Context, in *EventParameters, opts ...grpc.CallOption) (*WriteStatus, error)
}

type eventWriterClient struct {
	cc grpc.ClientConnInterface
}

func NewEventWriterClient(cc grpc.ClientConnInterface) EventWriterClient {
	return &eventWriterClient{cc}
}

func (c *eventWriterClient) WriteEvent(ctx context.Context, in *EventParameters, opts ...grpc.CallOption) (*WriteStatus, error) {
	out := new(WriteStatus)
	err := c.cc.Invoke(ctx, "/event_pb.EventWriter/WriteEvent", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// EventWriterServer is the server API for EventWriter service.
// All implementations must embed UnimplementedEventWriterServer
// for forward compatibility
type EventWriterServer interface {
	WriteEvent(context.Context, *EventParameters) (*WriteStatus, error)
	mustEmbedUnimplementedEventWriterServer()
}

// UnimplementedEventWriterServer must be embedded to have forward compatible implementations.
type UnimplementedEventWriterServer struct {
}

func (UnimplementedEventWriterServer) WriteEvent(context.Context, *EventParameters) (*WriteStatus, error) {
	return nil, status.Errorf(codes.Unimplemented, "method WriteEvent not implemented")
}
func (UnimplementedEventWriterServer) mustEmbedUnimplementedEventWriterServer() {}

// UnsafeEventWriterServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to EventWriterServer will
// result in compilation errors.
type UnsafeEventWriterServer interface {
	mustEmbedUnimplementedEventWriterServer()
}

func RegisterEventWriterServer(s grpc.ServiceRegistrar, srv EventWriterServer) {
	s.RegisterService(&EventWriter_ServiceDesc, srv)
}

func _EventWriter_WriteEvent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EventParameters)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EventWriterServer).WriteEvent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/event_pb.EventWriter/WriteEvent",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EventWriterServer).WriteEvent(ctx, req.(*EventParameters))
	}
	return interceptor(ctx, in, info, handler)
}

// EventWriter_ServiceDesc is the grpc.ServiceDesc for EventWriter service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var EventWriter_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "event_pb.EventWriter",
	HandlerType: (*EventWriterServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "WriteEvent",
			Handler:    _EventWriter_WriteEvent_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "event.proto",
}