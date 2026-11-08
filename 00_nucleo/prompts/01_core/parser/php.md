# Prompt: Parser de Manifestos PHP composer.json (parser/php)

**Camada**: L₁ — Núcleo  
**Criado em**: 2026-07-31  
**Arquivos gerados**: `01_core/parser/php.go`, `01_core/parser/php_test.go`  

---

## Contexto

Na etapa de Triagem do Bateia para o ecossistema PHP, o conteúdo de texto do manifesto `composer.json` precisa ser analisado para extrair o nome do pacote, versão e dependências (produção e desenvolvimento).

---

## Restrições Estruturais

- **L₁ (Núcleo)**: Função estritamente pura. Zero I/O, zero leitura de disco ou chamadas HTTP. Recebe a string de conteúdo e retorna a struct `*domain.Manifest`.
- Depende apenas de `01_core/domain`, `encoding/json` e `fmt`.
- Proibido importar $L_2$, $L_3$ ou $L_4$.

---

## Instrução

Criar a função:
`func ParseComposerJson(ref domain.CommitRef, content string) (*domain.Manifest, error)`

Regras do parsing:
1. Parsear o JSON em uma struct interna `ComposerJSON`.
2. Mapear o campo `name` para `ModuleName`.
3. Mapear o campo `require` como dependências do tipo `"direct"`.
4. Mapear o campo `require-dev` como dependências do tipo `"dev"`.
5. Preencher `Language: "php"`.
6. Manter `ConsultadoEm` zerado; o instante de coleta é fornecido explicitamente no dispatcher (ADR-0011).

---

## Critérios de Verificação

```
Dado um composer.json com "name": "laravel/laravel", "require": {"laravel/framework": "^10.0"}, "require-dev": {"phpunit/phpunit": "^10.0"}
Quando ParseComposerJson for executada
Então retorna um Manifest com ModuleName "laravel/laravel", Language "php", 1 dependência "direct" e 1 dependência "dev".
```

---

## Resultado Esperado

- `01_core/parser/php.go`: Implementação do parser de composer.json.
- `01_core/parser/php_test.go`: Suite de testes cobrindo composer.json.

---

## Histórico de Revisões

| Data | Motivo | Arquivos afetados |
|------|--------|-------------------|
| 2026-07-31 | Criação inicial Tekt L0 para PHP | php.go, php_test.go |

## Pureza temporal — 0128

Não consultar o relógio neste parser. ConsultadoEm permanece zero quando a coleta não foi informada. Preservar assinatura, parsing e erros; ParseManifestAt aplica o instante explícito recebido dos coletores L3. O ADR-0011 substitui qualquer obrigação anterior de gerar a hora atual neste componente.

## Ordem de dependências — 0128

Iterar os objetos JSON por chaves lexicograficamente ordenadas, usando sortedManifestKeys. Preservar grupos, tipos, versões, exclusões e política de erros. A ordem anterior dos mapas Go era indeterminada; a nova ordem não modifica artefatos históricos.
