package config

import (
	"context"
	"github.com/jfrog/jfrog-client-go/http/httpclient"
	"net/http"
	"time"

	"github.com/jfrog/jfrog-client-go/auth"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

type Config interface {
	GetCertificatesPath() string
	GetThreads() int
	IsDryRun() bool
	GetServiceDetails() auth.ServiceDetails
	GetLogger() log.Log
	IsInsecureTls() bool
	GetContext() context.Context
	GetDialTimeout() time.Duration
	GetOverallRequestTimeout() time.Duration
	GetHttpRetries() int
	GetHttpRetryWaitMilliSecs() int
	GetHttpClient() *http.Client
	GetKerberosDetails() httpclient.KerberosDetails
}

type servicesConfig struct {
	auth.ServiceDetails
	certificatesPath       string
	dryRun                 bool
	threads                int
	logger                 log.Log
	insecureTls            bool
	ctx                    context.Context
	dialTimeout            time.Duration
	overallRequestTimeout  time.Duration
	httpRetries            int
	httpRetryWaitMilliSecs int
	httpClient             *http.Client
	kerberosDetails        httpclient.KerberosDetails
}

func (config *servicesConfig) IsDryRun() bool {
	return config.dryRun
}

func (config *servicesConfig) GetCertificatesPath() string {
	return config.certificatesPath
}

func (config *servicesConfig) GetThreads() int {
	return config.threads
}

func (config *servicesConfig) GetServiceDetails() auth.ServiceDetails {
	return config.ServiceDetails
}

func (config *servicesConfig) GetLogger() log.Log {
	return config.logger
}

func (config *servicesConfig) IsInsecureTls() bool {
	return config.insecureTls
}

func (config *servicesConfig) GetContext() context.Context {
	return config.ctx
}

func (config *servicesConfig) GetDialTimeout() time.Duration {
	return config.dialTimeout
}

func (config *servicesConfig) GetOverallRequestTimeout() time.Duration {
	return config.overallRequestTimeout
}

func (config *servicesConfig) GetHttpRetries() int {
	return config.httpRetries
}

func (config *servicesConfig) GetHttpRetryWaitMilliSecs() int {
	return config.httpRetryWaitMilliSecs
}

func (config *servicesConfig) GetHttpClient() *http.Client {
	return config.httpClient
}

func (config *servicesConfig) GetKerberosDetails() httpclient.KerberosDetails {
	return config.kerberosDetails
}
