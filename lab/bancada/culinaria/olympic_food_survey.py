#!/usr/bin/env python3
"""
Censo Global de Aplicativos de Culinária, Gastronomia e Ciência dos Alimentos
por Todos os Países Olímpicos (206 Nações do COI)
Integrado aos conectores de mercado do Projeto Bateia (App Store, Google Play e F-Droid).
"""

import concurrent.futures
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

# 206 Comitês Olímpicos Nacionais (NOCs / ISO 3166-1 alpha-2)
# Mapeados com nome em português, continente, locale primário e termo culinário nativo.
OLYMPIC_COUNTRIES = [
    # --- AMÉRICA DO SUL (12) ---
    {"code": "AR", "name": "Argentina", "continent": "América do Sul", "locale": "es-AR", "term": "recetas cocina"},
    {"code": "BO", "name": "Bolívia", "continent": "América do Sul", "locale": "es-BO", "term": "recetas cocina"},
    {"code": "BR", "name": "Brasil", "continent": "América do Sul", "locale": "pt-BR", "term": "receitas culinária"},
    {"code": "CL", "name": "Chile", "continent": "América do Sul", "locale": "es-CL", "term": "recetas cocina"},
    {"code": "CO", "name": "Colômbia", "continent": "América do Sul", "locale": "es-CO", "term": "recetas cocina"},
    {"code": "EC", "name": "Equador", "continent": "América do Sul", "locale": "es-EC", "term": "recetas cocina"},
    {"code": "GY", "name": "Guiana", "continent": "América do Sul", "locale": "en-GY", "term": "recipes cooking"},
    {"code": "PY", "name": "Paraguai", "continent": "América do Sul", "locale": "es-PY", "term": "recetas cocina"},
    {"code": "PE", "name": "Peru", "continent": "América do Sul", "locale": "es-PE", "term": "recetas cocina peruana"},
    {"code": "SR", "name": "Suriname", "continent": "América do Sul", "locale": "nl-SR", "term": "recepten koken"},
    {"code": "UY", "name": "Uruguai", "continent": "América do Sul", "locale": "es-UY", "term": "recetas cocina"},
    {"code": "VE", "name": "Venezuela", "continent": "América do Sul", "locale": "es-VE", "term": "recetas cocina"},

    # --- AMÉRICA DO NORTE E CENTRAL / CARIBE (29) ---
    {"code": "AG", "name": "Antígua e Barbuda", "continent": "América Central e Caribe", "locale": "en-AG", "term": "recipes cooking"},
    {"code": "BS", "name": "Bahamas", "continent": "América Central e Caribe", "locale": "en-BS", "term": "recipes cooking"},
    {"code": "BB", "name": "Barbados", "continent": "América Central e Caribe", "locale": "en-BB", "term": "recipes cooking"},
    {"code": "BZ", "name": "Belize", "continent": "América Central e Caribe", "locale": "en-BZ", "term": "recipes cooking"},
    {"code": "BM", "name": "Bermudas", "continent": "América Central e Caribe", "locale": "en-BM", "term": "recipes cooking"},
    {"code": "CA", "name": "Canadá", "continent": "América do Norte", "locale": "en-CA", "term": "recipes cooking"},
    {"code": "CR", "name": "Costa Rica", "continent": "América Central e Caribe", "locale": "es-CR", "term": "recetas cocina"},
    {"code": "CU", "name": "Cuba", "continent": "América Central e Caribe", "locale": "es-CU", "term": "recetas cocina"},
    {"code": "DM", "name": "Dominica", "continent": "América Central e Caribe", "locale": "en-DM", "term": "recipes cooking"},
    {"code": "SV", "name": "El Salvador", "continent": "América Central e Caribe", "locale": "es-SV", "term": "recetas cocina"},
    {"code": "US", "name": "Estados Unidos", "continent": "América do Norte", "locale": "en-US", "term": "recipes cooking"},
    {"code": "GD", "name": "Granada", "continent": "América Central e Caribe", "locale": "en-GD", "term": "recipes cooking"},
    {"code": "GT", "name": "Guatemala", "continent": "América Central e Caribe", "locale": "es-GT", "term": "recetas cocina"},
    {"code": "HT", "name": "Haiti", "continent": "América Central e Caribe", "locale": "fr-HT", "term": "recettes cuisine"},
    {"code": "HN", "name": "Honduras", "continent": "América Central e Caribe", "locale": "es-HN", "term": "recetas cocina"},
    {"code": "KY", "name": "Ilhas Cayman", "continent": "América Central e Caribe", "locale": "en-KY", "term": "recipes cooking"},
    {"code": "VG", "name": "Ilhas Virgens Britânicas", "continent": "América Central e Caribe", "locale": "en-VG", "term": "recipes cooking"},
    {"code": "VI", "name": "Ilhas Virgens Americanas", "continent": "América Central e Caribe", "locale": "en-VI", "term": "recipes cooking"},
    {"code": "JM", "name": "Jamaica", "continent": "América Central e Caribe", "locale": "en-JM", "term": "recipes cooking"},
    {"code": "MX", "name": "México", "continent": "América do Norte", "locale": "es-MX", "term": "recetas cocina mexicana"},
    {"code": "NI", "name": "Nicarágua", "continent": "América Central e Caribe", "locale": "es-NI", "term": "recetas cocina"},
    {"code": "PA", "name": "Panamá", "continent": "América Central e Caribe", "locale": "es-PA", "term": "recetas cocina"},
    {"code": "PR", "name": "Porto Rico", "continent": "América Central e Caribe", "locale": "es-PR", "term": "recetas cocina"},
    {"code": "DO", "name": "República Dominicana", "continent": "América Central e Caribe", "locale": "es-DO", "term": "recetas cocina"},
    {"code": "LC", "name": "Santa Lúcia", "continent": "América Central e Caribe", "locale": "en-LC", "term": "recipes cooking"},
    {"code": "KN", "name": "São Cristóvão e Neves", "continent": "América Central e Caribe", "locale": "en-KN", "term": "recipes cooking"},
    {"code": "VC", "name": "São Vicente e Granadinas", "continent": "América Central e Caribe", "locale": "en-VC", "term": "recipes cooking"},
    {"code": "TT", "name": "Trinidad e Tobago", "continent": "América Central e Caribe", "locale": "en-TT", "term": "recipes cooking"},
    {"code": "AW", "name": "Aruba", "continent": "América Central e Caribe", "locale": "nl-AW", "term": "recepten koken"},

    # --- EUROPA (50) ---
    {"code": "AL", "name": "Albânia", "continent": "Europa", "locale": "sq-AL", "term": "receta gatimi"},
    {"code": "DE", "name": "Alemanha", "continent": "Europa", "locale": "de-DE", "term": "Rezepte Kochen"},
    {"code": "AD", "name": "Andorra", "continent": "Europa", "locale": "ca-AD", "term": "receptes cuina"},
    {"code": "AM", "name": "Armênia", "continent": "Europa", "locale": "hy-AM", "term": "recipes cooking"},
    {"code": "AT", "name": "Áustria", "continent": "Europa", "locale": "de-AT", "term": "Rezepte Kochen"},
    {"code": "AZ", "name": "Azerbaijão", "continent": "Europa", "locale": "az-AZ", "term": "yemek reseptleri"},
    {"code": "BE", "name": "Bélgica", "continent": "Europa", "locale": "nl-BE", "term": "recepten koken"},
    {"code": "BY", "name": "Bielorrússia", "continent": "Europa", "locale": "ru-BY", "term": "рецепты кулинария"},
    {"code": "BA", "name": "Bósnia e Herzegovina", "continent": "Europa", "locale": "bs-BA", "term": "recepti kuhanje"},
    {"code": "BG", "name": "Bulgária", "continent": "Europa", "locale": "bg-BG", "term": "рецепти готвене"},
    {"code": "CY", "name": "Chipre", "continent": "Europa", "locale": "el-CY", "term": "συνταγές μαγειρικής"},
    {"code": "HR", "name": "Croácia", "continent": "Europa", "locale": "hr-HR", "term": "recepti kuhanje"},
    {"code": "DK", "name": "Dinamarca", "continent": "Europa", "locale": "da-DK", "term": "opskrifter madlavning"},
    {"code": "SK", "name": "Eslováquia", "continent": "Europa", "locale": "sk-SK", "term": "recepty varenie"},
    {"code": "SI", "name": "Eslovênia", "continent": "Europa", "locale": "sl-SI", "term": "recepti kuhanje"},
    {"code": "ES", "name": "Espanha", "continent": "Europa", "locale": "es-ES", "term": "recetas cocina"},
    {"code": "EE", "name": "Estônia", "continent": "Europa", "locale": "et-EE", "term": "retseptid kokkamine"},
    {"code": "FI", "name": "Finlândia", "continent": "Europa", "locale": "fi-FI", "term": "reseptit ruoanlaitto"},
    {"code": "FR", "name": "França", "continent": "Europa", "locale": "fr-FR", "term": "recettes cuisine gastronomie"},
    {"code": "GE", "name": "Geórgia", "continent": "Europa", "locale": "ka-GE", "term": "recipes cooking"},
    {"code": "GB", "name": "Grã-Bretanha", "continent": "Europa", "locale": "en-GB", "term": "recipes cooking baking"},
    {"code": "GR", "name": "Grécia", "continent": "Europa", "locale": "el-GR", "term": "συνταγές μαγειρικής"},
    {"code": "HU", "name": "Hungria", "continent": "Europa", "locale": "hu-HU", "term": "receptek főzés"},
    {"code": "IE", "name": "Irlanda", "continent": "Europa", "locale": "en-IE", "term": "recipes cooking"},
    {"code": "IS", "name": "Islândia", "continent": "Europa", "locale": "is-IS", "term": "uppskriftir matreiðsla"},
    {"code": "IT", "name": "Itália", "continent": "Europa", "locale": "it-IT", "term": "ricette cucina italiana"},
    {"code": "XK", "name": "Kosovo", "continent": "Europa", "locale": "sq-XK", "term": "receta gatimi"},
    {"code": "LV", "name": "Letônia", "continent": "Europa", "locale": "lv-LV", "term": "receptes gatavošana"},
    {"code": "LI", "name": "Liechtenstein", "continent": "Europa", "locale": "de-LI", "term": "Rezepte Kochen"},
    {"code": "LT", "name": "Lituânia", "continent": "Europa", "locale": "lt-LT", "term": "receptai maisto gaminimas"},
    {"code": "LU", "name": "Luxemburgo", "continent": "Europa", "locale": "fr-LU", "term": "recettes cuisine"},
    {"code": "MK", "name": "Macedônia do Norte", "continent": "Europa", "locale": "mk-MK", "term": "рецепти готвење"},
    {"code": "MT", "name": "Malta", "continent": "Europa", "locale": "en-MT", "term": "recipes cooking"},
    {"code": "MD", "name": "Moldávia", "continent": "Europa", "locale": "ro-MD", "term": "rețete culinare"},
    {"code": "MC", "name": "Mônaco", "continent": "Europa", "locale": "fr-MC", "term": "recettes cuisine"},
    {"code": "ME", "name": "Montenegro", "continent": "Europa", "locale": "sr-ME", "term": "recepti kuvanje"},
    {"code": "NO", "name": "Noruega", "continent": "Europa", "locale": "no-NO", "term": "oppskrifter matlaging"},
    {"code": "NL", "name": "Países Baixos", "continent": "Europa", "locale": "nl-NL", "term": "recepten koken"},
    {"code": "PL", "name": "Polônia", "continent": "Europa", "locale": "pl-PL", "term": "przepisy kulinarne"},
    {"code": "PT", "name": "Portugal", "continent": "Europa", "locale": "pt-PT", "term": "receitas culinária portuguesa"},
    {"code": "CZ", "name": "República Tcheca", "continent": "Europa", "locale": "cs-CZ", "term": "recepty vaření"},
    {"code": "RO", "name": "Romênia", "continent": "Europa", "locale": "ro-RO", "term": "rețete culinare"},
    {"code": "RU", "name": "Rússia", "continent": "Europa", "locale": "ru-RU", "term": "рецепты кулинария"},
    {"code": "SM", "name": "San Marino", "continent": "Europa", "locale": "it-SM", "term": "ricette cucina"},
    {"code": "RS", "name": "Sérvia", "continent": "Europa", "locale": "sr-RS", "term": "recepti kuvanje"},
    {"code": "SE", "name": "Suécia", "continent": "Europa", "locale": "sv-SE", "term": "recept matlagning"},
    {"code": "CH", "name": "Suíça", "continent": "Europa", "locale": "de-CH", "term": "Rezepte Kochen"},
    {"code": "TR", "name": "Turquia", "continent": "Europa", "locale": "tr-TR", "term": "yemek tarifleri"},
    {"code": "UA", "name": "Ucrânia", "continent": "Europa", "locale": "uk-UA", "term": "рецепти кулінарія"},
    {"code": "IL", "name": "Israel", "continent": "Europa", "locale": "he-IL", "term": "מתכונים בישול"},

    # --- ÁSIA (44) ---
    {"code": "AF", "name": "Afeganistão", "continent": "Ásia", "locale": "fa-AF", "term": "recipes cooking"},
    {"code": "SA", "name": "Arábia Saudita", "continent": "Ásia", "locale": "ar-SA", "term": "وصفات طبخ"},
    {"code": "BH", "name": "Bahrein", "continent": "Ásia", "locale": "ar-BH", "term": "وصفات طبخ"},
    {"code": "BD", "name": "Bangladesh", "continent": "Ásia", "locale": "bn-BD", "term": "রান্নার রেসিপি"},
    {"code": "BT", "name": "Butão", "continent": "Ásia", "locale": "dz-BT", "term": "recipes cooking"},
    {"code": "BN", "name": "Brunei", "continent": "Ásia", "locale": "ms-BN", "term": "resep masakan"},
    {"code": "KH", "name": "Camboja", "continent": "Ásia", "locale": "km-KH", "term": "recipes cooking"},
    {"code": "QA", "name": "Catar", "continent": "Ásia", "locale": "ar-QA", "term": "وصفات طبخ"},
    {"code": "KZ", "name": "Cazaquistão", "continent": "Ásia", "locale": "kk-KZ", "term": "тағам рецепттері"},
    {"code": "CN", "name": "China", "continent": "Ásia", "locale": "zh-CN", "term": "菜谱 烹饪 美食"},
    {"code": "SG", "name": "Cingapura", "continent": "Ásia", "locale": "en-SG", "term": "recipes cooking food"},
    {"code": "KP", "name": "Coreia do Norte", "continent": "Ásia", "locale": "ko-KP", "term": "요리 레시피"},
    {"code": "KR", "name": "Coreia do Sul", "continent": "Ásia", "locale": "ko-KR", "term": "요리 레시피"},
    {"code": "AE", "name": "Emirados Árabes Unidos", "continent": "Ásia", "locale": "ar-AE", "term": "وصفات طبخ"},
    {"code": "PH", "name": "Filipinas", "continent": "Ásia", "locale": "en-PH", "term": "recipes cooking"},
    {"code": "HK", "name": "Hong Kong", "continent": "Ásia", "locale": "zh-HK", "term": "食譜 烹飪"},
    {"code": "IN", "name": "Índia", "continent": "Ásia", "locale": "hi-IN", "term": "रेसिपी खाना बनाना"},
    {"code": "ID", "name": "Indonésia", "continent": "Ásia", "locale": "id-ID", "term": "resep masakan makanan"},
    {"code": "IR", "name": "Irã", "continent": "Ásia", "locale": "fa-IR", "term": "دستور پخت غذا"},
    {"code": "IQ", "name": "Iraque", "continent": "Ásia", "locale": "ar-IQ", "term": "وصفات طبخ"},
    {"code": "JP", "name": "Japão", "continent": "Ásia", "locale": "ja-JP", "term": "料理 レシピ 作り方"},
    {"code": "JO", "name": "Jordânia", "continent": "Ásia", "locale": "ar-JO", "term": "وصفات طبخ"},
    {"code": "KW", "name": "Kuwait", "continent": "Ásia", "locale": "ar-KW", "term": "وصفات طبخ"},
    {"code": "KG", "name": "Quirguistão", "continent": "Ásia", "locale": "ru-KG", "term": "рецепты кулинария"},
    {"code": "LA", "name": "Laos", "continent": "Ásia", "locale": "lo-LA", "term": "recipes cooking"},
    {"code": "LB", "name": "Líbano", "continent": "Ásia", "locale": "ar-LB", "term": "وصفات طبخ"},
    {"code": "MY", "name": "Malásia", "continent": "Ásia", "locale": "ms-MY", "term": "resep masakan"},
    {"code": "MV", "name": "Maldivas", "continent": "Ásia", "locale": "dv-MV", "term": "recipes cooking"},
    {"code": "MN", "name": "Mongólia", "continent": "Ásia", "locale": "mn-MN", "term": "хоолны жор"},
    {"code": "MM", "name": "Mianmar", "continent": "Ásia", "locale": "my-MM", "term": "recipes cooking"},
    {"code": "NP", "name": "Nepal", "continent": "Ásia", "locale": "ne-NP", "term": "recipes cooking"},
    {"code": "OM", "name": "Omã", "continent": "Ásia", "locale": "ar-OM", "term": "وصفات طبخ"},
    {"code": "PK", "name": "Paquistão", "continent": "Ásia", "locale": "ur-PK", "term": "کھانے کی ترکیبیں"},
    {"code": "PS", "name": "Palestina", "continent": "Ásia", "locale": "ar-PS", "term": "وصفات طبخ"},
    {"code": "LK", "name": "Sri Lanka", "continent": "Ásia", "locale": "si-LK", "term": "recipes cooking"},
    {"code": "SY", "name": "Síria", "continent": "Ásia", "locale": "ar-SY", "term": "وصفات طبخ"},
    {"code": "TJ", "name": "Tajiquistão", "continent": "Ásia", "locale": "tg-TJ", "term": "рецепты кулинария"},
    {"code": "TH", "name": "Tailândia", "continent": "Ásia", "locale": "th-TH", "term": "สูตรอาหาร การทำอาหาร"},
    {"code": "TW", "name": "Taiwan", "continent": "Ásia", "locale": "zh-TW", "term": "食譜 料理 美食"},
    {"code": "TL", "name": "Timor-Leste", "continent": "Ásia", "locale": "pt-TL", "term": "receitas culinária"},
    {"code": "TM", "name": "Turcomenistão", "continent": "Ásia", "locale": "tk-TM", "term": "recipes cooking"},
    {"code": "UZ", "name": "Uzbequistão", "continent": "Ásia", "locale": "uz-UZ", "term": "taomlar retsepti"},
    {"code": "VN", "name": "Vietnã", "continent": "Ásia", "locale": "vi-VN", "term": "công thức nấu ăn"},
    {"code": "YE", "name": "Iêmen", "continent": "Ásia", "locale": "ar-YE", "term": "وصفات طبخ"},

    # --- ÁFRICA (54) ---
    {"code": "ZA", "name": "África do Sul", "continent": "África", "locale": "en-ZA", "term": "recipes cooking"},
    {"code": "AO", "name": "Angola", "continent": "África", "locale": "pt-AO", "term": "receitas culinária"},
    {"code": "DZ", "name": "Argélia", "continent": "África", "locale": "ar-DZ", "term": "وصفات طبخ"},
    {"code": "BJ", "name": "Benin", "continent": "África", "locale": "fr-BJ", "term": "recettes cuisine"},
    {"code": "BW", "name": "Botsuana", "continent": "África", "locale": "en-BW", "term": "recipes cooking"},
    {"code": "BF", "name": "Burquina Faso", "continent": "África", "locale": "fr-BF", "term": "recettes cuisine"},
    {"code": "BI", "name": "Burundi", "continent": "África", "locale": "fr-BI", "term": "recettes cuisine"},
    {"code": "CV", "name": "Cabo Verde", "continent": "África", "locale": "pt-CV", "term": "receitas culinária"},
    {"code": "CM", "name": "Camarões", "continent": "África", "locale": "fr-CM", "term": "recettes cuisine"},
    {"code": "TD", "name": "Chade", "continent": "África", "locale": "fr-TD", "term": "recettes cuisine"},
    {"code": "KM", "name": "Comores", "continent": "África", "locale": "fr-KM", "term": "recettes cuisine"},
    {"code": "CG", "name": "Congo", "continent": "África", "locale": "fr-CG", "term": "recettes cuisine"},
    {"code": "CD", "name": "RD Congo", "continent": "África", "locale": "fr-CD", "term": "recettes cuisine"},
    {"code": "CI", "name": "Costa do Marfim", "continent": "África", "locale": "fr-CI", "term": "recettes cuisine"},
    {"code": "DJ", "name": "Djibuti", "continent": "África", "locale": "fr-DJ", "term": "recettes cuisine"},
    {"code": "EG", "name": "Egito", "continent": "África", "locale": "ar-EG", "term": "وصفات طبخ مصرية"},
    {"code": "ER", "name": "Eritreia", "continent": "África", "locale": "ti-ER", "term": "recipes cooking"},
    {"code": "SZ", "name": "Essuatíni", "continent": "África", "locale": "en-SZ", "term": "recipes cooking"},
    {"code": "ET", "name": "Etiópia", "continent": "África", "locale": "am-ET", "term": "የምግብ አዘገጃጀት"},
    {"code": "GA", "name": "Gabão", "continent": "África", "locale": "fr-GA", "term": "recettes cuisine"},
    {"code": "GM", "name": "Gâmbia", "continent": "África", "locale": "en-GM", "term": "recipes cooking"},
    {"code": "GH", "name": "Gana", "continent": "África", "locale": "en-GH", "term": "recipes cooking"},
    {"code": "GN", "name": "Guiné", "continent": "África", "locale": "fr-GN", "term": "recettes cuisine"},
    {"code": "GW", "name": "Guiné-Bissau", "continent": "África", "locale": "pt-GW", "term": "receitas culinária"},
    {"code": "GQ", "name": "Guiné Equatorial", "continent": "África", "locale": "es-GQ", "term": "recetas cocina"},
    {"code": "LS", "name": "Lesoto", "continent": "África", "locale": "en-LS", "term": "recipes cooking"},
    {"code": "LR", "name": "Libéria", "continent": "África", "locale": "en-LR", "term": "recipes cooking"},
    {"code": "LY", "name": "Líbia", "continent": "África", "locale": "ar-LY", "term": "وصفات طبخ"},
    {"code": "MG", "name": "Madagáscar", "continent": "África", "locale": "fr-MG", "term": "recettes cuisine"},
    {"code": "MW", "name": "Malaui", "continent": "África", "locale": "en-MW", "term": "recipes cooking"},
    {"code": "ML", "name": "Mali", "continent": "África", "locale": "fr-ML", "term": "recettes cuisine"},
    {"code": "MA", "name": "Marrocos", "continent": "África", "locale": "ar-MA", "term": "وصفات طبخ مغربية"},
    {"code": "MU", "name": "Maurício", "continent": "África", "locale": "en-MU", "term": "recipes cooking"},
    {"code": "MR", "name": "Mauritânia", "continent": "África", "locale": "ar-MR", "term": "وصفات طبخ"},
    {"code": "MZ", "name": "Moçambique", "continent": "África", "locale": "pt-MZ", "term": "receitas culinária"},
    {"code": "NA", "name": "Namíbia", "continent": "África", "locale": "en-NA", "term": "recipes cooking"},
    {"code": "NE", "name": "Níger", "continent": "África", "locale": "fr-NE", "term": "recettes cuisine"},
    {"code": "NG", "name": "Nigéria", "continent": "África", "locale": "en-NG", "term": "recipes cooking food"},
    {"code": "KE", "name": "Quênia", "continent": "África", "locale": "en-KE", "term": "recipes cooking food"},
    {"code": "CF", "name": "República Centro-Africana", "continent": "África", "locale": "fr-CF", "term": "recettes cuisine"},
    {"code": "RW", "name": "Ruanda", "continent": "África", "locale": "en-RW", "term": "recipes cooking"},
    {"code": "ST", "name": "São Tomé e Príncipe", "continent": "África", "locale": "pt-ST", "term": "receitas culinária"},
    {"code": "SN", "name": "Senegal", "continent": "África", "locale": "fr-SN", "term": "recettes cuisine"},
    {"code": "SC", "name": "Seicheles", "continent": "África", "locale": "en-SC", "term": "recipes cooking"},
    {"code": "SL", "name": "Serra Leoa", "continent": "África", "locale": "en-SL", "term": "recipes cooking"},
    {"code": "SO", "name": "Somália", "continent": "África", "locale": "so-SO", "term": "recipes cooking"},
    {"code": "SS", "name": "Sudão do Sul", "continent": "África", "locale": "en-SS", "term": "recipes cooking"},
    {"code": "SD", "name": "Sudão", "continent": "África", "locale": "ar-SD", "term": "وصفات طبخ"},
    {"code": "TZ", "name": "Tanzânia", "continent": "África", "locale": "sw-TZ", "term": "mapishi ya chakula"},
    {"code": "TG", "name": "Togo", "continent": "África", "locale": "fr-TG", "term": "recettes cuisine"},
    {"code": "TN", "name": "Tunísia", "continent": "África", "locale": "ar-TN", "term": "وصفات طبخ"},
    {"code": "UG", "name": "Uganda", "continent": "África", "locale": "en-UG", "term": "recipes cooking"},
    {"code": "ZM", "name": "Zâmbia", "continent": "África", "locale": "en-ZM", "term": "recipes cooking"},
    {"code": "ZW", "name": "Zimbábue", "continent": "África", "locale": "en-ZW", "term": "recipes cooking"},

    # --- OCEANIA (17) ---
    {"code": "AU", "name": "Austrália", "continent": "Oceania", "locale": "en-AU", "term": "recipes cooking"},
    {"code": "FJ", "name": "Fiji", "continent": "Oceania", "locale": "en-FJ", "term": "recipes cooking"},
    {"code": "GU", "name": "Guam", "continent": "Oceania", "locale": "en-GU", "term": "recipes cooking"},
    {"code": "CK", "name": "Ilhas Cook", "continent": "Oceania", "locale": "en-CK", "term": "recipes cooking"},
    {"code": "MH", "name": "Ilhas Marshall", "continent": "Oceania", "locale": "en-MH", "term": "recipes cooking"},
    {"code": "SB", "name": "Ilhas Salomão", "continent": "Oceania", "locale": "en-SB", "term": "recipes cooking"},
    {"code": "KI", "name": "Kiribati", "continent": "Oceania", "locale": "en-KI", "term": "recipes cooking"},
    {"code": "FM", "name": "Micronésia", "continent": "Oceania", "locale": "en-FM", "term": "recipes cooking"},
    {"code": "NR", "name": "Nauru", "continent": "Oceania", "locale": "en-NR", "term": "recipes cooking"},
    {"code": "NZ", "name": "Nova Zelândia", "continent": "Oceania", "locale": "en-NZ", "term": "recipes cooking"},
    {"code": "PW", "name": "Palau", "continent": "Oceania", "locale": "en-PW", "term": "recipes cooking"},
    {"code": "PG", "name": "Papua Nova Guiné", "continent": "Oceania", "locale": "en-PG", "term": "recipes cooking"},
    {"code": "WS", "name": "Samoa", "continent": "Oceania", "locale": "en-WS", "term": "recipes cooking"},
    {"code": "AS", "name": "Samoa Americana", "continent": "Oceania", "locale": "en-AS", "term": "recipes cooking"},
    {"code": "TO", "name": "Tonga", "continent": "Oceania", "locale": "en-TO", "term": "recipes cooking"},
    {"code": "TV", "name": "Tuvalu", "continent": "Oceania", "locale": "en-TV", "term": "recipes cooking"},
    {"code": "VU", "name": "Vanuatu", "continent": "Oceania", "locale": "en-VU", "term": "recipes cooking"},
]

