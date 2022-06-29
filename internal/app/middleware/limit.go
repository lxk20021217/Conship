//Author: lxk20021217
//Date: 2022-06-17 23:49:22
//LastEditTime: 2022-06-24 15:16:37
//LastEditors: lxk20021217
//Description:
//FilePath: \Conship\internal\app\middleware\limit.go
//是谁总是天亮了才睡

package middleware

import (

	"time"

	"github.com/QinLiStudio/Conship/internal/app/configs"
	"github.com/QinLiStudio/Conship/pkg/error"
	"github.com/gin-gonic/gin"
)

func LimitRoute(limit int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.ClientIP() + c.FullPath()
		// 读取请求次数
		value, err := configs.REDISDB.Get(key).Int64()
		if err != nil {
			// 初始化请求次数
			configs.REDISDB.Set(key, 1, time.Hour)
			c.Next()
		} else if value < limit {
			// 请求次数 +1
			configs.REDISDB.Incr(key)
			c.Next()
		} else {
			// 请求次数超过限制
			error.ErrResponse(c, error.ErrTooManyRequests, error.TooManyRequests, "请求过于频繁，请一小时后再重试。", err)
			c.Abort()
		}
	}
}

/**
 * @description: 按时间单路由限流
 * @param {int} limit
 */
func LimitAverageRoute(limit int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.ClientIP() + c.FullPath()
		// 读取值并验证
		if err := configs.REDISDB.Get(key).Err(); err != nil {
			// 初始化请求次数
			configs.REDISDB.Set(key, 1, time.Hour/time.Duration(limit))
			c.Next()
		} else {
			error.ErrResponse(c, error.ErrTooManyRequests, error.TooManyRequests, "请求过于频繁，请一小时后再重试。", err)
			c.Abort()
		}
	}
}
