# Entrada limitada e publicação de dossiês

**Camada:** L3
**Proprietário:** `03_infra/host/storage.go`
**Contrato:** 0129. Possuir leitura por descritor não bloqueante de arquivo regular limitado, preparação de dossiê sem symlink/alias e escrita por temporário/rename. Extrair o comportamento dos helpers CLI sem alterar erros, permissões, limites, SameFile, ReadDir(1), cleanup ou precedência de erros. Não serializar documentos nem decidir política de reúso.

ReadDescriptor recebe o handle já aberto e validado pela entrada de exportação; não reabre por caminho nem fecha o handle do chamador. ReadBounded abre uma vez e delega a mesma validação por descritor.
