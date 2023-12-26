package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"
)

type Response struct {
	Token string `json:"token"`
}

// generateRandomString generates a random string of the specified length
func generateRandomString(length int) string {
	characters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rand.Seed(time.Now().UnixNano())
	result := make([]byte, length)
	for i := range result {
		result[i] = characters[rand.Intn(len(characters))]
	}
	return string(result)
}

func main() {
	// Open a file for writing
	file, err := os.Create("tokens.txt")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	// Create a writer for the file
	writer := bufio.NewWriter(file)

	// Close the writer when the program ends
	defer writer.Flush()

	// Ask the user for the number of tokens
	var numTokens int
	fmt.Print("Enter the number of tokens you need: ")
	fmt.Scan(&numTokens)

	for i := 0; i < numTokens; i++ {
		url := "https://api.discord.gx.games/v1/direct-fulfillment"
		method := "POST"
		payload := map[string]string{
			"partnerUserId": generateRandomString(64),
		}

		jsonPayload, err := json.Marshal(payload)
		if err != nil {
			fmt.Println("Error encoding JSON:", err)
			return
		}

		req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonPayload))

		if err != nil {
			fmt.Println("Error creating request:", err)
			return
		}

		// Set headers
		req.Header.Set("authority", "api.discord.gx.games")
		req.Header.Set("accept", "*/*")
		req.Header.Set("accept-language", "en-US,en;q=0.9")
		req.Header.Set("content-type", "application/json")
		req.Header.Set("origin", "https://www.opera.com")
		req.Header.Set("referer", "https://www.opera.com/")
		req.Header.Set("sec-ch-ua", "\"Opera GX\";v=\"105\", \"Chromium\";v=\"119\", \"Not?A_Brand\";v=\"24\"")
		req.Header.Set("sec-ch-ua-mobile", "?0")
		req.Header.Set("sec-ch-ua-platform", "\"Windows\"")
		req.Header.Set("sec-fetch-dest", "empty")
		req.Header.Set("sec-fetch-mode", "cors")
		req.Header.Set("sec-fetch-site", "cross-site")
		req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36 OPR/105.0.0.0")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error making request:", err)
			return
		}
		defer resp.Body.Close()

		var responseStruct Response
		err = json.NewDecoder(resp.Body).Decode(&responseStruct)
		if err != nil {
			fmt.Println("Error decoding response body:", err)
			return
		}

		// Save the token and construct the URL
		token := responseStruct.Token
		urlWithToken := fmt.Sprintf("https://discord.com/billing/partner-promotions/1180231712274387115/%s", token)

		// Print the token and URL to the console
		fmt.Printf("Token %d: %s\n", i+1, token)
		fmt.Printf("URL with Token %d: %s\n", i+1, urlWithToken)

		// Write the URL with the token to the file
		_, err = writer.WriteString(urlWithToken + "\n")
		if err != nil {
			fmt.Println("Error writing to file:", err)
			return
		}

		// Pause for a short duration to avoid rate limiting
		time.Sleep(time.Second)
	}

	fmt.Println("Tokens and URLs saved to tokens.txt")
}
