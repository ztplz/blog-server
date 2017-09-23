package models

import log "github.com/sirupsen/logrus"

const (
	qAddArticle = "INSERT INTO article (create_at, update_at, visit_count, reply_count, article_title, article_previewtext, article_content, top, category) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)"
)

// Article 文章的数据结构
type Article struct {
	ID                 uint   `db:"id" json:"id"`
	CreateAt           string `db:"create_at" json:"creat_at"`
	UpdateAt           string `db:"update_at" json:"update_at"`
	VisitCount         uint   `db:"visit_count" json:"visit_count"`
	ReplyCount         uint   `db:"reply._count" json:"reply_count"`
	ArticleTitle       string `db:"article_title" json:"article_title"`
	ArticlePreviewText string `db:"article_previewtext" json:"article_previewtext"`
	ArticleContent     string `db:"article_content" json:"article_content"`
	Top                bool   `db:"top" json:"top"`
	Category           string `db:"category" json:"category"`
	TagList            string `db:"tag_list" json:"tagList"`
}

// AddArticle 增加文章
func AddArticle(article *Article) error {
	stmt, err := DB.Prepare(qAddArticle)
	if err != nil {
		log.WithFields(log.Fields{
			"errorMsg": err.Error(),
		}).Info("Sql prepare failed")

		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(article.CreateAt, article.UpdateAt, article.VisitCount, article.ReplyCount, article.ArticleTitle, article.ArticlePreviewText, article.ArticleContent, article.Top, article.Category, article.TagList)
	if err != nil {
		log.WithFields(log.Fields{
			"errorMsg": err.Error(),
		}).Info("Sql Exec failed")
		return err
	}

	return nil
}
