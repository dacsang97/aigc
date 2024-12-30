package prompt

const defaultSystemTemplate = `Generate a commit message following the Conventional Commits standard:

<type>[optional scope]: <description>

[optional body]

[optional footer]

Rules:
1. Type MUST be one of:
   - feat: new feature (correlates with MINOR version)
   - fix: bug fix (correlates with PATCH version)
   - docs: documentation changes
   - style: formatting, missing semi colons, etc
   - refactor: refactoring code
   - perf: performance improvements
   - test: adding tests
   - chore: maintenance tasks

2. Scope is optional and should describe the section of code (e.g., feat(parser))
3. Description must be concise and in imperative mood (e.g., 'change' not 'changed')
4. Body should explain the motivation for the change and contrast with previous behavior
5. Breaking changes MUST be indicated by BREAKING CHANGE: in footer
6. A ! MAY be added before the : for breaking changes (e.g., feat!: breaking change)

IMPORTANT: Always generate the commit message in English, regardless of the input language.
Do not include any explanation in your response, only return the commit message content.`

const defaultUserTemplate = `User provided this commit message hint (which may be in any language):
%s

Please consider this message when generating the commit message. 
Understand the meaning and translate the intent to English if needed, 
but ensure the output follows the Conventional Commits standard and is in English.`

const defaultChangeMessageTemplate = `Analyze these file changes and generate a commit message:
"""
%s
"""

Guidelines:
1. Use appropriate type based on the changes (feat for new features, fix for bugs, etc.)
2. Add relevant scope if the changes are focused on a specific component
3. Description must be under 100 characters
4. Include breaking changes in footer with BREAKING CHANGE: prefix if any
5. Add detailed body explaining motivation and changes if significant
6. Use issue/PR references in footer if relevant

Return only the commit message without any extra content or backticks.

`
