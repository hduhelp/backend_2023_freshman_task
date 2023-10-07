// Package DelTx 更新表单实现id更新操作 同时利用Transaction进行并发控制
package DelTx

import (
	"HDUhelper_Todo/models"
)

func DelTx(id int) error {

	db := models.DbConnect()

	tx := db.Begin() // 开始事务

	// 设置事务的隔离级别为 SERIALIZABLE 以保证数据一致性
	tx = tx.Set("gorm:query_option", "SET TRANSACTION ISOLATION LEVEL SERIALIZABLE")

	err := tx.Where("Id = ?", id).Delete(models.ListItem{}).Error //完成delete操作
	if err != nil {
		tx.Rollback() // 发生错误，回滚事务
		return err
	}

	// 创建临时表
	err = tx.Exec("CREATE TEMPORARY TABLE temp_table AS SELECT * FROM list_items").Error
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Exec("DELETE FROM list_items").Error // 删除原始表中的所有数据
	if err != nil {
		tx.Rollback()
		return err
	}

	// 更新临时表中的自增字段值
	err = tx.Exec("UPDATE temp_table SET id = id - 1 WHERE id > ?", id).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	// 将临时表中的数据重新插入原始表中
	err = tx.Exec("INSERT INTO list_items SELECT * FROM temp_table").Error
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Exec("DROP TABLE temp_table").Error // 删除临时表
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit().Error // 提交delete操作
	if err != nil {
		return err
	}

	return nil
}
