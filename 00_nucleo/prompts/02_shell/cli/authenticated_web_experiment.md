# Contrato L0 — CLI do experimento web autenticado

## Superfície

Disponibilize uma operação explícita para o Passo 0020, distinta do experimento público do Passo 0019. Ela recebe por flags somente caminhos dos conjuntos N/F1/F2, diretório de saída, manifesto/autorização e parâmetros que devem coincidir com o orçamento congelado. O token não pode ser aceito por flag, stdin, arquivo de configuração ou texto de autorização.

A autenticação é habilitada apenas por uma flag inequívoca de execução autenticada e o adaptador lê o nome fixo da variável de ambiente no ponto de composição. Ausência da variável, autorização inválida ou uso do modo público deve falhar antes da rede. Ajuda, erros e saída nunca exibem o valor nem informam prefixo, sufixo, comprimento ou hash do token.

## Comportamento

- exigir exatamente os três arquivos de queries e um diretório de saída;
- exigir autorização explícita compatível com os sete campos congelados do Passo 0020;
- fixar 10 resultados, 20 segundos por tentativa, 10 minutos totais, concorrência 2, 36 primárias e uma repetição transitória máxima; opções divergentes são rejeitadas, não ajustadas;
- registrar a autorização redigida antes da primeira request;
- executar o plano intercalado do caso de uso;
- persistir de forma atômica os três snapshots e o receipt sanitizado somente após varredura de segredo e privacidade;
- em falha de segurança, não deixar snapshot parcial, receipt parcial ou arquivo temporário com material do run;
- retornar erro/veredito específico sem ecoar resposta ou headers brutos.

Saídas legíveis por máquina devem conter apenas caminhos, hashes, contagens redigidas e veredito. Não imprimir URLs com query string, headers, cookies, ambiente ou dumps de request/response.

## Replay

`--resume-dir` pode apontar para snapshots autenticados anteriores do mesmo
conjunto. Valide braço e hash, preserve apenas `SUCCESS`/`EMPTY` e execute de
novo casos opacos ou pulados. Grave um novo run completo sem alterar a origem.

A operação existente de replay deve aceitar o novo schema sem token, sem leitura de ambiente e sem provedor injetado. Ela deve falhar de modo não destrutivo se detectar segredo, privado, hash divergente, orçamento inválido ou estado desconhecido não suportado. Duas invocações sobre o mesmo snapshot devem produzir arquivos byte a byte idênticos.

## Códigos observáveis

O processo distingue sucesso pronto para julgamento, inconclusão por cota/exposição e falhas de segurança, orçamento e reprodutibilidade. `Unknown` não usa código de sucesso equivalente a resultado vazio. Mensagens podem nomear campo/estado e request ID redigido, mas nunca incluir payload bruto sensível.

## Passo 0131
Aplicar o ownership do ADR-0014 e passo 0131. Jornada/decisões em L1, transporte/renderização em L2, efeitos em L3, composição em L4. Preservar os contratos observáveis, formatos, hashes, versões, ordem, retomada e cancelamento; testes seguem o proprietário ou Lab para auditoria geral. Não conservar regra duplicada na casca/runtime.

## Consolidação 0131
Admissão, veredicto e agendamento pertencem a web_experiment_run.go em L1. Os auxiliares privados compatíveis delegam; L2 conserva flags, JSON estrito, filtro de dados sensíveis e publicação.
