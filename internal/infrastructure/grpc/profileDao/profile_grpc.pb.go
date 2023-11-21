// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.25.0
// source: profile.proto

package profileDao

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

// ProfileRepoClient is the client API for ProfileRepo service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ProfileRepoClient interface {
	FindProfileInStorage(ctx context.Context, in *ProfileName, opts ...grpc.CallOption) (*ProfileObj, error)
	InsertProfileToStorage(ctx context.Context, in *ProfileObj, opts ...grpc.CallOption) (*InsertResult, error)
	DeleteProfileFromStorage(ctx context.Context, in *ProfileName, opts ...grpc.CallOption) (*Empty, error)
	UpdateProfileInStorage(ctx context.Context, in *ProfileUpdate, opts ...grpc.CallOption) (*Empty, error)
	ListProfilesFromStorage(ctx context.Context, in *Page, opts ...grpc.CallOption) (*SliceOfProfileObj, error)
}

type profileRepoClient struct {
	cc grpc.ClientConnInterface
}

func NewProfileRepoClient(cc grpc.ClientConnInterface) ProfileRepoClient {
	return &profileRepoClient{cc}
}

func (c *profileRepoClient) FindProfileInStorage(ctx context.Context, in *ProfileName, opts ...grpc.CallOption) (*ProfileObj, error) {
	out := new(ProfileObj)
	err := c.cc.Invoke(ctx, "/profileDao.ProfileRepo/FindProfileInStorage", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *profileRepoClient) InsertProfileToStorage(ctx context.Context, in *ProfileObj, opts ...grpc.CallOption) (*InsertResult, error) {
	out := new(InsertResult)
	err := c.cc.Invoke(ctx, "/profileDao.ProfileRepo/InsertProfileToStorage", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *profileRepoClient) DeleteProfileFromStorage(ctx context.Context, in *ProfileName, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/profileDao.ProfileRepo/DeleteProfileFromStorage", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *profileRepoClient) UpdateProfileInStorage(ctx context.Context, in *ProfileUpdate, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/profileDao.ProfileRepo/UpdateProfileInStorage", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *profileRepoClient) ListProfilesFromStorage(ctx context.Context, in *Page, opts ...grpc.CallOption) (*SliceOfProfileObj, error) {
	out := new(SliceOfProfileObj)
	err := c.cc.Invoke(ctx, "/profileDao.ProfileRepo/ListProfilesFromStorage", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ProfileRepoServer is the server API for ProfileRepo service.
// All implementations must embed UnimplementedProfileRepoServer
// for forward compatibility
type ProfileRepoServer interface {
	FindProfileInStorage(context.Context, *ProfileName) (*ProfileObj, error)
	InsertProfileToStorage(context.Context, *ProfileObj) (*InsertResult, error)
	DeleteProfileFromStorage(context.Context, *ProfileName) (*Empty, error)
	UpdateProfileInStorage(context.Context, *ProfileUpdate) (*Empty, error)
	ListProfilesFromStorage(context.Context, *Page) (*SliceOfProfileObj, error)
	mustEmbedUnimplementedProfileRepoServer()
}

// UnimplementedProfileRepoServer must be embedded to have forward compatible implementations.
type UnimplementedProfileRepoServer struct {
}

func (UnimplementedProfileRepoServer) FindProfileInStorage(context.Context, *ProfileName) (*ProfileObj, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FindProfileInStorage not implemented")
}
func (UnimplementedProfileRepoServer) InsertProfileToStorage(context.Context, *ProfileObj) (*InsertResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method InsertProfileToStorage not implemented")
}
func (UnimplementedProfileRepoServer) DeleteProfileFromStorage(context.Context, *ProfileName) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteProfileFromStorage not implemented")
}
func (UnimplementedProfileRepoServer) UpdateProfileInStorage(context.Context, *ProfileUpdate) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateProfileInStorage not implemented")
}
func (UnimplementedProfileRepoServer) ListProfilesFromStorage(context.Context, *Page) (*SliceOfProfileObj, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListProfilesFromStorage not implemented")
}
func (UnimplementedProfileRepoServer) mustEmbedUnimplementedProfileRepoServer() {}

// UnsafeProfileRepoServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ProfileRepoServer will
// result in compilation errors.
type UnsafeProfileRepoServer interface {
	mustEmbedUnimplementedProfileRepoServer()
}

func RegisterProfileRepoServer(s grpc.ServiceRegistrar, srv ProfileRepoServer) {
	s.RegisterService(&ProfileRepo_ServiceDesc, srv)
}

func _ProfileRepo_FindProfileInStorage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ProfileName)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProfileRepoServer).FindProfileInStorage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/profileDao.ProfileRepo/FindProfileInStorage",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProfileRepoServer).FindProfileInStorage(ctx, req.(*ProfileName))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProfileRepo_InsertProfileToStorage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ProfileObj)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProfileRepoServer).InsertProfileToStorage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/profileDao.ProfileRepo/InsertProfileToStorage",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProfileRepoServer).InsertProfileToStorage(ctx, req.(*ProfileObj))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProfileRepo_DeleteProfileFromStorage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ProfileName)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProfileRepoServer).DeleteProfileFromStorage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/profileDao.ProfileRepo/DeleteProfileFromStorage",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProfileRepoServer).DeleteProfileFromStorage(ctx, req.(*ProfileName))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProfileRepo_UpdateProfileInStorage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ProfileUpdate)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProfileRepoServer).UpdateProfileInStorage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/profileDao.ProfileRepo/UpdateProfileInStorage",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProfileRepoServer).UpdateProfileInStorage(ctx, req.(*ProfileUpdate))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProfileRepo_ListProfilesFromStorage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Page)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProfileRepoServer).ListProfilesFromStorage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/profileDao.ProfileRepo/ListProfilesFromStorage",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProfileRepoServer).ListProfilesFromStorage(ctx, req.(*Page))
	}
	return interceptor(ctx, in, info, handler)
}

// ProfileRepo_ServiceDesc is the grpc.ServiceDesc for ProfileRepo service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ProfileRepo_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "profileDao.ProfileRepo",
	HandlerType: (*ProfileRepoServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "FindProfileInStorage",
			Handler:    _ProfileRepo_FindProfileInStorage_Handler,
		},
		{
			MethodName: "InsertProfileToStorage",
			Handler:    _ProfileRepo_InsertProfileToStorage_Handler,
		},
		{
			MethodName: "DeleteProfileFromStorage",
			Handler:    _ProfileRepo_DeleteProfileFromStorage_Handler,
		},
		{
			MethodName: "UpdateProfileInStorage",
			Handler:    _ProfileRepo_UpdateProfileInStorage_Handler,
		},
		{
			MethodName: "ListProfilesFromStorage",
			Handler:    _ProfileRepo_ListProfilesFromStorage_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "profile.proto",
}
