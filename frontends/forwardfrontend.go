package frontends

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/mailgun/oxy/forward"
	"github.com/mailgun/oxy/testutils"
	"github.com/rpheuts/routery/router"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type ForwardFrontend struct {
	config              *FrontendConfig
	routeRequestChannel chan *router.RouteRequest
	routes              []*router.RouteRequest
}

func (ff *ForwardFrontend) Initialize(config *FrontendConfig) error {
	ff.config = config

	return nil
}

func (ff *ForwardFrontend) Route(routeRequestChannel chan *router.RouteRequest) error {
	if !ff.config.Enabled {
		return nil
	}

	ff.routeRequestChannel = routeRequestChannel

	// Start watching route requests
	go ff.watchRouteRequests()

	// Start the forwarder
	go ff.watchWebRequests()

	return nil
}

func (ff *ForwardFrontend) watchRouteRequests() {
	for {
		event := <-ff.routeRequestChannel

		if !event.Remove {
			ff.routes = append(ff.routes, event)
			log.Printf("%v:%v: Received route-add request. %v\n", ff.config.Hostname, ff.config.Port, event)
		} else {
			ff.remove(event)
			log.Printf("%v:%v: Received route-remove request. %v\n", ff.config.Hostname, ff.config.Port, event)
		}

	}
}

func (ff *ForwardFrontend) remove(r *router.RouteRequest) bool {
	for i := range ff.routes {
		if ff.routes[i].Id == r.Id {
			copy(ff.routes[i:], ff.routes[i+1:])
			ff.routes[len(ff.routes)-1] = nil
			ff.routes = ff.routes[:len(ff.routes)-1]
			return true
		}
	}
	return false
}

func (ff *ForwardFrontend) watchWebRequests() {
	fwd, _ := forward.New()
	redirect := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		hostname := strings.Split(req.Host, ".")[0]

		for _, event := range ff.routes {
			if event.Hostname == hostname {
				req.URL = testutils.ParseURI(fmt.Sprintf("http://%v:%v", event.Endpoint, event.Port))
				fwd.ServeHTTP(w, req)
				log.Printf("%v:%v:Serving request. Hostname: %v Target: %v Port: %v\n", ff.config.Hostname, ff.config.Port, hostname, event.Endpoint, event.Port)
			}
		}
	})

	s := &http.Server{
		Addr:    fmt.Sprintf(":%v", ff.config.Port),
		Handler: redirect,
	}

	if ff.config.SSL {
		if ff.config.CA != "" {
			s.TLSConfig = ff.getCACert()
		}

		log.Printf("Listening on port: %v TLS: %v\n", ff.config.Port, ff.config.SSL)
		s.ListenAndServeTLS(ff.config.Cert, ff.config.Key)
	} else {
		log.Printf("Listening on port: %v TLS: %v\n", ff.config.Port, ff.config.SSL)
		s.ListenAndServe()
	}

}

func (ff *ForwardFrontend) getCACert() *tls.Config {
	mTLSConfig := &tls.Config{
		CipherSuites: []uint16{
			tls.TLS_RSA_WITH_RC4_128_SHA,
			tls.TLS_RSA_WITH_3DES_EDE_CBC_SHA,
			tls.TLS_RSA_WITH_AES_128_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_RC4_128_SHA,
			tls.TLS_RSA_WITH_AES_128_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
		},
	}

	mTLSConfig.PreferServerCipherSuites = true
	mTLSConfig.MinVersion = tls.VersionTLS10
	mTLSConfig.MaxVersion = tls.VersionTLS10

	certs := x509.NewCertPool()

	pemData, err := ioutil.ReadFile(ff.config.CA)
	if err != nil {
		// do error
	}
	certs.AppendCertsFromPEM(pemData)
	mTLSConfig.RootCAs = certs

	return mTLSConfig
}
