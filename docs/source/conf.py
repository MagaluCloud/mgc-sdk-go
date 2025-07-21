# Configuration file for the Sphinx documentation builder.
#
# For the full list of built-in configuration values, see the documentation:
# https://www.sphinx-doc.org/en/master/usage/configuration.html

# -- Project information -----------------------------------------------------
# https://www.sphinx-doc.org/en/master/usage/configuration.html#project-information

project = 'MGC SDK Go'
copyright = '2025, Magalu Cloud'
author = 'Magalu Cloud'
release = '0.3.45'

# -- General configuration ---------------------------------------------------
# https://www.sphinx-doc.org/en/master/usage/configuration.html#general-configuration

extensions = [
    'sphinx.ext.autodoc',
    'sphinx.ext.napoleon',
    'sphinx.ext.viewcode',
    'sphinx.ext.githubpages',
    'myst_parser',  # For Markdown support
    'sphinx_copybutton',
    'sphinx.ext.todo',
]

templates_path = ['_templates']
exclude_patterns = ['_build', 'Thumbs.db', '.DS_Store']

# -- Options for HTML output -------------------------------------------------
# https://www.sphinx-doc.org/en/master/usage/configuration.html#options-for-html-output

html_theme = 'sphinx_rtd_theme'
html_static_path = ['_static']
html_logo = '_static/logo.png'
html_favicon = '_static/favicon.ico'

# Theme configuration
html_theme_options = {
    'navigation_depth': 4,
    'titles_only': False,
    'collapse_navigation': False,
    'sticky_navigation': True,
    'includehidden': True,
    'prev_next_buttons_location': 'bottom',
    'style_external_links': True,
}

# MyST-Parser configuration
myst_enable_extensions = [
    "colon_fence",
    "deflist",
    "dollarmath",
    "html_image",
    "html_admonition",
    "replacements",
    "smartquotes",
    "substitution",
    "tasklist",
]

# Autodoc configuration
autodoc_default_options = {
    'members': True,
    'member-order': 'bysource',
    'special-members': '__init__',
    'undoc-members': True,
    'exclude-members': '__weakref__'
}

# Copybutton configuration
copybutton_prompt_text = ">>> |\.\.\. |\$ |In \[\d*\]: | (2, 5)\.\.\.: | (5, 8): "
copybutton_prompt_is_regexp = True

# Todo configuration
todo_include_todos = True

# Language configuration
language = 'en'

# Numbering configuration
numfig = True
numfig_format = {
    'figure': 'Figure %s',
    'table': 'Table %s',
    'code-block': 'Listing %s',
    'section': 'Section %s'
}

# Syntax highlighting configuration
highlight_language = 'go'
pygments_style = 'sphinx'
pygments_dark_style = 'monokai'

# Disable syntax highlighting for problematic code blocks
highlight_options = {
    'go': {
        'linenos': False,
        'hl_lines': [],
        'linenostart': 1,
    }
}

# Suppress warnings for syntax highlighting failures
suppress_warnings = ['misc.highlighting_failure']
