# Especificação do Conector Google Play Store

**Camada**: L3
**Identificador**: `google-play-store-apps`
**Versão**: `playstore-web-search-1.0.0`
**Responsabilidade**:
- Realizar descoberta de aplicativos Android no catálogo público da Google Play Store por geolocalização (`gl`) e idioma (`hl`).
- Extrair `packageId`, título, autor, score e URL canônica.
- Garantir orçamento observável de bytes (máx 2MB) e sanitização pura de strings (máx 200 caracteres UTF-8).
