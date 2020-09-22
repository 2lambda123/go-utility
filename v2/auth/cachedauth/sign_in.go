package cachedauth

import (
	"context"
	"sync"

	"github.com/pkg/errors"

	"github.com/SKF/go-utility/v2/auth"
)

var lock sync.Mutex
var tokens auth.Tokens

var config *Config

// Config is the configuration of the package
type Config struct {
	Stage string
}

// Configure will configure the package
func Configure(conf Config) {
	lock.Lock()
	defer lock.Unlock()

	config = &conf

	auth.Configure(auth.Config{Stage: conf.Stage})
}

// GetTokens will return the cached tokens
func GetTokens() auth.Tokens {
	lock.Lock()
	defer lock.Unlock()

	return tokens
}

// SignIn is thread safe and only returns new tokens if the old tokens are about to expire
func SignIn(ctx context.Context, username, password string) (err error) {
	lock.Lock()
	defer lock.Unlock()

	if config == nil {
		return errors.New("cachedauth is not configured")
	}

	const tokenExpireDurationDiff = 5
	if auth.IsTokenValid(tokens.AccessToken, tokenExpireDurationDiff) {
		return nil
	}

	tokens, err = auth.SignIn(ctx, username, password)
	if err != nil {
		tokens = auth.Tokens{}
		return err
	}

	return nil
}
