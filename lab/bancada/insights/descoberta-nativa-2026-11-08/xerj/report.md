# Ensaio XERJ independente — inventário cooperativo

Foram elaboradas duas alternativas de desenho para experimentar, com evidência parcial de mecanismos reais. **Nenhum aplicativo pronto que satisfaça o problema foi descoberto ou validado.** Não há base nesta pesquisa para escolher uma biblioteca de sincronização, prometer durabilidade em celulares baratos ou aprovar implantação. As mudanças de processo e os detalhes de aplicação abaixo são síntese deste ensaio; não são capacidades atribuídas aos projetos lidos.

Diretório completo: `/tmp/bateia-xerj-open-0154/xerj`. Entregas estruturadas: `alternatives.json`, `sources.json`, `coverage.json`. Testemunhos originais: `audit/NNN.json`, campo `stdout`. O baseline foi escrito antes da primeira chamada do produto e permaneceu congelado: SHA-256 `974fd615f3b14e4a0e97d772649f2bb2625e4b8a1e6691a0b4e3f73272813034`.

## O que a pesquisa conseguiu estabelecer

A descoberta partiu de problemas e capacidades: offline, inventário, replicação, conflitos, armazenamento móvel e histórico. Nenhum nome ou URL de projeto foi usado como atalho de memória. Depois da descoberta, nomes retornados pelo catálogo foram usados para selecionar o corpus de leitura. Foram lidos três corpora nativos, todos fixados em revisão de 2026-10-13:

- **caniuse**, pin `e951c6cd80296dab56d08a932f5f879e786fa774`, licença CC-BY-4.0. `caniuse/fulldata-json/data-2.0.json`, linha 1: “Method of storing data client-side, allows indexed database queries.” (`audit/025.json`, E1). `caniuse/data.json`, linha 1, Service Workers: “including hooks to enable bootstrapping of web applications while offline.” (`audit/027.json`, E2). Isto sustenta experimentar cliente web com cache e armazenamento local; não comprova sincronização, quota, durabilidade de fotos, suporte dos aparelhos concretos nem entrega em background. A linha 1 é o localizador real do JSON monolítico retornado; não se inventaram números mais precisos.
- **cassandra**, pin `ec2c276be145bb25085c5c35ac0da69113b676ff`, Apache-2.0. `cassandra/doc/modules/cassandra/pages/architecture/dynamo.adoc`, passagem iniciada na linha 251: “Updates / resolve according to the conflict resolution rule of last write wins.” A barra representa a quebra de linha literal registrada em E3 (`audit/029.json`). A mesma passagem explicita dependência dos relógios. **Inferência:** replicar um campo mutável não é conservar todas as observações. Uma correção pode ganhar e ocultar outra descrição; o desenho precisa tratar histórico separadamente.
- **bookkeeper**, pin `773a988ddf8b888a8869c6ae33f696935eba5cdc`, Apache-2.0. `bookkeeper/site3/website/docs/api/distributedlog-api.md`, passagem da linha 1: “A *log* in DistributedLog is an ordered, immutable sequence of *log records*.” (E4). Passagem da linha 36: “there is only one active writer for a log at any given time.” (E5). Ambas em `audit/034.json`; a segunda descreve também fencing durante partição. **Inferência:** histórico imutável e autoridade para confirmar são decisões diferentes. Não se inferiu que BookKeeper seja adequado a telefones ou barato para dois mantenedores.

Essas fontes são referências de mecanismos, não três alternativas conceituais nem uma lista de produtos recomendados. A identidade citada é corpus+pin, pois o catálogo não forneceu URL de repositório desses três projetos. O único link GitHub elegível observado, um polyfill de IndexedDB, não foi adquirido: não respondia ao conflito de inventário e uma leitura posterior indicava abandono de WebSQL. Nenhuma URL foi reconstruída por memória.

## Mudanças comuns propostas

Etiquetar individualmente as peças, incluindo itens visualmente idênticos. Registrar chegada, movimentação, separação, uso e correção como novas observações com autor, identificador da operação e referência à peça. Uma correção referencia o que corrige; não apaga o relato anterior. Horário do telefone ajuda a pessoa a entender o relato, mas não deve escolher sozinho qual observação é verdadeira.

Mostrar separadamente “salvo neste telefone”, “recebido”, “pedido pendente”, “confirmado” e idade da informação. Fotos podem sincronizar depois dos metadados, com pendência visível. Exportar histórico, inventário e manifesto de imagens, com originais acessíveis, é requisito de aceitação. Idempotência, fila, fotos e exportação são trabalho de implementação ainda não demonstrado pelas fontes.

## A1 — Observar livremente, confirmar depois

Qualquer voluntário registra operações offline. Descrições concorrentes permanecem visíveis; estados fisicamente incompatíveis geram revisão. Separações remotas são pedidos provisórios. Antes de viagem ou promessa definitiva, confirmar com a oficina que possui a peça. Uma projeção de estoque facilita consulta, mas o diário mantém o que as pessoas disseram.

É uma modificação material de B1 e da preservação de observações de B3. E1/E2 tornaram concreta uma possível interface web offline; E3 remove a confiança em “a replicação resolve”; E4 fundamenta a distinção entre diário e estado calculado. Nenhuma fonte demonstrou um CRDT de reservas, e esse mecanismo não é prometido.

