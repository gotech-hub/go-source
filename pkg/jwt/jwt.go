package jwt

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"go-source/pkg/utils"
	"strings"
	"time"
)

const (
	RS256 = "RS256"
	RS384 = "RS384"
	RS512 = "RS512"
	ES256 = "ES256"
	ES384 = "ES384"
	ES512 = "ES512"
	PS256 = "PS256"
	PS384 = "PS384"
	PS512 = "PS512"
	EdDSA = "EdDSA"
)

type JWT struct {
	alg        string
	privateKey interface{}
	publicKey  interface{}
	expire     time.Duration
	refExpire  time.Duration
}

func NewJWT(cfg JwtConfig) (*JWT, error) {
	_jwt := &JWT{
		alg:       cfg.Alg,
		expire:    cfg.Expire,
		refExpire: cfg.RefExpire,
	}

	switch cfg.Alg {
	case RS256, RS384, RS512, PS256, PS384, PS512:
		// parse private key
		pk, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(cfg.PrivateKey))
		if err != nil {
			return nil, err
		}

		_jwt.privateKey = pk

		// build public key from private key
		pubKey, err := utils.GetRSAPublicKeyFromPrivateKey(pk)
		if err != nil {
			return nil, err
		}

		_jwt.publicKey = pubKey
		break

	case ES256, ES384, ES512:
		// parse private key
		pk, err := jwt.ParseECPrivateKeyFromPEM([]byte(cfg.PrivateKey))
		if err != nil {
			return nil, err
		}

		_jwt.privateKey = pk
		_jwt.publicKey = &pk.PublicKey
		break

	case EdDSA:
		// parse private key
		pk, err := jwt.ParseEdPrivateKeyFromPEM([]byte(cfg.PrivateKey))
		if err != nil {
			return nil, err
		}

		// get public key from private key
		pkEd25519, ok := pk.(ed25519.PrivateKey)
		if !ok {
			return nil, fmt.Errorf("key is not a valid Ed25519 private key")
		}

		pubEd25519, ok := pkEd25519.Public().(ed25519.PublicKey)
		if !ok {
			return nil, fmt.Errorf("key is not a valid Ed25519 public key")
		}

		_jwt.privateKey = pk
		_jwt.publicKey = pubEd25519
		break

	default:
		return nil, fmt.Errorf("alg invalid: %s", cfg.Alg)
	}

	return _jwt, nil
}

func (j *JWT) SignToken(claims interface{}) (string, error) {
	byteClaim, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}

	var mapClaims jwt.MapClaims
	err = json.Unmarshal(byteClaim, &mapClaims)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.GetSigningMethod(j.alg), mapClaims)

	return token.SignedString(j.privateKey)
}

func (j *JWT) VerifyAndParseToken(tokenStr string, claims interface{}) error {
	// Parse jwt and check sign
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// check signing method
		alg := token.Method.Alg()
		if alg != j.alg {
			return nil, fmt.Errorf("signing method invalid: %v", alg)
		}

		return j.publicKey, nil
	})

	if err != nil {
		return err
	}

	// check token valid
	if !token.Valid {
		return fmt.Errorf("token invalid")
	}

	// parse claims
	byteClaim, err := json.Marshal(token.Claims)
	if err != nil {
		return err
	}

	err = json.Unmarshal(byteClaim, claims)
	if err != nil {
		return err
	}

	return nil
}

func (j *JWT) GetExpire() time.Duration {
	return j.expire
}

func (j *JWT) GetRefExpire() time.Duration {
	return j.refExpire
}

func ParseClaims(tokenStr string, claims interface{}) error {
	// split token
	splitToken := strings.Split(tokenStr, ".")

	if len(splitToken) != 3 {
		return fmt.Errorf("token invalid")
	}

	// decode claims
	byteClaim, err := base64.RawURLEncoding.DecodeString(splitToken[1])
	if err != nil {
		return err
	}

	err = json.Unmarshal(byteClaim, claims)
	if err != nil {
		return err
	}

	return nil
}

func ParseHeader(tokenStr string, header interface{}) error {
	// split token
	splitToken := strings.Split(tokenStr, ".")

	if len(splitToken) != 3 {
		return fmt.Errorf("token invalid")
	}

	// decode header
	byteHeader, err := base64.RawURLEncoding.DecodeString(splitToken[0])
	if err != nil {
		return err
	}

	err = json.Unmarshal(byteHeader, header)
	if err != nil {
		return err
	}

	return nil
}

func Verify(tokenStr, alg string, publicKey interface{}) error {
	// Parse jwt and check sign
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// check signing method
		if token.Method.Alg() != alg {
			return nil, fmt.Errorf("signing method invalid: %v", token.Method.Alg())
		}

		// parse public key
		switch alg {
		case RS256, RS384, RS512, PS256, PS384, PS512:
			pubKey, ok := publicKey.(*rsa.PublicKey)
			if !ok {
				return nil, fmt.Errorf("public key invalid")
			}

			return pubKey, nil

		case ES256, ES384, ES512:
			pubKey, ok := publicKey.(*ecdsa.PublicKey)
			if !ok {
				return nil, fmt.Errorf("public key invalid")
			}

			return pubKey, nil

		default:
			return nil, fmt.Errorf("alg invalid: %s", alg)
		}
	})

	if err != nil {
		return err
	}

	// check token valid
	if !token.Valid {
		return fmt.Errorf("token invalid")
	}

	return nil
}

func VerifyAndParseClaims(tokenStr, alg string, publicKey, claims interface{}) error {
	// Parse jwt and check sign
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// check signing method
		if token.Method.Alg() != alg {
			return nil, fmt.Errorf("signing method invalid: %v", token.Method.Alg())
		}

		// parse public key
		switch alg {
		case RS256, RS384, RS512, PS256, PS384, PS512:
			pubKey, ok := publicKey.(*rsa.PublicKey)
			if !ok {
				return nil, fmt.Errorf("public key invalid")
			}

			return pubKey, nil

		case ES256, ES384, ES512:
			pubKey, ok := publicKey.(*ecdsa.PublicKey)
			if !ok {
				return nil, fmt.Errorf("public key invalid")
			}

			return pubKey, nil

		default:
			return nil, fmt.Errorf("alg invalid: %s", alg)
		}
	})

	if err != nil {
		return err
	}

	// check token valid
	if !token.Valid {
		return fmt.Errorf("token invalid")
	}

	// parse claims
	byteClaim, err := json.Marshal(token.Claims)
	if err != nil {
		return err
	}

	err = json.Unmarshal(byteClaim, claims)
	if err != nil {
		return err
	}

	return nil
}
