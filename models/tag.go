package models

import (
	log "github.com/sirupsen/logrus"
)

// Tag 定义Tag的结构
type Tag struct {
	ID       uint   `db:"id" json:"id"`
	Color    string `db:"color" json:"color"`
	TagTitle string `db:"tag_title" json:"tag_title"`
	// ArticleID uint   `db:"article_id" json:"article_id"`
}

// sql 查询语句
const (
	qGetAllTag = "SELECT id, color, tag_title FROM tags"
	qAddTag    = "INSERT INTO tag (color, tag_title) VALUES (?, ?)"
)

// GetAllTag 获取全部标签
func GetAllTag() ([]Tag, error) {
	tags := make([]Tag, 0)

	// 从数据库查询所有不同的tag
	rows, err := DB.Query(qGetAllTag)
	if err != nil {
		log.WithFields(log.Fields{
			"errorMsg": err,
		}).Info("Query all tag failed")

		return nil, err
	}
	defer rows.Close()

	// 把结果写入到 tags 里
	for rows.Next() {
		id := 0
		color := ""
		tagTitle := ""
		err = rows.Scan(&id, &color, &tagTitle)
		tags = append(tags, Tag{
			ID:       uint(id),
			Color:    color,
			TagTitle: tagTitle,
		})
	}
	err = rows.Err()
	if err != nil {
		log.WithFields(log.Fields{
			"errorMsg": err,
		}).Info("Write result to tags failed")

		return nil, err
	}

	return tags, err
}

// AddTag 增加标签
func AddTag(color string, title string) error {
	// 向数据库插入标签
	stmt, err := DB.Prepare(qAddTag)
	if err != nil {
		log.WithFields(log.Fields{
			"errorMsg": err,
			"color":    color,
			"title":    title,
		}).Info("Sql prepare failed")

		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(color, title)
	if err != nil {
		log.WithFields(log.Fields{
			"errorMsg": err,
			"color":    color,
			"title":    title,
		}).Info("Insert tag to database failed")

		return err
	}

	return nil
}

// 修改标签