Atende por intenção às 72 horas de registro desconectado e reduz correções destrutivas. O preço é tolerar pedidos provisórios, conflitos e revisão humana. Não impede que duas pessoas façam promessas físicas inconsistentes se ignorarem a pendência. Metadados das 45 mil peças, fila e fotos precisam caber no pior telefone. Uma aplicação e um serviço simples são hipótese de manutenção; custos recorrentes e esforço dos dois desenvolvedores seguem desconhecidos.

Aplicável quando compromissos definitivos podem esperar confirmação e o volume de revisão for pequeno. Rejeitar se voluntários confundirem pedido com disponibilidade, se o histórico perder relatos ou se o trabalho diário de reconciliação ficar alto.

## A2 — Custódia por oficina e transferência com recibo

Cada peça pertence a uma custódia local. Todos podem observar offline, mas uma autoridade local explícita confirma separação e uso. Outra oficina envia pedido. No despacho a peça deixa de estar disponível no remetente e entra em trânsito; o destinatário só a oferece após aceite físico. Troca de responsável exige entrega de autoridade, não apenas expiração pelo relógio.

É uma modificação material de B2. E4/E5 fornecem a analogia de diário e escritor único, incluindo o problema de dois escritores durante partição. O protocolo de custódia proposto é síntese, não uma funcionalidade encontrada no BookKeeper. Ele só funciona se a autoridade local for de fato acessível e o trabalho físico seguir a entrega.

Troca conflitos posteriores por restrição sobre quem confirma. Dispensa ida diária a uma sede e mantém observações durante desconexão, mas pedidos remotos aguardam comunicação. Ausência do responsável pode bloquear transições definitivas; dois responsáveis offline podem produzir dupla autoridade. Detectar esse caso e suspender confirmação é requisito ainda não implementado. Não se afirma exclusividade entre celulares desconectados por simples regra de interface.

Aplicável se cada oficina aceitar responsabilidade local, recibos e troca de turno disciplinada. Rejeitar se a autoridade se tornar gargalo ou se a entrega produzir duas peças disponíveis no catálogo embora exista apenas uma peça física.

## Avaliação pequena antes de compromisso

Propor um piloto comparável de duas oficinas, seis telefones reais e 200 peças etiquetadas; carregar também 45 mil metadados representativos. Não foi executado neste ensaio. Durante 72 horas sem internet, registrar reservas e usos concorrentes, descrições e correções incompatíveis, relógios divergentes, reenvios, reinício de telefone e fotos pendentes. Na reconexão, variar ordem e duplicar entregas.

Critérios comuns: nenhum evento confirmado localmente desaparece; nenhuma correção apaga outra observação; reenvio não duplica uso; histórico e fotos exportados podem ser reconstruídos fora do cliente. Medir armazenamento, latência, bateria, dados transferidos, tempo de manutenção e entendimento dos voluntários. Fazer uma tentativa real de restaurar dados após falha, sem presumir que compatibilidade significa durabilidade.

Para A1, propor limite de 15 minutos de revisão por oficina/dia no cenário acordado e exigir que ninguém trate pedido pendente como confirmação. Para A2, simular responsável ausente, troca de turno, telefone perdido, tentativa de dupla autoridade e transferência sem aceite. Propor rejeição se mais de 10% das operações locais aguardarem autoridade por mais de 15 minutos. São limiares sugeridos para negociação, não dados observados. Comparar promessas incorretas, viagens evitadas e atrasos de operação, não apenas sucesso de sincronização.

## Cobertura e limites

B1 e B2 foram apresentados; B3 não foi promovido a alternativa independente porque a descoberta não forneceu evidência suficiente de convergência que preserve a semântica de inventário. Sua parte útil — manter observações — entrou em A1 sem alegação de convergência automática. B4, etiquetas/contagem/confirmação, permanece complemento das duas propostas: não houve evidência de planilha que satisfaça as 72 horas sem perda e sem reconciliação destrutiva.

Muitos resultados do catálogo foram tangenciais. A primeira descoberta falhou por catálogo ausente (`004`); sync obteve catálogo de 514 corpora (`006`). O build não tem cliente hospedado (`011`). Cassandra excedeu o limite de auto-pull de 50 MB (`016`), resolvido com pull explícito pelo produto (`018`). Uma tentativa de k=30 foi rejeitada e corrigida para 10 (`015`, `017`). Não se acessaram arquivos de fonte diretamente nem se executaram projetos. Testes de recuperação e de operação descritos acima são propostas, não resultados.

Permanecem pesquisáveis: aplicativo pronto de inventário offline, protocolo de reserva exclusiva, durabilidade/quota real em Android barato, fotos, exportação completa, custos e acessibilidade. A cobertura estruturada registra descartes, evidências e lacunas. Foram registradas 27 chamadas classificadas pelo runner como pesquisa e nove como setup, zero aquisições via acquire. Parada antes do teto: as três fontes úteis foram lidas e o baseline foi revisto; novos nomes lembrados não foram usados para preencher lacunas.

Não houve intervenção humana durante a execução. O nó automático próprio foi iniciado na porta 9300, PID 1201802, por pull nativo (`018`); encerrado somente esse nó (`035`); busca no endpoint próprio confirmou inacessibilidade (`036`). Nenhum timer, publicação ou feedback foi criado.
