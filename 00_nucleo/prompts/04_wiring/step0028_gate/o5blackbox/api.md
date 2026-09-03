# Prompt espelhado — API black-box O5

Defina uma API de oráculo que receba somente `contract_input` e `mutant_contract_output` do protocolo público O5. Não aceite ID de caso, hash de ataque/binding, marcador, expectativa ou resultado predefinido. Produza `Preserved`, `Violated` ou `Unknown` exclusivamente por predicados do contrato R3. Permita avaliação sequencial para determinismo e identidade de runtime.
