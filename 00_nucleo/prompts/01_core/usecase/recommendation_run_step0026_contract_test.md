# Prompt espelhado — oráculo L1 de execução durável 0026

## Autoridade

Teste A/B independente derivado do Passo 0026 e do contrato observável R1, sem leitura de implementação candidata.

## Obrigações observáveis

- Canonicalizar campos-conjunto e calcular hash idempotente sensível a mudanças semânticas.
- Validar chave e pedido antes de store, reserva de fila e geração de ID.
- Criar/confirmar exatamente um run sob concorrência; conflito e fila cheia não deixam reserva nem associação.
- Aplicar somente os arcos permitidos, CAS e progresso `planned=8` monotônico.
- Publicar toda a matriz requerida antes de `SUCCEEDED` e preservar `SUCCEEDED+INCONCLUSIVE`.
- Cancelar de forma idempotente e recuperar órfãos como `INTERRUPTED`, sem reescrever terminais.

## Fakes

Store, fila, gerador e executor são memória controlada, sem filesystem, goroutine interna, rede ou relógio. O fake não suaviza CAS nem transições; testemunhas concorrentes usam sincronização explícita.
