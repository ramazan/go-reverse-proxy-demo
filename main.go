package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func main() {
	e := echo.New()

	// create the reverse proxy
	url, _ := url.Parse("https://jsonplaceholder.typicode.com")
	proxy := httputil.NewSingleHostReverseProxy(url)

	reverseProxyRoutePrefix := "/user"
	routerGroup := e.Group(reverseProxyRoutePrefix)
	routerGroup.Use(func(handlerFunc echo.HandlerFunc) echo.HandlerFunc {
	    return func(context echo.Context) error {

		 req := context.Request()
		 res := context.Response().Writer

		 //may be some extra validations before sending request like Auth etc.
		 if req.Header.Get("X-Custom-Header") != "123" {
		     return echo.NewHTTPError(http.StatusUnauthorized, "Please provide valid credentials")
		 }

		  // Update the headers to allow for SSL redirection
		  req.Host = url.Host
		  req.URL.Host = url.Host
		  req.URL.Scheme = url.Scheme

		  //trim reverseProxyRoutePrefix
		  path := req.URL.Path
		  req.URL.Path = strings.TrimLeft(path, reverseProxyRoutePrefix)

		  // Note that ServeHttp is non blocking and uses a go routine under the hood
		  proxy.ServeHTTP(res, req)
                  return nil
	    }
	})

	e.Start(":2957")

}
