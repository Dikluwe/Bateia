# Controle independente de identidade live 0055

Papel B, A/B sem atestação física. Fonte: passo0055 e adendo-checkpoint congelado antes deste teste. Ler somente testes anteriores como fixture, sem implementação. Criar checkpoint por porta CLI live e verificar discovery_query=0055-v1. Remover ou alterar versão torna checkpoint inelegível; versão atual permite reuso. Teste legado M21ZeroLocaleGoldenIdentity cobre preservação replay e não deve ser alterado. Não executar antes de comunicar hashes ao coordenador.

Revisão r2: adendo-sintaxe0055 substitui explicitamente apenas emissão antiga contentquoted por `patterntype:regexp (a|b) AND (x|y)` e versão de plano/checkpoint por 0055-v2. Oráculo rejeita content/aspas. Correção do harness HTTP drena corpo e limita espera/cleanup de handlers sem mudar asserts de erro. R1 preservado em benchmark/0055/test-r1-*.txt. Nenhum resultado ou produção lido.

## Supersessão independente 0056

Contrato0056 exige discovery_query live0056-v1, preservando catálogo lexical0055-v2. Antes do código: atualizar somente versão corrente e adicionar0055-v2 como caso obsoleto explícito na fixture legada, mantendo os demais controles de replay/reexecução. Go baseline SHA256 5ee57d7614d0b502fb7cec0efc943b9b8435c07822cf61374065c72a3358ad2c; prompt baseline338680eba4bf65093f0eb84304d0e0f1bc9913e27770fd6e815db7452252472d. Nenhuma produção ou resultado lido. Freeze antes execução.
