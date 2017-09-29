/*
* 检查提交是否跟 redis 里或者 mysql里一些数据重复
*
* author: ztplz
* email: mysticzt@gmail.com
* github: https://github.com/ztplz
* create-at: 2017.08.15
 */

package utils

import (
	"github.com/ztplz/blog-server/models"
)

func CheckInHash(key string, field string, id uint) (bool, error) {
	b, err := models.RedisClient.HExists(string(id)).Result()
	if err != nil {
		categories, err := models.GetAllCategory()
		for _, category
	}
	
	return b, nil
}
