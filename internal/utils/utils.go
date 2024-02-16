package utils

import (
	"bufio"
	"crypto/md5"
	"encoding/csv"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
)

func LogError(msg string) {
	fmt.Println(color.RedString(msg))
}

func GetAllServices() []string {
	services := make([]string, 0, 6)
	services = append(services, "x")
	services = append(services, "instagram")
	services = append(services, "amazon")
	services = append(services, "facebook")
	services = append(services, "spotify")
	services = append(services, "github")
	return services
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
	cyan := color.New(color.FgHiCyan).SprintFunc()
	magenta := color.New(color.FgHiMagenta).SprintFunc()
	for i, line := range lines {
		switch {
		case i%3 == 0:
			fmt.Println(cyan(line))
		case i%3 == 1:
			fmt.Println(magenta(line))
		case i%3 == 2:
			fmt.Println(line)
		}
	}
	developerInfo := "Developed by: Aazar (https://www.github.com/44za12)"
	reachInfo := "Reach: https://aazar.me"

	fmt.Println()
	fmt.Println(color.HiCyanString(developerInfo))
	fmt.Println(color.HiMagentaString(reachInfo))
	fmt.Println()
}

func SaveResponse(htmlContent string, filename string) {
	_ = os.WriteFile(filename, []byte(htmlContent), 0644)
}

func randomUserAgent() string {
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
	result := make([]byte, length)
	result[0] = lettersAndUnderscores[rand.Intn(len(lettersAndUnderscores))]
	for i := 1; i < length; i++ {
		result[i] = lettersNumbersAndUnderscores[rand.Int63()%int64(len(lettersNumbersAndUnderscores))]
	}

	return string(result)
}

func HashString(str string) string {
	hash := md5.Sum([]byte(str))
	return fmt.Sprintf("%x", hash)
}

func StandardHeaders() map[string]string {
	return map[string]string{
		"User-Agent":                randomUserAgent(),
		"Accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8",
		"Accept-Language":           "en-US,en;q=0.5",
		"Accept-Encoding":           "gzip, deflate, br",
		"DNT":                       "1",
		"Connection":                "keep-alive",
		"Upgrade-Insecure-Requests": "1",
	}
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

	writer := csv.NewWriter(file)
	defer writer.Flush()
	serviceNames := GetAllServices()
	header := append([]string{"email"}, GetAllServices()...)
	if err := writer.Write(header); err != nil {
		return err
	}
	for email, services := range results {
		row := make([]string, 1, len(serviceNames)+1)
		row[0] = email
		for _, serviceName := range serviceNames {
			exists, ok := services[serviceName]
			if ok {
				row = append(row, fmt.Sprintf("%t", exists))
			} else {
				row = append(row, "error")
			}
		}
		if err := writer.Write(row); err != nil {
			return err
		}
	}
	return nil
}

func GetCurrentTimeStamp() int {
	now := time.Now()
	return int(now.UnixNano() / 1e6)
}

func GetCurrentTimeStampAsStr() string {
	now := time.Now()
	return strconv.Itoa(int(now.UnixNano() / 1e6))
}
