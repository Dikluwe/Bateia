# Laudo do Passo 0028 — mecanismos e hipóteses de melhoria

**Veredito:** `MECHANISM_IMPROVEMENT_SYNTHESIS_READY`  
**Atestação:** `EXECUTADO_SEM_ATESTACAO_FORTE_DE_ISOLAMENTO`  
**Contrato ativo:** R4, SHA-256 `0053b5eb34b563100d8165483937d6c8e9e8acf5b2d41a66f379802c37b97322`  
**Selo discriminatório ativo:** O23, SHA-256 `54a0ec1554db5688213eeb940d6057c6756f7efd6b4862099195d71404c6ad3d`  
**Candidato:** 35 arquivos, raiz SHA-256 `cc81bdec3eea946fba3939b36b495c6ca745684c6049badd74aa757b6b4064af`  
**Certificado:** SHA-256 `7b1098c57c6425a2adb2fb6d0d15fea9917502d201e3fb46ae0dfe5dc5ba6506`

## Conclusão

O candidato satisfaz o fragmento observável congelado do Passo 0028: extrai mecanismos arquiteturais somente de evidência avaliada e ligada a fontes imutáveis; conserva incerteza, conflitos e limites de licença; formula hipóteses sem convertê-las em decisão de adoção; publica os sete artifacts novos em execução durável; e acrescenta somente `bateia_start_improvement_synthesis` à superfície MCP.

A certificação vale exclusivamente para os bytes e hashes registrados no contrato R4, catálogo `tekt.mechanism-rules.v3`, protocolo O22, oracle O23, candidato integrado e predecessores 0025–0027.

## Cadeia discriminatória

O gate O22/O23 foi reexecutado após a inspeção do candidato em três ordens, cada uma em dois processos frescos. Foram observados:

- 160/160 negativos como `Violated`;
- 12/12 positivos como `Preserved`;
- 10/10 opacos como `Unknown`, fora do denominador;
- mutation score `1.0`;
- 1.092 processos frescos;
- digest estável `0a8ef723cfea8d65f964755087259ce3d9ef4d2b851a11f920a9ec1f4e945b1c`.

Os 42 arquivos do manifest protegido O23 foram rehashados sem divergência. Receipt discriminatório `0c3b49814ed7ef20017e655da5a26f5942b84d03f868d0d17c13b5cc7331081b` e selo `54a0ec1554db5688213eeb940d6057c6756f7efd6b4862099195d71404c6ad3d` permaneceram byte-identical.

## Gates finais

| Gate | Resultado |
|---|---|
| `go test -count=1 ./...`, duas execuções independentes | PASS |
| suítes protegidas com build tag `step0028_candidate` | PASS |
| race de L1, MCP, L3 e L4 | PASS |
| `go vet ./...` | PASS |
| build de `./04_wiring` | PASS; binário `8abaaf70817cbc30421cde9957b403d70b02934cece74746e837a2b152afb7d3` |
| `crystalline-lint .` | PASS com warnings não bloqueantes |
| `git diff --check` | PASS |
| JSON estrito, hashes e cadeia selada | PASS |
| replay do candidato duas vezes | PASS; bytes idênticos, zero rede |
| MCP compilado real via stdio duas vezes | PASS; bytes idênticos, dez tools |
| scan de raw content, canários e segredos | PASS |
| imports, I/O, catálogo embedded e runtime único | PASS |
| módulo aninhado `lab/bancada/history/benchmarks/0028/go.mod` | PASS |

As primeiras execuções não escaladas de `go test ./...` confirmaram a restrição do sandbox ao abrir listeners de `httptest`. As duas execuções contratuais foram então repetidas com loopback local habilitado e passaram integralmente. Nenhuma rede externa foi autorizada por essa escalada.

## Replay e superfície pública

O replay verificou `mechanism_catalog`, `mechanism_gaps`, `architecture_comparison`, `improvement_hypotheses`, `hypothesis_gaps`, `validation_plans` e `synthesis_provenance`. As duas projeções tiveram SHA-256 `b1b57868ddeaa1404edf00da543fe1c5f463b4a6b0cbbd0c608c8a50f86001b4`, sem transporte de rede, conteúdo bruto de terceiro, canário, segredo ou idempotency key publicada.

O binário real foi executado por stdio em dois roots duráveis novos. As duas respostas tiveram SHA-256 `d92d46b3a8d8c8a9e2dbcf13fc285ecb5ea50e9e34a2c2491272b8beab9c04d1`. A reconstrução das nove tools predecessoras resultou no hash congelado `ee630cff5eada0dee4afe82ad9513dea7b760f0beb7905e915d2714a9114a596`; a única nova tool foi a décima, e o getter recebeu somente os sete enums aditivos contratados.

## Arquitetura e refinamentos

L1 não introduziu I/O externo; L2 não importou L3/L4; L3 não importou camadas superiores; L4 permaneceu composição. O catálogo embedded mede exatamente `52172e289faf65e791f8fa99cad1eee75e349fcbcf78186bcf9446ea4a8f4142`. A fiação contém uma única construção de `FileStore` e uma de `LocalQueue`.

O adaptador adicional `03_infra/envfile/load.go` está corretamente em L3, sem dependência ascendente, preserva valores já presentes no ambiente e não participa do replay offline. O `go.mod` aninhado em `lab/bancada/history/benchmarks/0028` isola somente os dois programas públicos de calibração; as suítes protegidas e os testes de produção continuam enumerados e executados pelo módulo raiz.

Foram emitidos receipts separados preservando:

- `MCP_RECOMMENDATION_GATEWAY_READY` (0025);
- `MCP_DURABLE_RUNS_READY` (0026);
- `MULTISOURCE_MARKET_DISCOVERY_READY` (0027).

## Limitações e invalidação

O workspace foi compartilhado; portanto, não se reivindica isolamento temporal ou de filesystem forte. Os dez casos opacos permanecem `Unknown` por contrato. Nenhuma chamada nova a provedor público ou modelo externo foi necessária ou autorizada para o Passo 0028.

Qualquer alteração em hash causal, input protegido, candidato, catálogo, predecessor, identidade do runtime, replay, superfície MCP ou resultado de gate invalida este laudo e o certificado.
