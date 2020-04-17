package actions

import (
	"net/http"
	"testing"

	"github.com/gobuffalo/httptest"
	"github.com/stretchr/testify/require"
)

func Test_HomeHandler(t *testing.T) {
	r := require.New(t)
	w := httptest.New(App())

	res := w.HTML("/").Get()
	r.Equal(http.StatusOK, res.Code)
	r.Contains(res.Body.String(), "Welcome to Buffalo")
}
