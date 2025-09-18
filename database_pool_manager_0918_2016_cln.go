// 代码生成时间: 2025-09-18 20:16:34
package main

import (
    "database/sql"
    _ "github.com/go-sql-driver/mysql" 
    "fmt"
    "log"
)

// DatabaseConfig 定义数据库连接配置
type DatabaseConfig struct {
    User     string
    Pass     string
    Host     string
    Port     int
    DBName   string
}

// NewDatabasePool 创建并返回一个新的数据库连接池
func NewDatabasePool(config DatabaseConfig) (*sql.DB, error) {
    // 构造DSN（数据源名称）
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", config.User, config.Pass, config.Host, config.Port, config.DBName)

    // 打开数据库连接
    db, err := sql.Open("mysql", dsn)
    if err != nil {
        return nil, err
    }

    // 设置数据库连接池的参数
    db.SetMaxOpenConns(100)    // 设置最大打开的连接数
    db.SetMaxIdleConns(10)     // 设置连接池中的最大空闲连接数
    db.SetConnMaxLifetime(3600) // 设置连接的最大存活时间，单位秒

    // 测试数据库连接
    err = db.Ping()
    if err != nil {
        return nil, err
    }

    return db, nil
}

func main() {
    // 定义数据库配置
    config := DatabaseConfig{
        User:     "username",
        Pass:     "password",
        Host:     "localhost",
        Port:     3306,
        DBName:   "dbname",
    }

    // 创建数据库连接池
    db, err := NewDatabasePool(config)
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }
    defer db.Close()

    // 使用数据库连接池进行操作，例如查询
    // 这里只是一个示例，实际应用中应替换为具体的业务逻辑
    result, err := db.Query("SELECT * FROM some_table")
    if err != nil {
        log.Printf("Query failed: %v", err)
        return
    }
    defer result.Close()

    // 处理查询结果
    // ...

    // 打印成功消息
    fmt.Println("Database connection pool established successfully.")
}
