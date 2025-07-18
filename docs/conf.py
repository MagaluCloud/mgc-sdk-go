# Configuration file for the Sphinx documentation builder.
#
# This file only contains a selection of the most common options. For a full
# list see the documentation:
# https://www.sphinx-doc.org/en/master/usage/configuration.html

# -- Path setup --------------------------------------------------------------

# If extensions (or modules to document with autodoc) are in another directory,
# add these directories to sys.path here. If the directory is relative to the
# documentation root, use os.path.abspath to make it absolute, like shown here.
#
import os
import sys
import sphinx_rtd_theme

# Add the project root to the Python path
sys.path.insert(0, os.path.abspath('..'))

# -- Project information -----------------------------------------------------

project = 'MGC SDK Go'
copyright = '2025, Magalu Cloud'
author = 'Magalu Cloud'

version = '0.3'
release = '0.3.44'
