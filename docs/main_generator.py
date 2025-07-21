#!/usr/bin/env python3
"""
Documentation Generator for MGC SDK Go
Uses Sphinx with Markdown support to generate automatic documentation
"""

import os
import shutil
import subprocess
import sys
from pathlib import Path
import re
from typing import Dict

class DocumentationGenerator:
    def __init__(self, project_root: str = "../.."):
        self.project_root = Path(project_root).resolve() / ("v" + self.get_project_version())
        # check if project_root exists
        if not self.project_root.exists():
            self.project_root = Path(project_root).resolve() / "mgc-sdk-go"
        
        self.docs_dir = Path(__file__).parent.resolve()
        self.output_dir = self.docs_dir / "output"
        self.source_dir = self.docs_dir / "source"

        print("ls -la .")
        subprocess.run(["ls", "-la", "."])

        print(f"project_root: {self.project_root}")
        subprocess.run(["ls", "-la", self.project_root])

        print(f"docs_dir: {self.docs_dir}")
        subprocess.run(["ls", "-la", self.docs_dir])

        print(f"output_dir: {self.output_dir}")
        subprocess.run(["ls", "-la", self.output_dir])

        print(f"source_dir: {self.source_dir}")
        subprocess.run(["ls", "-la", self.source_dir])


        # Project configuration
        self.project_name = "MGC SDK Go"
        self.project_version = self.get_project_version()
        self.project_author = "Magalu Cloud"
        
        # Go modules to document
        self.go_modules = [
            "client", "compute", "blockstorage", "network", "kubernetes",
            "dbaas", "containerregistry", "sshkeys", "availabilityzones",
            "audit", "lbaas", "helpers"
        ]

    def get_project_version(self) -> str:
        """Captures the current project version from VERSION environment variable"""
        try:
            # Try to get version from environment variable
            version = os.environ.get("VERSION") 
            if not version:
                version = os.environ.get("READTHEDOCS_VERSION")
            
            if version and version.strip():
                version = version.strip()
                if version.startswith('v'):
                    version = version[1:]
                print(f"✅ Captured project version from VERSION env: {version}")
                return version
            else:
                print(f"⚠️  Error getting VERSION from environment, using default version")
                return "0.3.45"
                
        except Exception as e:
            print(f"⚠️  Error getting VERSION from environment: {e}, using default version")
            return "0.3.45"

    def clean_output_directory(self):
        """Completely cleans the output/ directory"""
        print("🧹 Cleaning output/ directory...")
        if self.output_dir.exists():
            shutil.rmtree(self.output_dir)
        self.output_dir.mkdir(exist_ok=True)
        print("✅ Output/ directory cleaned successfully")

    def create_sphinx_structure(self):
        """Creates the basic Sphinx structure"""
        print("📁 Creating Sphinx structure...")
        
        # Create source directory
        self.source_dir.mkdir(exist_ok=True)
        
        # Create directories for assets
        (self.source_dir / "_static").mkdir(exist_ok=True)
        (self.source_dir / "_templates").mkdir(exist_ok=True)

    
    def create_index_rst(self):
        """Creates the main index.rst file"""
        # Get list of available examples
        examples_dir = self.project_root / "cmd" / "examples"
        example_entries = []
        
        if examples_dir.exists():
            for example_dir in sorted(examples_dir.iterdir()):
                if example_dir.is_dir() and (example_dir / "main.go").exists():
                    example_name = example_dir.name
                    # Convert names like "availabilityzones" to "availability-zones"
                    display_name = example_name.replace('_', '-')
                    example_entries.append(f"   examples/{display_name}")
        
        examples_section = ""
        if example_entries:
            examples_section = f'''
.. toctree::
   :maxdepth: 1
   :caption: Examples:

{chr(10).join(example_entries)}
'''
        
        index_content = f'''# MGC Go SDK

Welcome to the MGC Go SDK documentation!

The MGC Go SDK provides a convenient way to interact with the Magalu Cloud API from Go applications.

.. toctree::
   :maxdepth: 2
   :caption: Content:

   introduction
   installation
   authentication
   regions
   usage
   error-handling
   contributing
   project-structure
{examples_section}
.. toctree::
   :maxdepth: 1
   :caption: Modules:

   modules/client
   modules/compute
   modules/blockstorage
   modules/network
   modules/kubernetes
   modules/dbaas
   modules/containerregistry
   modules/sshkeys
   modules/availabilityzones
   modules/audit
   modules/lbaas
   modules/helpers

Indices and tables
==================

* :ref:`genindex`
* :ref:`modindex`
* :ref:`search`
'''
        
        with open(self.source_dir / "index.rst", "w", encoding="utf-8") as f:
            f.write(index_content)
        print("✅ index.rst file created")

    def extract_readme_content(self) -> Dict[str, str]:
        """Extracts content from README.md and divides into sections"""
        readme_path = self.project_root / "README.md"
        if not readme_path.exists():
            print("⚠️  README.md not found")
            return {}
        
        with open(readme_path, "r", encoding="utf-8") as f:
            content = f.read()
        
        # Divide README into sections
        sections = {}
        lines = content.split('\n')
        current_section = "introduction"
        
        for line in lines:
            if line.startswith('## '):
                section_name = line[3:].strip().lower()
                section_name = re.sub(r'[^a-z0-9\s-]', '', section_name)
                section_name = re.sub(r'\s+', '-', section_name)
                current_section = section_name
                sections[current_section] = []
            elif line.startswith('# '):
                # Main title
                sections["introduction"] = [line]
                current_section = "introduction"
            else:
                if current_section not in sections:
                    sections[current_section] = []
                sections[current_section].append(line)
        
        # Convert lists to strings
        for section, lines in sections.items():
            sections[section] = '\n'.join(lines).strip()
        
        return sections

    def create_markdown_files(self, readme_sections: Dict[str, str]):
        """Creates Markdown files based on README content"""
        print("📝 Creating Markdown files...")
        
        # Section to file mapping
        section_files = {
            "introduction": "introduction.md",
            "installation": "installation.md",
            "authentication": "authentication.md",
            "regions": "regions.md",
            "global-services": "regions.md",  # Same file
            "project-structure": "project-structure.md",
            "usage-examples": "usage.md",
            "error-handling": "error-handling.md",
            "contributing": "contributing.md"
        }
        
        for section, content in readme_sections.items():
            if section in section_files and content.strip():
                filename = section_files[section]
                filepath = self.source_dir / filename
                
                # Add frontmatter for MyST and ensure proper H1 heading
                title = section.replace('-', ' ').title()
                
                # Process content to fix header levels
                lines = content.split('\n')
                processed_lines = []
                
                for line in lines:
                    if line.startswith('### '):
                        # Convert H3 to H2 to maintain hierarchy
                        processed_lines.append(line.replace('### ', '## ', 1))
                    else:
                        processed_lines.append(line)
                
                processed_content = '\n'.join(processed_lines)
                
                markdown_content = f"""---
title: {title}
---

# {title}

{processed_content}
"""
                
                with open(filepath, "w", encoding="utf-8") as f:
                    f.write(markdown_content)
                print(f"✅ Created {filename}")

    def create_modules_documentation(self):
        """Creates documentation for Go modules"""
        print("📚 Creating module documentation...")
        
        # Create modules directory
        modules_dir = self.source_dir / "modules"
        modules_dir.mkdir(exist_ok=True)
        
        # Create documentation for each module
        for module in self.go_modules:
            self.create_module_documentation(module, modules_dir)

    def create_module_documentation(self, module_name: str, modules_dir: Path):
        """Creates documentation for a specific module using go doc"""
        module_path = self.project_root / module_name
        
        if not module_path.exists():
            print(f"⚠️  Module {module_name} not found")
            return
        
        # Find Go files in the module (excluding test files)
        go_files = [f for f in module_path.glob("*.go") if not f.name.endswith("_test.go")]
        if not go_files:
            print(f"⚠️  No Go files found in {module_name}")
            return
        
        # Create module content with title
        module_content = f"""# {module_name.title()}

"""
        
        # Get package documentation using go doc
        package_doc = self.get_go_doc_for_package(module_name)
        if package_doc:
            module_content += f"{package_doc}\n\n"
        
        # Save module file
        module_file = modules_dir / f"{module_name}.md"
        with open(module_file, "w", encoding="utf-8") as f:
            f.write(module_content)
        
        print(f"✅ Documentation created for module {module_name} using go doc")

    def get_go_doc_for_package(self, module_name: str) -> str:
        """Gets package documentation using go doc"""
        try:
            # Change to project root directory
            original_cwd = os.getcwd()
            os.chdir(self.project_root)
            
            # Run go doc for the package
            result = subprocess.run([
                "go", "doc", "-all", module_name
            ], capture_output=True, text=True, timeout=30)
            
            os.chdir(original_cwd)
            
            if result.returncode == 0 and result.stdout.strip():
                return self.clean_go_doc_output(result.stdout)
            else:
                print(f"⚠️  No package documentation found for {module_name}")
                return ""
                
        except subprocess.TimeoutExpired:
            print(f"⚠️  Timeout getting documentation for {module_name}")
            return ""
        except Exception as e:
            print(f"⚠️  Error getting package documentation for {module_name}: {e}")
            return ""

    def clean_go_doc_output(self, doc_output: str) -> str:
        """Cleans and formats go doc output for Markdown"""
        if not doc_output:
            return ""
        
        lines = doc_output.split('\n')
        
        # Remove the first line
        if lines:
            lines = lines[1:]
        
        cleaned_lines = []
        in_code_block = False
        code_block_content = []
        constants_removed = False
        types_removed = False
        
        for line in lines:
            # Skip empty lines at the beginning
            if not cleaned_lines and not line.strip():
                continue
            
            # Remove "CONSTANTS" and "TYPES" strings (case-sensitive, only once each)
            stripped_line = line.strip()
            if not constants_removed and stripped_line == "CONSTANTS":
                constants_removed = True
                continue
            if not types_removed and stripped_line == "TYPES":
                types_removed = True
                continue
            
            # Check if this line starts a code block
            if (stripped_line.startswith('func ') or 
                stripped_line.startswith('type ') or 
                stripped_line.startswith('const ') or 
                stripped_line.startswith('var ') or
                stripped_line.startswith('package ') or
                stripped_line.startswith('import ')):
                
                # If we were in a code block, close it
                if in_code_block:
                    if code_block_content:
                        # Always use generic code block to avoid highlighting issues
                        cleaned_lines.append("```")
                        cleaned_lines.extend(code_block_content)
                        cleaned_lines.append("```")
                    in_code_block = False
                    code_block_content = []
                
                # Start new code block
                in_code_block = True
                code_block_content = [stripped_line]
                
            elif in_code_block:
                # Continue code block
                if line.strip():
                    code_block_content.append(line.strip())
                else:
                    # Empty line in code block
                    code_block_content.append("")
                    
            elif line.strip().startswith('//'):
                # Convert Go comments to Markdown
                comment = line.strip()[2:].strip()
                if comment:
                    # Escape any problematic characters in comments
                    comment = comment.replace("'", "\\'").replace('"', '\\"')
                    cleaned_lines.append(f"*{comment}*")
                    
            elif line.strip():
                # Regular text - escape problematic characters
                text = line.strip()
                text = text.replace("'", "\\'").replace('"', '\\"')
                cleaned_lines.append(text)
        
        # Close any remaining code block
        if in_code_block and code_block_content:
            # Always use generic code block to avoid highlighting issues
            cleaned_lines.append("```")
            cleaned_lines.extend(code_block_content)
            cleaned_lines.append("```")
        
        return '\n'.join(cleaned_lines).strip()

    def check_go_availability(self) -> bool:
        """Checks if Go is available and working"""
        try:
            result = subprocess.run(["go", "version"], capture_output=True, text=True, timeout=10)
            if result.returncode == 0:
                print(f"✅ Go found: {result.stdout.strip()}")
                return True
            else:
                print("❌ Go is not working properly")
                return False
        except FileNotFoundError:
            print("❌ Go not found. Please install Go to generate documentation.")
            return False
        except Exception as e:
            print(f"❌ Error checking Go: {e}")
            return False

    def create_examples_documentation(self):
        """Creates documentation for examples"""
        examples_dir = self.project_root / "cmd" / "examples"
        
        if not examples_dir.exists():
            print("⚠️  Examples directory not found")
            return
        
        # Create examples directory
        examples_docs_dir = self.source_dir / "examples"
        examples_docs_dir.mkdir(exist_ok=True)
        
        # Create individual example files
        for example_dir in sorted(examples_dir.iterdir()):
            if example_dir.is_dir():
                example_name = example_dir.name
                main_file = example_dir / "main.go"
                
                if main_file.exists():
                    # Convert names like "availabilityzones" to "availability-zones"
                    display_name = example_name.replace('_', '-')
                    
                    # Create example content
                    example_content = f"""# {example_name.title()}

Example usage of the `{example_name}` module.

"""
                    
                    # Try to extract comments from main.go
                    try:
                        with open(main_file, "r", encoding="utf-8") as f:
                            content = f.read()
                        
                        # Extract package comments
                        package_match = re.search(r'//\s*(.+?)(?:\n|$)', content)
                        if package_match:
                            example_content += f"{package_match.group(1).strip()}\n\n"
                    
                    except Exception:
                        pass
                    
                    example_content += f"**File:** `cmd/examples/{example_name}/main.go`\n\n"
                    example_content += "```go\n"
                    
                    # Include example code
                    try:
                        with open(main_file, "r", encoding="utf-8") as f:
                            lines = f.readlines()
                            example_content += ''.join(lines)
                    except Exception:
                        example_content += "// Error reading example file\n"
                    
                    example_content += "```\n"
                    
                    # Save individual example file
                    example_file = examples_docs_dir / f"{display_name}.md"
                    with open(example_file, "w", encoding="utf-8") as f:
                        f.write(example_content)
                    
                    print(f"✅ Created example documentation for {example_name}")
        
        print("✅ Examples documentation created")

    def create_project_structure_documentation(self):
        """Creates documentation for project structure"""
        structure_content = """# Project Structure

This section describes the organization of files and directories in the MGC Go SDK.

## Structure Overview

```
mgc-sdk-go/
├── client/         # Base client implementation and configuration
├── compute/        # Computing service API (instances, images, machine types)
├── blockstorage/   # Block storage service API
├── network/        # Network service API
├── kubernetes/     # Kubernetes service API
├── dbaas/          # Database as a Service API
├── containerregistry/ # Container Registry service API
├── sshkeys/        # SSH Keys service API
├── availabilityzones/ # Availability Zones service API
├── audit/          # Audit service API
├── lbaas/          # Load Balancer as a Service API
├── helpers/        # Utility functions
├── internal/       # Internal packages
└── cmd/            # Usage examples
```

## Module Descriptions

### client/
Contains the base HTTP client implementation and configurations for communicating with the Magalu Cloud API.

### compute/
Provides functionality to manage virtual instances, machine types, images, and snapshots.

### blockstorage/
Allows managing block storage volumes, snapshots, and volume types.

### network/
Provides functionality to manage VPCs, subnets, security groups, and other network resources.

### kubernetes/
Allows managing Kubernetes clusters, nodepools, and related configurations.

### dbaas/
Provides functionality to manage database instances, clusters, and configurations.

### containerregistry/
Allows managing container registries, repositories, and images.

### sshkeys/
Provides functionality to manage SSH keys.

### availabilityzones/
Allows querying available availability zones.

### audit/
Provides functionality to access audit logs and events.

### lbaas/
Allows managing load balancers and related configurations.

### helpers/
Contains reusable utility functions throughout the SDK.

### internal/
Contains internal packages not publicly exposed.

### cmd/
Contains practical examples of how to use each SDK module.
"""
        
        with open(self.source_dir / "project-structure.md", "w", encoding="utf-8") as f:
            f.write(structure_content)
        
        print("✅ Project structure documentation created")

    def create_requirements_txt(self):
        """Creates requirements.txt file for Python dependencies"""
        requirements = """sphinx>=7.0.0
sphinx-rtd-theme>=1.3.0
myst-parser>=2.0.0
sphinx-copybutton>=0.5.0
"""
        
        with open(self.docs_dir / "requirements.txt", "w", encoding="utf-8") as f:
            f.write(requirements)
        print("✅ requirements.txt file created")

    def create_makefile(self):
        """Creates Makefile to facilitate documentation generation"""
        makefile_content = """# Makefile for documentation generation

.PHONY: help clean html

help:
	@echo "Available targets:"
	@echo "  html       to make standalone HTML files"
	@echo "  clean      to clean output directory"

clean:
	rm -rf output/*

html:
	$(SPHINXBUILD) -b html $(ALLSPHINXOPTS) output/html
	@echo
	@echo "Build finished. The HTML pages are in output/html."

# Variables
SPHINXBUILD   = sphinx-build
SOURCEDIR     = source
BUILDDIR      = output

# Put it first so that "make" without argument is like "make help".
help:
	@$(SPHINXBUILD) -M help "$(SOURCEDIR)" "$(BUILDDIR)" $(SPHINXOPTS) $(O)

.PHONY: help Makefile

# Catch-all target: route all unknown targets to Sphinx using the new
# "make mode" option.  $(O) is meant as a shortcut for $(SPHINXOPTS).
%: Makefile
	@$(SPHINXBUILD) -M $@ "$(SOURCEDIR)" "$(BUILDDIR)" $(SPHINXOPTS) $(O)
"""
        
        with open(self.docs_dir / "Makefile", "w", encoding="utf-8") as f:
            f.write(makefile_content)
        print("✅ Makefile created")

    def install_dependencies(self):
        """Installs required Python dependencies"""
        print("📦 Installing Python dependencies...")
        try:
            subprocess.run([
                sys.executable, "-m", "pip", "install", "-r", 
                str(self.docs_dir / "requirements.txt")
            ], check=True, capture_output=True)
            print("✅ Dependencies installed successfully")
        except subprocess.CalledProcessError as e:
            print(f"❌ Error installing dependencies: {e}")
            print("Try installing manually: pip install -r requirements.txt")

    
    def run(self):
        """Executes the complete documentation generation process"""
        print("🚀 Starting documentation generation for MGC Go SDK")
        print("=" * 60)
        
        # 1. Check Go availability
        if not self.check_go_availability():
            print("⚠️  Continuing without Go documentation generation...")
        
        # 2. Clean output directory
        self.clean_output_directory()
        
        # 3. Create Sphinx structure
        self.create_sphinx_structure()
        
        # 4. Create configuration files
        self.create_requirements_txt()
        self.create_makefile()
        
        # 5. Extract README content
        readme_sections = self.extract_readme_content()
        
        # 6. Create Markdown files
        self.create_markdown_files(readme_sections)
        
        # 7. Create module documentation (with go doc if available)
        self.create_modules_documentation()
        
        # 8. Create examples documentation
        self.create_examples_documentation()
        
        # 9. Create project structure documentation
        self.create_project_structure_documentation()
        
        # 10. Create main index
        self.create_index_rst()
        
        # 11. Install dependencies
        self.install_dependencies()
                
        print("=" * 60)
        print("🎉 Documentation generation completed!")

def main():
    """Main function"""
    generator = DocumentationGenerator()
    generator.run()

if __name__ == "__main__":
    main()
