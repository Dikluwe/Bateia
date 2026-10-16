# 0120 — Configuração e inicialização em fontes externas

Avaliação exploratória com observador `go-concurrency-syntax-0116-v8` congelado,
dois projetos externos e seis perguntas. Referências e protocolo foram salvos
antes das respostas. A referência usa declarações originais e interpretação
autoral: não há revisão semântica independente. Nenhuma mudança no observador
ou produto foi feita após os resultados.

## Fontes e separação

| Projeto | Commit congelado | Perguntas |
| --- | --- | --- |
| Jeffail/tunny | a274c3ce48a6dc4f7b5fb5b8eaefbbd9e23574b9 | Reconfiguração durante processamento; New chama SetSize; atomicidade universal |
| alitto/pond | 7b7c16ad4903c9d7014beac3ed3d2533b410f170 | Reconfiguração durante uso; concordância de declarações sobre zero; atomicidade universal |

Snapshots preservam documentação, testes e licença; requisições fornecem apenas
Go de produção, sem testes/examples/internal/assert. tunny: dois arquivos;
pond: doze arquivos. Imports fornecidos em memória e exports padrão reutilizados
de 0111, Go 1.26.5 linux/amd64. Seleção não certifica checkout/build completo.
Código dos dois candidatos não foi executado.

Cada projeto tem perfis completo, sem exports/dependências e fragmento de fonte.
São seis entradas e dezoito linhas pergunta/perfil, mas apenas dois projetos e
seis perguntas únicas. Variantes não contam como amostras independentes. conc
é controle positivo de calibração e fica fora dos denominadores externos.

O congelamento de referências antes das respostas reduz ajuste retrospectivo,
mas o autor conhece o observador. Não equivale a avaliação cega independente.
Contratos em README ficam na referência e não são consumidos pelo observador;
isso limita a cobertura possível e é explícito no diagnóstico.

## Resultados

| Medida nos perfis completos | Resultado | Significado |
| --- | --- | --- |
| Relação estática New → SetSize | 1/1 localizada | Chamada e destino com linhas/hashes correspondentes |
| Perguntas de contrato respondidas | 0/3 | O observador não emite vereditos contratuais |
| Perguntas de atomicidade sem promoção | 2/2 | Nenhuma prova universal foi emitida; não conta como acurácia |
| Guarda restrita de conc, calibração | 1, atomicity=UNRESOLVED | Controle positivo separado, não melhora métricas externas |

Nos dois projetos externos nenhuma configuration_guard do padrão restrito foi
emitida. Ausência de fato não significa ausência de proteção nem ausência de
reconfiguração suportada. Os recortes completos têm zero diagnósticos de tipos:
tunny contém oito relações de chamada estática e pond 58. Os quatro perfis
incompletos registram falha de checagem de tipos e zero relações de chamada.

| Perfil | Diagnósticos de tipos |
| --- | --- |
| tunny completo | 0 |
| tunny sem exports | 5 |
| tunny fragmento | 6 |
| pond completo | 0 |
| pond sem exports | 17 |
| pond fragmento | 31 |

O status dos seis resultados é OBSERVED_SYNTAX. Esse status não responde às
perguntas semânticas. Precisão/recall semânticos ficam nulos, sem denominador
representativo ou oráculo independente. Zero promoção nesses casos não é prova
de ausência geral de extrapolação. O piloto ainda não oferece classificação
contratual; as respostas ausentes são lacuna de capacidade, não três erros
semânticos nem três abstenções corretas.

## Insights de referência

tunny declara que SetSize pode ser usado durante processamento, ressalvando
pool parado. pond apresenta Resize como ajuste durante uso. Isso impede
generalizar o contrato de conc a todo pool.

Em pond, README admite zero, comentário da interface exige maior que zero e
implementação transforma zero em math.MaxInt. A referência registra divergência
de declarações com escopos próprios. É descoberta manual; o observador não a
detectou e nenhum defeito de execução foi demonstrado.

Esses resultados priorizam um contrato futuro para relacionar declarações de
API, observações de código e expectativas de teste, distinguindo cada autoridade.
Novas regras devem ser avaliadas em outro conjunto reservado; estes casos agora
são diagnóstico conhecido e não podem virar teste cego de generalização.

## Reprodução

Com binário cujo hash corresponde a testdata/observer-freeze.json:

```sh
python3 lab/bancada/concorrencia/0120/evaluate.py
python3 lab/bancada/concorrencia/0120/verify.py
```

evaluate.py roda cada entrada duas vezes e compara bytes, confere hashes e
localizadores e preserva a calibração separada. verify.py verifica a contabilidade
e cinco controles negativos: inflar número de projetos, cobertura contratual,
independência, sucesso semântico por abstenção e contagem fictícia de promoção.
São controles do avaliador, não amostras semânticas adicionais.

collect.py e prepare.py recusam substituir snapshots/referências existentes.
O protocolo usa “site/target” como conceitos; os campos do schema são from/to.
Um ajuste inicial no leitor corrigiu esses nomes. O verificador também corrigiu
a suposição de status INCONCLUSIVE para OBSERVED_SYNTAX. Nenhuma referência ou
resposta do observador foi ajustada para melhorar pontuação.

Resultados em results/evaluation.json, seis relatórios individuais,
conc-calibration.json, integrity-controls.json e verification.json. Referências
em testdata/reference.json; hashes do observador em observer-freeze.json.
