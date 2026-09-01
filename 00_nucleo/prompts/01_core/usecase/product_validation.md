# Product validation use case

Valide domínio, dez IDs únicos e braços únicos com contagens coerentes. Calcule
conclusão por braço e exija pelo menos 80% em cada um. Se qualquer braço falhar,
produza `FAIL_DISCOVERY`, identifique a fase e mantenha evidência, decisão e
utilidade como `NOT_EVALUATED`. Nunca use média agregada para esconder um braço.
Ordenação da entrada não altera a saída.
