package model

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/gomodule/redigo/redis"
)

// 创建全局redis 连接池句柄
var RedisPool redis.Pool

//创建函数，  初始化Redis连接池
func InitRedis() {
	RedisPool = redis.Pool{
		MaxIdle:         20,
		MaxActive:       50,
		MaxConnLifetime: 60 * 5,
		IdleTimeout:     60,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", "127.0.0.1:6379")
		},
	}
}

// 查询redis中的是否与图片中一致
func CheckImgCode(uuid, imgCode string) bool {
	// 链接 redis 从redispool 中取出一个链接
	conn := RedisPool.Get()
	defer conn.Close()

	// 查询redis 数据
	code, err := redis.String(conn.Do("get", uuid))
	if err != nil {
		fmt.Println("查询错误:", err)
		return false
	}
	return code == imgCode
}

// 存储短信验证码到redis
func SaveSmsCode(phone, code string) error {
	conn := RedisPool.Get()
	defer conn.Close()
	_, err := conn.Do("setax", phone+"code", 60*3, code)
	return err
}

// 处理登录业务， 根据手机号/密码 获取用户名
func Login(mobile, pwd string) (string, error) {
	var user User

	//对参数pwd 做md5 hash 加密
	m5 := md5.New()
	m5.Write([]byte(pwd))
	pwd_hash := hex.EncodeToString(m5.Sum(nil))

	// 在user 表中查找
	err := GlobalConn.Where("name = ?", mobile).Where("password_hash = ?", pwd_hash).Select("name").Find(&user).Error
	if err != nil {
		fmt.Println("查找 失败：", err)
	}
	return user.Name, err
}

// 根据用户名在Mysql中查找用户信息
func GetInfo(userName string) (User, error) {
	var user User
	// mysql 中查找
	err := GlobalConn.Where("name = ?", userName).Find(&user).Error
	return user, err
}

//在Mysql中更新新用户名
func UpdateUserName(newname, oldname string) error {
	return GlobalConn.Model(new(User)).Where("name = ?", oldname).Update("name", newname).Error
}
