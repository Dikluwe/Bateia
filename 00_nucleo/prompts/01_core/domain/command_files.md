# 01_core/domain/command_files

**Camada:** L1
**Proprietário:** `01_core/domain/command_files.go`
**Contrato:** passo 0129 e ADR-0012.

Contratos portáveis de arquivos/raízes, flags declarativas e host de comandos. Nenhum tipo os.File/os.Root/syscall, I/O ou relógio concreto. Preservar io/fs como metadado e handles abstratos.

O host também oferece ReadBounded, PrepareDossier e WriteAtomic como operações de capacidade; leitura física, retenção de descritores e publicação pertencem ao adaptador.

O agregado CommandHost é apenas composição das portas estreitas CommandStorage, CommandExports, CommandTraceFiles, CommandWorkspace e CommandEnvironment, além de ExecutionRuntime. Expõe somente operações consumidas por comandos e seus contratos ativos; não replica a API completa do sistema operacional.
