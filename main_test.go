
package main

import (
	//"bytes"
	//"errors"
	"net/http"
	"net/http/httptest"
	"fmt"
	//"strings"
	
	"testing"
	"time"
	
	"github.com/labstack/echo"
	//"github.com/labstack/echo/middleware"
	
	"github.com/stretchr/testify/assert"
)

func Test_EchoHello(t *testing.T) {
    pl := fmt.Println
    
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/hello", nil)  //(echo.POST, "auth/logout", strings.NewReader(""))
	rec := httptest.NewRecorder()
	//c := e.NewContext(req, rec)
    //c.Handler()(c)

    e.GET("/hello", HelloHd)
    
    e.ServeHTTP(rec, req)
    
    time.Sleep(200 * time.Millisecond)
    
    pl(rec.Code, rec.Body.String())
    
	// Router
    assert.NotNil(t, e.Router())
    assert.Equal(t, http.StatusOK, rec.Code)
	//assert.Equal(t, http.StatusNotFound, rec.Code)
	
}


