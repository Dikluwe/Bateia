# Prompt espelhado — avaliador black-box O5

Avalie uma troca somente por bytes de `contract_input` e `mutant_contract_output`. Aplique invariantes globais de pureza, efeitos, vazamento, erro controlado e identidade, depois predicados específicos da porta. Falha normativa é `Violated`; opacidade contratual concreta e sem sobrealegação é `Unknown`; somente uma troca inteiramente decidível e conforme é `Preserved`. Em sequências, a mesma entrada exige saída canônica idêntica e as identidades de runtime não podem mudar. Não aceite bindings, IDs, hashes de ataque, expected ou branches por caso.
