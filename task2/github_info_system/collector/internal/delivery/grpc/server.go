package grpc

import (
	"context"
	"errors"
	"time"

	"local/github_info_system/collector/internal/domain"
	"local/github_info_system/collector/internal/usecase"
	"local/github_info_system/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	proto.UnimplementedCollectorServiceServer
	useCase *usecase.RepositoryUseCase
}

func NewServer(useCase *usecase.RepositoryUseCase) *Server {
	return &Server{
		useCase: useCase,
	}
}

func (s *Server) GetRepositoryInfo(ctx context.Context, req *proto.RepoRequest) (*proto.RepoResponse, error) {
	if req.Owner == "" || req.Repo == "" {
		return nil, status.Error(codes.InvalidArgument, "Owner and repo are required")
	}

	repo, err := s.useCase.GetRepositoryInfo(req.Owner, req.Repo)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrRepositoryNotFound):
			return nil, status.Error(codes.NotFound, err.Error())
		case errors.Is(err, domain.ErrGitHubAPI):
			return nil, status.Error(codes.Internal, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &proto.RepoResponse{
		Name:        repo.Name,
		Description: repo.Description,
		Stars:       int32(repo.Stars),
		Forks:       int32(repo.Forks),
		CreatedAt:   repo.CreatedAt.Format(time.RFC3339),
	}, nil
}
