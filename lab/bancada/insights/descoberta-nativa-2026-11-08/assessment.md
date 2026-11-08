# Parecer anônimo

As famílias operacionais centrais são compartilhadas: registrar sem rede e confirmar depois; ou confirmar por custódia local e transferir com recibo. Não há alternativa completa exclusiva demonstrada. Existe em X uma opção complementar parcialmente exclusiva: merge explícito de conjuntos/texto adequados, com identidade de réplica. Isso não autoriza declarar vencedor nem transforma três títulos em três famílias.

Li somente `problem.txt`, `rubric.md`, `packet-X.json` e `packet-Y.json`. Não consultei workspace, mapping, fontes externas ou outros ensaios. Os literais e testemunhos são aceitos conforme a verificação fornecida. Pins dos corpora hospedados são identidades reportadas pelo produto; nenhum código/teste de fonte foi executado e este parecer não verificou Git novamente.

## Notas por candidata

Notas em ordem relevância/evidência/especificidade, sem soma. Todas atingem o limiar de candidata relevante, sob as condições descritas.

| Candidata | Notas | Fundamento e limite decisório |
|---|---|---|
| X-A1 | 4/3/2 | Captura não espera confirmação central; outbox, histórico e reserva transacional explícitos. O contrato “saves the data unconditionally” impede usar force update para saldo/reserva; “Store the cache for 24 hours” exige dedupe além de 72h e de replay de backup. Fotos, export e custo são requisitos, não capacidades demonstradas. |
| X-A2 | 3/3/2 | Colaboração seletiva complementa uma política de reserva. “The default is last-write-wins” e o teste “both concurrent runs survive with distinct identities” tornam concretas modelagem e identidade de réplica. Transporte, manutenção e exclusividade de reserva ficam abertos; não é solução completa independente. |
| X-A3 | 4/3/2 | Custódia com banco transacional na rede local; telefone fora dela observa e aguarda confirmação. `BEGIN IMMEDIATE` e teste `['reserved', 'conflict']` sustentam serialização no mesmo banco, não entre cópias desconectadas. Handoff, fotos e export completo ainda precisam ser construídos. |
| Y-A1 | 4/3/2 | Observações persistentes, pedidos provisórios e confirmação antes de viagem; piloto inclui metadados completos, recuperação e entendimento. LWW/relógios são evidência específica de limite. IndexedDB/Service Workers apoiam uma interface candidata, não durabilidade, sync ou fila implementada. |
| Y-A2 | 3/3/2 | Custódia define responsabilidade, trânsito e aceite; piloto testa ausência/troca de responsável. “there is only one active writer for a log at any given time” sustenta o problema de autoridade. Coordenação local concreta falta: detectar duas sequências depois não prova impedir duas confirmações anteriores. |

Evidência 3 significa que documentação/teste específico permite conferir condição ou limite que muda decisão. Em Y isso vem sobretudo de mecanismos e contraexemplos; não significa cliente pronto. Em X o teste local também não demonstra o protocolo distribuído proposto. Localizadores, trechos integrais usados e lacunas de cada candidata estão em `assessment.json`.

## Pertinência das fontes

Em X, synchronize_cache e datum oferecem componentes de implementação diretamente pertinentes ao cliente: outbox e merge, respectivamente. Não oferecem o sistema integrado pedido. Benchstock é muito próximo ao domínio e transfere a regra de reserva serializada, mas é protótipo por SKU e banco único; não atende diretamente 60 voluntários, peças individuais, fotos e sync. Seu histórico de reserva ausente e CSV sem movimentos mudam concretamente a decisão: não adotar intacto e exigir exportação/reconstrução completa. Inventory transfere operação em rede isolada sem CDN, não gravação em qualquer telefone desconectado. RecallWeave é analogia de replay/supersessão, limitada por escritor único e compactação; não deve ser motor de sync. Hot Stuff ajuda a rejeitar adoção individual por Android 13 e risco de perda com aparelho. Os outros READMEs e o viewer incidental não sustentam alternativas adicionais.

Em Y, caniuse é referência de capacidade de plataforma: armazenamento e inicialização offline tornam um cliente web investigável, mas o contexto é insuficiente para durabilidade/quota/fotos. Cassandra transfere um contraexemplo útil e preciso: LWW com relógios não preserva automaticamente duas observações. Não prova que toda replicação falha nem recomenda infraestrutura distribuída. BookKeeper transfere as distinções entre log imutável e escritor único/fencing, não custódia física ou coordenação offline em Android. As fontes de Y são pertinentes como mecanismos/limites, menos próximas de implementação cliente/inventário. Nenhum pacote encontrou solução pronta comprovada para o problema completo.

