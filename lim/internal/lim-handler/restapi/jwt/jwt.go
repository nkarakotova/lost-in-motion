package jwt

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"lim/cmd/registry"
	"lim/internal/lim-core/errors/servicesErrors"
	"lim/internal/lim-core/models"
	"lim/internal/lim-handler/restapi/operations/auth"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/golang-jwt/jwt"
)

var secretKey = []byte("secret-key")

func SignupHandlerFunc(params auth.SignupParams, app registry.App) middleware.Responder {
	telephone := params.ClientInfo.Telephone
	password := params.ClientInfo.Password
	name := params.ClientInfo.Name
	email := params.ClientInfo.Email

	client := models.Client{
		Telephone: *telephone,
		Password: *password,
		Name: *name,
		Mail: email.String(),
	}

	ctx := context.Background()
	err := app.Services.ClientService.Create(ctx, &client)
	if err != nil && err == servicesErrors.ClientAlreadyExists {
		return middleware.Error(http.StatusConflict, "Client already exists")
	} else if err != nil && err == servicesErrors.ClientMailIncorrect {
		return middleware.Error(http.StatusUnprocessableEntity, "Client mail incorrect")
	} else if err != nil && err == servicesErrors.ClientTelephoneIncorrect {
		return middleware.Error(http.StatusUnprocessableEntity, "Client telephone incorrect")
	} else if err != nil {
		return middleware.Error(http.StatusInternalServerError, "Could not create client")
	}

	tokenString, err := createToken(*telephone, "client")
	if err != nil {
		return middleware.Error(http.StatusInternalServerError, "Could not create token")
	}

	return middleware.ResponderFunc(func(rw http.ResponseWriter, p runtime.Producer) {
		rw.WriteHeader(http.StatusOK)
		fmt.Fprint(rw, tokenString)
	})
}

func LoginHandlerFunc(params auth.LoginParams, app registry.App) middleware.Responder {
	telephone := params.ClientLoginInfo.Telephone
	password := params.ClientLoginInfo.Password

	if *telephone == "0000000000" && *password == "12345" {
		tokenString, err := createToken(*telephone, "admin")
		if err != nil {
			return middleware.Error(http.StatusInternalServerError, "Could not create token")
		}
		return middleware.ResponderFunc(func(rw http.ResponseWriter, p runtime.Producer) {
			rw.WriteHeader(http.StatusOK)
			fmt.Fprint(rw, tokenString)
		})
	} else if *telephone == "0000000000" {
		return middleware.Error(http.StatusUnauthorized, "Authorization error")
	}

	ctx := context.Background()
	_, err := app.Services.ClientService.Login(ctx, *telephone, *password)
	if err != nil && err == servicesErrors.ClientDoesNotExists {
		return middleware.Error(http.StatusNotFound, "Client does not exists")
	} else if err != nil && err == servicesErrors.InvalidPassword {
		return middleware.Error(http.StatusUnauthorized, "Authorization error")
	} else if err != nil {
		return middleware.Error(http.StatusInternalServerError, "Could not login")
	}

	tokenString, err := createToken(*telephone, "client")
	if err != nil {
		return middleware.Error(http.StatusInternalServerError, "Could not create token")
	}
	return middleware.ResponderFunc(func(rw http.ResponseWriter, p runtime.Producer) {
		rw.WriteHeader(http.StatusOK)
		fmt.Fprint(rw, tokenString)
	})
}

func GetAuthenticatedUser(r *http.Request) (string, string, error) {
	authHeader := r.Header.Get("token")
	if authHeader == "" {
		return "", "", fmt.Errorf("missing Authorization header")
	}

	token, err := jwt.Parse(authHeader, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return "", fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})
	if err != nil || !token.Valid {
		return "", "", fmt.Errorf("invalid token")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		telephone := claims["telephone"].(string)
		role := claims["role"].(string)
		return telephone, role, nil
	} else {
		return "", "", fmt.Errorf("invalid token claims")
	}
}

func createToken(telephone, role string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"telephone": telephone,
		"role":     role,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyToken(tokenString string) (string, string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return "", "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		telephone := claims["telephone"].(string)
		role := claims["role"].(string)
		return telephone, role, nil
	}

	return "", "", fmt.Errorf("invalid token")
}
