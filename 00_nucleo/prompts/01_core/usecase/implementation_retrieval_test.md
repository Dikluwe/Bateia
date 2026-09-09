# Passo 0044 — testes segregados

Stubs puros sem I/O: provedor sem porta retorna não avaliação explícita, erro é explícito, seleção prioriza hit de caminho neutro preserva README e limite, mutação de bytes invalida hit previamente indexado.

Regime A/B, FShared: notattested. Spec: 00_nucleo/passos/0044-recuperar-implementacoes-por-simbolos.md. Proibida leitura de corpos candidatos e oráculos antigos. Congelar hashes antes da primeira execução.

Clarificação pública de integração: limit é teto de arquivos de código; README/manifesto são contexto adicional conforme limites finais existentes. Fixture usa limite1, máximo3 totais com README/go.mod, hit deve vencer concorrente de nome pertinente.
