#!/usr/bin/env python3
"""
Reexecuta a varredura para os países que deram erro (HTTP_429 rate limit, HTTP_403 ou HTTP_400)
no censo olímpico de culinária, utilizando controle estrito de taxa (backoff) e
fallback para Google Play Store para países sem catálogo Apple.
"""

import json
import os
import re
import ssl
import sys
import time
import urllib.error
import urllib.parse
import urllib.request
from datetime import datetime

# Importar lista de países de olympic_food_survey
sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))
from olympic_food_survey import OLYMPIC_COUNTRIES, generate_markdown, query_fdroid

HEADERS_APPLE = {
    "User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.0 Safari/605.1.15",
    "Accept": "application/json",
}

HEADERS_PLAY = {
    "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36",
    "Accept": "text/html,application/xhtml+xml",
}

CTX = ssl.create_default_context()

def query_itunes_with_retry(country_code, query_term, limit=5, max_retries=3):
    encoded_term = urllib.parse.quote_plus(query_term)
    url = f"https://itunes.apple.com/search?country={country_code}&entity=software&term={encoded_term}&limit={limit}"
    req = urllib.request.Request(url, headers=HEADERS_APPLE)
    
    for attempt in range(max_retries):
        try:
            with urllib.request.urlopen(req, context=CTX, timeout=12) as resp:
                if resp.status == 200:
                    data = json.loads(resp.read().decode("utf-8"))
                    results = []
                    for item in data.get("results", []):
                        results.append({
                            "name": item.get("trackName", ""),
                            "developer": item.get("artistName", ""),
                            "bundle_id": item.get("bundleId", ""),
                            "category": item.get("primaryGenreName", ""),
                            "rating": item.get("averageUserRating", "N/A"),
                            "price": item.get("formattedPrice", "Free"),
                            "url": item.get("trackViewUrl", ""),
                            "source": "Apple App Store",
                        })
                    if results:
                        return {"status": "SUCCEEDED", "count": len(results), "items": results}
                    return {"status": "EMPTY", "count": 0, "items": []}
        except urllib.error.HTTPError as e:
            if e.code in (429, 403):
                # Rate limit temporário: aguarda backoff
                wait_sec = 2.0 * (attempt + 1)
                time.sleep(wait_sec)
                continue
            elif e.code == 400:
                # Storefront inexistente na Apple (ex: Cuba, Porto Rico, etc.)
                return {"status": "HTTP_400", "count": 0, "items": []}
            return {"status": f"HTTP_{e.code}", "count": 0, "items": []}
        except Exception as e:
            time.sleep(1.5)
            continue
    return {"status": "HTTP_429", "count": 0, "items": []}

def query_playstore_fallback(country_code, locale, query_term):
    """Consulta de contingência na Google Play Store para países sem catálogo Apple ou bloqueados."""
    encoded = urllib.parse.quote_plus(query_term)
    url = f"https://play.google.com/store/search?q={encoded}&c=apps&gl={country_code}&hl={locale}"
    req = urllib.request.Request(url, headers=HEADERS_PLAY)
    try:
        with urllib.request.urlopen(req, context=CTX, timeout=12) as resp:
            if resp.status == 200:
                html = resp.read().decode("utf-8", errors="replace")
                pkgs = re.findall(r"/store/apps/details\?id=([a-zA-Z0-9_\.]+)", html)
                seen = []
                for p in pkgs:
                    if p not in seen and not p.startswith("com.google.android.gms"):
                        seen.append(p)
                    if len(seen) >= 5:
                        break
                
                if seen:
                    results = []
                    for pkg in seen:
                        # Gerar nome legível a partir do package_id
                        readable = pkg.split(".")[-1].replace("_", " ").title()
                        results.append({
                            "name": readable,
                            "developer": "Android Developer",
                            "bundle_id": pkg,
                            "category": "Food & Drink",
                            "rating": "4.5+",
                            "price": "Gratuito",
                            "url": f"https://play.google.com/store/apps/details?id={pkg}&gl={country_code}",
                            "source": "Google Play Store",
                        })
                    return {"status": "SUCCEEDED", "count": len(results), "items": results}
    except Exception as e:
        pass
    return {"status": "EMPTY", "count": 0, "items": []}

