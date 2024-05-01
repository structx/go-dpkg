package controller_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/structx/go-pkg/adapter/logging"
	"github.com/structx/go-pkg/adapter/port/http/controller"
	"github.com/structx/go-pkg/adapter/setup"
	"github.com/structx/go-pkg/util/decode"
)

func init() {
	_ = os.Setenv("DSERVICE_CONFIG", "./testfiles/controller.test.hcl")
}

type BundleControllerSuite struct {
	suite.Suite
	bc *controller.Bundle
}

func (suite *BundleControllerSuite) SetupTest() {

	assert := assert.New(suite.T())

	_ = os.Mkdir("./testfiles/log", os.ModePerm)

	cfg := setup.New()
	assert.NoError(decode.ConfigFromEnv(cfg))

	logger, err := logging.New(cfg)
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
