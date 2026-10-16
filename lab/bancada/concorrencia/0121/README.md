# 0121 — Documentação ligada a evidência estática

Capacidade experimental pura, separada do observador v8 congelado. Inspect
recebe fonte/ambiente v8 e documentos/associações opcionais e produz uma projeção
documental com o relatório original preservado. Nenhuma dependência produtiva.

## O que faz

Go Doc anexado por AST a funções, métodos e métodos de interfaces ganha texto
exato, hash/localizador físico e dono sintático. Interface e implementação
permanecem entidades distintas. api_owners inventaria donos também sem Doc,
permitindo apontar um README a uma declaração exata.

Para README, bindings precisa apontar documento, intervalo de linhas e owner
Location exata. Caminhos, bytes/hash, dono e posições são validados. Essa
associação recebe CALLER_PROVIDED_BINDING; Go Doc recebe AST_ATTACHED_GO_DOC.
Validar âncora não autentica a interpretação da associação fornecida.

As declarações usam autoridade DOCUMENTATION_TEXT_ONLY. Evidências no corpo
recebem LEXICAL_CONTAINMENT_ONLY: observações, origem das relações de chamada e
guarda/escrita de configuração. Contenção pode incluir closures; não prova
alcance, execução, concordância com texto ou validade de contrato. Não há
classificação de linguagem natural, contradição automática ou garantia de
atomicidade. Arquivos com erro de parsing não geram declarações; falha de
tipos preserva documentação sintática e mantém as lacunas do observador.

## Resultados em fontes conhecidas

| Fonte | Donos AST | Grupos Go Doc | Associações fornecidas | Ocorrências de links |
| --- | --- | --- | --- | --- |
| tunny | 22 | 14 | 1 | 37 |
| pond | 127 | 69 | 2 | 9 |
| conc | 51 | 49 | 0 | 77 |

São 132 grupos documentais anexados automaticamente, três associações externas
fornecidas e 123 ocorrências de links. Um mesmo fato pode aparecer em mais de
uma associação; contagem de links não mede novos fatos ou qualidade semântica.
Os relatórios v8 correspondem integralmente aos recibos anteriores.

WithMaxGoroutines de conc reúne Go Doc com configuração/guarda do corpo,
mantendo atomicity=UNRESOLVED. SetSize de tunny reúne comentário e README
associado à declaração concreta. Em pond, comentário BasePool.Resize e README
sobre zero estão localizados separadamente; o método *pool.Resize recebe outra
associação fornecida. Não houve ligação nominal automática de implementação
de interface. A divergência continua sendo interpretação manual.

Fontes são diagnósticos já conhecidos de 0111/0120. Não constituem holdout,
referência independente ou medida de precisão semântica. Candidatos não foram
executados. As três perguntas contratuais de 0120 continuam sem veredito;
esta fatia prepara evidência rastreável para interpretação futura.

## Uso e reprodução

```sh
go build -o /tmp/bateia-declarations-0121 ./lab/bancada/concorrencia/0121/cmd/declarations
python3 lab/bancada/concorrencia/0121/replay.py
/tmp/bateia-declarations-0121 --input lab/bancada/concorrencia/0121/testdata/pond-request.json
```

--output aceita somente arquivo novo, publicado com temporário 0600/hardlink.
JSON rejeita duplicatas, desconhecidos, múltiplos documentos e nesting excessivo.
Limite serializado CLI 64 MiB por incluir exports; documentos até 32, 256 KiB
cada e 2 MiB totais/serializados, até mil bindings e dez mil donos; soma de
declarações/links até dez mil. Limites de fontes/ambiente v8 permanecem.

replay.py refaz duas execuções por entrada, confere identidade dos resultados,
preservação do observador e documentação em fragmentos com falhas de tipos.
Gera entradas/recibos atuais sem substituir evidência histórica. Testes cobrem
associações erradas, comentários soltos/trailing, tipos/parsing incompletos,
diretivas line, ordem, closures, JSON e publicação exclusiva.

Recibos em results/replay-receipt.json e results/verification.json; projeções
individuais em tunny.json, pond.json e conc.json. observer_report conserva os
fatos antigos; declarations e api_owners são a nova vista documental.
