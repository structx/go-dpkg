package controller_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/trevatk/go-pkg/http/controller"
	"github.com/trevatk/go-pkg/logging"
)

func init() {
	_ = os.Setenv("LOG_LEVEL", "DEBUG")
	_ = os.Setenv("LOG_PATH", "pkg.log")
}

type BundleControllerSuite struct {
	suite.Suite
	bc *controller.Bundle
}

func (suite *BundleControllerSuite) SetupTest() {

	assert := assert.New(suite.T())

	logger, err := logging.NewLoggerFromEnv()
	assert.NoError(err)

	suite.bc = controller.NewBundle(logger)
}

func (suite *BundleControllerSuite) TestHealth() {

	assert := assert.New(suite.T())

	testcases := []struct {
		expected int
	}{
		{
			expected: http.StatusOK,
		},
	}

	for _, testcase := range testcases {

		rr := httptest.NewRecorder()

		request, err := http.NewRequest(http.MethodGet, "/health", nil)
		assert.NoError(err)

		suite.bc.Health(rr, request)

		assert.Equal(testcase.expected, rr.Code)
	}
}

func TestBundleControllerSuite(t *testing.T) {
	suite.Run(t, new(BundleControllerSuite))
}
