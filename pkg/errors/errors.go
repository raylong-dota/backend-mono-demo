package auth

import "github.com/go-kratos/kratos/v2/errors"

// 系统语意错误 (判断标准和具体业务无关！！！)
var ErrAuthFail = errors.New(401, "Authentication failed", "Missing token or token incorrect")
