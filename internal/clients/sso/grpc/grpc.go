package grpc

import (
	"context"
	"fmt"
	"log/slog"

	ssov1 "github.com/kirill-dolgii/protos/gen/go/sso"
	"github.com/kirill-dolgii/url-shortner/internal/clients/sso/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	api    ssov1.AuthClient
	logger *slog.Logger
}

func New(
	ctx context.Context,
) (*Client, error) {
	const op = "grpc.New"

	cc, err := grpc.DialContext(insecure.NewCredentials())

	if err != nil {
		return nil, fmt.Errorf("%s, %w", op, err)
	}

	return &Client{}
}
