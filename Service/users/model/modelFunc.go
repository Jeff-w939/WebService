package model

import (
	//"Project3_WebService/Web/model"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis"
	//"go.etcd.io/bbolt"
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
	_, err := conn.Do("setex", phone+"_code", 60*10, code)
	return err
}

//校验短信验证码
func CheckSmsCode(phone, code string) error {
	// 链接redis
	conn := RedisPool.Get()

	//从redis中， 根据key 获取value --短信验证码
	smsCode, err := redis.String(conn.Do("get", phone+"_code"))
	if err != nil {
		fmt.Println("redis get phone_code err :", err)
		return err
	}
	// 验证码匹配
	if smsCode != code {
		return errors.New("验证码对不上，匹配失败噢！")
	}
	//匹配成功
	return nil
}

// 注册用户信息， 写入Mysql
func RegisterUser(mobile, pwd string) error {
	var user User
	user.Name = mobile   // 默认使用手机号作为用户名
	user.Mobile = mobile // 新加的
	// 使用MD5 对pwd 加密
	m5 := md5.New() // 初始化对象
	m5.Write([]byte(pwd))
	//不使用额外的密匙
	pwd_hash := hex.EncodeToString(m5.Sum(nil))

	user.Password_hash = pwd_hash
	return GlobalConn.Create(&user).Error
}
