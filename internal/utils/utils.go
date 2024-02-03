package utils

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/fatih/color"
)

func LogError(msg string) {
	fmt.Println(color.RedString(msg))
}

func PrintBanner() {
	msg := `   _____           .__ .__         _________.__                    __   .__      
  /     \  _____   |__||  |       /   _____/|  |    ____   __ __ _/  |_ |  |__   
 /  \ /  \ \__  \  |  ||  |       \_____  \ |  |  _/ __ \ |  |  \\   __\|  |  \  
/    Y    \ / __ \_|  ||  |__     /        \|  |__\  ___/ |  |  / |  |  |   Y  \ 
\____|__  /(____  /|__||____/    /_______  /|____/ \___  >|____/  |__|  |___|  / 
        \/      \/                       \/            \/                    \/  
                                                                                 `
	lines := strings.Split(msg, "\n")
	red := color.New(color.FgRed).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	blue := color.New(color.FgBlue).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	for i, line := range lines {
		switch {
		case i%4 == 0:
			fmt.Println(red(line))
		case i%4 == 1:
			fmt.Println(green(line))
		case i%4 == 2:
			fmt.Println(blue(line))
		case i%4 == 3:
			fmt.Println(yellow(line))
		}
	}
}

func RandomUserAgent() string {
	userAgents := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36 Edg/120.0.2210.144",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; Xbox; Xbox One) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36 Edge/44.18363.8131",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/122.0",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:115.0) Gecko/20100101 Firefox/115.0",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36 OPR/106.0.0.0",
		"Mozilla/5.0 (Windows NT 10.0; WOW64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36 OPR/106.0.0.0",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36 Edg/120.0.2210.144",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 14.3; rv:109.0) Gecko/20100101 Firefox/122.0",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 14.3; rv:115.0) Gecko/20100101 Firefox/115.0",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 14_3) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.2 Safari/605.1.15",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 14_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36 OPR/106.0.0.0",
	}
	index := rand.Intn(len(userAgents))
	return userAgents[index]
}

func RandomString(length int) string {
	const lettersAndUnderscores = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_"
	const lettersNumbersAndUnderscores = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_0123456789"

	if length == 0 {
		return ""
	}

	// Ensure the first character is not a number
	result := make([]byte, length)
	result[0] = lettersAndUnderscores[rand.Intn(len(lettersAndUnderscores))]

	// Fill the rest of the string with any allowed character
	for i := 1; i < length; i++ {
		result[i] = lettersNumbersAndUnderscores[rand.Int63()%int64(len(lettersNumbersAndUnderscores))]
	}

	return string(result)
}

func NewHttpClient(proxyURL string) (*http.Client, error) {
	if proxyURL == "" {
		return &http.Client{}, nil
	}

	proxy, err := url.Parse(proxyURL)
	if err != nil {
		return nil, fmt.Errorf("invalid proxy URL: %w", err)
	}

	transport := &http.Transport{
		Proxy: http.ProxyURL(proxy),
	}

	return &http.Client{
		Transport: transport,
	}, nil
}

func LoadEmailsFromFile(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var emails []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		emails = append(emails, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return emails, nil
}

func LoadProxiesFromFile(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var proxies []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		proxies = append(proxies, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return proxies, nil
}

func OutputResultsToFile(filePath string, results map[string]map[string]bool) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // SetIndent for pretty printing
	if err := encoder.Encode(results); err != nil {
		return err
	}

	return nil
}
