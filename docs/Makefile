# Minimal makefile for Sphinx documentation
#

# You can set these variables from the command line, and also
# from the environment for the first two.
SPHINXOPTS    ?=
SPHINXBUILD   ?= sphinx-build
SOURCEDIR     = .
BUILDDIR      = _build

# Put it first so that "make" without argument is like "make help".
help:
	@$(SPHINXBUILD) -M help "$(SOURCEDIR)" "$(BUILDDIR)" $(SPHINXOPTS) $(O)

.PHONY: help Makefile

# Catch-all target: route all unknown targets to Sphinx using the new
# "make mode" option.  $(O) is meant as a shortcut for $(SPHINXOPTS).
%: Makefile
	@$(SPHINXBUILD) -M $@ "$(SOURCEDIR)" "$(BUILDDIR)" $(SPHINXOPTS) $(O)

# Custom targets
clean:
	rm -rf $(BUILDDIR)/*

html: clean
	$(SPHINXBUILD) -b html $(SPHINXOPTS) $(SOURCEDIR) $(BUILDDIR)/html
	@echo
	@echo "Build finished. The HTML pages are in $(BUILDDIR)/html."

serve: html
	@echo "Starting local server at http://localhost:8000"
	@echo "Press Ctrl+C to stop"
	python -m http.server 8000 --directory $(BUILDDIR)/html

install:
	pip install sphinx sphinx-rtd-theme

check:
	@echo "Checking for broken links..."
	$(SPHINXBUILD) -b linkcheck $(SPHINXOPTS) $(SOURCEDIR) $(BUILDDIR)/linkcheck
	@echo
	@echo "Link check finished. Look for any errors in the above output " \
	      "or in $(BUILDDIR)/linkcheck/output.txt."

spelling:
	@echo "Checking spelling..."
	$(SPHINXBUILD) -b spelling $(SPHINXOPTS) $(SOURCEDIR) $(BUILDDIR)/spelling
	@echo
	@echo "Spell check finished. Look for any errors in the above output " \
	      "or in $(BUILDDIR)/spelling/output.txt." 