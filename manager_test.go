package tx_test

import (
	"context"
	"errors"
	"testing"

	"github.com/koaim/tx"
)

var (
	repoErr = errors.New("error")
)

type repo struct{ okCalls int }

func (r *repo) Ok()          { r.okCalls++ }
func (r *repo) Error() error { return repoErr }

func TestInTx_OK(t *testing.T) {
	t.Parallel()

	r := &repo{}
	mngr := tx.NewFakeManager(r)

	err := mngr.InTx(t.Context(), func(ctx context.Context, r *repo) error {
		r.Ok()
		return nil
	})

	if err != nil {
		t.Fatal(err)
	}
	if r.okCalls != 1 {
		t.Fatalf("\nwant: %v\nhave: %v", 1, r.okCalls)
	}
}

func TestInTx_Err(t *testing.T) {
	t.Parallel()

	r := repo{}
	mngr := tx.NewFakeManager(r)

	err := mngr.InTx(t.Context(), func(ctx context.Context, r repo) error {
		return r.Error()
	})

	if !errors.Is(err, repoErr) {
		t.Fatalf("\nwant: %q\nhave: %q", repoErr, err)
	}
}
