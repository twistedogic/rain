package archive

import (
	"context"
	"testing"
	"time"
)

func Test_Client_GetAllSymbols(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()
	symbols, err := New().GetAllSymbols(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(symbols) == 0 {
		t.Fatal(symbols)
	}
}
