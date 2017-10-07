package models

import (
	log "github.com/sirupsen/logrus"
)

const (
	qUpdateVisitorCount = "UPDATE visitor_count SET count = count + 1 WHERE id = 1"
	qGetAllVisitCount   = "SELECT count FROM visitor_count WHERE id =1"
)

// CountVistor 数据库访问人数加 1
func CountVistor() error {
	// sql语句预处理
	stmt, err := DB.Prepare(qUpdateVisitorCount)
	if err != nil {
		log.WithFields(log.Fields{
			"errorMsg": err,
		}).Fatal("Sql prepare failed")
	}
	defer stmt.Close()

	_, err = stmt.Exec()
	if err != nil {
		log.WithFields(log.Fields{
			"errorMsg": err,
		}).Fatal("Increase visit count failed")
	}

	return err
}

// GetAllVisitorCount 获取数据所有访客人数
func GetAllVisitorCount() (uint, error) {
	var count uint

	err := DB.QueryRow(qGetAllVisitCount).Scan(&count)
	if err != nil {
		log.WithFields(log.Fields{
			"errorMsg": err,
		}).Info("Query visitor count failed")

		return 0, err
	}

	return count, nil
}
