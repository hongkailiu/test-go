package http

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	log "github.com/sirupsen/logrus"
)

func getToken(localID string, key interface{}) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"localID": localID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	return token.SignedString(key)
}

func getLocalIDFromToken(tokenString string, key interface{}) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return key, nil
	})
	if err != nil {
		log.Warnf("found error when jwt.Parse(), %s", err.Error())
		return "", err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		localID := claims["localID"].(string)
		//https://github.com/dgrijalva/jwt-go/issues/224
		exp := claims["exp"].(float64)
		log.WithFields(log.Fields{"localID": localID, "exp": exp}).Debug("found in token")
		//expInt64, err := strconv.ParseInt(exp, 10, 64)
		if err != nil {
			log.Warnf("found error when getLocalID(s), %s", err.Error())
			return "", err
		}
		if time.Now().Unix() > int64(exp) {
			log.Warnf("token is expired, %s", exp)
			return "", err
		}
		return localID, nil
	}
	return "", fmt.Errorf("unexpected error when getLocalIDFromToken: %s", tokenString)
}
