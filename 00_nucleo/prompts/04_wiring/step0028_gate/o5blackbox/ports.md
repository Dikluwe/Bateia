# Prompt proprietário — predicados black-box O5 das portas do Passo 0028

Implemente em `04_wiring/step0028_gate/o5blackbox/ports.go` os predicados
externos e binding-free das portas `SYNTHESIS_START`, `SYNTHESIS_EXECUTE`,
`RUN_READ`, `ARTIFACT_READ`, `RUN_CANCEL`, `TOOLS_LIST`, `REPLAY` e
`REFINEMENT_AUDIT`.

Use como autoridades semânticas somente o contrato observável 0028 revisão R3,
o catálogo fechado `tekt.mechanism-rules.v2` e o protocolo executável O5. Não
importe implementação candidata, descrições de mutação, expectativas por caso,
bindings privados, registros ou vereditos.

Os predicados devem ser funções determinísticas apenas de `contract_input` e do
`mutant_contract_output` observado. Valide JSON estrito, shapes e enums fechados,
ordem e schema das dez tools, pré-validação antes da admissão, lifecycle e status
disjuntos, publicação atômica dos cinco artefatos de síntese, bytes/hash/size de
artefato, replay offline e identidade única do runtime. Rejeite efeitos externos,
vazamento de campos não seguros, publicação parcial, confusão entre lifecycle e
status e alegação de refinamento sem os três receipts congelados.

Classifique como `Unknown` somente quando a própria construção observada conserva
uma opacidade contratual real, como source histórica indisponível sem releitura,
status de síntese `INCONCLUSIVE`, target/novelty/transfer `UNKNOWN`, licença
`UNRECOGNIZED`, teto inconclusivo ou classe de mecanismo fora do catálogo. Nunca
converta `Unknown` em sucesso nem o use para tolerar uma violação estrutural.

Preserve mensagens de falha genéricas e sem binding por caso. O arquivo pertence
à camada L4 e não deve introduzir I/O, rede, relógio, concorrência nem dependência
de L1–L3.

## Suplemento O8 / R4

Replay e o envelope causal usam somente o hash público do catálogo v3. Todas as
outras relações de porta calibradas em O6 permanecem inalteradas.
