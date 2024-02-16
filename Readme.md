# MailSleuth

MailSleuth is an extremely quick and efficient email OSINT (Open Source Intelligence) tool designed to check the presence of email addresses across various social media platforms and other web services. It supports single email checks, bulk processing from files, proxy usage for anonymity, and configurable concurrency for performance optimization.

## Features

- **Single Email Check**: Quickly verify the presence of an email address on supported platforms.
- **Bulk Email Processing**: Process multiple email addresses from a file, ideal for large-scale investigations.
- **Proxy Support**: Use HTTP proxies to anonymize requests, supporting both single proxy and proxy list from a file.
- **Concurrency Control**: Set the number of concurrent operations to balance speed and resource usage.
- **Output Customization**: Save the processing results to a JSON file for further analysis or reporting.
- **Service Listing**: Easily list all supported services for email checks.

# Installation

```sh
go install -v github.com/44za12/mailsleuth/cmd/mailsleuth@latest
```

## Usage

### Basic Command

```bash
mailsleuth --email "user@example.com"
```

OR

```bash
mailsleuth -e "user@example.com"
```

### Using Proxies

Single proxy:

```bash
mailsleuth --email "user@example.com" --proxy "http://user:pass@host:port"
```

OR 

```bash
mailsleuth -e "user@example.com" -p "http://user:pass@host:port"
```

Proxy list from a file:

```bash
mailsleuth --emails-file "emails.txt" --proxies-file "proxies.txt"
```

OR

```bash
mailsleuth -E "emails.txt" -P "proxies.txt"
```

### Bulk Processing

Process multiple emails from a file and save results to a CSV file:

```bash
mailsleuth --emails-file "emails.txt" --output "results.csv"
```

OR

```bash
mailsleuth -E "emails.txt" -o "results.csv"
```

### Concurrency Control

Set the concurrency limit (default is 10):

```bash
mailsleuth --emails-file "emails.txt" --concurrency 20
```

OR

```bash
mailsleuth -E "emails.txt" -c 20
```

### List Supported Services

```bash
mailsleuth --list
```

OR

```bash
mailsleuth -l
```

## Supported Services

- Instagram
- Spotify
- X (Twitter)
- Amazon
- Facebook
- Github
- Kommo
- Axonaut
- HubSpot
- Insightly
- Nimble
- Wordpress
- Voxmedia
- Gravatar
- AnyDo
- LastPass
- Zoho
- Outlook

## Contributing

Contributions are welcome! Please feel free to submit a pull request or open an issue for any bugs, features, or improvements.

## Inspiration

MailSleuth is highly inspired by the functionality and core concepts of [holehe](https://github.com/megadose/holehe) and [mosint](https://github.com/alpkeskin/mosint). These tools have paved the way in email OSINT by offering powerful features for uncovering email usage across various platforms. MailSleuth aims to build upon their foundation by addressing some of the gaps and introducing additional features that were missing, such as:

- Enhanced proxy support for improved anonymity and circumvention of rate-limiting issues.
- Advanced concurrency control to optimize the speed and efficiency of bulk email processing.
- A more extensive list of supported services, continuously updated to include new platforms as they become popular.
- An intuitive command-line interface that simplifies the process of conducting email investigations.

We acknowledge the groundbreaking work done by the developers of holehe and mosint and are grateful for the inspiration they've provided. MailSleuth seeks to complement these tools by offering a broader set of features and capabilities for the OSINT community.