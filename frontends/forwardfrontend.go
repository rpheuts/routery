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
	"github.com/rpheuts/routery/authentication"
	"github.com/rpheuts/routery/config"
	"encoding/base64"
)

type ForwardFrontend struct {
	config              *FrontendConfig
	routeryConfig       *config.RouteryConfig
	routeRequestChannel chan *router.RouteRequest
	routes              []*router.RouteRequest
}

func (ff *ForwardFrontend) Initialize(config *FrontendConfig, routeryConfig *config.RouteryConfig) error {
	ff.config = config
	ff.routeryConfig = routeryConfig

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

	proxy := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
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
		Handler: ff.basicAuth(proxy),
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
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
		},
	}

	mTLSConfig.PreferServerCipherSuites = true
	mTLSConfig.MinVersion = tls.VersionTLS10
	mTLSConfig.MaxVersion = tls.VersionTLS12

	certs := x509.NewCertPool()

	pemData, err := ioutil.ReadFile(ff.config.CA)
	if err != nil {
		// do error
	}
	certs.AppendCertsFromPEM(pemData)
	mTLSConfig.RootCAs = certs

	return mTLSConfig
}

func (ff *ForwardFrontend) basicAuth(pass http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Authorization"] == nil || len(r.Header["Authorization"]) <= 0 {
			w.Header().Add("WWW-Authenticate", "Basic")
			http.Error(w, "", 401)
			return
		}

		auth := strings.SplitN(r.Header["Authorization"][0], " ", 2)

		if len(auth) != 2 || auth[0] != "Basic" {
			http.Error(w, "bad syntax", http.StatusBadRequest)
			return
		}

		payload, _ := base64.StdEncoding.DecodeString(auth[1])
		pair := strings.SplitN(string(payload), ":", 2)

		log.Printf("Auth received for user: %v\n", pair[0])
		if len(pair) != 2 || !authentication.Authenticate(ff.routeryConfig, pair[0], pair[1]) {
			http.Error(w, "authorization failed", http.StatusUnauthorized)
			return
		}

		pass(w, r)
	}
}
