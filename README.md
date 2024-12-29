# AIGC - AI-powered Git Commit Message Generator

AIGC is a command-line tool that uses AI to automatically generate meaningful Git commit messages based on your code changes.

## Features

- 🤖 AI-powered commit message generation
- 🔑 Configurable API key settings
- 🎯 Support for different AI models
- 📝 Detailed logging system
- 🔄 Optional automatic push after commit

## Installation

```bash
# Clone the repository
git clone https://github.com/dacsang97/aigc.git
cd aigc

# Install dependencies
go mod download

# Build the project
go build -o aigc
```

## Configuration

Before using AIGC, you need to configure your API key:

```bash
# Set your API key
aigc config YOUR_API_KEY
```

You can view your current configuration by running:

```bash
aigc config
```

## Usage

### Basic Commit

```bash
aigc commit
```

### Commit and Push

```bash
aigc commit --push
# or
aigc commit -p
```

### Debug Mode

```bash
aigc --debug commit
```

### Change AI Model

```bash
aigc --model "your-preferred-model" commit
```

## Configuration File

AIGC stores its configuration in `~/.aigc/config.yaml` with the following structure:

```yaml
api_key: your-api-key
model: google/gemini-flash-1.5-8b
debug: false
```

## Logs

Logs are stored in `~/.aigc/log/` directory with daily rotation.

## License

[MIT License](LICENSE)
