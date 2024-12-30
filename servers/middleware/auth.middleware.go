package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ngdangkietswe/swe-gateway-service/domain"
	"github.com/ngdangkietswe/swe-gateway-service/utils"
	"github.com/ngdangkietswe/swe-go-common-shared/cache"
	"github.com/ngdangkietswe/swe-go-common-shared/config"
	"github.com/ngdangkietswe/swe-go-common-shared/constants"
	"github.com/ngdangkietswe/swe-go-common-shared/util"
	"github.com/ngdangkietswe/swe-protobuf-shared/generated/auth"
	"github.com/ngdangkietswe/swe-protobuf-shared/generated/common"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"log"
	"net/http"
	"strings"
	"time"
)

type AuthMiddleware struct {
	authConn *grpc.ClientConn
	cache    *cache.RedisCache
}

// ShouldSkip is a middleware function that checks if the request should skip the auth middleware
func (a AuthMiddleware) ShouldSkip(ctx *gin.Context) bool {
	publicRoutes := []string{
		"/swagger",
		"/api/v1/auth/login",
		"/api/v1/auth/register",
	}

	for _, route := range publicRoutes {
		if strings.HasPrefix(ctx.Request.URL.Path, route) {
			return true
		}
	}

	return false
}

// Handle is a middleware function that checks if the request has a valid token
func (a AuthMiddleware) Handle(ctx *gin.Context) {
	token := utils.GetTokenFromReq(ctx)
	if token == "" {
		ctx.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "Token is required"})
		ctx.Abort()
		return
	}

	claims, err := util.ParseToken(token, config.GetString("JWT_SECRET", ""))
	if err != nil {
		ctx.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		ctx.Abort()
		return
	}

	claimsUser := (*claims)["user"].(map[string]interface{})
	userId := claimsUser["user_id"].(string)

	md := metadata.Pairs(constants.AuthorizationHeader, fmt.Sprintf("%s %s", constants.TokenPrefix, token))
	newCtx := metadata.NewOutgoingContext(context.Background(), md)

	// Get and cache user permission to the context.
	userPermission, err := a.getAndCacheUserPermission(newCtx, userId)
	if err != nil {
		ctx.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "Unknown error"})
		ctx.Abort()
		return
	}

	grpcUserPermission, err := json.Marshal(userPermission)
	if err != nil {
		ctx.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "Unknown error"})
		ctx.Abort()
		return
	}

	ctx.Request.Header.Add(constants.GrpcMetadataUserPermission, string(grpcUserPermission))

	ctx.Next()
	return
}

// getAndCacheUserPermission is a function that gets the user permission from the cache.
// If it doesn't exist, it will get it from the auth service.
func (a AuthMiddleware) getAndCacheUserPermission(ctx context.Context, userId string) (*domain.UserPermission, error) {
	var (
		userPermission       *domain.UserPermission
		permissionOfUserResp *auth.PermissionOfUserResp
		permissions          []*domain.Permission
	)

	cacheKey := fmt.Sprintf("%s_%s", constants.UserPermissionCacheKeyPrefix, userId)
	if err := a.cache.Get(cacheKey, &userPermission); err != nil {
		log.Printf("Error getting user permission from cache: %v", err)
		authClient := auth.NewPermissionInternalServiceClient(a.authConn)

		permissionOfUserResp, err = authClient.PermissionOfUser(ctx, &common.IdReq{
			Id: userId,
		})
		if err != nil {
			log.Printf("Error getting user permission: %v", err)
			return nil, err
		}

		for _, permission := range permissionOfUserResp.GetData().Permissions {
			permissions = append(permissions, &domain.Permission{
				Action:   permission.Action.Name,
				Resource: permission.Resource.Name,
			})
		}

		userPermission = &domain.UserPermission{
			Permissions: permissions,
		}

		if err = a.cache.Set(cacheKey, userPermission, 30*time.Minute); err != nil {
			return nil, err
		}
	}

	return userPermission, nil
}

func (a AuthMiddleware) AsMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if a.ShouldSkip(ctx) {
			ctx.Next()
			return
		}

		a.Handle(ctx)
	}
}

func NewAuthMiddleware(
	authConn *grpc.ClientConn,
	cache *cache.RedisCache) Middleware {
	return &AuthMiddleware{
		authConn: authConn,
		cache:    cache,
	}
}
