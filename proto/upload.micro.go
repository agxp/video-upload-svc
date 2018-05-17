// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: proto/upload.proto

/*
Package video_upload is a generated protocol buffer package.

It is generated from these files:
	proto/upload.proto

It has these top-level messages:
	Request
	Response
	UploadRequest
	UploadResponse
	PropertyRequest
	PropertyResponse
	UploadFinishRequest
	UploadFinishResponse
*/
package video_upload

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	client "github.com/micro/go-micro/client"
	server "github.com/micro/go-micro/server"
	context "context"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ client.Option
var _ server.Option

// Client API for Upload service

type UploadClient interface {
	S3Request(ctx context.Context, in *Request, opts ...client.CallOption) (*Response, error)
	UploadFile(ctx context.Context, in *UploadRequest, opts ...client.CallOption) (*UploadResponse, error)
	WriteVideoProperties(ctx context.Context, in *PropertyRequest, opts ...client.CallOption) (*PropertyResponse, error)
	UploadFinish(ctx context.Context, in *UploadFinishRequest, opts ...client.CallOption) (*UploadFinishResponse, error)
}

type uploadClient struct {
	c           client.Client
	serviceName string
}

func NewUploadClient(serviceName string, c client.Client) UploadClient {
	if c == nil {
		c = client.NewClient()
	}
	if len(serviceName) == 0 {
		serviceName = "video_upload"
	}
	return &uploadClient{
		c:           c,
		serviceName: serviceName,
	}
}

func (c *uploadClient) S3Request(ctx context.Context, in *Request, opts ...client.CallOption) (*Response, error) {
	req := c.c.NewRequest(c.serviceName, "Upload.S3Request", in)
	out := new(Response)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *uploadClient) UploadFile(ctx context.Context, in *UploadRequest, opts ...client.CallOption) (*UploadResponse, error) {
	req := c.c.NewRequest(c.serviceName, "Upload.UploadFile", in)
	out := new(UploadResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *uploadClient) WriteVideoProperties(ctx context.Context, in *PropertyRequest, opts ...client.CallOption) (*PropertyResponse, error) {
	req := c.c.NewRequest(c.serviceName, "Upload.WriteVideoProperties", in)
	out := new(PropertyResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *uploadClient) UploadFinish(ctx context.Context, in *UploadFinishRequest, opts ...client.CallOption) (*UploadFinishResponse, error) {
	req := c.c.NewRequest(c.serviceName, "Upload.UploadFinish", in)
	out := new(UploadFinishResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Upload service

type UploadHandler interface {
	S3Request(context.Context, *Request, *Response) error
	UploadFile(context.Context, *UploadRequest, *UploadResponse) error
	WriteVideoProperties(context.Context, *PropertyRequest, *PropertyResponse) error
	UploadFinish(context.Context, *UploadFinishRequest, *UploadFinishResponse) error
}

func RegisterUploadHandler(s server.Server, hdlr UploadHandler, opts ...server.HandlerOption) {
	s.Handle(s.NewHandler(&Upload{hdlr}, opts...))
}

type Upload struct {
	UploadHandler
}

func (h *Upload) S3Request(ctx context.Context, in *Request, out *Response) error {
	return h.UploadHandler.S3Request(ctx, in, out)
}

func (h *Upload) UploadFile(ctx context.Context, in *UploadRequest, out *UploadResponse) error {
	return h.UploadHandler.UploadFile(ctx, in, out)
}

func (h *Upload) WriteVideoProperties(ctx context.Context, in *PropertyRequest, out *PropertyResponse) error {
	return h.UploadHandler.WriteVideoProperties(ctx, in, out)
}

func (h *Upload) UploadFinish(ctx context.Context, in *UploadFinishRequest, out *UploadFinishResponse) error {
	return h.UploadHandler.UploadFinish(ctx, in, out)
}
