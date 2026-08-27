# Idea evidence diagnosis CLI

Expor `bateia idea evidence diagnose` como comando estritamente offline. Exigir
assessment, bundles e julgamentos; aceitar catálogo versionado opcional; escrever
o relatório de cobertura em JSON para stdout ou arquivo.

Erros devem identificar a entrada inválida. O comando não acessa rede, não lê
tokens e não modifica os artefatos congelados.