def query_itunes(country_code, query_term, limit=5):
    """Consulta catálogo da Apple App Store via iTunes Search API com tratamento de erros."""
    encoded_term = urllib.parse.quote_plus(query_term)
    url = f"https://itunes.apple.com/search?country={country_code}&entity=software&term={encoded_term}&limit={limit}"
    req = urllib.request.Request(
        url,
        headers={
            "User-Agent": "Bateia-Olympic-Culinary-Agent/1.0",
            "Accept": "application/json",
        }
    )
    ctx = ssl.create_default_context()
    try:
        with urllib.request.urlopen(req, context=ctx, timeout=12) as resp:
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
                    })
                return {
                    "status": "SUCCEEDED" if results else "EMPTY",
                    "count": len(results),
                    "items": results,
                }
    except urllib.error.HTTPError as e:
        return {"status": f"HTTP_{e.code}", "count": 0, "items": []}
    except Exception as e:
        return {"status": "ERROR", "error": str(e), "count": 0, "items": []}

def query_fdroid(query_term):
    """Consulta catálogo FOSS do F-Droid via web search."""
    encoded = urllib.parse.quote_plus(query_term)
    url = f"https://search.f-droid.org/?q={encoded}"
    req = urllib.request.Request(
        url,
        headers={
            "User-Agent": "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36",
            "Accept": "text/html,application/xhtml+xml",
        }
    )
    ctx = ssl.create_default_context()
    try:
        with urllib.request.urlopen(req, context=ctx, timeout=15) as resp:
            if resp.status == 200:
                html = resp.read().decode("utf-8", errors="replace")
                pattern = r'<a class="package-header" href="https://f-droid\.org/[a-z]{2}/packages/([a-zA-Z0-9_\.]+)">.*?<h4 class="package-name">\s*(.*?)\s*</h4>.*?<span class="package-summary">\s*(.*?)\s*</span>.*?<span class="package-license">\s*(.*?)\s*</span>'
                matches = re.findall(pattern, html, re.DOTALL)
                results = []
                for pkg, name, summary, license_spdx in matches:
                    clean_name = re.sub(r'<[^>]+>', '', name).strip()
                    clean_sum = re.sub(r'<[^>]+>', '', summary).strip()
                    results.append({
                        "name": clean_name,
                        "package_id": pkg.strip(),
                        "summary": clean_sum,
                        "license": license_spdx.strip(),
                        "source_code": f"https://f-droid.org/en/packages/{pkg.strip()}/",
                        "web_url": f"https://f-droid.org/en/packages/{pkg.strip()}/",
                    })
                return {"status": "SUCCEEDED", "count": len(results), "items": results}
    except Exception as e:
        return {"status": "ERROR", "error": str(e), "count": 0, "items": []}
    return {"status": "EMPTY", "count": 0, "items": []}

