#!/usr/bin/env bash
# ============================================================
# DockSea — Instalador / Atualizador (Linux .deb)
# Uso:  curl -fsSL https://raw.githubusercontent.com/Josemarcio15/dockSea/dev/install.sh | bash
# ============================================================
set -euo pipefail

REPO="Josemarcio15/dockSea"
API="https://api.github.com/repos/${REPO}/releases"
TMP_DEB="/tmp/docksea_latest.deb"

# ── Cores ───────────────────────────────────────────────────
RESET="\033[0m"
BOLD="\033[1m"
GREEN="\033[32m"
CYAN="\033[36m"
YELLOW="\033[33m"
RED="\033[31m"

info()    { printf "${CYAN}${BOLD}[dockSea]${RESET} %s\n" "$*"; }
success() { printf "${GREEN}${BOLD}✓${RESET} %s\n" "$*"; }
warn()    { printf "${YELLOW}${BOLD}⚠${RESET}  %s\n" "$*"; }
error()   { printf "${RED}${BOLD}✗ Erro:${RESET} %s\n" "$*" >&2; exit 1; }

check_deps() {
    for cmd in curl dpkg; do
        command -v "$cmd" &>/dev/null || error "Dependência não encontrada: '$cmd'. Instale e tente novamente."
    done
}

fetch_release() {
    info "Buscando versão mais recente em github.com/${REPO} ..."

    local response
    response=$(curl -fsSL -H "User-Agent: docksea-installer" "$API") \
        || error "Falha ao consultar a API do GitHub. Verifique sua conexão."

    TAG=$(printf '%s' "$response" | grep -o '"tag_name": *"[^"]*"' | head -1 | cut -d'"' -f4)
    [[ -z "$TAG" ]] && error "Não foi possível determinar a versão mais recente."

    DEB_URL=$(printf '%s' "$response" \
        | grep -o '"browser_download_url": *"[^"]*\.deb"' \
        | head -1 \
        | cut -d'"' -f4)

    [[ -z "$DEB_URL" ]] && error "Nenhum pacote .deb encontrado no release '${TAG}'."

    success "Última versão: ${TAG}"
}

download_deb() {
    info "Baixando ${DEB_URL##*/} ..."
    curl -fL --progress-bar "$DEB_URL" -o "$TMP_DEB" \
        || error "Falha no download. Tente novamente."
    success "Download concluído."
}

elevate_cmd() {
    if [ "${EUID:-$(id -u)}" -eq 0 ]; then
        "$@"
    elif command -v sudo &>/dev/null; then
        sudo "$@"
    elif command -v pkexec &>/dev/null; then
        pkexec "$@"
    else
        error "Necessário permissões de root para instalar o pacote .deb."
    fi
}

install_deb() {
    info "Instalando pacote .deb..."
    if command -v apt-get &>/dev/null; then
        elevate_cmd apt-get install -y --allow-downgrades "$TMP_DEB" || {
            warn "Tentando corrigir dependências ausentes..."
            elevate_cmd apt-get install -f -y
        }
    else
        elevate_cmd dpkg -i "$TMP_DEB"
    fi

    rm -f "$TMP_DEB"
    success "DockSea ${TAG} instalado com sucesso! 🚀"
    info "Você já pode abrir o DockSea pelo menu de aplicativos ou executando 'DockSea'."
}

main() {
    check_deps
    fetch_release
    download_deb
    install_deb
}

main
