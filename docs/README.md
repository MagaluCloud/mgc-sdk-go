# MGC Go SDK Documentation

This directory contains the documentation for the MGC Go SDK, built with Sphinx and designed to be hosted on ReadTheDocs.

## Structure

```
docs/
├── index.rst              # Main documentation page
├── installation.rst       # Installation guide
├── authentication.rst     # Authentication guide
├── configuration.rst      # Client configuration
├── services/              # Service-specific documentation
│   └── index.rst         # Services overview
├── examples.rst           # Code examples
├── error-handling.rst     # Error handling guide
├── contributing.rst       # Contributing guidelines
├── conf.py               # Sphinx configuration
├── Makefile              # Build commands
└── README.md             # This file
```

## Building the Documentation

### Prerequisites

- Python 3.7 or higher
- pip

### Installation

1. Install Sphinx and the ReadTheDocs theme:

```bash
make install
```

Or manually:

```bash
pip install sphinx sphinx-rtd-theme
```

### Building

Build the HTML documentation:

```bash
make html
```

The built documentation will be available in `_build/html/`.

### Local Development

To serve the documentation locally:

```bash
make serve
```

This will start a local server at http://localhost:8000.

### Other Commands

- `make clean` - Clean build directory
- `make check` - Check for broken links
- `make spelling` - Check spelling
- `make help` - Show all available commands

## ReadTheDocs Integration

This documentation is configured to work with ReadTheDocs. The key configuration is in `conf.py`:

- Uses the `sphinx_rtd_theme`
- Configured for GitHub integration
- Includes autodoc extensions for Go code documentation

## Adding New Documentation

1. Create a new `.rst` file in the appropriate directory
2. Add it to the table of contents in `index.rst`
3. Follow the existing documentation style
4. Build and test locally before committing

## Documentation Style Guide

- Use clear, concise language
- Include code examples for all major features
- Follow the existing structure and formatting
- Use proper RST syntax
- Include cross-references where appropriate

## Contributing

When contributing to the documentation:

1. Follow the existing style and structure
2. Test your changes locally
3. Ensure all links work correctly
4. Update the table of contents if needed
5. Submit a pull request with a clear description

## Troubleshooting

### Common Issues

1. **Import errors**: Make sure the Go code is in the Python path
2. **Theme not found**: Install `sphinx_rtd_theme`
3. **Build errors**: Check RST syntax and ensure all files exist

### Getting Help

- Check the Sphinx documentation
- Look at existing `.rst` files for examples
- Open an issue on GitHub for documentation-specific problems 