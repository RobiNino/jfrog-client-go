package httpclient

import (
	"context"
	"crypto/tls"
	krb5Client "github.com/jcmturner/gokrb5/v8/client"
	krb5Cofig "github.com/jcmturner/gokrb5/v8/config"
	"github.com/jcmturner/gokrb5/v8/keytab"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"net"
	"net/http"
	"time"

	"github.com/jfrog/jfrog-client-go/auth/cert"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
)

var DefaultDialTimeout = 30 * time.Second

func ClientBuilder() *httpClientBuilder {
	builder := &httpClientBuilder{}
	builder.SetDialTimeout(DefaultDialTimeout)
	return builder
}

type httpClientBuilder struct {
	certificatesDirPath   string
	clientCertPath        string
	clientCertKeyPath     string
	insecureTls           bool
	ctx                   context.Context
	dialTimeout           time.Duration
	overallRequestTimeout time.Duration
	retries               int
	retryWaitMilliSecs    int
	httpClient            *http.Client
	kerberosDetails       KerberosDetails
}

func (builder *httpClientBuilder) SetCertificatesPath(certificatesPath string) *httpClientBuilder {
	builder.certificatesDirPath = certificatesPath
	return builder
}

func (builder *httpClientBuilder) SetClientCertPath(certificatePath string) *httpClientBuilder {
	builder.clientCertPath = certificatePath
	return builder
}

func (builder *httpClientBuilder) SetClientCertKeyPath(certificatePath string) *httpClientBuilder {
	builder.clientCertKeyPath = certificatePath
	return builder
}

func (builder *httpClientBuilder) SetInsecureTls(insecureTls bool) *httpClientBuilder {
	builder.insecureTls = insecureTls
	return builder
}

func (builder *httpClientBuilder) SetHttpClient(httpClient *http.Client) *httpClientBuilder {
	builder.httpClient = httpClient
	return builder
}

func (builder *httpClientBuilder) SetKerberosDetails(kerberosDetails KerberosDetails) *httpClientBuilder {
	builder.kerberosDetails = kerberosDetails
	return builder
}

func (builder *httpClientBuilder) SetContext(ctx context.Context) *httpClientBuilder {
	builder.ctx = ctx
	return builder
}

func (builder *httpClientBuilder) SetDialTimeout(dialTimeout time.Duration) *httpClientBuilder {
	builder.dialTimeout = dialTimeout
	return builder
}

func (builder *httpClientBuilder) SetOverallRequestTimeout(overallRequestTimeout time.Duration) *httpClientBuilder {
	builder.overallRequestTimeout = overallRequestTimeout
	return builder
}

func (builder *httpClientBuilder) SetRetries(retries int) *httpClientBuilder {
	builder.retries = retries
	return builder
}

func (builder *httpClientBuilder) SetRetryWaitMilliSecs(retryWaitMilliSecs int) *httpClientBuilder {
	builder.retryWaitMilliSecs = retryWaitMilliSecs
	return builder
}

func (builder *httpClientBuilder) AddClientCertToTransport(transport *http.Transport) error {
	if builder.clientCertPath != "" {
		certificate, err := cert.LoadCertificate(builder.clientCertPath, builder.clientCertKeyPath)
		if err != nil {
			return err
		}
		transport.TLSClientConfig.Certificates = []tls.Certificate{certificate}
	}
	return nil
}

func (builder *httpClientBuilder) Build() (*HttpClient, error) {
	kerberosClient, err := builder.createKerberosClientIfNeeded()
	if err != nil {
		return nil, err
	}

	if builder.httpClient != nil {
		// Using a custom http.Client, pass-though.
		return &HttpClient{client: builder.httpClient, ctx: builder.ctx, retries: builder.retries, retryWaitMilliSecs: builder.retryWaitMilliSecs, kerberosClient: kerberosClient}, nil
	}

	var transport *http.Transport

	if builder.certificatesDirPath == "" {
		transport = builder.createDefaultHttpTransport()
		//#nosec G402 -- Insecure TLS allowed here.
		transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: builder.insecureTls}
	} else {
		transport, err = cert.GetTransportWithLoadedCert(builder.certificatesDirPath, builder.insecureTls, builder.createDefaultHttpTransport())
		if err != nil {
			return nil, errorutils.CheckErrorf("failed creating HttpClient: " + err.Error())
		}
	}
	err = builder.AddClientCertToTransport(transport)
	return &HttpClient{client: &http.Client{Transport: transport, Timeout: builder.overallRequestTimeout}, ctx: builder.ctx, retries: builder.retries, retryWaitMilliSecs: builder.retryWaitMilliSecs, kerberosClient: kerberosClient}, err
}

func (builder *httpClientBuilder) createDefaultHttpTransport() *http.Transport {
	return &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   builder.dialTimeout,
			KeepAlive: 20 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
}

func (builder *httpClientBuilder) createKerberosClientIfNeeded() (*krb5Client.Client, error) {
	log.Debug(">>KERBEROS>> Initializing Kerberos client...")
	krb5Details := builder.kerberosDetails
	if krb5Details.Krb5ConfigPath == "" || (krb5Details.Password == "" && krb5Details.KeytabPath == "") {
		log.Debug(">>KERBEROS>> kerberos details missing, skipping Kerberos client initialization...")
		return nil, nil
	}

	krbConf, err := krb5Cofig.Load(krb5Details.Krb5ConfigPath)
	if err != nil {
		log.Debug(">>KERBEROS>> Error encountered when loading Krb5 config from path: ", krb5Details.Krb5ConfigPath, "error: ", err)
		return nil, err
	}

	var cl *krb5Client.Client
	if krb5Details.Password != "" {
		log.Debug(">>KERBEROS>> Initializing Kerberos client with password...")
		cl = krb5Client.NewWithPassword(krb5Details.Username, krb5Details.Realm, krb5Details.Password, krbConf, krb5Client.DisablePAFXFAST(true))
	} else if krb5Details.KeytabPath != "" {
		log.Debug(">>KERBEROS>> Initializing Kerberos client with keytab...")
		ktFromFile, err := keytab.Load(krb5Details.KeytabPath)
		if err != nil {
			log.Debug(">>KERBEROS>> Error encountered when loading keytab from path: ", krb5Details.KeytabPath, "error: ", err)
			return nil, err
		}
		cl = krb5Client.NewWithKeytab(krb5Details.Username, krb5Details.Realm, ktFromFile, krbConf, krb5Client.DisablePAFXFAST(true))
	}

	err = cl.Login()
	if err != nil {
		log.Debug(">>KERBEROS>> Error encountered when trying to log in the client. error: ", err)
		return nil, err
	}
	log.Debug(">>KERBEROS>> Done initializing Kerberos client...")
	return cl, nil
}

type KerberosDetails struct {
	Krb5ConfigPath string `json:"krb5ConfigPath,omitempty"`
	Username       string `json:"username,omitempty"`
	Realm          string `json:"realm,omitempty"`
	Password       string `json:"password,omitempty"`
	KeytabPath     string `json:"keytabPath,omitempty"`
}