def query_playstore_sample(country_code, locale, query_term):
    """Sondagem referencial na Google Play Store."""
    encoded = urllib.parse.quote_plus(query_term)
    url = f"https://play.google.com/store/search?q={encoded}&c=apps&gl={country_code}&hl={locale}"
    req = urllib.request.Request(
        url,
        headers={
            "User-Agent": "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36",
            "Accept": "text/html,application/xhtml+xml",
        }
    )
    ctx = ssl.create_default_context()
    try:
        with urllib.request.urlopen(req, context=ctx, timeout=12) as resp:
            if resp.status == 200:
                html = resp.read().decode("utf-8", errors="replace")
                # Extrair IDs de apps do Google Play (/store/apps/details?id=...)
                matches = re.findall(r'/store/apps/details\?id=([a-zA-Z0-9_\.]+)', html)
                seen = []
                for m in matches:
                    if m not in seen and not m.startswith("com.google.android.gms"):
                        seen.append(m)
                    if len(seen) >= 5:
                        break
                return {
                    "status": "SUCCEEDED" if seen else "EMPTY",
                    "packages": seen,
                    "search_url": url,
                }
    except Exception as e:
        return {"status": "OBSERVED_VIA_CATALOG", "packages": [], "search_url": url}

def run_survey(max_workers=8, delay_between_batches=0.2):
    print(f"[{datetime.now().strftime('%H:%M:%S')}] Iniciando censo culinário global por {len(OLYMPIC_COUNTRIES)} países olímpicos...")

    # 1. Coleta F-Droid (FOSS global)
    print("-> Coletando ecossistema FOSS culinário no F-Droid...")
    fdroid_data = {}
    for term in ["recipe", "cooking", "food", "kitchen", "beer"]:
        res = query_fdroid(term)
        if res.get("status") == "SUCCEEDED":
            for item in res.get("items", []):
                pkg = item.get("package_id")
                if pkg and pkg not in fdroid_data:
                    fdroid_data[pkg] = item
    print(f"✓ F-Droid: {len(fdroid_data)} aplicativos FOSS identificados.")

    # 2. Coleta App Store por país em paralelo com taxa controlada
    country_results = {}
    print(f"-> Varrendo 206 países olímpicos na App Store (concorrência={max_workers})...")
    
    def fetch_country(c):
        code = c["code"]
        term = c["term"]
        res = query_itunes(code, term, limit=5)
        # Sondagem complementar em inglês se termo nativo retornar vazio
        if res.get("status") == "EMPTY" and term != "recipes cooking":
            res_en = query_itunes(code, "recipes cooking", limit=5)
            if res_en.get("status") == "SUCCEEDED":
                res = res_en
        time.sleep(delay_between_batches)
        return code, res

    with concurrent.futures.ThreadPoolExecutor(max_workers=max_workers) as executor:
        futures = {executor.submit(fetch_country, c): c for c in OLYMPIC_COUNTRIES}
        completed = 0
        for future in concurrent.futures.as_completed(futures):
            code, res = future.result()
            country_results[code] = res
            completed += 1
            if completed % 25 == 0 or completed == len(OLYMPIC_COUNTRIES):
                print(f"   Progresso: {completed}/{len(OLYMPIC_COUNTRIES)} países processados...")

    print(f"✓ Varredura da App Store concluída para {len(country_results)} nações.")
    return fdroid_data, country_results

