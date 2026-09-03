# Prompt espelhado — subjects variants O3

Decodifique uma tabela protegida com exatamente M001–M160, vinculada por hash ao diagnóstico sanitizado O2→O3. Cada entrada cria um `VariantSubject` independente, baseado no subject de referência, que altera exatamente um campo observável declarado para seu ID. Rejeite campos não suportados, IDs duplicados, schemas divergentes e documentos com bytes JSON adicionais. Nenhuma variante pode ler `Expected`, candidato ou ataques.
