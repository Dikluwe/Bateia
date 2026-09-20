# Veredito independente B — passo 0051

Regime A/B executado sem atestação de isolamento: filesystem e ferramentas compartilhados. B recebeu contrato, regras, tipos públicos domain e assinaturas das APIs; não leu corpos candidatos ou testes anteriores. Busca de assinaturas expôs nomes de testes anteriores, sem corpos. Não há selo de refinamento ou alegação de independência física.

A primeira execução foi posterior ao congelamento `tests-freeze-r1.json`: núcleo passou; os dois testes Markdown falharam porque observavam todo o relatório, incluindo a seção geral Arquivos Extraídos. A saída preservada em `independent-r1.log` mostra o hit documental correto e a fonte geral também presente. O contrato restringe o trecho do hit documental, não proíbe a seção geral. A revisão r2 recorta a observação antes dessa seção, sem mudar valores esperados; original do teste preservado em `renderer-tests-r1.go.txt`. A revisão também move mecanicamente o teste do núcleo para parser_test para atender à linhagem do arquivo documentação do parser. Integrador aprovou correção de escopo; não houve leitura da implementação para a revisão.

Após novo congelamento em `tests-freeze.json`, comando `GOCACHE=/tmp/bateia-0048-cache go test ./01_core/parser ./02_shell/report -run Documentation0051 -count=1` passou, conforme `independent-r2.log`.

O recorte distingue função neutra localizada por Doc, método e frase negativa, escopos inelegíveis (pacote, solto, corpo, literal, tipo, gerado, declaração incompleta), prioridade/pesos do código, exato, duplicação, limite de resultados, empate por caminho, projeção documental e mista, adulteração de hash/nome/declaração/grupo Doc/subgrupo, releitura alterada, JSON e renderização com limite de vinte linhas.

Limitações: não executa budget de bytes remoto, não prova ganho universal PT/EN, não mede seleção conjunta nem ação/objeto/condições, não testa diretamente todo o catálogo de capacidades. A frase negativa prova apenas que a projeção continua documental. Gates globais e ensaios em fonte real pertencem ao integrador. A aprovação deste recorte não é prova de cumprimento de requisitos pelo código documentado.
