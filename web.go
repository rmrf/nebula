package nebula

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// WebHostInfo for web api result for return
type WebHostInfo struct {
	Name string `json:"name"`
	IPs  string `json:"ips"`
}

func startWeb(hostMap *HostMap, addr, token string) {
	type PostBody struct {
		Token string `json:"token""  binding:"required"`
	}

	r := gin.New()

	r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// your custom format
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))
	r.Use(gin.Recovery())
	r.POST("/api/hostmap", func(gc *gin.Context) {
		var pBody PostBody
		if err := gc.ShouldBindJSON(&pBody); err != nil {
			gc.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if pBody.Token != token {
			log.Printf("Wrong token: %s", token)
			gc.JSON(http.StatusForbidden, gin.H{"message": "bad token"})
			return
		}
		var webHostsInfo []WebHostInfo
		for _, info := range hostMap.Hosts {
			cert := info.GetCert().Details
			ips := fmt.Sprintf("%s", cert.Ips)
			hi := WebHostInfo{Name: cert.Name, IPs: ips}
			webHostsInfo = append(webHostsInfo, hi)
		}
		gc.JSON(http.StatusOK, webHostsInfo)
		return
	})
	s := &http.Server{Addr: addr,
		Handler:      r,
		ReadTimeout:  6 * time.Second,
		WriteTimeout: 6 * time.Second}
	s.ListenAndServe()

}
