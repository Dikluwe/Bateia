# Parecer preliminar — melhorar OpenCascade ou substituir

## Entrada

O Bateia decompôs a decisão em B-Rep, booleanas/tolerâncias, STEP/IGES,
tesselação, constraints, paralelismo, API moderna, extensibilidade e testes
geométricos reproduzíveis.

Foram executadas duas rodadas públicas no GitHub:

- consultas longas: 2 `SUCCESS`, 10 `EMPTY`, seis projetos encontrados;
- consultas curtas: 12 `SUCCESS`, com candidatos de extensão e alternativa.

O contraste mostra que as consultas longas eram restritivas demais. Os dois
snapshots foram preservados para não apagar essa evidência.

## Veredito inicial

**Direção recomendada: `EXTEND + COMPOSE`, não substituição imediata.**

O conjunto encontrado não sustenta um substituto único que cubra, com evidência
equivalente, a combinação de B-Rep, interoperabilidade industrial e maturidade
do OCCT. A linha de menor risco é manter o OCCT como núcleo de compatibilidade e
modernizar suas bordas.

## Componentes para melhorar o OCCT

- `bschwind/opencascade-rs`: referência para uma API Rust mais segura;
- `tpaviot/pythonocc-core`, `CadQuery/cadquery` e `AutodeskAILab/occwl`:
  referências para API Python mais expressiva e paramétrica;
- `donalffons/opencascade.js`: referência para distribuição WebAssembly;
- `fougue/mayo`: referência para visualização e conversão CAD desacopladas;
- `MeshInspector/MeshLib` e `cutdigital/mcut`: componentes especializados para
  reparo, interseção e booleanas de malha;
- `manifold-csg`: controle diferencial para robustez de CSG em malhas.

Esses projetos não devem ser tratados automaticamente como prova de qualidade;
o snapshot comprova apenas descoberta e descrição pública.

## Alternativas por mudança de representação

- `ricosjp/truck`: candidato mais direto a kernel CAD moderno em Rust; requer
  auditoria de B-Rep, tolerâncias, formatos e maturidade;
- `hannobraun/fornjot`: B-Rep em Rust, mas a própria descrição informa que não
  está mais em desenvolvimento;
- `leap71/PicoGK` + `LEAP71_ShapeKernel`: opção para geometria implícita/engenharia
  computacional, não substituição transparente de um kernel B-Rep industrial;
- `MeshLib`, `mcut` e `manifold-csg`: alternativas especializadas para malhas e
  CSG, não para toda a superfície STEP/IGES do OCCT;
- `gcad3d/gcad3d`: aplicação CAD paramétrica, não necessariamente um kernel
  substituível como biblioteca.

## Melhorias com maior retorno

1. Criar uma fachada moderna e versionada sobre OCCT, com ownership, erros
   explícitos e unidades/tolerâncias tipadas.
2. Montar um corpus diferencial de falhas geométricas: OCCT versus Manifold,
   MeshLib e MCUT para booleanas/tesselação onde as representações coincidirem.
3. Isolar STEP/IGES em uma fronteira de interoperabilidade testável, evitando
   contaminar a API de domínio com tipos legados.
4. Tornar tesselação, healing e conversão tarefas incrementais, canceláveis e
   paralelizáveis, com cache por hash de shape e parâmetros.
5. Separar kernel, constraints, visualização e automação em componentes para
   permitir substituição gradual.
6. Manter `truck` como trilha experimental: reproduzir o mesmo corpus e medir
   cobertura antes de considerar migração.

## Próximo gate

Para elevar este parecer de descoberta para decisão de reúso, coletar bundles em
commits imutáveis dos candidatos principais e comprovar:

- suporte real a B-Rep e tolerâncias;
- booleanas robustas em casos degenerados;
- STEP/IGES e preservação de metadados;
- licença e atividade de manutenção;
- testes existentes e API pública reutilizável;
- custo de migração a partir de tipos OCCT.

Até esse gate, a decisão permanece **`EXTEND + COMPOSE` com alternativas em
avaliação**, e não “trocar OpenCascade”.
