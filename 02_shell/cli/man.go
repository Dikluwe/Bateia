/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/man.md
 * @layer L2
 * @updated 2026-09-19
 */

package cli

import "fmt"

// GenerateManPage gera o manual no formato roff para a manpage UNIX bateia(1).
func GenerateManPage() string {
	return `.TH BATEIA 1 "August 2026" "Bateia 1.4.0" "User Commands"
.SH NAME
bateia \- Motor estatistico de busca de codigo e analise arquitetural em Go (Tekt v1.4)
.SH SYNOPSIS
.B bateia
[\fIOPTIONS\fR]
.SH DESCRIPTION
.B Bateia
e um motor de busca de codigo e analise arquitetural baseado no Funil de Reducao de 3 Estagios (Descoberta -> Triagem -> Extracao). Permite buscar padroes arquiteturais no GitHub/Sourcegraph de forma idempotente, imutavel por repo@commit e com estrito controle de cota de API.
.SH OPTIONS
.TP
\fB\-query\fR \fISTRING\fR
Termo ou padrao de busca arquitetural.
.TP
\fB\-lang\fR \fISTRING\fR
Filtrar por linguagem de programacao (go, rust, python, typescript).
.TP
\fB\-max\fR \fINUMBER\fR
Numero maximo de repositorios a extrair no estagio 3 (default: 5).
.TP
\fB\-depth\fR \fINUMBER\fR
Profundidade de busca recursiva no grafo de dependencias (default: 1).
.TP
\fB\-output\fR \fIPATH\fR
Caminho do relatorio gerado (default: relatorio_bateia.md).
.TP
\fB\-man\fR
Exibe esta manpage no formato roff UNIX.
.SH AUTHOR
Escrito pela equipe Bateia (Five-Layer Lattice Tekt v1.4).
`
}

// PrintManPage imprime a manpage no stdout.
func PrintManPage() {
	fmt.Print(GenerateManPage())
}
