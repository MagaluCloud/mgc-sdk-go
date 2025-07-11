# Scripts de Documentação

Este diretório contém scripts para gerar documentação automaticamente baseada no código Go usando `go doc`.

## Scripts Disponíveis

### 1. `generate-docs.sh` (Bash)

Script em Bash para gerar documentação RST baseada nos comentários Go.

**Uso:**
```bash
# Tornar executável
chmod +x scripts/generate-docs.sh

# Executar
./scripts/generate-docs.sh
```

**Funcionalidades:**
- Extrai documentação de todos os pacotes do SDK
- Gera arquivos RST para cada pacote
- Cria índice da API
- Atualiza o índice principal da documentação
- Gera HTML usando Sphinx (se disponível)

### 2. `generate_docs.py` (Python)

Script em Python mais robusto para geração de documentação.

**Uso:**
```bash
# Gerar apenas RST
python3 scripts/generate_docs.py .

# Gerar RST e HTML
python3 scripts/generate_docs.py . --html
```

**Funcionalidades:**
- Parse mais robusto da saída do `go doc`
- Melhor tratamento de erros
- Suporte a timeouts
- Geração opcional de HTML

## Uso via Makefile

O projeto inclui comandos no Makefile para facilitar a geração de documentação:

```bash
# Gerar documentação usando script Bash
make docs

# Gerar documentação usando script Python
make docs-python

# Apenas construir HTML (após gerar RST)
make docs-html
```

## Workflow Automatizado

O projeto inclui um workflow do GitHub Actions (`.github/workflows/docs.yml`) que:

1. **Gera documentação automaticamente** quando há mudanças no código Go
2. **Constrói HTML** usando Sphinx
3. **Publica no GitHub Pages** quando há push para `main`
4. **Prepara para ReadTheDocs** quando há releases

## Estrutura Gerada

A documentação é gerada no diretório `docs/api/` com a seguinte estrutura:

```
docs/
├── api/
│   ├── index.rst          # Índice da API
│   ├── Client.rst         # Documentação do cliente
│   ├── Audit.rst          # Documentação do módulo audit
│   ├── BlockStorage.rst   # Documentação do módulo blockstorage
│   ├── Compute.rst        # Documentação do módulo compute
│   ├── ContainerRegistry.rst
│   ├── DBaaS.rst
│   ├── Kubernetes.rst
│   ├── LoadBalancer.rst
│   ├── Network.rst
│   ├── SSHKeys.rst
│   └── Helpers.rst
├── _build/                # Arquivos HTML gerados
└── index.rst              # Índice principal (atualizado)
```

## Requisitos

### Para os Scripts
- Go 1.21+
- Bash (para o script shell)
- Python 3.8+ (para o script Python)

### Para HTML
- Sphinx
- sphinx-rtd-theme

**Instalação:**
```bash
pip install sphinx sphinx-rtd-theme
```

## Personalização

### Adicionar Novos Pacotes

Para adicionar um novo pacote à documentação automática:

1. **No script Bash** (`generate-docs.sh`):
   ```bash
   packages=(
       # ... pacotes existentes ...
       "github.com/magalucloud/mgc-sdk-go/novo-pacote:NovoPacote"
   )
   ```

2. **No script Python** (`generate_docs.py`):
   ```python
   self.packages = [
       # ... pacotes existentes ...
       ("github.com/magalucloud/mgc-sdk-go/novo-pacote", "NovoPacote"),
   ]
   ```

### Modificar Formato de Saída

Os scripts geram arquivos RST por padrão. Para modificar o formato:

1. **RST**: Formato padrão, compatível com Sphinx e ReadTheDocs
2. **Markdown**: Modificar os scripts para gerar `.md` em vez de `.rst`
3. **HTML direto**: Usar `go doc -html` e processar a saída

## Troubleshooting

### Erro: "go doc: no such file or directory"
- Verifique se o Go está instalado e no PATH
- Execute `go version` para confirmar

### Erro: "Sphinx não encontrado"
- Instale o Sphinx: `pip install sphinx sphinx-rtd-theme`
- Ou use apenas a geração RST (sem HTML)

### Documentação vazia
- Verifique se os comentários Go seguem o padrão `// Package`, `// Function`, etc.
- Execute `go doc <pacote>` manualmente para testar

### Timeout na geração
- O script Python tem timeout de 30s por pacote
- Aumente o timeout no código se necessário

## Contribuição

Para contribuir com melhorias nos scripts:

1. Mantenha compatibilidade com Go 1.21+
2. Adicione testes para novas funcionalidades
3. Documente mudanças no README
4. Teste com diferentes versões do Go

## Licença

Os scripts seguem a mesma licença do projeto principal. 