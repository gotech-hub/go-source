package jwt

import "time"

type JwtConfig struct {
	Alg        string        `env:"ALG,required,notEmpty"` // support: RS256, RS384, RS512, ES256, ES384, ES512, PS256, PS384, PS512
	PrivateKey string        `env:"PRIVATE_KEY,required,notEmpty"`
	Expire     time.Duration `env:"EXPIRE"`
	RefExpire  time.Duration `env:"REF_EXPIRE"`
}
