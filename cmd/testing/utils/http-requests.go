package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"sync"
	"time"
)

const defaultEndpoint = "http://cloudy-pics-asg-lb-258440012.us-east-1.elb.amazonaws.com:3000"
const defaultTimeout = 3 * time.Second

func timeHTTPRequest(url string, withBody bool, bearer string, method string, body string) (string, int, error) {
	// Add the default endpoint if the URL is relative
	if url[0] == '/' {
		url = defaultEndpoint + url
	}

	client := http.Client{
		Timeout: defaultTimeout,
	}

	// Data to be sent in the request body
	var jsonData []byte
	if body != "" {
		data := map[string]string{
			"username": body,
		}

		// Marshal the data into JSON
		var err error
		jsonData, err = json.Marshal(data)
		if err != nil {
			fmt.Println("Error marshaling data:", err)
			return "", 0, err
		}
	}

	// Create a new request
	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return "", 0, err
	}

	// Set the authorization header
	req.Header.Set("Authorization", "Bearer "+bearer)

	// Set the content type if there is a body
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}

	start := time.Now()
	resp, err := client.Do(req)
	if err != nil {
		// T = Timeout
		fmt.Fprintf(os.Stderr, "\033[1;31mT\033[0m")
		fmt.Println(err)
		os.Stderr.Sync()
		return "", 0, err
	}
	defer resp.Body.Close()

	// Check that status code is 200
	if !(resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated || resp.StatusCode == http.StatusNoContent) {
		// F = Failed request
		fmt.Fprintf(os.Stderr, "\033[1;31mF\033[0m")
		os.Stderr.Sync()
		return "", 0, fmt.Errorf("status code %d", resp.StatusCode)
	}

	// Calculate the time taken for the request in milliseconds
	duration := time.Since(start).Milliseconds()

	// Get the response body as a string, if needed
	var responseBody string
	if withBody {
		var responseMap map[string]interface{}
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", 0, err
		}
		err = json.Unmarshal(bodyBytes, &responseMap)
		if err != nil {
			return "", 0, err
		}
		if id, ok := responseMap["id"].(string); ok {
			responseBody = id
		} else {
			return "", 0, fmt.Errorf("ID field not found or not a string")
		}
	}

	// Print the success, as a green dot
	fmt.Fprintf(os.Stderr, "\033[32m.\033[0m")
	os.Stderr.Sync()

	return responseBody, int(duration), nil
}

func timeHTTPPostFile(url string, filePath string, fileName string, bearer string) (string, int, error) {
	// Add the default endpoint if the URL is relative
	if url[0] == '/' {
		url = defaultEndpoint + url
	}

	client := http.Client{
		Timeout: defaultTimeout,
	}

	file, err := os.Open(filePath)
	if err != nil {
		return "", 0, fmt.Errorf("Error: could not open file %s", filePath)
	}
	defer file.Close()

	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	// Add the file to the request
	part, err := writer.CreateFormFile("photo", fileName)
	if err != nil {
		return "", 0, fmt.Errorf("Error: could not create form file")
	}

	if _, err := io.Copy(part, file); err != nil {
		return "", 0, fmt.Errorf("Error: could not copy file to part")
	}

	if err := writer.Close(); err != nil {
		return "", 0, fmt.Errorf("Error: could not close writer")
	}

	// Create the request (without sending it)
	req, err := http.NewRequest("POST", url, &requestBody)
	if err != nil {
		return "", 0, fmt.Errorf("Error: could not create request")
	}

	// Set the content type
	req.Header.Set("Content-Type", "image/jpeg")

	// Set the authorization header
	req.Header.Set("Authorization", "Bearer "+bearer)

	// Send the request
	start := time.Now()
	resp, err := client.Do(req)
	if err != nil {
		// T = Timeout
		fmt.Fprintf(os.Stderr, "\033[1;31mT\033[0m")
		os.Stderr.Sync()
		return "", 0, err
	}
	defer resp.Body.Close()

	// Check that status code is 200
	if resp.StatusCode != http.StatusCreated {
		// F = Failed request
		fmt.Fprintf(os.Stderr, "\033[1;31mF\033[0m")
		os.Stderr.Sync()
		return "", 0, fmt.Errorf("status code %d", resp.StatusCode)
	}

	// Calculate the time taken for the request in milliseconds
	duration := time.Since(start).Milliseconds()

	var responseMap map[string]interface{}
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", 0, err
	}
	err = json.Unmarshal(bodyBytes, &responseMap)
	if err != nil {
		return "", 0, err
	}

	// Extract the ID field
	var responseBody string
	if id, ok := responseMap["id"].(string); ok {
		responseBody = id
	} else {
		return "", 0, fmt.Errorf("ID field not found or not a string")
	}

	// Print the success, as a green dot
	fmt.Fprintf(os.Stderr, "\033[32m.\033[0m")
	os.Stderr.Sync()

	return responseBody, int(duration), nil
}

func TimeHTTPRequest(url string, bearer string, method string, body string) (int, error) {
	_, duration, err := timeHTTPRequest(url, false, bearer, method, body)
	return duration, err
}

func TimeHTTPRequestWithBody(url string, bearer string, method string, body string) (string, int, error) {
	return timeHTTPRequest(url, true, bearer, method, body)
}

func TimeHTTPRequestWaiting(url string, wg *sync.WaitGroup) (int, error) {
	defer wg.Done()
	return TimeHTTPRequest(url, "", "GET", "")
}
