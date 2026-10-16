# 0109 Observador de concorrência no Lab

**Data**: 2026-10-06  
**Estado**: IMPLEMENTADO E VERIFICADO NO RECORTE SINTÁTICO  
**Contrato**: [passo 0109](../passos/0109-observar-propriedades-de-concorrencia-no-lab.md)  
**Evidência**: [bancada](../../lab/bancada/concorrencia/0109/README.md)

## Resultado

Foi materializado um observador puro Go AST no Lab, com CLI isolada observe
e evaluate. Recebe arquivos/hash/origem/build declarados e produz observações
por propriedade, função, condições sintáticas e posição física. Não usa a
classificação recebida para decidir propriedades nem executa fontes. Nenhum
import produtivo do Lab ou mudança no catálogo foi introduzido.

O recorte localiza forma de buffer, envio/recebimento, select/default, go
statement e candidato Done. Autoridade permanece SYNTAX_ONLY. Limites de
workers/admissão e conclusão de cancelamento permanecem lacunas explícitas.
Hash garante vínculo com bytes; build e origem não são autenticados.

Comentários/literais não fabricam observações; make sombreado não autoriza
buffer; closures têm escopo próprio e //line não corrompe localização física.
Parsing parcial conserva arquivos válidos. Entrada tem limites, integridade e
JSON estrito. Saída nova é publicada completa sem substituir arquivos ou aliases.

## Bancada e fonte real

Os 15 casos de 12 famílias conferiram com perfis e lacunas congelados; testes
de expectativas erradas e observações extras realmente produziram falha.
Essa concordância verifica o contrato sintático no corpus autoral, não uma
taxa de qualidade semântica independente. Variantes não são amostras independentes.

O ensaio real usou dois arquivos congelados de sourcegraph/conc no SHA
5f936abd7ae87036af1f75c95fb9d0daaf00116b, com hashes conferidos ao importar
a evidência do 0106. Produziu sete observações, todas SYNTAX_ONLY, e lacunas
de build/tipos/limites/cancelamento. Replay foi idêntico byte a byte.

O piloto não resolve make(limiter,n), nem transforma o wrapper p.handle.Go
em go statement. Também não interpreta default como rejeição: esse ramo
da fonte pode iniciar worker. Esses resultados evidenciam a necessidade de
resolução de tipos e fluxo antes de estudar propriedades semânticas. Nenhuma
fonte externa foi compilada/executada e nenhum efeito no alvo foi certificado.

## Verificação

Suíte completa go test ./... count=1, testes focados com race, vet, build do
produto e da CLI experimental passaram. Lint cristalino e diff check passaram.
Assinatura Git foi desativada apenas no processo das fixtures; loopback foi
permitido. Builds ficaram em /tmp. Gates selados com módulo próprio não
foram reproduzidos. Nenhum artefato histórico foi reescrito.

## Limites e próximo corte

Não há seleção de build, typecheck, alcance, aliases, fluxo interprocedural,
identidade de recurso ou avaliação da viabilidade de caminhos. O observador
não certifica mecanismo, ausência de defeito ou benefício. A comparação
do corpus não verifica a verdade semântica das condições registradas.

A próxima fatia deve congelar referência revisada separadamente e comparar
propriedades condicionais sobre um conjunto selecionado de canais/operadores.
O piloto permanece no Lab até que essa evidência justifique uma promoção
com contrato proprietário Tekt e compatibilidade do produto.
