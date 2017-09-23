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
	qGetAllCategory = "SELECT category FROM category_list"
	qDeleteCategory = "DELETE FROM category_list WHERE category = ?"
	qUpdateCategory = "UPDATE category_list SET category = ? WHERE category = ?"
)

// AddCategory 增加分类名
func AddCategory(category string) error {
	stmt, err := DB.Prepare(qAddCategory)
	if err != nil {
		log.WithFields(log.Fields{
			"errorMsg": err,
		}).Info("Sql prepare failed")

		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(category)
	if err != nil {
		log.WithFields(log.Fields{
			"errorMsg": err,
		}).Info("Sql Exec failed")

		return err
	}

	return nil
}

// 获取全部分类名
func GetAllCategory() ([]string, error) {
	categories := make([]string, 0)
	rows, err := DB.Query(qGetAllCategory)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var category string
		err = rows.Scan(&category)
		log.Println(category)
		categories = append(categories, category)
	}
	err = rows.Err()
	if err != nil {
		log.Println(err)
	}

	log.Println(categories)
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
