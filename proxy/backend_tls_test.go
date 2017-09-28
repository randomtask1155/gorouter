package proxy_test

import (
	"crypto/tls"
	"crypto/x509"
	"net/http"

	"code.cloudfoundry.org/gorouter/common/uuid"
	"code.cloudfoundry.org/gorouter/test_util"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Backend TLS", func() {
	var registerConfig test_util.RegisterConfig

	freshProxyCACertPool := func() *x509.CertPool {
		var err error
		caCertPool, err = x509.SystemCertPool()
		Expect(err).NotTo(HaveOccurred())
		return caCertPool
	}
	// createCertAndAddCA creates a signed cert with a root CA and adds the CA
	// to the specified cert pool
	createCertAndAddCA := func(cn test_util.CertNames, cp *x509.CertPool) test_util.CertChain {
		certChain := test_util.CreateSignedCertWithRootCA(cn)
		cp.AddCert(certChain.CACert)
		return certChain
	}

	registerAppAndTest := func() *http.Response {
		ln := test_util.RegisterHandler(r, "test", func(conn *test_util.HttpConn) {
			req, err := http.ReadRequest(conn.Reader)
			if err != nil {
				conn.WriteResponse(test_util.NewResponse(http.StatusInternalServerError))
				return
			}
			err = req.Body.Close()
			Expect(err).ToNot(HaveOccurred())
			conn.WriteResponse(test_util.NewResponse(http.StatusOK))
		}, registerConfig)
		defer ln.Close()

		conn := dialProxy(proxyServer)

		conn.WriteLines([]string{
			"GET / HTTP/1.1",
			"Host: test",
		})

		resp, _ := conn.ReadResponse()
		return resp
	}

	BeforeEach(func() {
		var err error

		privateInstanceId, _ := uuid.GenerateUUID()
		// Clear proxy's CA cert pool
		proxyCertPool := freshProxyCACertPool()

		// Clear backend app's CA cert pool
		backendCACertPool := x509.NewCertPool()

		backendCertChain := createCertAndAddCA(test_util.CertNames{CommonName: privateInstanceId}, proxyCertPool)
		clientCertChain := createCertAndAddCA(test_util.CertNames{CommonName: "gorouter"}, backendCACertPool)

		backendTLSConfig := backendCertChain.AsTLSConfig()
		backendTLSConfig.ClientCAs = backendCACertPool

		conf.Backends.ClientAuthCertificate, err = tls.X509KeyPair(clientCertChain.CertPEM, clientCertChain.PrivKeyPEM)
		Expect(err).NotTo(HaveOccurred())

		registerConfig = test_util.RegisterConfig{
			TLSConfig:  backendTLSConfig,
			InstanceId: privateInstanceId,
			AppId:      "app-1",
		}
	})

	Context("when the backend does not require a client certificate", func() {
		It("makes an mTLS connection with the backend", func() {
			resp := registerAppAndTest()
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
		})
	})
	Context("when the backend requires a client certificate", func() {
		BeforeEach(func() {
			registerConfig.TLSConfig.ClientAuth = tls.RequireAndVerifyClientCert
		})

		It("makes an mTLS connection with the backend", func() {
			resp := registerAppAndTest()
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
		})
		Context("when the gorouter presents certs that the backend does not trust", func() {
			BeforeEach(func() {
				registerConfig.TLSConfig.ClientCAs = x509.NewCertPool()
			})
			It("returns a HTTP 496 status code", func() {
				resp := registerAppAndTest()
				Expect(resp.StatusCode).To(Equal(496))
			})
		})
		Context("when the gorouter does not present certs", func() {
			BeforeEach(func() {
				conf.Backends.ClientAuthCertificate = tls.Certificate{}
			})
			It("returns a HTTP 496 status code", func() {
				resp := registerAppAndTest()
				Expect(resp.StatusCode).To(Equal(496))
			})
		})
	})

	Context("when the backend instance certificate is signed with an invalid CA", func() {
		BeforeEach(func() {
			var err error
			caCertPool, err = x509.SystemCertPool()
			Expect(err).ToNot(HaveOccurred())
		})

		It("returns a HTTP 526 status code", func() {
			resp := registerAppAndTest()
			Expect(resp.StatusCode).To(Equal(526))
		})
	})

	Context("when the backend instance id does not match the common name on the backend's cert", func() {
		BeforeEach(func() {
			registerConfig.InstanceId = "foo-instance"
		})

		It("returns a HTTP 503 Service Unavailable error", func() {
			resp := registerAppAndTest()
			Expect(resp.StatusCode).To(Equal(http.StatusServiceUnavailable))
		})
	})

	Context("when the backend instance returns a cert that only has a DNS SAN", func() {
		BeforeEach(func() {
			proxyCertPool := freshProxyCACertPool()
			backendCertChain := createCertAndAddCA(test_util.CertNames{
				SANs: test_util.SubjectAltNames{DNS: registerConfig.InstanceId},
			}, proxyCertPool)
			registerConfig.TLSConfig = backendCertChain.AsTLSConfig()

		})

		It("returns a successful 200 OK response from the backend", func() {
			resp := registerAppAndTest()
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
		})
	})

	Context("when the backend instance returns a cert that has a matching CommonName but non-matching DNS SAN", func() {
		BeforeEach(func() {
			proxyCertPool := freshProxyCACertPool()
			backendCertChain := createCertAndAddCA(test_util.CertNames{
				CommonName: registerConfig.InstanceId,
				SANs:       test_util.SubjectAltNames{DNS: "foo"},
			}, proxyCertPool)
			registerConfig.TLSConfig = backendCertChain.AsTLSConfig()
		})

		It("returns a HTTP 503 Service Unavailable error", func() {
			resp := registerAppAndTest()
			Expect(resp.StatusCode).To(Equal(http.StatusServiceUnavailable))
		})
	})

	Context("when the backend instance returns a cert that has a non-matching CommonName but matching DNS SAN", func() {
		BeforeEach(func() {
			proxyCertPool := freshProxyCACertPool()
			backendCertChain := createCertAndAddCA(test_util.CertNames{
				CommonName: "foo",
				SANs:       test_util.SubjectAltNames{DNS: registerConfig.InstanceId},
			}, proxyCertPool)
			registerConfig.TLSConfig = backendCertChain.AsTLSConfig()
		})

		It("returns a successful 200 OK response from the backend", func() {
			resp := registerAppAndTest()
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
		})
	})

	Context("when the backend instance returns a cert that has a matching CommonName but non-matching IP SAN", func() {
		BeforeEach(func() {
			proxyCertPool := freshProxyCACertPool()
			backendCertChain := createCertAndAddCA(test_util.CertNames{
				CommonName: registerConfig.InstanceId,
				SANs:       test_util.SubjectAltNames{IP: "192.0.2.1"},
			}, proxyCertPool)
			registerConfig.TLSConfig = backendCertChain.AsTLSConfig()
		})

		It("returns a successful 200 OK response from the backend (only works for Go1.8 and before)", func() {
			resp := registerAppAndTest()
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
		})
	})

	Context("when the backend instance returns a cert that has a non-matching CommonName but matching IP SAN", func() {
		BeforeEach(func() {
			proxyCertPool := freshProxyCACertPool()
			backendCertChain := createCertAndAddCA(test_util.CertNames{
				CommonName: "foo",
				SANs:       test_util.SubjectAltNames{IP: "127.0.0.1"},
			}, proxyCertPool)
			registerConfig.TLSConfig = backendCertChain.AsTLSConfig()
		})

		It("returns with a HTTP 503 Service Unavailable error (possible route integrity failure)", func() {
			resp := registerAppAndTest()
			Expect(resp.StatusCode).To(Equal(http.StatusServiceUnavailable))
		})
	})

	Context("when the backend registration does not include instance id", func() {
		BeforeEach(func() {
			registerConfig.InstanceId = ""
		})

		It("fails to validate (backends registering with a tls_port MUST provide a name that we can validate on their server certificate)", func() {
			resp := registerAppAndTest()
			Expect(resp.StatusCode).To(Equal(http.StatusServiceUnavailable))
		})
	})

	Context("when the backend is only listening for non TLS connections", func() {
		BeforeEach(func() {
			registerConfig.IgnoreTLSConfig = true
		})
		It("returns a HTTP 525 SSL Handshake error", func() {
			resp := registerAppAndTest()
			Expect(resp.StatusCode).To(Equal(525))
		})
	})
})