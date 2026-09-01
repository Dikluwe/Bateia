# Contrato L0 — domínio do experimento web autenticado

## Autoridade

Este prompt materializa o Passo 0020 para a camada de domínio. A implementação deve ampliar os tipos do snapshot observável sem alterar a intenção congelada do Passo 0019. Credencial é capacidade efêmera da infraestrutura, nunca dado de domínio persistível.

## Vocabulário obrigatório

Modele explicitamente:

- autorização com `provider=github`, `purpose=public_repository_search`, `authentication=read_only_token`, `private_repositories_allowed=false`, `planned_queries=36` e `authorized=true`;
- orçamento congelado com braços `N`, `F1`, `F2`, 12 queries por braço, 36 requests primárias, máximo 10 resultados, timeout por tentativa de 20 segundos, timeout total de 10 minutos, concorrência máxima 2 e no máximo uma repetição por falha transitória;
- tentativa observável, separada da query lógica, com `attempt_number` 1 ou 2 e indicação de repetição contabilizada;
- contadores de rate limit redigidos: limite, restante, usado e instante de reset, todos opcionais; nenhum header bruto;
- visibilidade de resultado suficiente para distinguir público, privado e desconhecido;
- receipt do run com contadores planejados, iniciados, repetidos, concluídos, desconhecidos e pulados por orçamento, além de margem segura declarada e motivo de parada.

Os estados externos são enumeração fechada:

```text
SUCCESS
EMPTY
UNKNOWN_RATE_LIMIT
UNKNOWN_TIMEOUT
UNKNOWN_PROVIDER_RESPONSE
REJECTED_PRIVATE
SKIPPED_BUDGET
```

`Unknown` é informação material. Não pode ser representado por `EMPTY`, lista vazia, sucesso, ausência de erro ou valor zero. Resultados de tentativa desconhecida são vazios, mas seu estado permanece um dos três `UNKNOWN_*`.

## Invariantes estruturais

1. O plano contém exatamente 36 queries lógicas: 12 por braço, IDs não vazios e únicos dentro do run.
2. A ordem primária é exatamente `N1, F1-1, F2-1, ..., N12, F1-12, F2-12`; cada item possui posição global imutável de 1 a 36.
3. Cada query possui exatamente uma tentativa primária e, somente após falha transitória, no máximo uma segunda tentativa da mesma query. Assim, há 36 requests primárias e no máximo 72 tentativas de rede; ambas as contagens são observáveis. Retry não cria query nem altera braço, texto, limite ou ordem primária.
4. Nenhum braço recebe query adicional. Uma repetição apenas repete uma falha transitória daquele mesmo item e é contabilizada separadamente.
5. Todo item planejado termina em uma resposta/tentativa observável ou em `SKIPPED_BUDGET`; nunca desaparece do snapshot.
6. `SUCCESS` exige ao menos um resultado público válido. `EMPTY` exige resposta válida do provedor com zero resultados. Resultado privado produz `REJECTED_PRIVATE` e não é persistido como hit. Visibilidade ausente ou ambígua produz `UNKNOWN_PROVIDER_RESPONSE`.
7. Cada tentativa preserva braço, query ID, request ID, instante, status HTTP opcional, estado, request hash, response hash e rate limit redigido. Cada hit público preserva rank, URL original e canônica, título, descrição e visibilidade pública.
8. Request IDs são únicos por tentativa. Associação entre run, braço, query, posição e tentativa é imutável.
9. Hashes são SHA-256 hexadecimais de serialização canônica documentada. O hash de resposta inclui identidade, estado, HTTP, rate limit redigido e hits, mas exclui o próprio hash. O hash do snapshot inclui autorização redigida, orçamento, plano e todas as tentativas, mas exclui o próprio hash.
10. Token, valor do token, `Authorization`, cookies e derivados capazes de reconstruir credencial são campos proibidos em todo valor serializável, hash, erro e receipt.

## Validação e veredito

Validadores de domínio devem rejeitar snapshots com estado desconhecido, associação inválida, hashes malformados, retry excedente, resultado não público, orçamento divergente ou campo sensível. Falha fechada: opacidade vira `Unknown` ou invalidação, nunca sucesso.

O veredito do run usa precedência:

1. vazamento de token ou conteúdo privado: `FAIL_SECURITY`;
2. orçamento excedido: `FAIL_BUDGET`;
3. hash ou replay divergente: `FAIL_REPRODUCIBILITY`;
4. cota impede exposição comparável: `INCONCLUSIVE_RATE_LIMIT`;
5. exposição não pode ser provada comparável: `UNKNOWN_PROVIDER_EXPOSURE`;
6. integridade passa, cada braço conclui ao menos 80%, diferença máxima de conclusão é 10 pontos percentuais e nenhuma estratégia recebe query extra: `READY_FOR_BLINDED_JUDGMENT`.

Conclusão significa `SUCCESS` ou `EMPTY`. `UNKNOWN_*`, `REJECTED_PRIVATE` e `SKIPPED_BUDGET` não contam como concluídos.

## Compatibilidade

Preserve replay de snapshots anteriores quando o schema os identificar. O novo schema deve ser versionado e não reinterpretar `UNKNOWN` legado como vazio ou sucesso.
