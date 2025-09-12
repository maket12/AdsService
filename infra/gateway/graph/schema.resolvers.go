package graph

import (
	authpb "AdsService/authservice/proto"
	"AdsService/infra/gateway/graph/model"
	"AdsService/infra/gateway/internal/authctx"
	userpb "AdsService/userservice/proto"
	"context"
	"fmt"
	"github.com/99designs/gqlgen/graphql"
	"google.golang.org/protobuf/types/known/emptypb"
	"io"
	"net/http"
)

func (r *mutationResolver) Register(ctx context.Context, email string, password string) (*model.AuthPayload, error) {
	if r.C.Auth == nil {
		return nil, fmt.Errorf("auth service is unvailable")
	}

	res, err := r.C.Auth.Register(ctx, &authpb.RegisterRequest{Email: email, Password: password})
	if err != nil {
		return nil, err
	}
	return &model.AuthPayload{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
	}, nil
}

func (r *mutationResolver) Login(ctx context.Context, email string, password string) (*model.AuthPayload, error) {
	if r.C.Auth == nil {
		return nil, fmt.Errorf("auth service is unvailable")
	}

	res, err := r.C.Auth.Login(ctx, &authpb.LoginRequest{Email: email, Password: password})
	if err != nil {
		return nil, err
	}
	return &model.AuthPayload{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
	}, nil
}

func (r *mutationResolver) UpdateProfile(ctx context.Context, name *string, phone *string) (*model.UserProfile, error) {
	if r.C.Auth == nil {
		return nil, fmt.Errorf("user service is unvailable")
	}

	ctx = authctx.MetadataFromRequest(ctx.Value("httpRequest").(*http.Request))
	p, err := r.C.User.UpdateProfile(ctx, &userpb.UpdateProfileRequest{
		Name:  getOrEmpty(name),
		Phone: getOrEmpty(phone),
	})
	if err != nil {
		return nil, err
	}
	return mapUserProfile(p), nil
}

func (r *mutationResolver) UploadMyPhoto(ctx context.Context, file graphql.Upload) (string, error) {
	if r.C.User == nil {
		return "", fmt.Errorf("user service is unavailable")
	}

	fileBytes, err := io.ReadAll(file.File)
	if err != nil {
		return "", fmt.Errorf("failed to read upload: %w", err)
	}

	ctx = authctx.MetadataFromRequest(ctx.Value("httpRequest").(*http.Request))

	res, err := r.C.User.UploadPhoto(ctx, &userpb.UploadPhotoRequest{
		Data:        fileBytes,
		Filename:    file.Filename,
		ContentType: file.ContentType,
	})
	if err != nil {
		return "", err
	}

	return res.PhotoId, nil
}

func (r *queryResolver) Me(ctx context.Context) (*model.UserProfile, error) {
	if r.C.Auth == nil {
		return nil, fmt.Errorf("user service is unvailable")
	}

	ctx = authctx.MetadataFromRequest(ctx.Value("httpRequest").(*http.Request))
	p, err := r.C.User.GetProfile(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, err
	}
	return mapUserProfile(p), nil
}

func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
