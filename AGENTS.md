# Agent Instructions

## Version Control

This project uses **Jujutsu (jj)** for version control instead of Git.

### Why Jujutsu?

Jujutsu provides:
- Atomic commits with automatic conflict resolution
- Intuitive branching and rebasing workflows
- Better ergonomics for complex history manipulation
- Simplified collaboration patterns

### Required Commands

Use **only** these jj commands for version control:

```bash
jj status      # Check current working copy status
jj diff        # View changes
jj commit      # Create atomic commit
```

Do not use `git` commands in this project.

### Commit Format

All commits must use **Conventional Commit** format:

```
<type>(<scope>): <subject>

<body>

<footer>
```

#### Types
- `feat`: New feature
- `fix`: Bug fix
- `refactor`: Code refactoring
- `test`: Tests or test improvements
- `docs`: Documentation
- `chore`: Build, dependencies, or tooling
- `perf`: Performance improvements
- `ci`: CI/CD configuration

#### Examples

```
feat(cli): add color output support

Implement ANSI color codes for terminal output.

Closes #123
```

```
fix(core): handle nil pointer in image generation

Previously the code would panic if title was empty.
Now it validates input before processing.
```

```
refactor: simplify hex color parsing

Extract hexToRGB into separate utility function.
```

### Test-Driven Development

Follow **TDD (Red-Green-Refactor)** for all changes:

1. **Red**: Write a failing test first
2. **Verify**: Run the test to confirm it fails
3. **Green**: Write the minimum implementation to pass the test
4. **Refactor**: Clean up code while keeping tests green

Do not skip the verification step - always confirm the test fails before implementing.

### Atomic Commits

- Make one logical change per commit
- Keep commits small and focused
- Each commit should be independently testable
- Avoid mixing refactoring with feature changes
- Use jj's automatic conflict resolution when needed

### Workflow

1. Make changes to files
2. `jj status` - verify changes
3. `jj diff` - review what changed
4. `jj commit -m "type(scope): description"` - create atomic commit
5. Continue working in the automatically created new working copy

### No Git

⚠️ **Important**: Do not use any git commands including:
- `git add`, `git commit`, `git push`
- `git status`, `git diff`
- `git branch`, `git checkout`
- `.git` directory manipulation

This project is managed exclusively with Jujutsu.
