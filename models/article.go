package models

import (
	log "github.com/sirupsen/logrus"
)

const (
	// ArticleTitleLengthMax 文章最多标题数
	ArticleTitleLengthMax = 100

	// ArticlePreviewTextLengthMax 文章预览最多字符数
	ArticlePreviewTextLengthMax = 600

	qGetArticleCount  = "SELECT COUNT(*) as count FROM article"
	qGetAllArticle    = "SELECT * FROM article"
	qAddArticle       = "INSERT INTO article (create_at, update_at, visit_count, reply_count, article_title, article_previewtext, article_content, top, category, tag_list) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	qGetArticleByPage = `SELECT id, create_at, update_at, visit_count, reply_count, article_title, article_previewtext, article_content, top, category, tag_list 
						FROM article 
							ORDER BY id DESC LIMIT ?, ?`
	qGetArticleByID = "SELECT id, create_at, update_at, visit_count, reply_count, article_title, article_previewtext, article_content, top, category, tag_list FROM article WHERE id = ?"
)

// Article 文章的数据结构
type Article struct {
	ID                 uint   `db:"id" json:"id"`
	CreateAt           string `db:"create_at" json:"create_at"`
	UpdateAt           string `db:"update_at" json:"update_at"`
	VisitCount         uint   `db:"visit_count" json:"visit_count"`
	ReplyCount         uint   `db:"reply._count" json:"reply_count"`
	ArticleTitle       string `db:"article_title" json:"article_title"`
	ArticlePreviewText string `db:"article_previewtext" json:"article_previewtext"`
	ArticleContent     string `db:"article_content" json:"article_content"`
	Top                bool   `db:"top" json:"top"`
	Category           uint   `db:"category" json:"category"`
	TagList            string `db:"tag_list" json:"tag_list"`
}

// AddArticle 增加文章
func AddArticle(article *Article) (int64, error) {
	stmt, err := DB.Prepare(qAddArticle)
	if err != nil {
		log.WithFields(log.Fields{
			"errorMsg": err,
			"article":  *article,
		}).Info("Sql prepare failed")

		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(article.CreateAt, article.UpdateAt, article.VisitCount, article.ReplyCount, article.ArticleTitle, article.ArticlePreviewText, article.ArticleContent, article.Top, article.Category, article.TagList)
	if err != nil {
		log.WithFields(log.Fields{
			"errorMsg": err,
			"article":  *article,
		}).Info("Sql Exec failed")
		return 0, err
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		log.WithFields(log.Fields{
			"errorMsg": err,
			"article":  *article,
		}).Info("LastInsertId Exec failed")

		return 0, err
	}

	return lastID, nil
}

// GetArticleByID 根据 id 查询
func GetArticleByID(id uint64) (*Article, error) {
	var article Article

	row := DB.QueryRow(qGetArticleByID, uint(id))
	err := row.Scan(&article.ID, &article.CreateAt, &article.UpdateAt, &article.VisitCount, &article.ReplyCount, &article.ArticleTitle, &article.ArticlePreviewText, &article.ArticleContent, &article.Top, &article.Category, &article.TagList)
	if err != nil {
		log.WithFields(log.Fields{
			"errorMsg": err,
		}).Info("Sql query failed")

		return nil, err
	}

	return &article, nil
}

// GetArticleByPage 分页查询
func GetArticleByPage(limit int64, page int64) (*[]Article, error) {
	var article Article
	var articles []Article

	// 查询区间
	crows, err := DB.Query(qGetArticleCount)
	if err != nil {
		log.WithFields(log.Fields{
			"errorMsg": err,
		}).Info("Query article count failed")

		return nil, err
	}
	defer crows.Close()

	var counts int64

	for crows.Next() {
		err = crows.Scan(&counts)
	}
	if err != nil {
		log.WithFields(log.Fields{
			"errorMsg": err,
		}).Info("Scan article count failed")

		return nil, err
	}

	var qoffset int64
	var qlimit int64

	if limit*page > counts {
		qoffset = limit*(page-1) + 1
		qlimit = counts - limit*(page-1)
	}

	qoffset = limit*(page-1) + 1
	qlimit = limit

	rows, err := DB.Query(qGetArticleByPage, qoffset, qlimit)
	// rows, err := DB.Query(qGetArticleByPage, 0, 30)
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

// GetAllArticle 取出所有博文
func GetAllArticle() (*[]Article, error) {
	var article Article
	var articles []Article

	// 查询所有博文
	rows, err := DB.Query(qGetAllArticle)
	if err != nil {
		log.WithFields(log.Fields{
			"errorMsg": err,
		}).Info("DB query all article failed")

		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		// 这里以后优化性能
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
