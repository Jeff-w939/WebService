package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

type Student struct {
	Id   int
	Name string
	Age  int
}

var GlobalConn *gorm.DB

func main() {
	// 连接数据库mysql, gorm框架建好的Mysql数据库连接conn, 就是一个连接池的句柄
	conn, err := gorm.Open("mysql", "root:199791@tcp(127.0.0.1:3306)/WService")
	if err != nil {
		fmt.Println("gorm.open err:", err)
		return
	}

	GlobalConn = conn
	GlobalConn.DB().SetMaxIdleConns(10)
	GlobalConn.DB().SetConnMaxLifetime((100))

	//表示不创建复数表名
	GlobalConn.SingularTable(true)
	// 用gorm创建数据库表
	fmt.Println(GlobalConn.AutoMigrate(new(Student)).Error)

	// 插入数据
	InsertData()
}

// 插入数据
func InsertData() {
	// 先创建数据
	var stu Student
	stu.Name = "wangbaba"
	stu.Age = 24

	//再插入数据
	fmt.Println(GlobalConn.Create(&stu).Error)
}

//查找数据
func SearchData() {
	var stu []Student

	////查询第一条全部信息
	//GlobalConn.First(&stu)
	//
	////查询指定字段的第一条信息
	//GlobalConn.Select("name", "age").First(&stu)

	//查询指定字段的所有信息
	GlobalConn.Select("name", "age").Find(&stu)

	//查询带有条件的字段 where
	GlobalConn.Select("name", "age").Where("name = ?", "wangbaba").Find(&stu)
}

//更新数据
func UpdateData() {
	//更新一个字段 update
	fmt.Println(GlobalConn.Model(new(Student)).Where("name = ?", "wangbaba").Update("name", "wangguo").Error)

	//更新多个字段 updates
	fmt.Println(GlobalConn.Model(new(Student)).Where("name = ?", "wangbaba").Updates(map[string]interface{}{"name": "wang", "age": 22}).Error)
}

//删除数据
func DeleteData() {
	// 软删除， 并不是真正的删除数据
	fmt.Println(GlobalConn.Where("name = ?", "wangbaba").Delete(new(Student)).Error)
}
