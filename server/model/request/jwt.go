package request

import (
	"server/model/appTypes"

	"github.com/gofrs/uuid"
	jwt "github.com/golang-jwt/jwt/v4"
)

// 都是自定义的claims
// JwtCustomClaims 结构体用于存储JWT的自定义Claims，继承自BaseClaims，并包含标准的JWT注册信息
type JwtCustomClaims struct {
	BaseClaims           // 基础Claims，包含用户ID、UUID和角色ID
	jwt.RegisteredClaims // 标准JWT声明，例如过期时间、发行者等
}

// JwtCustomRefreshClaims 结构体用于存储刷新Token的自定义Claims，包含用户ID和标准的JWT注册信息
type JwtCustomRefreshClaims struct {
	UserID               uint // 用户ID，用于与刷新Token相关的身份验证
	jwt.RegisteredClaims      // 标准JWT声明
}

// BaseClaims 结构体用于存储基本的用户信息，作为JWT的Claim部分
type BaseClaims struct {
	UserID uint            // 用户ID，标识用户唯一性
	UUID   uuid.UUID       // 用户的UUID，唯一标识用户
	RoleID appTypes.RoleID // 用户角色ID，表示用户的权限级别
}

// UUID（Universally Unique Identifier）是一种用于标识信息的唯一标识符。它的特点如下：
//
// 1. **唯一性**：UUID 通过算法生成，确保在全球范围内几乎不会重复。
// 2. **格式**：通常是一个 128 位的数字，表示为 32 个十六进制字符，分为 5 组，格式为 `8-4-4-4-12`，例如：`550e8400-e29b-41d4-a716-446655440000`。
// 3. **用途**：常用于 分布式系统 中唯一标识实体（如用户、资源、会话等），避免 ID 冲突。
//
// 在代码中，`UUID` 用于唯一标识用户，确保用户身份的唯一性和安全性。
