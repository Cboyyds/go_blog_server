package database

import (
	"context"
	"server/global"
	"server/model/elasticsearch"

	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/scriptlanguage"
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

// Comment 评论表
type Comment struct {
	global.MODEL
	ArticleID string    `json:"article_id"`                     // 文章 ID
	PID       *uint     `json:"p_id"`                           // 父评论 ID
	PComment  *Comment  `json:"-" gorm:"foreignKey:PID"`        // 序列化时不包含这个json字段
	Children  []Comment `json:"children" gorm:"foreignKey:PID"` // 子评论，写PID是因为它也是个评论，要从关联同一个表（自引用）, GORM 中使用   foreignKey   标签时，GORM 会尝试在数据库中创建外键约束。
	UserUUID  uuid.UUID `json:"user_uuid" gorm:"type:char(36)"` // 用户 uuid
	//  github.com/gofrs/uuid   是一个在GitHub上的Go语言库，它提供了生成和操作UUID（Universally Unique Identifier，通用唯一识别码）的功能。UUID是一个128位的数值，通常用16个十六进制数字表示，用于确保在分布式系统中的唯一性。
	User    User   `json:"user" gorm:"foreignKey:UserUUID;references:UUID"` // 关联的用户
	Content string `json:"content"`                                         // 内容
}

// AfterCreate 钩子，创建后调用
func (c *Comment) AfterCreate(_ *gorm.DB) error {
	source := "ctx._source.comments += 1"
	script := types.Script{Source: &source, Lang: &scriptlanguage.Painless}
	_, err := global.ESClient.Update(elasticsearch.ArticleIndex(), c.ArticleID).Script(&script).Do(context.TODO())
	return err
}

// AfterDelete 钩子，删除后调用
func (c *Comment) BeforeDelete(_ *gorm.DB) error {
	var articleID string
	if err := global.DB.Model(&c).Pluck("article_id", &articleID).Error; err != nil {
		return err
	}
	source := "ctx._source.comments -= 1"
	script := types.Script{Source: &source, Lang: &scriptlanguage.Painless}
	_, err := global.ESClient.Update(elasticsearch.ArticleIndex(), articleID).Script(&script).Do(context.TODO())
	return err
}
