# Bancada0146 — qualidade dos prompts de exploração

Fatia documental de 28 prompts, sem modificar Go produtivo ou testes existentes. A intenção antecedente desta mudança está no [passo0146](../../../../00_nucleo/passos/0146-consolidar-prompts-de-exploracao.md); a consolidação do comportamento é posterior ao código0142–0145 e não reclassifica a autoria histórica.

## Artefatos e reprodução

- `baseline.json`: SHA256 de arquivos existentes no início desta fatia, incluindo produto, testes, regras, L0 e bancadas0142–0145. Não é um baseline de Git HEAD: o checkout já continha trabalho anterior.
- `before/`: cópias integrais dos 28 prompts substituídos, com nomes que codificam o caminho original. Textos históricos continuam consultáveis; não são fontes executáveis disfarçadas em L0.
- `consolidate.py`: materialização documental delimitada dos prompts a partir das obrigações locais declaradas; não escreve Go, regras ou bancadas antigas.
- `validate.py`: verifica links/âncoras, estrutura, vínculo entre critérios e prompts, alterações permitidas e preservação de Go/regras/evidência histórica. Controles negativos artificiais demonstram somente essas guardas estruturais; não julga semântica ou utilidade investigativa.
- `probe_generation.py`: contraexemplo **sintético/offline** usando o binário preservado de0145; cria entrada válida de geração 999999, aplica feedback e tenta ler o resultado. Não representa uma exploração real, não consulta serviço e não altera os artefatos de0145. Executar apenas em diretório sem saídas do ensaio anterior, pois a publicação protege saídas existentes.
- `generation-probe/result.json`: feedback aceito produz geração 1000000 e a próxima leitura rejeita a memória. Identidade do binário e recibos ficam no ensaio; a matriz registra MEM-04 como divergência pendente.
- `final/`: logs e recibo do gate canônico, emitidos para esta árvore sem substituir evidência antiga.

Rodar `python3 lab/bancada/prompts/0146/validate.py` e `BATEIA_VERIFY_OUTPUT=lab/bancada/prompts/0146/final make verify`. O gate ajusta assinatura Git apenas no processo e fixtures podem exigir loopback. Build é escrito nesta bancada em vez de /dev/shm, conforme o runner canônico. Testes verdes existentes não demonstram satisfação de MEM-04 nem dos controles explicitamente classificados como parciais na matriz.

A mudança não avalia ganho de criatividade, novidade ou qualidade de insights. Ela torna essas distinções, os contratos correntes e as lacunas acessíveis aos próximos autores de código/testes.
