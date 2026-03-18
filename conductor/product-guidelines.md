# Product Guidelines: termapp

## Branding and Voice
- **Identity**: Professional, efficient, and developer-centric.
- **Voice**: Technical but approachable. Use active voice and concise language.
- **Tone**: Helpful, direct, and authoritative on best practices for CLI design.

## User Experience (UX) Principles
- **Efficiency First**: Minimize keystrokes and reduce cognitive load for terminal users.
- **Predictability**: Commands should follow consistent naming and behavior patterns (e.g., `-h` or `help` for assistance).
- **Graceful Failure**: Error messages should be clear, actionable, and avoid technical jargon where possible.
- **Context-Awareness**: Leverage the stage-based architecture to provide relevant feedback and completions.

## Design Patterns for CLI
- **Progressive Disclosure**: Show only what's necessary at each stage; hide complex options behind flags or sub-commands.
- **Feedback Loops**: Provide immediate visual feedback for user actions (e.g., changes in the prompt or explicit confirmation messages).
- **Accessibility**: Support keyboard-only navigation (default) and consider color blindness in terminal output.

## Documentation and Prose Style
- **Code Examples**: Always provide runnable code snippets in documentation.
- **Terminology**: Use "Stage" for application states and "Stack" for navigation context consistently.
