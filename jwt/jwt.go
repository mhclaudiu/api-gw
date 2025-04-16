package jwt

import (
	"api-gw/functions"
	"api-gw/json"
	"api-gw/logging"
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	_jwt "github.com/golang-jwt/jwt/v5"
)

func Auth(next http.HandlerFunc, auth bool) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		if !auth {

			ctx := context.WithValue(r.Context(), "user", "N/A")
			next(w, r.WithContext(ctx))

			return
		}

		rcvAuth := r.Header.Get("Authorization")
		clientAddr := functions.RemoteHost(r)

		if !strings.HasPrefix(rcvAuth, "Bearer ") {

			json.Write(w, APIxOBJ_Json_Status_Rsp{
				Status: APIxOBJ_Json_Status{
					Err: true,
					Msg: "Sorry .. Invalid Authorization ..",
				},
			}, http.StatusUnauthorized)

			appLog.Add(logging.Entry{
				Event: fmt.Sprintf("Client: %s | Invalid Authorization", functions.RemoteHost(r)),
				Code:  logging.CONST_CODE_ERROR,
			})

			return
		}

		rcvToken := strings.TrimPrefix(rcvAuth, "Bearer ")
		token, err := _jwt.Parse(rcvToken, func(t *_jwt.Token) (interface{}, error) {

			return secret, nil
		})

		if err != nil || !token.Valid {

			json.Write(w, APIxOBJ_Json_Status_Rsp{
				Status: APIxOBJ_Json_Status{
					Err: true,
					Msg: "Sorry .. Invalid Token ..",
				},
			}, http.StatusForbidden)

			appLog.Add(logging.Entry{
				Event: fmt.Sprintf("Client: %s | Invalid Token: %s", clientAddr, rcvToken),
				Code:  logging.CONST_CODE_WARNING,
			})

			return
		}

		claims := token.Claims.(_jwt.MapClaims)

		appLog.Add(logging.Entry{
			Event: fmt.Sprintf("Client: %s - User: %s Successfully Authenticated", clientAddr, claims["user"]),
			Code:  logging.CONST_CODE_INFO,
		})

		ctx := context.WithValue(r.Context(), "user", claims["user"])
		next(w, r.WithContext(ctx))
	}
}

func CreateToken() (string, error) {

	claims := _jwt.MapClaims{
		"exp":  time.Now().Add(time.Hour * 1).Unix(),
		"user": "test",
	}

	token := _jwt.NewWithClaims(_jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}
