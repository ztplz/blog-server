package models

import (
	"log"
)

// 定义Tag 结构
type Tag struct {
	ID       uint   `db:"id" json:"id"`
	Color    string `db:"color" json:"color"`
	TagTitle string `db:"tag_title" json:"tag_title"`
}

const (
	qGetAllTag = "SELECT * FROM tag"
	qAddTag    = "INSERT INTO tag (color, tag_title) VALUES (?, ?)"
)

// 获取全部标签
func GetAllTag() ([]Tag, error) {
	tags := make([]Tag, 0)
	rows, err := DB.Query(qGetAllTag)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		// tag := new(Tag)
		id := 0
		color := ""
		tagTitle := ""
		err = rows.Scan(&id, &color, &tagTitle)
		// err = rows.Scan(tag.ID, tag.Color, tag.TagTitle)
		// err = rows.Scan(tag)
		tags = append(tags, Tag{
			ID:       uint(id),
			Color:    color,
			TagTitle: tagTitle,
		})
		// tags = append(tags, &tag)
	}
	err = rows.Err()
	if err != nil {
		log.Println(err)
	}

	return tags, err
}

// 增加标签
func AddTag(color string, title string) error {
	stmt, err := DB.Prepare(qAddTag)
	defer stmt.Close()
	if err != nil {
		return err
	}
	_, err = stmt.Exec(color, title)
	if err != nil {
		return err
	}

	return nil
}
