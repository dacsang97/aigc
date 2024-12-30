# AIGC - AI-powered Git Commit Message Generator

AIGC is a command-line tool that uses AI to automatically generate meaningful Git commit messages based on your code changes.

## Features

- 🤖 AI-powered commit message generation
- 🔑 Support for multiple AI providers (OpenAI, OpenRouter, Groq, and custom providers)
- 🎯 Configurable models and API endpoints
- 📝 Detailed logging system
- 🔄 Optional automatic push after commit
- 🌍 Multi-language commit message hints

## Installation

If you have Go installed, you can install AIGC by running:

```bash
go install github.com/dacsang97/aigc@latest
```

If you don't have Go installed, you can download the binary from [here](https://github.com/dacsang97/aigc/releases).

## Configuration

Before using AIGC, you need to configure your settings.

### Provider Configuration

AIGC supports multiple AI providers:

```bash
# OpenAI
aigc config --provider openai --api-key YOUR_API_KEY --model gpt-4

# OpenRouter
aigc config --provider openrouter --api-key YOUR_API_KEY --model google/gemini-flash-1.5-8b

# Groq (OpenAI-compatible)
aigc config --provider custom --api-key YOUR_API_KEY --model llama3-8b-8192 --endpoint https://api.groq.com/openai/v1/chat/completions

# Other OpenAI-compatible providers
aigc config --provider custom --api-key YOUR_API_KEY --model your-model --endpoint https://your-api-endpoint/v1/chat/completions
```

You can get API keys from:

- OpenAI: https://platform.openai.com/api-keys
- OpenRouter: https://openrouter.ai/keys
- Groq: https://console.groq.com/keys

### Other Settings

```bash
# Enable debug mode
aigc config --debug true

# View current configuration
aigc config
```

### Project-Specific Rules

You can create a `.aigcrules` file in your project directory to provide additional context and rules for commit message generation. These rules will be automatically loaded when running `aigc` commands.

Example `.aigcrules`:

```yaml
- This is a monorepo project with packages: api, web, docs
- Always use package name as scope when changes are package-specific
- Include performance impact for any database-related changes
- Reference Jira ticket number in footer if available
```

## Usage

### Basic Commit

```bash
# Let AI generate commit message automatically
aigc commit

# Provide your own commit message hint (in any language)
aigc commit -m "Thêm tính năng đăng nhập qua Google"
aigc commit --message "修复登录页面的样式问题"
aigc commit -m "Add user authentication feature"

# The AI will combine your message with the code changes
# to generate a more accurate and contextual commit message
```

### Commit and Push

```bash
# Generate and push
aigc commit --push
aigc commit -p

# Generate with your message and push
aigc commit -m "Add new feature" --push
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
provider:
  provider: openrouter # openai, openrouter, or custom
  model: google/gemini-flash-1.5-8b
  api_key: your-api-key
  endpoint: "" # optional, for custom providers
debug: false
rules: ""
```

## Logs

Logs are stored in `~/.aigc/log/` directory with daily rotation.

## License

[MIT License](LICENSE)
