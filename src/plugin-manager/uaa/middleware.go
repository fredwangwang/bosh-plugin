package uaa

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"strings"
	"time"
)

type tokenKey struct {
	KeyId string `json:"kid"`
	Value string `json:"value"`
}

type tokenKeys struct {
	Keys []tokenKey `json:"keys"`
}

type GinUAAMiddleware struct {
	uaaUrl     string
	scopes     []string
	publicKeys map[string]*rsa.PublicKey
}

func (u *GinUAAMiddleware) FetchTokenKeys() error {
	if u.publicKeys == nil {
		u.publicKeys = map[string]*rsa.PublicKey{}
	}

	resp, err := http.Get(u.uaaUrl + "/token_keys")
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("response from tokey key is not valid")
	}
	defer resp.Body.Close()

	keys := tokenKeys{}
	if err := json.NewDecoder(resp.Body).Decode(&keys); err != nil {
		return fmt.Errorf("failed to unmarshal token keys")
	}

	log.Printf("%#v\n", keys)

	for _, key := range keys.Keys {
		block, _ := pem.Decode([]byte(key.Value))
		if block == nil {
			return fmt.Errorf("failed to parse PEM block")
		}

		publicKeyInter, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return errors.Wrap(err, "failed to parse public key")
		}

		publicKey, ok := publicKeyInter.(*rsa.PublicKey)
		if !ok {
			return fmt.Errorf("failed to get valid RSA public key")
		}

		u.publicKeys[key.KeyId] = publicKey
	}

	log.Printf("%#v\n", u.publicKeys)
	return nil
}

func (u *GinUAAMiddleware) LoadOrFetchPublicKey(keyId string) (*rsa.PublicKey, error) {
	publicKey, ok := u.publicKeys[keyId]
	if ok {
		return publicKey, nil
	}

	if err := u.FetchTokenKeys(); err != nil {
		log.Println(err)
	}

	log.Println(u.publicKeys)

	publicKey, ok = u.publicKeys[keyId]
	if ok {
		return publicKey, nil
	}

	return nil, fmt.Errorf("public key %s not found", keyId)
}

func (u *GinUAAMiddleware) JWtFromHeader(c *gin.Context) (string, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("auth header is empty")
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", fmt.Errorf("auth header is invalid")
	}
	return parts[1], nil
}

func (u *GinUAAMiddleware) ParseToken(c *gin.Context) (*jwt.Token, error) {
	token, err := u.JWtFromHeader(c)
	if err != nil {
		return nil, err
	}

	return jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod("RS256") != t.Method {
			return nil, fmt.Errorf("invalid signing algorithm")
		}
		keyId := t.Header["kid"].(string)
		return u.LoadOrFetchPublicKey(keyId)
	})
}

func (u *GinUAAMiddleware) GetClaimsFromJWT(c *gin.Context) (jwt.MapClaims, error) {
	token, err := u.ParseToken(c)
	if err != nil {
		return nil, err
	}

	claims := token.Claims.(jwt.MapClaims)
	return claims, nil
}

func (u *GinUAAMiddleware) UAAJWTAuthentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, err := u.GetClaimsFromJWT(c)
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		if claims["exp"] == nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, "the token is malformed")
			return
		}

		if exp, ok := claims["exp"].(float64); !ok {
			c.AbortWithStatusJSON(http.StatusBadRequest, "the token is malformed")
			return
		} else {
			if int64(exp) < time.Now().Unix() {
				log.Println("token expired")
				c.AbortWithStatus(http.StatusNotFound)
				return
			}
		}

		scopes := claims["scope"].([]interface{})
		for _, scope := range scopes {
			for _, desiredScope := range u.scopes {
				if scope.(string) == desiredScope {
					c.Next()
					return
				}
			}
		}

		log.Printf("scopes %v does not match desired scopes %v\n", scopes, u.scopes)
		c.AbortWithStatus(404)
	}
}

func UAAJWTMiddleware(uaaUrl string, allowedScopes []string) gin.HandlerFunc {
	middleware := &GinUAAMiddleware{
		uaaUrl: uaaUrl,
		scopes: allowedScopes,
	}

	return middleware.UAAJWTAuthentication()
}
