# 0122 — Pesquisa de modelos de decisão

O Bateia foi usado para coleta pública com commits fixados:

| Repositório | Commit | Itens | Estado |
| --- | --- | --- | --- |
| NandhaKishorM/laya | 3cf26cbcb18725dbc2d127bb8bb2c4c43243ae63 | 29 | Parcial; orçamento deixou alvos sem examinar |
| 0xBakeer/arbiter | bf71358774c85aeff2051ec7b1dcc4577a87706d | 23 | Parcial; orçamento deixou alvos sem examinar |
| apolinario/decision-index | 9eb2dbe2a358004c8782c66e40a83ac07b953fec | 1 | Falhas HTTP 403; item é metadata, não README |

O README do Decision Index foi obtido diretamente de raw.githubusercontent.com
no SHA fixado e registrado como fallback separado. Não muda complete/failures
do coletor. 53 itens do Bateia mais um documento de fallback tiveram hashes
conferidos; são fontes documentais, não avaliação de modelos.

plans.json conserva o pedido; results/github-bundles.json contém bytes e alvos
coletados/pulados/falhos; collection-review.json resume o escopo; verification.json
registra os checks do projeto. APIs públicas sem tokens, sem instalação/clone
ou execução de candidatos. Nenhuma mudança no produto ou nos observadores.

## Diagnóstico preparado

pilot-requests.json contém quatro pares alegação/declaração conhecidos, usando
fontes congeladas de conc/tunny/pond. Referência autoral, não independente.
Request compatível com o formato state/questions/model do servidor nativo Laya;
alias english exige associação a pesos fixados antes de inferência.

Não foi preenchido min_confidence: limiar depende de calibração. As três opções
são supports, contradicts e undetermined. Cada documento é tratado separadamente;
as declarações divergentes de pond não são misturadas numa promessa de runtime.

Para extrair um pedido de diagnóstico (não executa modelo):

```sh
python3 -c 'import json; p=json.load(open("lab/bancada/decisoes/0122/pilot-requests.json")); print(json.dumps(p["cases"][0]["request"], indent=2))'
python3 lab/bancada/decisoes/0122/verify.py
```

prepare_pilot.py recusa substituir payloads existentes. verify.py confere hashes,
referências de commits, estado parcial/falhas e shape dos pedidos. Não verifica
mérito semântico nem mede acurácia. Binário usado na coleta:
/tmp/bateia-product-0121, em diretório temporário sem .env.

Plano de uso, alternativas e critérios no
[relatório](../../../../00_nucleo/relatorios/pesquisa-modelos-de-decisao-0122.md).
