# Passo 0079 — Container durável em L3

## Objetivo

Transferir para `03_infra/runcontainer` a construção coordenada de store, fila e serviço
de execuções. `04_wiring` deve somente montar adaptadores e controlar o ciclo de vida.

## Contrato

O container expõe o serviço e as operações `Recover`, `Start` e `Shutdown`. Falhas de
construção ou inicialização liberam os recursos já abertos. O encerramento tenta fechar
fila e store, preservando a precedência de erros.

## Aceite

- existe uma única construção produtiva do agregado;
- recuperação ocorre antes do início da fila;
- o mesmo executor atende caminhos síncrono e assíncrono;
- as fachadas públicas e os formatos persistidos permanecem compatíveis.