## Equivalência e distinções parciais

**X-A1 e Y-A1:** mesmo compromisso operacional de observar offline e deixar promessa remota pendente. X explicita autoridade central transacional para reserva; Y confirma com quem tem a peça e deixa a unicidade dessa confirmação menos formalizada. Essa diferença importa para o piloto, mas não basta para chamar o fluxo inteiro de exclusivo.

**X-A3 e Y-A2:** mesma custódia, responsabilidade local, indisponibilidade em trânsito e aceite. X dá implementação mais concreta de reserva local em Wi-Fi; Y destaca responsável e sequências concorrentes. A exclusividade depende de autoridade acessível e entrega real em ambos. Marcas ou a mera presença de SQLite não criam família adicional.

**X-A2 e a parte de relatos/revisão de Y-A1:** preservação de relatos e disputa física são compartilhadas. A diferença restante é real, embora parcial: X mantém uma opção de merge automático limitado de conjuntos/texto apropriados; Y retém eventos/descrições e revisão, e declara não ter fonte suficiente para promover convergência automática independente. Conservar duas observações não demonstra mesclar texto colaborativo convergentemente. Por isso registro **X-A2-merge como exclusiva observada complementar**, relevante sob notas 3/3/2, com confiança moderada. Ela precisa combinar política de disponibilidade, transporte e demais requisitos; não é uma exclusiva de solução ponta a ponta. Se a unidade for somente o fluxo operacional completo, o resultado é zero exclusivas integrais.

O baseline Y-B3 já imaginava merge convergente, mas não é alternativa final apoiada equivalente; deve ser reconhecido como ideia anterior, não ignorado. A ausência final não demonstra incapacidade do produto. Cliente web em Y, projeto de log, vários nomes de bibliotecas e app individual descartado em X são exclusivas apenas aparentes de famílias úteis.

## O que a pesquisa acrescentou ao baseline

Nenhuma família nova foi demonstrada. Ambos já pensavam em diário, preservação de relatos e custódia/recibos. Não basta repetir a classificação “materialmente modificada” declarada pelos próprios pacotes: a revisão comprovável é mais estreita.

Em X, force update e retenção de dedupe alteram o contrato de integração e acrescentam replay após 96h/backup ao ensaio. LWW versus estruturas explícitas e identidades distintas restringem a opção colaborativa; o teste sustenta investigá-la, não medi-la. A fronteira banco único/réplicas separadas e rede local/telefone isolado restringe onde a reserva pode ser confirmada. Lacunas de CSV/histórico e limites do Android alteram adoção e aceitação. Isso vai além de citar uma ideia lembrada, mas não é descoberta de custódia ou diário. A contribuição à exclusiva parcial de merge é sua concretização e limitação, não a existência da família — presente nos dois baselines.

Em Y, a pesquisa oferece contraexemplo concreto à sobrescrita e retira a convergência genérica de alternativa independente; o baseline já previa preservar conflitos. Escritor único/fencing reforça entrega explícita de autoridade e a proibição de substituí-la por relógio. Capacidade web concretiza uma interface a ensaiar. O ganho é restrição/rejeição de garantias indevidas e critérios de teste; não um aplicativo descoberto ou um protocolo de custódia encontrado. Detalhes de fluxo físico e pilotos são síntese dos autores, não efeitos observados nas fontes.

## Avaliação e incerteza

Ambos propõem ensaio pequeno de 72h com conflitos, reenvios, fotos, exportação/restauração e métricas de trabalho. X inclui replay após 72h, backup/compactação e falha da rede/autoridade. Y dá critérios mais explícitos para entendimento de pendência e proporção/tempo de operações bloqueadas. Esses aspectos se complementam. Os limiares de 15/30 minutos, 10% e uma hora são propostas negociáveis, não resultados.

Faltam durabilidade nos aparelhos reais, carga de 45 mil peças, fotos/saída de dados completas, custo e esforço para dois generalistas, e segurança do handoff. Nenhuma redução de viagem/promessa incorreta foi medida. A pesquisa permite escolher hipóteses e rejeitar integrações perigosas antes de um piloto; não fundamenta contratação, implantação ou superioridade de campo. A comparação não usa quantidade de chamadas, títulos ou marcas e não infere identidades dos produtos.
