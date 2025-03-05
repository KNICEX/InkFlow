package dao

import (
	"context"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"io"
	"strings"
	"testing"
	"time"
)

type EsTestSuite struct {
	suite.Suite
	es *elasticsearch.Client
}

func initEs(t *testing.T) *elasticsearch.Client {
	c, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{"http://localhost:9200"},
	})
	require.NoError(t, err)
	return c
}
func TestSuite(t *testing.T) {
	suite.Run(t, new(EsTestSuite))
}

func (s *EsTestSuite) SetupTest() {
	s.es = initEs(s.T())
}

func (s *EsTestSuite) TestCreateIndex() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := s.es.Indices.Create("user_idx",
		s.es.Indices.Create.WithContext(ctx),
		s.es.Indices.Create.WithBody(strings.NewReader(userIndexMapping)),
	)
	require.NoError(s.T(), err)
	require.NotNil(s.T(), res)
	bs, err := io.ReadAll(res.Body)
	require.NoError(s.T(), err)
	fmt.Println(string(bs))
}

func (s *EsTestSuite) TestGetDoc() {

}