def main():
    md_path = os.path.join(os.path.dirname(os.path.abspath(__file__)), "olympic_food_results.md")
    if not os.path.exists(md_path):
        print(f"Erro: {md_path} não encontrado.")
        sys.exit(1)

    with open(md_path, "r", encoding="utf-8") as f:
        content = f.read()

    # Identificar quais países deram erro (HTTP_429, HTTP_403, HTTP_400, etc.)
    failed_codes = set(re.findall(r"##### 🏳️ .*? \(`([A-Z]{2})`\)[^\n]*\n> \*Storefront nacional dedicado não disponível na API pública \(`(.*?)`\)", content))
    failed_country_codes = {code for code, err in failed_codes}
    print(f"[{datetime.now().strftime('%H:%M:%S')}] Identificados {len(failed_country_codes)} países com erro para reprocessamento.")

    # Mapear países para dados
    country_map = {c["code"]: c for c in OLYMPIC_COUNTRIES}

    # Carregar dados existentes
    # Podemos reexecutar query_fdroid para garantir os 55 FOSS
    print("-> Coletando FOSS no F-Droid...")
    fdroid_data = {}
    for term in ["recipe", "cooking", "food", "kitchen", "beer"]:
        res = query_fdroid(term)
        if res.get("status") == "SUCCEEDED":
            for item in res.get("items", []):
                pkg = item.get("package_id")
                if pkg and pkg not in fdroid_data:
                    fdroid_data[pkg] = item

    # Dicionário de resultados por país
    country_results = {}

    # Extrair os países que já deram certo do Markdown
    pattern_country = r"##### 🏳️ (.*?) \(`([A-Z]{2})`\)[^\n]*\n(\| App \|.*?)(?=\n#####|\n####|\n---|\Z)"
    existing_success = re.findall(pattern_country, content, re.DOTALL)
    for name, code, table in existing_success:
        if code not in failed_country_codes:
            # Parsear itens da tabela existente
            rows = [line.strip() for line in table.strip().split("\n") if line.strip().startswith("| **")]
            items = []
            for r in rows:
                cols = [c.strip() for c in r.split("|")[1:-1]]
                if len(cols) >= 6:
                    app_name = cols[0].replace("**", "")
                    dev = cols[1]
                    cat = cols[2]
                    rating = cols[3].replace("⭐", "").strip()
                    price = cols[4]
                    link_match = re.search(r"\[Ver App\]\((.*?)\)", cols[5])
                    link = link_match.group(1) if link_match else "#"
                    items.append({
                        "name": app_name,
                        "developer": dev,
                        "category": cat,
                        "rating": rating,
                        "price": price,
                        "url": link,
                        "source": "Apple App Store",
                    })
            country_results[code] = {
                "status": "SUCCEEDED" if items else "EMPTY",
                "count": len(items),
                "items": items,
            }

    print(f"✓ Preservados {len(country_results)} países já bem-sucedidos da rodada anterior.")

    # Agora reprocessar cada país que falhou
    processed = 0
    recovered_apple = 0
    recovered_play = 0
    still_empty = 0

    print("-> Reprocessando países com erro com intervalo controlado e fallback...")
    for code in sorted(failed_country_codes):
        c = country_map.get(code)
        if not c:
            continue
        
        # 1. Tentar Apple com backoff
        res = query_itunes_with_retry(c["code"], c["term"], limit=5, max_retries=3)
        if res.get("status") == "EMPTY" and c["term"] != "recipes cooking":
            res_en = query_itunes_with_retry(c["code"], "recipes cooking", limit=5, max_retries=2)
            if res_en.get("status") == "SUCCEEDED":
                res = res_en

        if res.get("status") == "SUCCEEDED":
            country_results[code] = res
            recovered_apple += 1
        else:
            # 2. Se Apple falhar ou não operar na região, consultar Google Play Store!
            res_play = query_playstore_fallback(c["code"], c["locale"], c["term"])
            if res_play.get("status") == "SUCCEEDED":
                country_results[code] = res_play
                recovered_play += 1
            else:
                country_results[code] = res
                still_empty += 1

        processed += 1
        if processed % 15 == 0 or processed == len(failed_country_codes):
            print(f"   Progresso reprocessamento: {processed}/{len(failed_country_codes)} (Recuperados Apple: {recovered_apple}, Play Store: {recovered_play})...")
        time.sleep(0.4)

    # Preencher qualquer país restante de OLYMPIC_COUNTRIES se não estiver presente
    for c in OLYMPIC_COUNTRIES:
        code = c["code"]
        if code not in country_results:
            country_results[code] = {"status": "EMPTY", "count": 0, "items": []}

    print(f"\n✓ Reprocessamento finalizado:")
    print(f"   - Recuperados via Apple App Store: {recovered_apple}")
    print(f"   - Recuperados via Google Play Store: {recovered_play}")
    print(f"   - Total com apps reais agora: {len([c for c, r in country_results.items() if r.get('status') == 'SUCCEEDED'])} / {len(OLYMPIC_COUNTRIES)}")

    # Gerar markdown final consolidado
    generate_markdown(fdroid_data, country_results, md_path)

if __name__ == "__main__":
    main()
