package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {

	// initialization of config file
	viper.SetConfigName("settings") // name of config file (without extension)
	viper.AddConfigPath(".")        // optionally look for config in the working directory
	err := viper.ReadInConfig()     // Find and read the config file
	if err != nil {                 // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n ", err))
	}

	// init and run HTTP server
	router := gin.Default()

	//POST func to serve file list from requested folder
	router.POST("/notify", func(c *gin.Context) {

		//Read POST data
		groupname := c.PostForm("group")
		message := c.PostForm("message")

		mobilereceivers := viper.GetStringSlice("GROUPS." + groupname + ".mobiles")

		for _, mobileno := range mobilereceivers {

			//make a web request to SMS vendor URL
			smsurl := "http://www.smsjust.com/sms/user/urlsms.php?username=mkclos_trans&pass=Trans@123&senderid=MKCLTD&msgtype=TXT&dltentityid=1201158047881908712&dltheaderid=1205158079076000975&dest_mobileno=" + mobileno + "&message=" + url.QueryEscape(message)

			fmt.Printf(smsurl)

			resp, err := http.Get(smsurl)

			if err != nil {
				log.Fatalln(err)
			}
			defer resp.Body.Close()
		}

		//Responce to POST
		if err != nil {
			//responce if error in JSON format
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "ERROR",
				"error":  err,
			})
		} else {
			//Responce files list with JSON
			c.JSON(http.StatusOK, gin.H{
				"status": "OK",
				"Mobile": mobilereceivers,
				"mess":   message,
			})
		}
	})

	router.Run(":10000")

}
