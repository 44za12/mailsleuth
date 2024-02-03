package main

import (
	"fmt"
	"log"
	"os"

	"github.com/44za12/mailsleuth/internal/processor"
	"github.com/44za12/mailsleuth/internal/utils"

	"github.com/urfave/cli/v2"
)

func main() {
	utils.PrintBanner()
	app := &cli.App{
		Name:  "MailSleuth",
		Usage: "An extremely quick and efficient email osint tool.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "email",
				Aliases: []string{"e"},
				Usage:   "Single email to check",
			},
			&cli.StringFlag{
				Name:    "proxy",
				Aliases: []string{"p"},
				Usage:   "Single proxy to use (Format: http://HOST:PORT or http://USER:PASS@HOST:PORT)",
			},
			&cli.StringFlag{
				Name:    "emails-file",
				Aliases: []string{"E"},
				Usage:   "File containing list of emails (one email per line)",
			},
			&cli.StringFlag{
				Name:    "proxies-file",
				Aliases: []string{"P"},
				Usage:   "File containing list of proxies (One proxy per line; Each proxy format: http://HOST:PORT or http://USER:PASS@HOST:PORT)",
			},
			&cli.IntFlag{
				Name:    "concurrency",
				Aliases: []string{"c"},
				Value:   10,
				Usage:   "Concurrency limit for processing emails",
			},
			&cli.StringFlag{
				Name:    "output",
				Aliases: []string{"o"},
				Usage:   "Output file for results (valid only with --emails-file) Example: results.json",
			},
			&cli.BoolFlag{
				Name:    "list",
				Aliases: []string{"l"},
				Usage:   "List available services",
			},
		},
		Action: func(c *cli.Context) error {
			if c.NArg() == 0 && c.NumFlags() == 0 {
				cli.ShowAppHelp(c)
				return nil
			}
			if c.Bool("list") {
				fmt.Println("Available services: [instagram, spotify, x]")
				return nil
			}
			email := c.String("email")
			proxy := c.String("proxy")
			emailsFile := c.String("emails-file")
			proxiesFile := c.String("proxies-file")
			concurrency := c.Int("concurrency")
			outputFile := c.String("output")

			if outputFile != "" && emailsFile == "" {
				return cli.Exit("--output is valid only with --emails-file", 1)
			}

			if email != "" {
				processor := processor.Processor{Email: email, Proxy: proxy}
				err := processor.Process()
				if err != nil {
					return cli.Exit(fmt.Sprintf("Error processing email: %v", err), 1)
				}
			}

			if emailsFile != "" {
				emails, err := utils.LoadEmailsFromFile(emailsFile)
				if err != nil {
					return cli.Exit(fmt.Sprintf("Error loading emails: %v", err), 1)
				}
				proxies, err := utils.LoadProxiesFromFile(proxiesFile)
				if err != nil {
					return cli.Exit(fmt.Sprintf("Error loading proxies: %v", err), 1)
				}

				processorMany := processor.ProcessorMany{Emails: emails, Proxies: proxies, ConcurrencyLimit: concurrency}
				results, err := processorMany.Process()
				if err != nil {
					return cli.Exit(fmt.Sprintf("Error processing emails: %v", err), 1)
				}
				if outputFile != "" {
					err := utils.OutputResultsToFile(outputFile, results)
					if err != nil {
						return cli.Exit(fmt.Sprintf("Error writing to output file: %v", err), 1)
					}
				} else {
					for email, result := range results {
						fmt.Printf("Results for %s: %v\n", email, result)
					}
				}
			}

			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
