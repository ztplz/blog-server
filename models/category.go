package models

import (
	"errors"

	log "github.com/sirupsen/logrus"
)

// Category 文章的数据结构
type Category struct {
	ID       uint   `db:"id" json:"id"`
	Category string `db:"category" json:"category"`
}

const (
	qAddCategory    = "INSERT INTO category_list (category) VALUES (?)"
	qGetAllCategory = "SELECT id, category FROM category_list"
	qDeleteCategory = "DELETE FROM category_list WHERE category = ?"
	qUpdateCategory = "UPDATE category_list SET category = ? WHERE category = ?"
)

// AddCategory 增加分类名
func AddCategory(category string) (int64, error) {
	stmt, err := DB.Prepare(qAddCategory)
	if err != nil {
		log.WithFields(log.Fields{
			"errorMsg": err,
			"category": category,
		}).Info("Sql prepare failed")

		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(category)
	if err != nil {
		log.WithFields(log.Fields{
			"errorMsg": err,
			"category": category,
		}).Info("Sql Exec failed")

		return 0, err
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		log.WithFields(log.Fields{
			"errorMsg": err,
			"category": category,
		}).Info("LastInsertId Exec failed")

		return 0, err
	}

	return lastID, nil
}

// GetAllCategory 获取全部分类名
func GetAllCategory() ([]Category, error) {
	var categories []Category

	rows, err := DB.Query(qGetAllCategory)
	if err != nil {
		log.WithFields(log.Fields{
			"errorMsg": err,
		}).Info("DB query all category failed")

		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var category Category
		err = rows.Scan(&category.ID, &category.Category)
		categories = append(categories, category)
	}
	err = rows.Err()
	if err != nil {
		log.WithFields(log.Fields{
			"errorMsg": err,
		}).Info("Rows scan failed")
	}

	return categories, err
}

// 删除某个分类名，删除分类名时务必删除关联 article 表的 category 字段相关的分类名
func DeleteCategory(category string) error {
	stmt, err := DB.Prepare(qDeleteCategory)
	defer stmt.Close()
	if err != nil {
		log.Println(err)
		return err
	}

	res, err := stmt.Exec(category)
	if err != nil {
		log.Println(err)
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		log.Println(err)
		return err
	}

	if rows == 0 {
		return errors.New("No Category Delete")
	}

	return nil
}

// 修改某个分类名
func UpdateCategory(category string, key string) error {
	stmt, err := DB.Prepare(qUpdateCategory)
	defer stmt.Close()
	if err != nil {
		log.Println(err)
		return err
	}

	res, err := stmt.Exec(key, category)
	if err != nil {
		log.Println(err)
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		log.Println(err)
		return err
	}

	if rows == 0 {
		return errors.New("No Category Update")
	}

	return nil
}
