package auth

import (
	"context"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/golang-jwt/jwt"
)

const HeaderKey = "x-md-global-uid"
const tokenType = "Bearer"

var jwtKey = []byte("greeter")
var openMethod = make(map[string]bool)

// 新增需要开放的方法
// 注意 key 值
func init() {
	openMethod["/helloworld.v1.Greeter/SayHello"] = true
}

// CreateJWT 创建JWT
// 传入的参数根据项目进行定制
// JWT 中的信息建议是非敏感信息
func CreateJWT(userId int64) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userId,
	})
	return claims.SignedString(jwtKey)
}

// CheckJWT 检查JWT是否合法
func CheckJWT(jwtToken string) (map[string]interface{}, error) {
	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		result := make(map[string]interface{}, 2)
		result["user_id"] = claims["user_id"]
		return result, nil
	} else {
		return nil, errAuthTypeError
	}
}

// NewAuthMiddleware 权限拦截校验
func NewAuthMiddleware() func(handler middleware.Handler) middleware.Handler {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			// 判断方法是否 Open
			tr, ok := transport.FromServerContext(ctx)
			if ok && openMethod[tr.Operation()] {
				return handler(ctx, req)
			}

			var jwtToken string
			jwtToken = tr.RequestHeader().Get("Authorization")
			if jwtToken == "" {
				// 缺少可认证的token，返回错误
				return nil, errAuthFail
			}
			token, err := CheckJWT(jwtToken[len(tokenType)+1:])
			if err != nil {
				// 缺少合法的token，返回错误
				return nil, errAuthFail
			}
			ctx = context.WithValue(ctx, HeaderKey, token["user_id"])
			return handler(ctx, req)
		}
	}
}