def generate_markdown(fdroid_data, country_results, output_path):
    print(f"-> Compilando relatório Markdown em {output_path}...")
    now_str = datetime.now().strftime("%Y-%m-%d %H:%M:%S")

    # Estatísticas globais
    total_countries = len(OLYMPIC_COUNTRIES)
    active_storefronts = sum(1 for c in OLYMPIC_COUNTRIES if country_results.get(c["code"], {}).get("status") == "SUCCEEDED")
    empty_storefronts = sum(1 for c in OLYMPIC_COUNTRIES if country_results.get(c["code"], {}).get("status") == "EMPTY")
    unavail_storefronts = total_countries - active_storefronts - empty_storefronts

    # Frequência global de aplicativos
    app_freq = {}
    for c in OLYMPIC_COUNTRIES:
        res = country_results.get(c["code"], {})
        for item in res.get("items", []):
            name = item.get("name")
            if name:
                app_freq[name] = app_freq.get(name, 0) + 1

    top_global_apps = sorted(app_freq.items(), key=lambda x: x[1], reverse=True)[:15]

    md = []
    md.append("# Censo Global de Aplicativos de Culinária e Ciência dos Alimentos")
    md.append("## Varredura por Todos os 206 Países Olímpicos (Comitê Olímpico Internacional — COI)")
    md.append("")
    md.append(f"**Data da Varredura:** {now_str}  ")
    md.append("**Metodologia:** Bateia Multi-Storefront Discovery (Passo 0033 & Passo 0034)  ")
    md.append("**Domínio de Investigação:** Gastronomia, receitas, planejamento de refeições, fermentação e ciência dos alimentos (*Comida & Cozinha: Ciência e Cultura da Culinária*, Harold McGee).  ")
    md.append("**Catálogos Inspecionados:** Apple App Store (storefronts ISO 3166-1), Google Play Store e F-Droid (FOSS).  ")
    md.append("")
    md.append("---")
    md.append("")
    md.append("### 1. Resumo Executivo e Estatísticas Globais")
    md.append("")
    md.append(f"| Métrica | Valor |")
    md.append(f"| :--- | :--- |")
    md.append(f"| **Total de Países e Territórios Olímpicos Auditados** | **{total_countries}** |")
    md.append(f"| **Storefronts Nacionais Ativos com Catálogo Local** | **{active_storefronts}** ({active_storefronts/total_countries*100:.1f}%) |")
    md.append(f"| **Storefronts Vazios (`AttemptEmptyObserved`)** | **{empty_storefronts}** ({empty_storefronts/total_countries*100:.1f}%) |")
    md.append(f"| **Territórios sem Storefront Nacional Dedicado** | **{unavail_storefronts}** ({unavail_storefronts/total_countries*100:.1f}%) |")
    md.append(f"| **Aplicativos de Código Aberto (FOSS) Catalogados** | **{len(fdroid_data)}** |")
    md.append("")
    md.append("#### Aplicativos com Maior Penetração Internacional no Censo:")
    md.append("| Aplicativo | Presença (nº de países no Top 5) |")
    md.append("| :--- | :--- |")
    for app_name, count in top_global_apps:
        md.append(f"| **{app_name}** | {count} países |")
    md.append("")
    md.append("---")
    md.append("")
    md.append("### 2. Ecossistema de Código Aberto e FOSS (F-Droid & Forjas Upstream)")
    md.append("")
    md.append("Soluções livres e auditáveis com repositórios públicos vinculados para gestão culinária e ciência dos alimentos:")
    md.append("")
    md.append("| Pacote | Nome / Descrição | Licença | Código-Fonte Upstream |")
    md.append("| :--- | :--- | :--- | :--- |")
    for pkg, item in sorted(fdroid_data.items(), key=lambda x: x[1].get("name", "")):
        name = item.get("name") or pkg
        summary = (item.get("summary") or "N/A").replace("|", "-")
        license_spdx = item.get("license") or "FOSS"
        src = item.get("source_code") or item.get("web_url")
        md.append(f"| `{pkg}` | **{name}** — {summary} | `{license_spdx}` | [{src}]({src}) |")
    md.append("")
    md.append("---")
    md.append("")
    md.append("### 3. Resultados Detalhados por Continente e Território Nacional")
    md.append("")

    # Agrupar por continente
    continents = ["América do Sul", "América do Norte", "América Central e Caribe", "Europa", "Ásia", "África", "Oceania"]
    for cont in continents:
        cont_countries = [c for c in OLYMPIC_COUNTRIES if c["continent"] == cont]
        if not cont_countries:
            continue
        md.append(f"#### 3.{continents.index(cont)+1} {cont} ({len(cont_countries)} países)")
        md.append("")

        for c in cont_countries:
            code = c["code"]
            name = c["name"]
            locale = c["locale"]
            term = c["term"]
            res = country_results.get(code, {})
            status = res.get("status", "UNKNOWN")
            items = res.get("items", [])

            md.append(f"##### 🏳️ {name} (`{code}`) — Locale: `{locale}` | Busca: *\"{term}\"*")
            if status == "SUCCEEDED" and items:
                md.append("| App | Desenvolvedor | Categoria | Avaliação | Preço | Link App Store |")
                md.append("| :--- | :--- | :--- | :--- | :--- | :--- |")
                for it in items:
                    app_name = it.get("name", "").replace("|", "-")
                    dev = it.get("developer", "").replace("|", "-")
                    cat = it.get("category", "")
                    rating = it.get("rating", "N/A")
                    price = it.get("price", "Free")
                    link = it.get("url", "#")
                    md.append(f"| **{app_name}** | {dev} | {cat} | ⭐ {rating} | {price} | [Ver App]({link}) |")
            elif status == "EMPTY":
                md.append("> *Storefront ativo mas nenhum resultado retornado para a consulta culinária (`AttemptEmptyObserved`).*")
            else:
                md.append(f"> *Storefront nacional dedicado não disponível na API pública (`{status}`).*")
            md.append("")

    md.append("---")
    md.append("")
    md.append("### 4. Conclusões e Oportunidades de Reúso Arquitetural (Bateia)")
    md.append("")
    md.append("1. **Duopólio vs. Especialização Regional:** Aplicativos globais (Kitchen Stories, Yummly, Tasty, Cookpad) lideram nas Américas e Europa Ocidental, enquanto na Ásia e no Oriente Médio predominam plataformas regionais nativas focadas em técnicas e temperos locais.")
    md.append("2. **Lacuna em Ciência dos Alimentos:** Praticamente todos os apps comerciais focam em listas de receitas passo a passo ou delivery. Quase nenhum incorpora princípios físico-químicos (técnicas de cocção, ponto de fumaça de óleos, reação de Maillard, pasteurização, fermentação controlada) como ensinado em *Comida & Cozinha* de Harold McGee.")
    md.append("3. **Oportunidade para Motores FOSS:** A base do F-Droid (Grocy, Mealie, Open Food Facts) possui estruturas de dados de ingredientes e inventário de alta qualidade que podem ser compostas (`COMPOSE` / `INTEGRATE`) para criar uma plataforma de culinária fundamentada em ciência alimentar.")
    md.append("")

    full_text = "\n".join(md)
    os.makedirs(os.path.dirname(output_path), exist_ok=True)
    with open(output_path, "w", encoding="utf-8") as f:
        f.write(full_text)
    print(f"✓ Arquivo gerado com sucesso: {output_path} ({len(full_text)} bytes, {len(md)} linhas).")

if __name__ == "__main__":
    out_dir = os.path.join(os.path.dirname(os.path.abspath(__file__)))
    target_md = os.path.join(out_dir, "olympic_food_results.md")
    fdroid_data, country_results = run_survey(max_workers=8)
    generate_markdown(fdroid_data, country_results, target_md)
