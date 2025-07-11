#!/usr/bin/env python3
"""
Script para gerar documentação automaticamente baseada no go doc
Este script extrai a documentação dos comentários Go e gera arquivos RST para o Sphinx
"""

import os
import sys
import subprocess
import re
from pathlib import Path
from typing import List, Dict, Optional, Tuple

class GoDocGenerator:
    def __init__(self, project_root: str):
        self.project_root = Path(project_root)
        self.docs_dir = self.project_root / "docs"
        self.api_docs_dir = self.docs_dir / "api"
        
        # Criar diretório se não existir
        self.api_docs_dir.mkdir(exist_ok=True)
        
        # Pacotes para documentar
        self.packages = [
            ("github.com/magalucloud/mgc-sdk-go/client", "Client"),
            ("github.com/magalucloud/mgc-sdk-go/audit", "Audit"),
            ("github.com/magalucloud/mgc-sdk-go/availabilityzones", "AvailabilityZones"),
            ("github.com/magalucloud/mgc-sdk-go/blockstorage", "BlockStorage"),
            ("github.com/magalucloud/mgc-sdk-go/compute", "Compute"),
            ("github.com/magalucloud/mgc-sdk-go/containerregistry", "ContainerRegistry"),
            ("github.com/magalucloud/mgc-sdk-go/dbaas", "DBaaS"),
            ("github.com/magalucloud/mgc-sdk-go/kubernetes", "Kubernetes"),
            ("github.com/magalucloud/mgc-sdk-go/lbaas", "LoadBalancer"),
            ("github.com/magalucloud/mgc-sdk-go/network", "Network"),
            ("github.com/magalucloud/mgc-sdk-go/sshkeys", "SSHKeys"),
            ("github.com/magalucloud/mgc-sdk-go/helpers", "Helpers"),
        ]

    def run_go_doc(self, package_path: str, args: List[str] = None) -> str:
        """Executa go doc e retorna a saída"""
        if args is None:
            args = []
        
        cmd = ["go", "doc"] + args + [package_path]
        try:
            result = subprocess.run(
                cmd, 
                capture_output=True, 
                text=True, 
                cwd=self.project_root,
                timeout=30
            )
            return result.stdout
        except subprocess.TimeoutExpired:
            print(f"⚠️  Timeout ao executar go doc para {package_path}")
            return ""
        except subprocess.CalledProcessError as e:
            print(f"⚠️  Erro ao executar go doc para {package_path}: {e}")
            return ""

    def parse_go_doc_output(self, output: str) -> Dict:
        """Parseia a saída do go doc"""
        lines = output.split('\n')
        
        # Encontrar seção do pacote
        package_doc = ""
        functions = []
        types = []
        constants = []
        
        in_package = False
        in_functions = False
        in_types = False
        in_constants = False
        
        for line in lines:
            line = line.strip()
            
            if line.startswith('package '):
                in_package = True
                in_functions = False
                in_types = False
                in_constants = False
                continue
                
            if line == 'FUNCTIONS':
                in_package = False
                in_functions = True
                in_types = False
                in_constants = False
                continue
                
            if line == 'TYPES':
                in_package = False
                in_functions = False
                in_types = True
                in_constants = False
                continue
                
            if line == 'CONSTANTS':
                in_package = False
                in_functions = False
                in_types = False
                in_constants = True
                continue
                
            if line == 'VARIABLES':
                break
                
            if in_package and line:
                package_doc += line + '\n'
            elif in_functions and line.startswith('func '):
                functions.append(line)
            elif in_types and line.startswith('type '):
                types.append(line)
            elif in_constants and line.startswith('const '):
                constants.append(line)
        
        return {
            'package_doc': package_doc.strip(),
            'functions': functions,
            'types': types,
            'constants': constants
        }

    def generate_rst_content(self, package_path: str, package_name: str) -> str:
        """Gera conteúdo RST para um pacote"""
        print(f"📝 Gerando documentação para {package_name}...")
        
        # Obter documentação completa
        full_doc = self.run_go_doc(package_path, ["-all"])
        parsed = self.parse_go_doc_output(full_doc)
        
        # Obter documentação resumida para listagem
        short_doc = self.run_go_doc(package_path, ["-short"])
        short_parsed = self.parse_go_doc_output(short_doc)
        
        # Gerar RST
        rst_content = []
        
        # Título
        rst_content.append(f"{package_name}")
        rst_content.append("=" * len(package_name))
        rst_content.append("")
        
        # Documentação do pacote
        if parsed['package_doc']:
            rst_content.append("Package Documentation")
            rst_content.append("-------------------")
            rst_content.append("")
            rst_content.append(".. code-block:: go")
            rst_content.append("")
            for line in parsed['package_doc'].split('\n'):
                if line.strip():
                    rst_content.append(f"   {line}")
            rst_content.append("")
            rst_content.append("")
        
        # Lista de funções
        if short_parsed['functions']:
            rst_content.append("Functions")
            rst_content.append("---------")
            rst_content.append("")
            for func in short_parsed['functions']:
                func_name = func.split('(')[0].split()[-1]
                rst_content.append(f"- :func:`{func_name}`")
            rst_content.append("")
        
        # Lista de tipos
        if short_parsed['types']:
            rst_content.append("Types")
            rst_content.append("-----")
            rst_content.append("")
            for type_def in short_parsed['types']:
                type_name = type_def.split()[1]
                rst_content.append(f"- :type:`{type_name}`")
            rst_content.append("")
        
        # Lista de constantes
        if short_parsed['constants']:
            rst_content.append("Constants")
            rst_content.append("---------")
            rst_content.append("")
            for const in short_parsed['constants']:
                const_name = const.split()[1]
                rst_content.append(f"- :const:`{const_name}`")
            rst_content.append("")
        
        # Documentação detalhada
        rst_content.append("Detailed Documentation")
        rst_content.append("---------------------")
        rst_content.append("")
        
        # Funções detalhadas
        for func in parsed['functions']:
            func_name = func.split('(')[0].split()[-1]
            rst_content.append(f".. function:: {func}")
            rst_content.append("")
            
            # Obter documentação específica da função
            func_doc = self.run_go_doc(f"{package_path}.{func_name}")
            if func_doc.strip():
                rst_content.append("   " + func_doc.strip().replace('\n', '\n   '))
            rst_content.append("")
        
        # Tipos detalhados
        for type_def in parsed['types']:
            type_name = type_def.split()[1]
            rst_content.append(f".. type:: {type_def}")
            rst_content.append("")
            
            # Obter documentação específica do tipo
            type_doc = self.run_go_doc(f"{package_path}.{type_name}")
            if type_doc.strip():
                rst_content.append("   " + type_doc.strip().replace('\n', '\n   '))
            rst_content.append("")
        
        return '\n'.join(rst_content)

    def generate_api_index(self):
        """Gera o índice da API"""
        print("📋 Gerando índice da API...")
        
        index_content = [
            "API Reference",
            "============",
            "",
            "This section contains the complete API reference for the MGC Go SDK.",
            "",
            ".. toctree::",
            "   :maxdepth: 2",
            "   :caption: API Packages",
            ""
        ]
        
        for _, package_name in self.packages:
            index_content.append(f"   {package_name}")
        
        index_content.append("")
        
        index_file = self.api_docs_dir / "index.rst"
        index_file.write_text('\n'.join(index_content))
        print(f"✅ Índice da API gerado: {index_file}")

    def update_main_index(self):
        """Atualiza o índice principal para incluir a API"""
        main_index = self.docs_dir / "index.rst"
        
        if not main_index.exists():
            print("⚠️  Índice principal não encontrado")
            return
        
        content = main_index.read_text()
        
        if "api/index" not in content:
            # Adicionar API Reference após Examples
            lines = content.split('\n')
            new_lines = []
            added = False
            
            for line in lines:
                new_lines.append(line)
                if line.strip() == "Examples" and not added:
                    new_lines.append("   api/index")
                    added = True
            
            main_index.write_text('\n'.join(new_lines))
            print("✅ API Reference adicionada ao índice principal")
        else:
            print("ℹ️  API Reference já está no índice principal")

    def generate_html(self):
        """Gera documentação HTML usando Sphinx"""
        try:
            print("🌐 Gerando documentação HTML...")
            subprocess.run(
                ["make", "html"], 
                cwd=self.docs_dir, 
                check=True,
                capture_output=True
            )
            print("✅ Documentação HTML gerada em docs/_build/html")
        except subprocess.CalledProcessError:
            print("⚠️  Erro ao gerar HTML. Verifique se o Sphinx está instalado.")
        except FileNotFoundError:
            print("⚠️  Sphinx não encontrado. Instale com: pip install sphinx sphinx-rtd-theme")

    def generate_all(self):
        """Gera toda a documentação"""
        print("🚀 Iniciando geração de documentação...")
        
        # Gerar documentação para cada pacote
        for package_path, package_name in self.packages:
            rst_content = self.generate_rst_content(package_path, package_name)
            rst_file = self.api_docs_dir / f"{package_name}.rst"
            rst_file.write_text(rst_content)
            print(f"✅ Documentação gerada: {rst_file}")
        
        # Gerar índice da API
        self.generate_api_index()
        
        # Atualizar índice principal
        self.update_main_index()
        
        # Gerar HTML (opcional)
        if "--html" in sys.argv:
            self.generate_html()
        
        print("🎉 Documentação gerada com sucesso!")
        print(f"📁 Arquivos gerados em: {self.api_docs_dir}")

def main():
    if len(sys.argv) < 2:
        print("Uso: python generate_docs.py <project_root> [--html]")
        sys.exit(1)
    
    project_root = sys.argv[1]
    generator = GoDocGenerator(project_root)
    generator.generate_all()

if __name__ == "__main__":
    main() 