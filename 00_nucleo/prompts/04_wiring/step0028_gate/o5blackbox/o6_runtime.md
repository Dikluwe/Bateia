# Prompt proprietário — refinamento O6 das portas runtime do Passo 0028

Produza `04_wiring/step0028_gate/o5blackbox/o6_runtime.go` como camada de
refinamento segregada do oráculo O5. Não altere `ports.go`. Exponha entry points
O6 prefixados para `SYNTHESIS_START`, `RUN_READ`, `RUN_CANCEL`, `REPLAY`,
`REFINEMENT_AUDIT` e `TOOLS_LIST`, permitindo integração posterior explícita.

Derive toda decisão somente do contrato R3, do catálogo fechado v2, do protocolo
executável O5, dos artefatos O5 próprios e dos diagnósticos públicos sanitizados
autorizados. Não leia nem codifique conteúdo, descrição, hash, expectativa ou
identidade privada por caso.

Para start, valide shape requerido, texto UTF-8 e controles, limites públicos,
IDs, refs, sets, source pairs e policy antes de qualquer efeito. Modele exatamente
os prefixos `reserve -> id -> create`, distinguindo reuse, criação completa,
rejeição de fila e falhas pré-fila; sucesso não pode terminar num prefixo parcial.

Para leitura e cancelamento, valide a representação integral do run, a disjunção
dos três status, cinco estágios de synthesis, progress e terminalidade. Não exija
status materializado para terminal degradado, mas exija-o em `SUCCEEDED`.
Cancelamento não resolve, reserva, gera ID, cria run nem publica, e uma transição
de cancelamento observada deve avançar a versão CAS.

Replay valida integridade de inputs congelados, permanece offline e exige hashes
dos sete artefatos produzidos e seus índices. Não compare esses hashes com o mapa
de hashes das sources do request: autoridades causais distintas não são
intercambiáveis. Preserve `Unknown` somente quando o conteúdo congelado expõe
estado Unknown, licença não reconhecida, budget esgotado, baseline ausente ou
compatibilidade não medida.

Refinement deve rejeitar sobrealegação: booleans de receipt seguem os hashes
congelados, layers seguem a evidência fornecida e `runtime_unique` requer as cinco
identidades observadas. Receipt, layer ou identidade ausente com resultado
conservador é `Unknown`, não `Violated`; alegação positiva sem evidência é
violação.

Valide tools por constraints normativas observáveis, sem comparar serialização,
ordem de propriedades ou keywords JSON Schema irrelevantes. Preserve exatamente
os dez nomes, a posição da tool 10, sua descrição, constraints causais e a enum
aditiva de artefatos.

## Suplemento O9

Antes da fila, rejeite refs de proposta que cruzem namespaces de mecanismo e
produto, categorias referenciadas sem source declarada, múltiplos `run_kind` para
o mesmo `run_id` e associação incoerente entre hash e artifact. No cancelamento,
uma transição CAS avança exatamente uma versão; um estado terminal observado
permanece estável e idempotente.

## Suplemento O11

Vincule este refinamento a R4, catálogo v3 e protocolo O8. Um start não pode
materializar evidência nem findings na observação de publicação: ele devolve
somente receipt ou erro controlado. Em sucesso, o novo `run_id` não pode colidir
com nenhum `run_id` histórico declarado em `source_artifacts`; a identidade da
source permanece histórica e imutável.

Cancel terminal aceita exatamente o estado na versão esperada ou uma única
transição CAS para a versão seguinte. Qualquer segundo avanço viola
idempotência. Preserve todos os predicados anteriores de resolução pré-fila,
admissão, lifecycle, replay e tools.
