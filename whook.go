package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"strings"

	"github.com/tidusant/c3m-common/c3mcommon"
	"github.com/tidusant/c3m-common/log"
	"github.com/tidusant/c3m-common/mycrypto"
	"github.com/tidusant/chadmin-repo/models"
	rpsex "github.com/tidusant/chadmin-repo/session"
	//"io"

	"net/http"
	//	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

func init() {

}

func main() {
	var port int
	var debug bool
	var mytoken string
	//fmt.Println(mycrypto.Encode("abc,efc", 5))
	flag.IntVar(&port, "port", 8181, "help message for flagname")
	flag.BoolVar(&debug, "debug", false, "Indicates if debug messages should be printed in log files")
	flag.StringVar(&mytoken, "mytoken", "xHrldoad34KDOSkkm", "mytoken")
	flag.Parse()

	logLevel := log.DebugLevel
	if !debug {
		logLevel = log.InfoLevel
		gin.SetMode(gin.ReleaseMode)
	}
	log.Debugf("encode %s", mycrypto.Encode2("uo,ghtk,5a52f89a7b4b30ed5ecfce9d"))

	log.SetOutputFile(fmt.Sprintf("portal-"+strconv.Itoa(port)), logLevel)
	defer log.CloseOutputFile()
	log.RedirectStdOut()

	log.Infof("running with port:" + strconv.Itoa(port))

	//init config

	router := gin.Default()

	router.POST("/:name/:encode", func(c *gin.Context) {
		strrt := ""
		if c.Param("name") == mytoken {
			c.Header("Access-Control-Allow-Origin", "*")
			// requestDomain := c.Request.Header.Get("Origin")
			// if requestDomain == "" {
			// 	requestDomain = "http://" + c.Request.Host
			// }

			// allowDomain := c3mcommon.CheckDomain(requestDomain)

			//if allowDomain != "" {
			//c.Header("Access-Control-Allow-Origin", allowDomain)
			c.Header("Access-Control-Allow-Headers", "access-control-allow-origin, access-control-allow-headers,access-control-allow-credentials")
			c.Header("Access-Control-Allow-Credentials", "true")
			log.Debugf("check request:%s", c.Request.URL.Path)
			if rpsex.CheckRequest(c.Request.URL.Path, c.Request.UserAgent(), c.Request.Referer(), c.Request.RemoteAddr, "POST") {
				strrt = myRoute(c, "")
			} else {
				log.Debugf("check request error")
			}
			// } else {
			// 	log.Debugf("Not allow " + requestDomain)
			// }

		}
		if strrt == "" {
			strrt = c3mcommon.Fake64()
		}
		c.String(http.StatusOK, strrt)

	})

	router.Run(":" + strconv.Itoa(port))

}

func myRoute(c *gin.Context, rpcname string) string {
	encode := mycrypto.Decode2(c.Param("encode"))

	if encode == "" {
		return ""
	}
	sdata := strings.Split(encode, ",")
	if len(sdata) < 3 {
		return ""
	}
	//data := c.PostForm("label_id")
	var args struct {
		LabelID    string `json:"label_id"`
		StatusID   string `json:"status_id"`
		PartnerID  string `json:"partner_id"`
		ActionTime string `json:"action_time"`
		ReasonCode string `json:"reason_code"`
		Reason     string `json:"reason"`
		Weight     string `json:"weight"`
		Fee        string `json:"fee"`
	}

	args.LabelID = c.PostForm("label_id")
	args.StatusID = c.PostForm("status_id")
	args.PartnerID = c.PostForm("partner_id")
	args.ActionTime = c.PostForm("action_time")
	args.ReasonCode = c.PostForm("reason_code")
	args.Reason = c.PostForm("reason")
	args.Weight = c.PostForm("weight")
	args.Fee = c.PostForm("fee")

	//userIP, _, _ := net.SplitHostPort(c.Request.RemoteAddr)
	info, _ := json.Marshal(args)
	strrt := string(info)
	//save into db
	var whook models.Whook
	whook.Name = sdata[1]
	whook.Action = sdata[0]
	whook.Data = strrt
	whook.ShopID = sdata[2]
	whook.Status = 0
	rpsex.SaveWhook(whook)
	log.Debugf("action %s, partner %s, data %v", sdata[0], sdata[1], sdata[2])
	return ""
}
