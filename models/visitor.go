package models

import (
	"time"
  
	log "github.com/sirupsen/logrus"
)

const (
	qUpdateVisitorCount = "INSERT INTO visitor_count (date, count) VALUES (?, ?)"
	qGetAllVisitCount   = "SELECT count FROM visitor_count WHERE id =1"
)

// CountVistor 数据库按日期记录访问人数加
func CountVistor(count uint) error {
	// sql语句预处理
	stmt, err := DB.Prepare(qUpdateVisitorCount)
	if err != nil {
		log.WithFields(log.Fields{
			"errorMsg": err,
		}).Fatal("Sql prepare failed")
	}
	defer stmt.Close()

	date := time.Now().Format("2006-01-02")
	_, err = stmt.Exec(date, count)
	if err != nil {
		log.WithFields(log.Fields{
			"errorMsg": err,
		}).Fatal("Insert visit count failed")
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
