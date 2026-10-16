# 0118 — Alegações e contrato de uso de conc

Revisão autoral limitada à fonte congelada de `sourcegraph/conc`, commit
`5f936abd7ae87036af1f75c95fb9d0daaf00116b`, coletada no passo 0106.
O registro `claims.json` liga oito alegações a quinze referências com caminhos,
linhas físicas, hashes do conteúdo e hashes dos trechos. Documentação e testes
do candidato foram lidos; seu código e seus testes não foram executados.

| Alegação | Conclusão da revisão |
| --- | --- |
| Configurar antes de Go | Ordem demonstrada pelo exemplo e documentação |
| Configurar depois de Go | Documentação declara panic; teste lido espera rejeição |
| Reconfigurar depois de Go e Wait | Contraria a expectativa do teste lido |
| Reutilizar depois de Wait | Teste lido espera outro lote de tarefas; não reconfigura |
| Configurar simultaneamente ao primeiro Go é suportado | Não estabelecido nas fontes revisadas |
| Guarda anterior demonstra atomicidade | Inferência sem suporte; observação mantém UNRESOLVED |
| Intercalação do modelo demonstra defeito no candidato | Transferência sem suporte |
| tasks mede goroutines já iniciadas | Comentário e predicado possuem escopos diferentes |

`Wait` chama `init`, fecha `tasks` e reinicia `initOnce`, mas não torna `tasks`
nulo. A guarda consulta `tasks != nil`. Por isso, permissão de reutilização
para tarefas não implica permissão de reconfiguração. O caminho Wait antes de
Go merece uma obrigação própria de validação: não foi executado nesta revisão
e não constitui diagnóstico de defeito.

O recorte não encontrou declaração explícita autorizando ou proibindo a
sobreposição entre configuração e primeiro Go. Silêncio não estabelece nenhuma
das duas. O modelo 0117 omite canais, sync.Once e workers e permanece evidência
restrita a seu próprio contrato.

## Reprodução e alcance da verificação

Executar a partir da raiz:

```sh
python3 lab/bancada/concorrencia/0118/verify.py
```

`prepare.py` registra a revisão autoral e recusa substituir `claims.json`.
`verify.py` confere integridade/localizadores, origem e invariantes do registro.
Cinco controles rejeitam adulteração do trecho, promoção do modelo a defeito,
suporte inferido do silêncio, declaração falsa de execução e declaração falsa
de precisão independente. As classificações esperadas são autorais e congeladas;
esses controles não verificam automaticamente a interpretação semântica nem
medem a qualidade do classificador.

O recibo está em `results/review-receipt.json`; a verificação do projeto e os
hashes estão em `results/verification.json`. Produto e observador v8 permanecem
sem mudança nesta fatia. Referência semântica independente continua pendente.
