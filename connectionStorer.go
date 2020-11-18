package main

import "context"

type connectionStorer interface {
	GetConnectionIDs(ctx context.Context) ([]string, error)
	AddConnectionID(ctx context.Context, connectionID string) error
	MarkConnectionIDDisconnected(ctx context.Context, connectionID string) error
}