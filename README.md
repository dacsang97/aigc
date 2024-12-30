# AIGC - AI-powered Git Commit Message Generator

AIGC is a command-line tool that uses AI to automatically generate meaningful Git commit messages based on your code changes.

## Features

- 🤖 AI-powered commit message generation
- 🔑 Configurable API key settings
- 🎯 Support for different AI models
- 📝 Detailed logging system
- 🔄 Optional automatic push after commit

## Installation

If you have Go installed, you can install AIGC by running:

```bash
go install github.com/dacsang97/aigc
```

If you don't have Go installed, you can download the binary from [here](https://github.com/dacsang97/aigc/releases).

## Configuration

Before using AIGC, you need to configure your settings.
If you need an API key, you can get it from [here](https://openrouter.ai/settings/keys).

```bash
# Set your API key
aigc config --apikey YOUR_API_KEY

# Set AI model
aigc config --model "your-preferred-model"

# Enable debug mode
aigc config --debug true

# Add rules for better commit messages
aigc config --rules "This is a monorepo with packages: api, web, docs","Always include scope for package changes","Include performance impact if relevant"

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

The rules can be any text that helps the AI understand your project's context and requirements better. They will be used alongside the global rules configured with `aigc config --rules`.

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
aigc commit -m "Cho phép user input message, AI sẽ dựa vào đó để generate commit" --push
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
