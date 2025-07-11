#!/bin/bash

# Script para gerar documentação automaticamente baseada no go doc
# Este script extrai a documentação dos comentários Go e gera arquivos markdown/rst

set -e

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Diretórios
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
DOCS_DIR="$PROJECT_ROOT/docs"
API_DOCS_DIR="$DOCS_DIR/api"

echo -e "${GREEN}Gerando documentação da API...${NC}"

# Criar diretório de documentação da API se não existir
mkdir -p "$API_DOCS_DIR"

# Função para gerar documentação de um pacote
generate_package_docs() {
    local package_path="$1"
    local package_name="$2"
    local output_file="$API_DOCS_DIR/${package_name}.rst"
    
    echo -e "${YELLOW}Gerando documentação para $package_name...${NC}"
    
    # Gerar documentação usando go doc
    {
        echo "$package_name"
        echo "=" * ${#package_name}
        echo ""
        
        # Extrair documentação do pacote
        package_doc=$(go doc "$package_path" 2>/dev/null | head -20)
        if [[ -n "$package_doc" ]]; then
            echo "Package Documentation"
            echo "-------------------"
            echo ""
            echo ".. code-block:: go"
            echo ""
            echo "$package_doc" | sed 's/^/   /'
            echo ""
            echo ""
        fi
        
        # Listar todas as funções, tipos e constantes
        # Funções
        functions=$(go doc -short "$package_path" 2>/dev/null | grep "^func" || true)
        if [[ -n "$functions" ]]; then
            echo "Functions"
            echo "---------"
            echo ""
            echo "$functions" | while read -r func; do
                if [[ -n "$func" ]]; then
                    func_name=$(echo "$func" | awk '{print $2}' | cut -d'(' -f1)
                    echo "- :func:\`$func_name\`"
                fi
            done
            echo ""
        fi
        
        # Tipos
        types=$(go doc -short "$package_path" 2>/dev/null | grep "^type" || true)
        if [[ -n "$types" ]]; then
            echo "Types"
            echo "-----"
            echo ""
            echo "$types" | while read -r type_def; do
                if [[ -n "$type_def" ]]; then
                    type_name=$(echo "$type_def" | awk '{print $2}')
                    echo "- :type:\`$type_name\`"
                fi
            done
            echo ""
        fi
        
        # Constantes
        constants=$(go doc -short "$package_path" 2>/dev/null | grep "^const" || true)
        if [[ -n "$constants" ]]; then
            echo "Constants"
            echo "---------"
            echo ""
            echo "$constants" | while read -r const; do
                if [[ -n "$const" ]]; then
                    const_name=$(echo "$const" | awk '{print $2}')
                    echo "- :const:\`$const_name\`"
                fi
            done
            echo ""
        fi
        
        # Exemplo de uso
        echo "Example Usage"
        echo "-------------"
        echo ""
        echo ".. code-block:: go"
        echo ""
        package_import="${package_path#./}"
        echo "   import \"github.com/magalucloud/mgc-sdk-go/$package_import\""
        echo ""
        echo "   // Use the $package_name package"
        echo "   // See the examples directory for complete examples"
        echo ""
        
    } > "$output_file"
    
    echo -e "${GREEN}✓ Documentação gerada: $output_file${NC}"
}

# Lista de pacotes para documentar (usando caminhos relativos)
packages=(
    "./client:Client"
    "./audit:Audit"
    "./availabilityzones:AvailabilityZones"
    "./blockstorage:BlockStorage"
    "./compute:Compute"
    "./containerregistry:ContainerRegistry"
    "./dbaas:DBaaS"
    "./kubernetes:Kubernetes"
    "./lbaas:LoadBalancer"
    "./network:Network"
    "./sshkeys:SSHKeys"
    "./helpers:Helpers"
)

# Gerar documentação para cada pacote
for package in "${packages[@]}"; do
    IFS=':' read -r package_path package_name <<< "$package"
    generate_package_docs "$package_path" "$package_name"
done

# Gerar índice da API
echo -e "${YELLOW}Gerando índice da API...${NC}"

cat > "$API_DOCS_DIR/index.rst" << 'EOF'
API Reference
============

This section contains the complete API reference for the MGC Go SDK.

.. toctree::
   :maxdepth: 2
   :caption: API Packages

   Client
   Audit
   AvailabilityZones
   BlockStorage
   Compute
   ContainerRegistry
   DBaaS
   Kubernetes
   LoadBalancer
   Network
   SSHKeys
   Helpers

EOF

echo -e "${GREEN}✓ Índice da API gerado: $API_DOCS_DIR/index.rst${NC}"

# Atualizar o índice principal para incluir a API
if grep -q "api/index" "$DOCS_DIR/index.rst"; then
    echo -e "${YELLOW}Índice principal já contém referência à API${NC}"
else
    # Adicionar API Reference ao índice principal
    sed -i '/Examples/a\   api/index' "$DOCS_DIR/index.rst"
    echo -e "${GREEN}✓ API Reference adicionada ao índice principal${NC}"
fi

# Gerar documentação HTML usando Sphinx (se disponível)
if command -v sphinx-build &> /dev/null; then
    echo -e "${YELLOW}Gerando documentação HTML...${NC}"
    cd "$DOCS_DIR"
    make html
    echo -e "${GREEN}✓ Documentação HTML gerada em $DOCS_DIR/_build/html${NC}"
else
    echo -e "${YELLOW}Sphinx não encontrado. Instale com: pip install sphinx sphinx-rtd-theme${NC}"
fi

echo -e "${GREEN}✓ Documentação gerada com sucesso!${NC}"
echo -e "${YELLOW}Arquivos gerados em: $API_DOCS_DIR${NC}" 