# Contributing to Labs

Thank you for your interest in improving the Labs repository! This guide explains how to contribute and maintain the educational materials.

## Philosophy

All contributions should align with our core principles:

- **Production-Ready First:** Every example must work in real environments, not toy implementations
- **ELI5 Methodology:** Complex concepts explained through analogies and simple language, never sacrificing depth
- **Extensive Documentation:** Inline comments explaining both "what" and "why"
- **Progressive Complexity:** Materials build from fundamentals to advanced patterns
- **Real-World Focus:** All labs and examples reflect production scenarios

## Types of Contributions

### 1. Adding New Courses
New course additions should follow the standard directory structure:

```
course-name/
├── README.md                    # Course overview, prerequisites, time estimates
├── QUICKSTART.md               # 30-minute quick introduction
├── 01-fundamentals/            # Basic concepts with ELI5 explanations
├── 02-core-skills/             # Core practical skills with working examples
├── 03-advanced-patterns/       # Production patterns and architectures
├── 04-labs/                    # Hands-on exercises with scenarios
├── 05-reference/               # Lookup materials, FAQs, troubleshooting
└── scripts/                    # Setup and utility scripts
```

Before starting a new course:
1. Verify it aligns with backend engineering, SRE, or DevOps domains
2. Check existing courses to avoid duplication
3. Consider prerequisites and how it integrates with other courses
4. Open an issue to discuss the new course topic

### 2. Improving Existing Courses

Improvements welcome for:
- **Clarity:** Simplifying explanations or fixing confusing sections
- **Accuracy:** Correcting outdated information or technical errors
- **Examples:** Adding working code or making examples more comprehensive
- **Coverage:** Expanding sections with missing topics
- **Organization:** Improving structure or flow

### 3. Adding Working Examples

Working examples are the heart of these materials. When adding examples:

**Requirements:**
- Complete, runnable code (not pseudo-code)
- Minimal external dependencies where possible
- Extensive inline comments explaining implementation
- README explaining setup and expected output
- `*.example` files for configuration that needs user input
- Scripts to automate setup where applicable

**Structure:**
```
example-name/
├── README.md              # What this example demonstrates
├── setup.sh               # Automation script (if needed)
├── main.<ext>             # Primary implementation
├── config.example.<ext>   # Configuration template
└── docs/                  # Additional documentation
```

### 4. Fixing Issues

Issues welcome for:
- **Bugs:** Examples that don't work or outdated information
- **Clarity:** Confusing explanations
- **Missing Information:** Gaps in course coverage
- **Better Patterns:** Architectural improvements

When reporting issues, include:
- Clear description of the problem
- Steps to reproduce (for bugs)
- Current behavior and expected behavior
- Affected course/section/example

## Submission Process

### For Small Changes (typos, clarifications, minor fixes)
1. Fork the repository
2. Create a branch: `git checkout -b fix/issue-description`
3. Make your changes
4. Submit a pull request with a clear description

### For Medium Changes (new examples, expanded lessons)
1. Open an issue first to discuss your plans
2. Wait for feedback from maintainers
3. Fork the repository and create a feature branch
4. Implement your changes following the structure
5. Test your work thoroughly
6. Submit a pull request

### For Large Changes (new courses, significant restructuring)
1. Open an issue with detailed proposal
2. Include: course overview, learning objectives, structure outline
3. Wait for approval before investing significant effort
4. Follow the submission process for medium changes

## Code Style & Documentation

### Lesson Files
- Use clear markdown formatting
- Include headers for major sections
- Use code blocks with language specification
- Include ELI5 explanations with analogies
- Add practical examples and warnings where appropriate

Example structure:
```markdown
# Lesson Title

## Overview
[Brief description of what you'll learn]

## The Concept Explained Simply
[ELI5 explanation with analogy]

## How It Works
[Technical explanation]

## Practical Example
[Working code example]

## Common Mistakes
[Gotchas and how to avoid them]

## Next Steps
[What comes next]
```

### Code Examples
- Extensive inline comments
- Clear variable names
- Error handling included
- Production patterns demonstrated
- README explaining the example

Example code style:
```go
// PackageDescription explains what this package does and why
package example

import (
    "fmt"
    // Clear comment explaining why this import is needed
    "github.com/some/package"
)

// FunctionName does X with detailed explanation
// It's used in production because of Y reason
// Example: see the README
func FunctionName(param string) (result string, err error) {
    // Real error handling, not just panic()
    if param == "" {
        return "", fmt.Errorf("param cannot be empty")
    }
    
    // Implementation with comments
    return result, nil
}
```

### Configuration Files
- Include example versions (e.g., `terraform.tfvars.example`)
- Document each setting with comments
- Explain why defaults exist
- Show both minimal and production configurations

## Testing & Validation

Before submitting:

1. **Read through your changes** - Does it make sense?
2. **Test the examples** - Do they actually run?
3. **Check the structure** - Does it follow patterns?
4. **Verify links** - Are all references correct?
5. **Run the setup scripts** - Do they work as documented?

For code examples:
```bash
cd example-directory
bash setup.sh           # If included
# Run the example
# Verify expected output
# Test error cases
```

## Writing Quality Examples

### What Makes a Good Example?

✅ **Good Examples:**
- Demonstrate a single, clear concept
- Include all necessary setup
- Work without external services (or clearly document requirements)
- Have production-grade error handling
- Include both happy path and error cases
- Are extensively commented
- Have a clear README explaining what/why

❌ **Avoid:**
- Toy implementations or "learning only" code
- Hardcoded values that should be configurable
- Skipping error handling
- Unexplained magic numbers or patterns
- Dependencies on external services without clear alternatives
- Overly complex examples that confuse rather than clarify

## Review Process

All submissions go through review:

1. **Initial Check:** Does this align with our philosophy and structure?
2. **Technical Review:** Is the content accurate and production-ready?
3. **Clarity Review:** Is it understandable to the target audience?
4. **Testing:** Do examples work as documented?
5. **Feedback:** Maintainers will suggest improvements

This process ensures quality and consistency across all materials.

## Questions?

Before contributing, you might find answers in:
- Existing courses (for structure and style)
- The main README (for philosophy)
- This CONTRIBUTING guide (for process)

Open an issue to ask questions before getting started!

## Recognition

All contributors will be recognized in:
- Course-specific acknowledgments
- Annual CHANGELOG updates
- GitHub contributor metrics

Thank you for helping make these materials better for everyone! 🙏

---

## Quick Reference

**Directory Structure Template:**
```bash
course-name/
├── README.md
├── QUICKSTART.md
├── 01-fundamentals/
│   ├── lesson-01-*.md
│   ├── lesson-02-*.md
│   └── examples/
│       └── example-1/
├── 02-core-skills/
│   ├── lesson-*.md
│   └── examples/
├── 03-advanced-patterns/
│   ├── lesson-*.md
│   └── examples/
├── 04-labs/
│   ├── lab-*.md
│   └── solutions/
├── 05-reference/
│   ├── common-mistakes.md
│   ├── troubleshooting.md
│   ├── faq.md
│   └── additional-resources.md
└── scripts/
    └── setup.sh
```

**File Naming:**
- Lessons: `lesson-##-kebab-case.md`
- Labs: `lab-##-kebab-case.md`
- Examples: `kebab-case/`
- Scripts: `kebab-case.sh`

**Always Include:**
- Extensive documentation
- Working examples
- Setup instructions
- Prerequisites
- Expected time commitment
