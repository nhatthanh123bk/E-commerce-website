package repository

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/blogs/app/model"
	"github.com/blogs/db"
	"github.com/blogs/helper"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis/v7"
	"github.com/labstack/echo/v4"
	"github.com/twinj/uuid"
)

// CreateToken create token detail from user id. Afterward the one will
// be used to create metadata token and stored into Redis. At the same time,
// a couple of tokens will be sent to end users.
func CreateToken(id uint) (*model.TokenDetail, error) {
	td := &model.TokenDetail{}
	td.AtExpires = time.Now().Add(time.Hour * 72).Unix()
	td.AccessUuid = uuid.NewV4().String()

	td.RfExpires = time.Now().Add(time.Hour * 168).Unix()
	td.RefreshUuid = td.AccessUuid + strconv.Itoa(int(id))

	var err error
	//create an access token
	atClaims := jwt.MapClaims{
		"authorized":  true,
		"access_uuid": td.AccessUuid,
		"user_id":     id,
		"exp":         td.AtExpires,
	}
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte("ACCESS_SECRET"))
	if err != nil {
		return nil, err
	}

	// create a refresh token
	rtClaims := jwt.MapClaims{
		"authorized":   true,
		"refresh_uuid": td.RefreshUuid,
		"user_id":      id,
		"exp":          td.RfExpires,
	}
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte("REFRESH_SECRET"))
	if err != nil {
		return nil, err
	}

	return td, nil
}

// StoreTokenIntoRedis store token metadata into Redis. These token metadata
// will be deleted automatically when being expired.
func StoreTokenIntoRedis(td *model.TokenDetail, userId uint) error {
	at := time.Unix(td.AtExpires, 0)
	rt := time.Unix(td.RfExpires, 0)
	now := time.Now()

	accessErr := db.RedisClient.Set(td.AccessUuid, strconv.Itoa(int(userId)), at.Sub(now)).Err()
	if accessErr != nil {
		return accessErr
	}

	refreshErr := db.RedisClient.Set(td.RefreshUuid, strconv.Itoa(int(userId)), rt.Sub(now)).Err()
	if refreshErr != nil {
		return refreshErr
	}

	return nil
}

func ExtractToken(req *http.Request) (string, error) {
	bearerToken := req.Header.Get("Authorization")
	if bearerToken == "" {
		helper.Logger.Errorw("required token!")
		return "", fmt.Errorf("need access token to authenticate")
	}
	if strArr := strings.Split(bearerToken, " "); len(strArr) == 2 {
		return strArr[1], nil
	}

	return "", fmt.Errorf("bearer token is invalid")
}

func DecodeToken(req *http.Request) (*jwt.Token, error) {
	strToken, err := ExtractToken(req)
	if strToken == "" {
		return nil, err
	}
	token, err := jwt.Parse(strToken, func(token *jwt.Token) (interface{}, error) {
		return []byte("ACCESS_SECRET"), nil
	})
	if err != nil {
		helper.Logger.Errorw("Can't parse token!", "token:", strToken, "error:", err)
		return nil, err
	}

	helper.Logger.Infow("parsed token", "token:", token)
	return token, nil
}

func ExtractTokenMetadata(req *http.Request) (*model.AccessTokenMetadata, error) {
	token, err := DecodeToken(req)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessUuid, ok := claims["access_uuid"].(string)
		if !ok {
			helper.Logger.Errorw("access_uuid is invalid", "accessUuid:", accessUuid)
			return nil, fmt.Errorf("access_uuid is invalid")
		}
		userId, ok := claims["user_id"].(float64)
		if !ok {
			helper.Logger.Errorw("user_id is invalid", "userId:", userId)
			return nil, fmt.Errorf("user_id is invalid")
		}

		return &model.AccessTokenMetadata{AccessUuid: accessUuid, UserId: int(userId)}, nil
	}

	return nil, fmt.Errorf("token isn't valid")
}

//GetCurrentUserId gets the current userId from access token
func GetCurrentUserId(req *http.Request) (int, error) {
	tokenMetadata, err := ExtractTokenMetadata(req)
	if err != nil {
		return 0, err
	}

	return int(tokenMetadata.UserId), nil
}

