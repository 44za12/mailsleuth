package processor

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"

	"github.com/44za12/mailsleuth/internal/utils"
	"github.com/44za12/mailsleuth/pkg/shopping/amazon"
	"github.com/44za12/mailsleuth/pkg/social/instagram"
	"github.com/44za12/mailsleuth/pkg/social/spotify"
	"github.com/44za12/mailsleuth/pkg/social/x"
	"github.com/cheggaaa/pb/v3"
	"github.com/fatih/color"
)

type Processor struct {
	Email   string
	Proxy   string
	Verbose bool
}

type ProcessorMany struct {
	Emails           []string
	Proxies          []string
	ConcurrencyLimit int
	Verbose          bool
}

func (p *Processor) Process() error {
	if p.Verbose {
		msg := fmt.Sprintf("Checking for email: %s with proxy: %s", p.Email, p.Proxy)
		fmt.Println(color.YellowString(msg))
	}
	if p.Email == "" {
		return errors.New("email required but not provided")
	}
	client, err := utils.NewHttpClient(p.Proxy)
	if err != nil {
		return fmt.Errorf("failed to create HTTP client: %v", err)
	}

	results, err := processSingleEmail(p.Email, client, p.Verbose)
	if err != nil {
		return err
	}
	for service, exists := range results {
		if exists {
			fmt.Printf(color.GreenString("%s account exists for: %s\n"), service, p.Email)
		} else {
			fmt.Printf(color.HiBlueString("%s account does not exist for: %s\n"), service, p.Email)
		}
	}
	return nil
}

func (pm *ProcessorMany) Process() (map[string]map[string]bool, error) {
	if len(pm.Emails) == 0 {
		return nil, errors.New("no emails provided")
	}

	results := make(map[string]map[string]bool)
	var mu sync.Mutex
	var wg sync.WaitGroup
	progressBar := pb.StartNew(len(pm.Emails))
	sem := make(chan struct{}, pm.ConcurrencyLimit)

	for i, email := range pm.Emails {
		wg.Add(1)
		sem <- struct{}{}

		go func(email string, proxyIndex int) {
			defer wg.Done()
			defer func() { <-sem }()
			defer progressBar.Increment()
			proxy := ""
			if len(pm.Proxies) > 0 {
				proxy = pm.Proxies[proxyIndex%len(pm.Proxies)]
			}

			client, err := utils.NewHttpClient(proxy)
			if err != nil {
				if pm.Verbose {
					msg := fmt.Sprintf("Error creating HTTP client for %s: %v\n", email, err)
					utils.LogError(msg)
				}
				return
			}

			emailResults, err := processSingleEmail(email, client, pm.Verbose)
			if err != nil {
				if pm.Verbose {
					msg := fmt.Sprintf("Error processing email %s: %v\n", email, err)
					utils.LogError(msg)
				}
				return
			}

			mu.Lock()
			results[email] = emailResults
			mu.Unlock()
		}(email, i)
	}

	wg.Wait()
	progressBar.Finish()
	return results, nil
}

func processSingleEmail(email string, client *http.Client, verbose bool) (map[string]bool, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	results := make(map[string]bool)
	var mu sync.Mutex
	var wg sync.WaitGroup

	services := map[string]func(string, *http.Client) (bool, error){
		"instagram": instagram.Check,
		"x":         x.Check,
		"spotify":   spotify.Check,
		"amazon":    amazon.Check,
	}

	for name, checkFunc := range services {
		wg.Add(1)
		go func(name string, checkFunc func(string, *http.Client) (bool, error)) {
			defer wg.Done()

			exists, err := checkFunc(email, client)
			if err != nil {
				if verbose {
					utils.LogError(fmt.Sprintf("Error checking %s for %s: %v\n", name, email, err))
				}
				return
			}

			mu.Lock()
			results[name] = exists
			mu.Unlock()
		}(name, checkFunc)
	}

	wg.Wait()

	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	return results, nil
}
