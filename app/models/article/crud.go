package article

import (
	"goblog/pkg/model"
	"goblog/pkg/types"
)

// Article 文章模型
type Article struct {
	ID    uint64
	Title string
	Body  string
}

// Get 通过 ID 获取文章
func Get(idstr string) (Article, error) {
	var article Article
	id := types.StringToUint64(idstr)
	if err := model.DB.First(&article, id).Error; err != nil {
		return article, err
	}

	return article, nil
}
