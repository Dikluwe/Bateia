# 0112 — Contrastes de relações

Corpus autoral de perfis completos de relações. Testar padrão positivo e nomes
renomeados; método Go sem goroutine, receive em outro canal/campo, launch em
outro case, chamada dentro de closure e função variável devem retirar vínculos
correspondentes. Defer condicional conserva condição, defer launch conserva
modo DEFERRED; código lexicamente após return continua contenção sem alegar
alcance. Erro de tipos suprime relações. Hash/posição em ambas extremidades e
ordem estável. Fonte real deve ter duas chamadas a WaitGroup.Go mas somente
uma em case de envio do limiter; default permanece ramo distinto. Relações
não promovem WORKER/ADMISSION_BOUND_UNRESOLVED.
