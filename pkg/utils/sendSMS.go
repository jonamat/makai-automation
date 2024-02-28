package utils

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
)

func SendSMS(message string) error {
	urlEncodedMessage := url.QueryEscape(message)
	openWRTUrl := os.Getenv("OPENWRT_URL")
	phoneNumber := os.Getenv("SMS_PHONE_NUMBER")
	username := os.Getenv("SMS_USERNAME")
	password := os.Getenv("SMS_PASSWORD")

	req, err := http.NewRequest("GET", openWRTUrl+"/cgi-bin/sms_send?username="+username+"&password="+password+"&number="+phoneNumber+"&text="+urlEncodedMessage, nil)
	if err != nil {
		fmt.Println("Error during login request creation", err.Error())
		return err
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("Error during login request execution", err.Error())
		return err
	}

	if res.StatusCode != http.StatusOK {
		fmt.Println("Error during login request execution", res.Status)
		return fmt.Errorf("error during login request execution: %s", res.Status)
	}

	return nil
}
