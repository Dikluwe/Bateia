# Laudo 0019 — Experimento público da gramática funcional

**Veredito**: `INCONCLUSIVE_RATE_LIMIT`
**Estado**: `EXECUTADO`

## Execução

Foram enviadas 36 consultas congeladas à API pública de busca de repositórios
do GitHub, sem token, cookies ou acesso a repositórios privados. A ordem foi
intercalada por posição: N, F1 e F2.

| Braço | Sucesso | Rate-limit/Unknown | Conclusão | Hits orgânicos únicos |
|---|---:|---:|---:|---:|
| N | 4/12 | 8/12 | 33,3% | 15 |
| F1 | 3/12 | 9/12 | 25,0% | 0 |
| F2 | 3/12 | 9/12 | 25,0% | 1 |

Todos os snapshots passaram na verificação de hashes. Cada replay foi executado
duas vezes e os pares foram idênticos byte a byte.

## Por que não houve julgamento

O gate exigia ao menos 50% das queries concluídas em cada braço. Nenhum braço
alcançou esse limite. Julgar os resultados parciais confundiria qualidade da
formulação com exposição desigual ao provedor.

O registry competitivo permaneceu fechado. Precisão, recall e ganho permanecem
`null`; a gramática v2 não foi promovida nem rejeitada.

## Observação do provedor

Uma resposta pública trouxe uma descrição de repositório anormalmente grande,
contendo texto de documentação e exemplos literais de nomes de secrets. Não era
credencial local e foi preservada no snapshot bruto porque faz parte da resposta
do provedor. Esse hit não foi julgado nem promovido.

## Conclusão

O experimento terminou corretamente como inconclusivo por rate-limit. A próxima
execução não deve repetir chamadas até obter resultado favorável. Um novo passo
deve escolher entre uma janela pública agendada, autorização explícita para um
token somente leitura, ou outro índice observável, mantendo os três braços e
budgets idênticos.

O protocolo Tekt foi executado sem atestação de isolamento: uma mesma autoridade
conversacional coordenou contrato, execução e veredito.
