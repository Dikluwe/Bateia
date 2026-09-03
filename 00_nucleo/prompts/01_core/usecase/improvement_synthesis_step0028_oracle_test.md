# Oráculo independente de síntese de melhoria do Passo 0028

**Camada:** L1  
**Arquivo proprietário:** `01_core/usecase/improvement_synthesis_step0028_oracle_test.go`

Valide os fixtures protegidos e ofereça um harness que receba artifacts capturados
pelo verificador, sem depender de construtor concreto. Cubra produto-only, target
UNKNOWN, licença não reconhecida, invariância a popularidade/novidade, composição,
contraevidência, source hash pré-fila e replay.

Toda hipótese preserva assumptions, risks, arrays explícitos, busca limitada e
plano falsificável. UNKNOWN nunca vira ABSENT/SUPPORTED; licença desconhecida
nunca habilita COMPONENT_REUSE; composição sem witness conserva as cinco
compatibilidades desconhecidas; `AdoptionDisposition` é sempre
`NOT_AN_ADOPTION_DECISION`. Não execute plano nem realize I/O externo.
