package processor

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/44za12/mailsleuth/internal/requestor"
	"github.com/44za12/mailsleuth/internal/utils"
	"github.com/44za12/mailsleuth/pkg/entertainment/spotify"
	"github.com/44za12/mailsleuth/pkg/programming/github"
	"github.com/44za12/mailsleuth/pkg/shopping/amazon"
	"github.com/44za12/mailsleuth/pkg/social/facebook"
	"github.com/44za12/mailsleuth/pkg/social/instagram"
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
		if p.Proxy != "" {
			msg := fmt.Sprintf("Checking for email: %s with proxy: %s\n", p.Email, p.Proxy)
			fmt.Println(color.YellowString(msg))
		} else {
			msg := fmt.Sprintf("Checking for email: %s\n", p.Email)
			fmt.Println(color.YellowString(msg))
		}
	}
	if p.Email == "" {
		return errors.New("email required but not provided")
	}
	results, err := processSingleEmail(p.Email, p.Proxy, p.Verbose)
	if err != nil {
		return err
	}
	for service, exists := range results {
		if exists {
			fmt.Printf(color.GreenString("✅ %s account exists for: %s\n"), service, p.Email)
		} else {
			fmt.Printf(color.HiRedString("❌ %s account does not exist for: %s\n"), service, p.Email)
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

			emailResults, err := processSingleEmail(email, proxy, pm.Verbose)
			if err != nil {
				if pm.Verbose {
					msg := fmt.Sprintf("error processing email %s: %v\n", email, err)
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

func processSingleEmail(email string, proxy string, verbose bool) (map[string]bool, error) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	results := make(map[string]bool)
	var mu sync.Mutex
	var wg sync.WaitGroup

	services := map[string]func(string, *requestor.Requestor) (bool, error){
		"instagram": instagram.Check,
		"x":         x.Check,
		"spotify":   spotify.Check,
		"amazon":    amazon.Check,
		"facebook":  facebook.Check,
		"github":    github.Check,
	}

	for name, checkFunc := range services {
		wg.Add(1)
		go func(name string, checkFunc func(string, *requestor.Requestor) (bool, error)) {
			defer wg.Done()
			requestorObj, err := requestor.NewRequestor(email, proxy)
			if err != nil {
				utils.LogError("error getting new requestor")
				return
			}
			exists, err := checkFunc(email, requestorObj)
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
