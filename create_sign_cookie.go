// create_sign_cookie.go
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/service/cloudfront/sign"
)

func main() {
	// Key Id
	keyID := "XXXXXXXXX"    
                                            
        // CloudFront private key                                                 
	privKeyPath := "./your_peivete_key.pem"   
 
        // CloudFront Resouce URL                                                                   
	url := "https://xxxxxxxx.cloudfront.net/512MB.zip" 

	expireAt := time.Now().Add(24 * time.Hour)

	privKey, err := sign.LoadPEMPrivKeyFile(privKeyPath)
	if err != nil {
		log.Fatalf("Load private key from %s failed\n", privKeyPath)
	}

	s := sign.NewCookieSigner(keyID, privKey)

	policy := &sign.Policy{
		Statements: []sign.Statement{
			{
				// Read the provided documentation on how to set this
				// correctly, you'll probably want to use wildcards.
				Resource: "http*://*.cloudfront.net/*",
				Condition: sign.Condition{
					// Optional IP source address range
					//IPAddress: &sign.IPAddress{SourceIP: "192.0.2.0/24"},
					// Optional date URL is not valid until
					//DateGreaterThan: &sign.AWSEpochTime{time.Now().Add(30 * time.Minute)},
					// Required date the URL will expire after
					DateLessThan: &sign.AWSEpochTime{expireAt},
				},
			},
		},
	}

	// Get Signed cookies for a resource that will expire in 1000 hour
	signedCookies, err := s.SignWithPolicy(policy)
	if err != nil {
		fmt.Errorf("failed to SignWithPolicy errorMsg: %s", err.Error())
	}

	for _, c := range signedCookies {
		fmt.Printf("%s: %s, %s, %s, %t\n", c.Name, c.Value, c.Path, c.Domain, c.Secure)
	}

	fmt.Printf("access with signedCookies: %s\ncontent: %s", signedCookies, httpGetWithCookie(url, signedCookies))
}

func httpGet(url string) string {
	res, err := http.Get(url)

	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	data, _ := ioutil.ReadAll(res.Body)

	return fmt.Sprintf("%s", data)
}

func httpGetWithCookie(url string, cookies []*http.Cookie) string {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	for _, c := range cookies {
		req.AddCookie(c)
	}

	client := &http.Client{}
	res, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	data, _ := ioutil.ReadAll(res.Body)

	return fmt.Sprintf("%d", len(data))
}
