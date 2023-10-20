package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {
	args := os.Args[1:]
	if len(args) < 2 {
		printHelp()
		return
	}
	countFlagstr := args[0]
	args = args[1:]
	if len(args) < 2 {
		fmt.Println("Please provide a valid search query.")
		return
	}
	countFlag, interr := strconv.Atoi(countFlagstr)
	if interr != nil {
		fmt.Println("Please provide a valid count.")
		return
	}
	args = args[2:]
	queryFlag := strings.Join(args, " ")
	queryFlag = url.QueryEscape(queryFlag)
	client := &http.Client{}
	// Create a new context for the API request.
	ctx := context.Background()

	// Create the Unsplash API request URL.
	url := fmt.Sprintf("https://api.unsplash.com/search/photos?query=%s&per_page=%d", queryFlag, countFlag)

	// Make the API request.
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.Fatalf("Failed to create HTTP request: %v", err)
	}

	// Set the Unsplash API authorization header.
	req.Header.Set("Authorization", "Client-ID kodpBdLnV7SCBziObkhKkVXxN3AnIO4ac4-gt92vneQ")

	// Send the API request.
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to send HTTP request: %v", err)
	}

	defer resp.Body.Close()

	// Decode the API response.
	var results struct {
		Results []struct {
			Urls struct {
				Regular string `json:"regular"`
			} `json:"urls"`
		} `json:"results"`
	}

	err = json.NewDecoder(resp.Body).Decode(&results)
	if err != nil {
		log.Fatalf("Failed to decode API response: %v", err)
	}

	// Download the images to the current directory.
	for i := 0; i < countFlag; i++ {
		photo := results.Results[i]

		// Get the image URL.
		imageURL := photo.Urls.Regular

		// Download the image.
		resp, err := client.Get(imageURL)
		if err != nil {
			log.Fatalf("Failed to download image: %v", err)
		}

		defer resp.Body.Close()

		// Create the file path for the downloaded image.
		filePath := filepath.Join(".", fmt.Sprintf("%d.jpg", i))

		// Write the image to the file.
		file, err := os.Create(filePath)
		if err != nil {
			log.Fatalf("Failed to create file: %v", err)
		}

		defer file.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("Failed to read image body: %v", err)
		}

		_, err = file.Write(body)
		if err != nil {
			log.Fatalf("Failed to write image to file: %v", err)
		}
		cmd := exec.Command("catimg", filePath)
		cmd.StdoutPipe()
	}

	fmt.Println("Successfully downloaded images!")
}

func printHelp() {
	fmt.Println("Usage: i-need <count> images of <query>")
	fmt.Println("Example: i-need 5 images of cats")
}