// DeleteToken delete token stored in redis
func DeleteToken(req *http.Request) error {
	tokenMetadata, err := ExtractTokenMetadata(req)
	if err != nil {
		return err
	}

	// delete
	_, err = db.RedisClient.Del(tokenMetadata.AccessUuid).Result()
	if err != nil {
		return err
	}

	return nil
}

// when access token is expired, we can automatically utilize refresh
// token to reproduce a new pairs token, instead of having to login again
// over again.
func GenerateTokenFromRefreshToken(c echo.Context) (*model.Token, error) {
	// get refresh token from request
	rfTokenString := c.QueryParam("refresh_token")
	if rfTokenString == "" {
		helper.Logger.Errorw("require refresh token")
		return &model.Token{}, helper.ErrUnAuthorization
	}

	// decode refresh token
	token, err := jwt.Parse(rfTokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("REFRESH_SECRET"), nil
	})
	if err != nil {
		helper.Logger.Errorw("refresh token's invalid", "token:", rfTokenString)
		return &model.Token{}, helper.ErrUnAuthorization
	}
	claim, ok := token.Claims.(jwt.MapClaims)
	helper.Logger.Infow("Metadata", "uuid:", claim["refresh_uuid"], "user_id:", claim["user_id"])
	if ok && token.Valid {
		rfUuid, ok := claim["refresh_uuid"].(string)
		if !ok {
			helper.Logger.Infow("refresh_uuid's invalid", "refresh_uuid:", claim["refresh_uuid"])
			return &model.Token{}, helper.ErrUnAuthorization
		}

		userId, ok := claim["user_id"].(float64)
		if !ok {
			helper.Logger.Infow("user_id's invalid", "user_id:", claim["user_id"])
			return &model.Token{}, helper.ErrUnAuthorization
		}
		tokenMetadata := model.RefreshTokenMetadata{
			RefreshUuid: rfUuid,
			UserId:      int(userId),
		}
		// authenticate refresh token based on rfUuid, userId
		isAuthenticated, err := AuthenticateRefreshToken(tokenMetadata)

		if err == redis.Nil {
			helper.Logger.Errorw("metadata doesn't exist in redis", "refresh uuid:", tokenMetadata.RefreshUuid)
			return &model.Token{}, helper.ErrUnAuthorization
		}
		if err != nil {
			helper.Logger.Errorw("Error occurred when operating redis", "error:", err)
			return &model.Token{}, helper.ErrInternal
		}

		if !isAuthenticated {
			helper.Logger.Errorw("Unauthenticated!", "refresh token:", rfTokenString)
			return &model.Token{}, helper.ErrUnAuthorization
		}
		// here, the credential was authenticated successfully , we're going to delete available
		// refresh token created before in redis.
		_, err = db.RedisClient.Del(tokenMetadata.RefreshUuid).Result()
		if err != nil {
			helper.Logger.Errorw("Can't delete the available refresh", "refresh uuid:", tokenMetadata.RefreshUuid)
			return &model.Token{}, helper.ErrInternal
		}
		// create a new pairs
		td, err := CreateToken(uint(tokenMetadata.UserId))
		if err != nil {
			helper.Logger.Errorw("Can't create detailed token", "userid:", tokenMetadata.UserId)
			return &model.Token{}, helper.ErrInternal
		}
		err = StoreTokenIntoRedis(td, uint(tokenMetadata.UserId))
		if err != nil {
			helper.Logger.Errorw("Can't stored token into redis", "userid:", tokenMetadata.UserId)
			return &model.Token{}, helper.ErrInternal
		}
		tokens := model.Token{
			"access_token":  td.AccessToken,
			"refresh_token": td.RefreshToken,
		}
		return &tokens, nil
	}
	return &model.Token{}, nil
}

func AuthenticateRefreshToken(tokenMetadata model.RefreshTokenMetadata) (bool, error) {
	userId, err := db.RedisClient.Get(tokenMetadata.RefreshUuid).Result()
	if err == redis.Nil {
		return false, err
	}
	if err != nil {
		return false, err
	}
	if strconv.Itoa(tokenMetadata.UserId) != userId {
		return false, nil
	}

	return true, nil
}
