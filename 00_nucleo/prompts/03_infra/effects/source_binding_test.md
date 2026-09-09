# Testes independentes de identidade de fonte — passo 0063

Autor B, contexto novo, sem atestação física; não ler produção candidata, testes existentes, fixtures ou relatórios. Bases: contrato congelado `00_nucleo/passos/0063-vincular-fonte-ao-efeito.md`, regras `.agents/AGENTS.md` e API pública domain.

Criar testes blackbox por API pública do runner. Construir fontes sintéticas nos próprios testes. Calcular SHA-256 dos bytes integrais independentemente. Exercitar pedido sem hash, hash correto e em maiúsculas, rejeição antecipada de hash inválido, mismatch sem execução, alteração exclusivamente em helper, fonte fora do recorte, símbolo ausente e erro de parse. Para fontes aceitas, conferir declaração selecionada e limites byte/linha físicos, preservando Unicode, CRLF e diretivas `//line`. Conferir JSON público incluindo declaração null quando indisponível. Não depender de nomes internos de metadata além de execution.source_sha256 e execution.stop_reason congelados.

Não executar testes durante autoria. Aplicar gofmt e congelar SHA-256 das saídas antes de entregar ao integrador.
