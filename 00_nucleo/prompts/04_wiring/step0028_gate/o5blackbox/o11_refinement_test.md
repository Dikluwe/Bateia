# Prompt — controles independentes O11 do Passo 0028

Produza controles de regressão binding-free para os refinamentos O11 usando
somente contrato R4, catálogo v3, protocolo O8, artefatos próprios do oráculo e
diagnósticos públicos sanitizados O10.

Exercite relações concretas, nunca IDs de caso ou expectativas importadas:

- `SYNTHESIS_START` não publica artifact nem finding que transforme proposta em
  evidência e o receipt não reutiliza identidade de source histórica;
- uma incerteza real de extração continua `Unknown` quando coexiste com mecanismo
  resolvido;
- `ARTIFACT_READ` com entries mistas reconhece baseline/compatibilidade realmente
  não resolvidos, mas não usa apenas `novelty|transfer_mode=UNKNOWN` como proxy;
- mecanismo documental R4 usa o RuleID documental efetivo no preimage e não pode
  ser elevado acima de `PARTIAL`;
- cancel terminal aceita no máximo uma transição CAS.

Valide também o manifest protegido O11, JSON estrito, repetição, shuffle, race,
vet e ausência de dependência em case ID, attack hash, record, expected, marker
ou candidato. Não leia nem materialize ataques, records/configs, candidate,
receipt, seal, certificado ou veredito.

Após a invalidação O11→O13, o manifest O11 permanece histórico e byte-estável;
ele não deve continuar exigindo que o harness corrente tenha hashes O11.
