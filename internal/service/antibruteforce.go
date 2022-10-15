package service

import (
	"context"
	"errors"
	"net"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/make-it-git/otus-antibruteforce/internal/contract"
	api "github.com/make-it-git/otus-antibruteforce/pkg/antibruteforce/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Service struct {
	api.UnimplementedAntiBruteforceServer
	storage contract.NetAddrStorage
	bucket  contract.LeakyBucket
}

func NewService(storage contract.NetAddrStorage, bucket contract.LeakyBucket) *Service {
	return &Service{
		storage: storage,
		bucket:  bucket,
	}
}

func errInvalidArg(err error) error {
	return status.Error(codes.InvalidArgument, err.Error())
}

func errInternal(err error) error {
	return status.Error(codes.Internal, err.Error())
}

var errInvalidIPAddress = errors.New("invalid ip address")

func accepted() *api.AuthCheckResponse {
	return &api.AuthCheckResponse{Accepted: true}
}

func declined() *api.AuthCheckResponse {
	return &api.AuthCheckResponse{Accepted: false}
}

func (s *Service) BlackListExtend(_ context.Context, request *api.SubnetAddress) (*empty.Empty, error) {
	_, ipNet, err := net.ParseCIDR(request.GetSubnetAddress())
	if err != nil {
		return nil, errInvalidArg(err)
	}

	err = s.storage.BlackListAdd(ipNet.String())
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *Service) BlackListRemove(_ context.Context, request *api.SubnetAddress) (*empty.Empty, error) {
	_, ipNet, err := net.ParseCIDR(request.GetSubnetAddress())
	if err != nil {
		return nil, errInvalidArg(err)
	}

	err = s.storage.BlackListRemove(ipNet.String())
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *Service) WhiteListAdd(_ context.Context, request *api.SubnetAddress) (*empty.Empty, error) {
	_, ipNet, err := net.ParseCIDR(request.GetSubnetAddress())
	if err != nil {
		return nil, errInvalidArg(err)
	}

	err = s.storage.WhiteListAdd(ipNet.String())
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *Service) WhiteListRemove(_ context.Context, request *api.SubnetAddress) (*empty.Empty, error) {
	_, ipNet, err := net.ParseCIDR(request.GetSubnetAddress())
	if err != nil {
		return nil, errInvalidArg(err)
	}

	err = s.storage.WhiteListRemove(ipNet.String())
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *Service) ClearLists(_ context.Context, _ *empty.Empty) (*empty.Empty, error) {
	if err := s.storage.ClearLists(); err != nil {
		return nil, errInvalidArg(err)
	}

	return &emptypb.Empty{}, nil
}

func (s *Service) AuthCheck(_ context.Context, request *api.AuthCheckRequest) (*api.AuthCheckResponse, error) {
	ip := net.ParseIP(request.GetIp())
	if len(ip) == 0 {
		return nil, errInvalidArg(errInvalidIPAddress)
	}

	statusIP, err := s.storage.GetStatus(ip)
	if err != nil {
		return nil, errInternal(err)
	}

	switch statusIP {
	case contract.Blacklisted:
		return declined(), nil
	case contract.Whitelisted:
		return accepted(), nil
	case contract.Unknown:
		break
	}

	err = s.bucket.Try(request.GetLogin(), request.GetPassword(), request.GetIp())
	if errors.Is(err, contract.ErrDeclined) {
		return declined(), nil
	} else if err != nil {
		return nil, errInternal(err)
	}

	return accepted(), nil
}
