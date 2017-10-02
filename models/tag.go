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
	qAddTag    = "INSERT INTO tags (color, tag_title) VALUES (?, ?)"
)

// GetAllTag 获取全部标签
func GetAllTag() (*[]Tag, error) {
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
		var id uint
		var color string
		var tagTitle string
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

	return &tags, err
}

// AddTag 增加标签
func AddTag(color string, title string) (int64, error) {
	// 向数据库插入标签
	stmt, err := DB.Prepare(qAddTag)
	if err != nil {
		log.WithFields(log.Fields{
			"errorMsg": err,
			"color":    color,
			"title":    title,
		}).Info("Sql prepare failed")

		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(color, title)
	if err != nil {
		log.WithFields(log.Fields{
			"errorMsg": err,
			"color":    color,
			"title":    title,
		}).Info("Sql exec failed")

		return 0, err
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		log.WithFields(log.Fields{
			"errorMsg": err,
			"color":    color,
			"title":    title,
		}).Info("LastInsertId Exec failed")

		return 0, err
	}

	log.Info(lastID)

	return lastID, nil
}

// 修改标签
