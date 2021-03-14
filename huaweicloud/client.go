/**
 * @Author: jie.an
 * @Description:
 * @File:  client.go
 * @Version: 1.0.0
 * @Date: 2021/3/14 18:30
 */

package huaweicloud

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/basic"
	"os"
)

type Auth struct {
	UsedRegion string
	Cred       basic.Credentials
}

func NewAuth(region string, opts ...AuthOption) *Auth {
	a := &Auth{
		UsedRegion: region,
	}
	if len(opts) == 0 {
		// use env Auth as default Auth
		WithEnv()(a)
	} else {
		// Loop through each option
		for _, opt := range opts {
			// Call the option giving the instantiated
			// *client as the argument
			opt(a)
		}
	}

	// return the modified client instance
	return a
}

type AuthOption func(*Auth)

func WithEnv() AuthOption {
	return func(a *Auth) {
		a.Cred = basic.NewCredentialsBuilder().
			WithAk(os.Getenv("HW_ACCESS_KEY")).
			WithSk(os.Getenv("HW_SECRET_KEY")).
			Build()
	}
}
