package aliyun

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth/credentials/provider"
)

type Provider struct {
	isEnvProvider     bool
	isAKSKProvider    bool
	isProfileProvider bool
	UsedRegion        string
	Pro               provider.Provider
}

func NewProvider(region string, opts ...ProviderOption) *Provider {
	const (
		EnvProvider     = true
		AKSKProvider    = false
		ProfileProvider = false
	)
	// default provider is EnvProvider
	p := &Provider{
		isEnvProvider:     EnvProvider,
		isAKSKProvider:    AKSKProvider,
		isProfileProvider: ProfileProvider,
		UsedRegion:        region,
	}

	if len(opts) == 0 {
		// use env provider as default provider
		InitProviderWithEnv()(p)
	} else {
		// Loop through each option
		for _, opt := range opts {
			// Call the option giving the instantiated
			// *Provider as the argument
			opt(p)
		}
	}

	// return the modified house instance
	return p
}

type ProviderOption func(*Provider)

func InitProviderWithEnv() ProviderOption {
	return func(p *Provider) {
		p.Pro = provider.NewEnvProvider()
	}
}
