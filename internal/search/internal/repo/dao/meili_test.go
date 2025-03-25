package dao

import (
	"fmt"
	"github.com/meilisearch/meilisearch-go"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

func initMeili() meilisearch.ServiceManager {
	return meilisearch.New("http://127.0.0.1:7700")
}

type Article struct {
	Id         int64     `json:"id"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	Tags       []string  `json:"tags"`
	LikeCnt    int64     `json:"like_cnt"`
	CommentCnt int64     `json:"comment_cnt"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func TestCreateIndex(t *testing.T) {
	cli := initMeili()
	res, err := cli.CreateIndex(&meilisearch.IndexConfig{
		Uid:        "test",
		PrimaryKey: "id",
	})
	require.NoError(t, err)
	fmt.Println(res)
	time.Sleep(time.Second * 3)

	idx, err := cli.GetIndex("test")
	require.NoError(t, err)
	fmt.Println(idx)

	_, err = cli.DeleteIndex("test")
	require.NoError(t, err)
}

type MeiliTestSuite struct {
	suite.Suite
	cli meilisearch.ServiceManager
}

func (s *MeiliTestSuite) SetupSuite() {
	s.cli = initMeili()
	_, err := s.cli.CreateIndex(&meilisearch.IndexConfig{
		Uid:        "test",
		PrimaryKey: "id",
	})
	require.NoError(s.T(), err)
	time.Sleep(time.Second * 2)
}

func TestMeiliSuite(t *testing.T) {
	suite.Run(t, new(MeiliTestSuite))
}
func (s *MeiliTestSuite) TearDownSuite() {
	_, err := s.cli.DeleteIndex("test")
	require.NoError(s.T(), err)
}

func (s *MeiliTestSuite) TestAddDocuments() {
	res, err := s.cli.Index("test").AddDocuments([]Article{
		{
			Id:         1,
			Title:      "what is web3",
			Content:    "test1 content",
			Tags:       []string{"web3", "eth", "btc"},
			LikeCnt:    100,
			CommentCnt: 10,
			CreatedAt:  time.Now().Add(-time.Hour * 24),
			UpdatedAt:  time.Now(),
		},
		{
			Id:         2,
			Title:      "how to draw with sai2",
			Content:    "test2 content",
			Tags:       []string{"sai", "draw", "art", "painting"},
			LikeCnt:    200,
			CommentCnt: 20,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			Id:    3,
			Title: "中文怎么样呢?",
			Content: "今天，我们来到了一个新的世界，" +
				"在这里，我们可以自由地表达自己，" +
				"我们可以用自己的语言，我想和你聊一聊关于未来的事情。",
			Tags:       []string{"中文", "哲学", "语言"},
			LikeCnt:    300,
			CommentCnt: 30,
			CreatedAt:  time.Now().Add(-time.Hour),
			UpdatedAt:  time.Now(),
		},
	})
	require.NoError(s.T(), err)
	fmt.Println(res)
	time.Sleep(time.Second * 2)

	fmt.Println("search \"web3\" ...")
	searchRes, er := s.cli.Index("test").Search("web3", &meilisearch.SearchRequest{})
	require.NoError(s.T(), er)
	fmt.Printf("search res: %+v\n", searchRes)

	fmt.Println("search \"自由\" ...")
	searchRes, er = s.cli.Index("test").Search("自由", &meilisearch.SearchRequest{})
	require.NoError(s.T(), er)
	fmt.Printf("search res: %+v\n", searchRes)

	fmt.Println("search \"哲学\" ...")
	searchRes, er = s.cli.Index("test").Search("哲学", &meilisearch.SearchRequest{})
	require.NoError(s.T(), er)
	fmt.Printf("search res: %+v\n", searchRes)

	fmt.Println("search \"未来 汽车\" ...")
	searchRes, er = s.cli.Index("test").Search("未来 汽车", &meilisearch.SearchRequest{})
	require.NoError(s.T(), er)
	fmt.Printf("search res: %+v\n", searchRes)
}
