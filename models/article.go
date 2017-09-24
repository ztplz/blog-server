package models

import log "github.com/sirupsen/logrus"

const (
	// ArticleTitleLengthMax 文章最多标题数
	ArticleTitleLengthMax = 100

	// ArticlePreviewTextLengthMax 文章预览最多字符数
	ArticlePreviewTextLengthMax = 300

	qAddArticle       = "INSERT INTO article (create_at, update_at, visit_count, reply_count, article_title, article_previewtext, article_content, top, category, tag_list) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	qGetArticleByPage = `SELECT id, create_at, update_at, visit_count, reply_count, article_title, article_previewtext, article_content, top, category, tag_list 
						FROM article 
							ORDER BY id DESC LIMIT ?, ?`
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
	TagList            string `db:"tag_list" json:"tag_list"`
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

// GetArticleByPage 分页查询
func GetArticleByPage(page int64, limit int64) (*[]Article, error) {
	var article Article
	var articles []Article
	// 查询偏移量
	offset := (page - 1) * limit

	rows, err := DB.Query(qGetArticleByPage, offset, limit)
	if err != nil {
		log.WithFields(log.Fields{
			"errorMsg": err,
			"page":     page,
			"limit":    limit,
		}).Info("DB query article failed")

		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		// 这里以后优化性能
		// var id 					uint
		// var createAt 			string
		// var updateAt 			string
		// var visitCount 			uint
		// var replyCount 			uint
		// var articleTitle 		string
		// var articlePreviewtext 	string
		// var articleContent 		string
		// var top 				bool
		// var category 			string
		// var tagList 			string

		err = rows.Scan(
			&article.ID,
			&article.CreateAt,
			&article.UpdateAt,
			&article.VisitCount,
			&article.ReplyCount,
			&article.ArticleTitle,
			&article.ArticlePreviewText,
			&article.ArticleContent,
			&article.Top,
			&article.Category,
			&article.TagList)
		articles = append(articles, article)
	}
	err = rows.Err()
	if err != nil {
		log.WithFields(log.Fields{
			"errorMsg": err,
		}).Info("Rows scan failed")

		return nil, err
	}

	return &articles, nil
}
