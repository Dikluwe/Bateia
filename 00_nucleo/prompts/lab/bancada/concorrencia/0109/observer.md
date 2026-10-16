# Prompt 0109 — observer

**Regime**: Lab, sem promoção para L1–L4  
**Data**: 2026-10-06

Observador puro sobre Go AST em arquivos congelados. Receber schema 1, origem declarada e arquivos path/content/hash; exigir SHA256 e caminhos relativos únicos, até 32 arquivos, 256KiB cada e 2MiB totais. Não usar rede, disco, relógio, tipos carregados ou executar fontes. Preservar identidade por bytes e posição física ignorando diretivas //line. Registrar sintaxe de capacidade literal de canal, envios/recebimentos, formato select/default/Done e go, por função e condições sintáticas. Nomes make sombreados não autorizam atribuição. Closures têm escopo próprio. Não inferir mecanismo, admissão limitada, limite de workers ou conclusão de cancelamento; emitir lacunas explícitas. Parsing parcial deve preservar outros arquivos. Ordem dos arquivos não muda relatório. Autoridade SYNTAX_ONLY e configuração de build apenas declarada.

Identificadores de função incluem posição física e receptor quando presente;
textos/localizadores permanecem UTF-8 válido. Não aceitar conteúdo inválido como
fonte JSON íntegra.

## Evolução 0110, versão v2

Preservar recibos v1. Agrupar ASTs em file set compartilhado por diretório e
pacote. Checar tipos com go/types em memória, sem importer/IO, tamanhos fixos
amd64; não usar GOARCH declarado para alegar seleção de build. Qualquer erro
ou parsing incompleto no grupo impede novas atribuições por tipos. Emitir
SUPPLIED_PACKAGE_TYPECHECK_FAILED nesse grupo; manter TYPECHECK_NOT_PERFORMED
como lacuna histórica da trilha sintática, esclarecendo seu escopo.
Em make de tipo não literal, requerer objeto builtin make, expressão de tipo,
Underlying chan e chamada tipada em pacote sem erros. Emitir
channel_type_identity=RESOLVED_CHANNEL_TYPE com autoridade
GO_TYPES_SUPPLIED_PACKAGE; channel_buffer_syntax continua SYNTAX_ONLY.
Emitir NAMED_CHANNEL_TYPE_UNRESOLVED para tentativas make(T,...) não resolvidas,
sem marcar slice/map validamente tipados como canal. Preservar identidade por
arquivo/posição/hash, limites globais e Done sem promoção. Esta checagem não
certifica compilabilidade do repositório, build real ou limite de trabalho.

Esclarecimento antes da verificação: na v2 substituir a lacuna histórica
TYPECHECK_NOT_PERFORMED por SYNTAX_PROPERTIES_NOT_TYPECHECKED. A trilha
sintática permanece disponível mesmo em fontes inválidas; a checagem restrita
serve somente para as novas atribuições. Preservar corpus v1; criar cópia v2
com esse contrato explícito. Parsing sem nome de pacote invalida todos os grupos
do mesmo diretório para resolução.

0111 v3: ObserveWithEnvironment aceita ambiente validado externo em memória;
Observe permanece sem ambiente. Report inclui environment_sha256 e autoridade
distinta nas atribuições de tipos. Hash/escopo declara fornecimento de dados,
não autentica fonte nem build. Dependência falha invalida dependentes.

0111: type_diagnostics opcionais localizam falhas de checagem por bytes e
posição física; texto limitado a 300 runes, até 100 por grupo. São diagnóstico
da entrada fornecida e não veredito de compilabilidade do projeto.

0112 v4: Report inclui relations opcionais do contrato flow_relations. Perfis
anteriores mantidos. Relações e lacunas são independentes da força de mecanismo;
status sintático não certifica fluxo. Reportar limites de alcance, execução de
callbacks e identidade de instâncias explicitamente.

0113 v5: callback_relations opcionais separados. Identidade dos parâmetros,
passagem e invocação estrutural conservam CALLBACK_EXECUTION_UNRESOLVED e
não atestam execução. ExecutionScope identifica closure/função de cada uso.

0114 v6: call_site_reachability/defer_order opcionais separam contenção lexical
e alcance limitado do site. Não promover relações/callbacks ou limites globais.

0115 v7: channel_identity inclui audit de declaração, bases e eventos de
write/exposição/capture. Estado de instância sempre inconclusivo neste recorte.

0116 v8: configuration_guards opcionais localizam guarda/predicado/panic antes
da escrita e writes de sentinela por declaração; atomicidade permanece desconhecida.
