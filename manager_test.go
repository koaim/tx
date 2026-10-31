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

type repo struct{}

func (r *repo) Ok()          {}
func (r *repo) Error() error { return repoErr }

func TestInTx_OK(t *testing.T) {
	t.Parallel()

	r := repo{}
	mngr := tx.NewFakeManager(r)

	err := mngr.InTx(t.Context(), func(ctx context.Context, r repo) error {
		r.Ok()
		return nil
	})

	if err != nil {
		t.Fatal(err)
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
		t.Fatalf("\nwant = %v\ngot  = %v", repoErr, err)
	}
}
