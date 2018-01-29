package main

import (
	"encoding/json"
	"flag"

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

	//logLevel := log.DebugLevel
	if !debug {
		//logLevel = log.InfoLevel
		gin.SetMode(gin.ReleaseMode)
	}
	log.Debugf("encode %s, action: %s", mycrypto.Encode2("ghtk"), mycrypto.Encode2("updateorderstatus"))

	// log.SetOutputFile(fmt.Sprintf("portal-"+strconv.Itoa(port)), logLevel)
	// defer log.CloseOutputFile()
	// log.RedirectStdOut()

	log.Infof("running with port:" + strconv.Itoa(port))

	//init config

	router := gin.Default()

	router.POST("/:name/:action/:partner", func(c *gin.Context) {
		strrt := ""
		if c.Param("name") == mytoken {
			requestDomain := c.Request.Header.Get("Origin")
			if requestDomain == "" {
				requestDomain = "http://" + c.Request.Host
			}

			allowDomain := c3mcommon.CheckDomain(requestDomain)

			c.Header("Access-Control-Allow-Origin", "*")
			if allowDomain != "" {
				c.Header("Access-Control-Allow-Origin", allowDomain)
				c.Header("Access-Control-Allow-Headers", "access-control-allow-origin, access-control-allow-headers,access-control-allow-credentials")
				c.Header("Access-Control-Allow-Credentials", "true")
				log.Debugf("check request:%s", c.Request.URL.Path)
				if rpsex.CheckRequest(c.Request.URL.Path, c.Request.UserAgent(), c.Request.Referer(), c.Request.RemoteAddr, "POST") {
					strrt = myRoute(c, "")
				} else {
					log.Debugf("check request error")
				}
			} else {
				log.Debugf("Not allow " + requestDomain)
			}

		}
		if strrt == "" {
			strrt = c3mcommon.Fake64()
		}
		c.String(http.StatusOK, strrt)

	})

	router.Run(":" + strconv.Itoa(port))

}

func myRoute(c *gin.Context, rpcname string) string {
	action := mycrypto.Decode2(c.Param("action"))
	partner := mycrypto.Decode2(c.Param("partner"))
	if partner == "" || action == "" {
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
	whook.Name = partner
	whook.Action = action
	whook.Data = strrt
	whook.Status = 0
	rpsex.SaveWhook(whook)
	log.Debugf("action %s, partner %s, data %v", action, partner, strrt)
	return ""
}
