# Kinikan CLI

## Description

Kinikan CLI is a tool designed for developers leveraging Docker Compose in their cloud applications. It focuses on automating the creation of add-ons or plugins for various Platform as a Service (PaaS) providers, such as Heroku and Railway. By analyzing your Docker Compose file, Kinikan identifies necessary services and seamlessly sets up corresponding add-ons, simplifying the initial setup process. This allows developers to concentrate more on building their applications, knowing that the foundational cloud services are efficiently handled by Kinikan.

## Features

- Automated creation of add-ons/plugins based on Docker Compose.
- Support for multiple PaaS providers.
- Interactive platform selection.
- Simplified API key management.

## To-Do

- [ ] Support for more PaaS providers.
- [ ] Automate deployment process.

## Installation

From the Releases tab on GitHub, download the appropriate release archive for your OS and architecture. Extract the archive and move the `kinikan` binary to a directory included in your system's PATH. This makes the `kinikan` command available from any terminal.

## Usage

Before using Kinikan, you must set the appropriate environment variable depending on the PaaS provider you are using:

- For Heroku, set `HEROKU_API_KEY`.
- For Railway, set `RAILWAY_API_KEY`.

These can be set in your shell as follows:

```bash
# For Heroku
export HEROKU_API_KEY='your_heroku_api_key'

# For Railway
export RAILWAY_API_KEY='your_railway_token'
```

## Running Kinikan CLI

The basic command structure is:

```bash
# To run Kinikan
kinikan run [--platform] [--filePath]

# Flags:
# --platform: Specify the PaaS provider (e.g., heroku, fly). If not specified, Kinikan will prompt you to choose.
# --filePath: Path to your Docker Compose file. If not specified, Kinikan defaults to the Docker Compose file in the root directory.
```

For example:

```bash
# Run Kinikan specifying the platform and file path
kinikan run --platform heroku --filePath /path/to/docker-compose.yml

# Run Kinikan interactively without specifying flags
kinikan run
```

Remember, `--platform` and `--filePath` are optional. If not specified, Kinikan will interactively prompt for the necessary information or use sensible defaults.
