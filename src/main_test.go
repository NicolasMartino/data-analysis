package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHandleGet(t *testing.T) {
	expected := "some data"

	svr := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, expected)
		}))

	defer svr.Close()

	res, err := handleGet(svr.URL)

	fmt.Printf("%v", res.Status)
	require.NoError(t, err)
	require.Equal(t, res.Status, http.StatusOK)
	require.Equal(t, res.Body, expected)
	require.Equal(t, res.RequestUrl, svr.URL)
}
