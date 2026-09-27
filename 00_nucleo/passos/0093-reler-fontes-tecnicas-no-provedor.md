# Passo 0093 — Reler fontes técnicas no provedor

## Ganho

O passo 0092 verifica a consistência interna de um bundle, mas um JSON pode
ser fabricado com bytes e hashes coerentes. O novo comando `idea evidence
recheck` não recebe bundles: relê os arquivos pelo provedor GitHub no commit
fixado, recalcula os hashes e aplica os verificadores existentes.

O relatório identifica essa observação como `verification_basis:
PROVIDER_RECHECK`. Isso distingue releitura da origem configurada de mera
conferência de um bundle fornecido. Não constitui assinatura digital ou
certificado de autenticidade de um JSON posteriormente exportado.

## Entrada e execução

A entrada é um `TechnicalSourceRecheckRequest` com `batch` (o mesmo lote de
achados do passo 0092) e `candidates` (candidatos explícitos). Cada candidato
deve informar `owner`, `repo` e `commit_id` com SHA hexadecimal completo de
40 caracteres. `url` é opcional; se presente, deve corresponder a
`https://github.com/<owner>/<repo>`. Cada achado usa a chave
`<owner>/<repo>@<commit>` em `candidate`, o mesmo commit em `source.revision`
e o digest esperado em `source.digest`.

```sh
bateia idea evidence recheck --input releitura.json --as-of 2026-09-26T00:00:00Z --output releitura-report.json --bundles-output bundles-relidos.json
```

O executável configura o provedor GitHub existente; `GITHUB_TOKEN` é opcional
e pode ser necessário para repositórios privados ou limites da API.
`--timeout` é de um minuto por padrão e aceita até cinco minutos.

O relatório inclui revisão, bundles e falhas operacionais. `--bundles-output`
exporta apenas os bundles para o fluxo já existente:

```sh
bateia idea decide --idea ideia.json --assessment assessment.json --policy reuse-policy.json --technical-evidence achados.json --technical-bundles bundles-relidos.json --evidence-as-of 2026-09-26T00:00:00Z
```

Esse último comando continua sendo uma verificação offline. A indicação de
origem na releitura não protege os arquivos exportados contra adulteração;
para uma observação nova, execute `recheck` novamente. Não há mudança no modo
simples de pesquisa e busca.

## Barreiras

- Somente suporte atual de `LICENSE_ID` e `SBOM_COMPONENT` com técnica
  correspondente pode gerar releitura. Achados vencidos, contraditos ou sem
  verificador não provocam consultas.
- A árvore na revisão é conferida e todos os documentos raiz de licença
  são incluídos. Não se omite uma segunda licença para criar certeza falsa.
- Revisão divergente, extração stale, caminhos extras/duplicados, arquivos
  ausentes e falhas do provedor tornam o bundle incompleto e impedem promoção.
- O hash é calculado dos bytes devolvidos pelo provedor e comparado ao digest
  esperado. Divergência deixa o achado não verificado.
- Limites: dez candidatos, cem achados, vinte arquivos por candidato e quatro
  MiB de conteúdo admitido no total. O limite de bytes é aplicado após a
  extração; não limita o tráfego nem a memória do cliente remoto existente.
- Entradas e saídas não podem apontar ao mesmo arquivo. O comando não executa
  código dos candidatos. O provedor existente pode usar clone temporário como
  fallback de leitura.

Worktrees locais não são aceitos neste comando: seus bytes atuais não provam
o conteúdo de um commit declarado. Identificação de licença não é análise
jurídica; componente declarado em SBOM não prova presença no binário. Apenas
o gate de licença anteriormente inconclusivo pode ser preenchido no fluxo de
decisão existente; os demais critérios continuam necessários.

## Testes

Provedor simulado verifica leitura fixada, hash recalculado, conteúdo
divergente, extração incompleta, revisão errada, stale, ambiguidade de licença,
falhas de árvore/leitura, orçamento de bytes, cancelamento e rejeição de
alvos não fixados. Testes da CLI verificam exportação dos bundles e proteção
contra sobrescrita da entrada. Não exigem acesso à rede.
