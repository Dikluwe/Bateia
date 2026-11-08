# run_capabilities

Containers internos L1 de descoberta, extração e síntese. A configuração é validada pela capacidade. A extração recebe o run investigativo, agrega mecanismos e lacunas, ordena com as regras existentes, propaga erro sem publicação e não promove evidência. Catálogo copiado na configuração; lifecycle e portas do store continuam no serviço durável.

## Passo 0131
Aplicar o ownership do ADR-0014 e passo 0131. Jornada/decisões em L1, transporte/renderização em L2, efeitos em L3, composição em L4. Preservar os contratos observáveis, formatos, hashes, versões, ordem, retomada e cancelamento; testes seguem o proprietário ou Lab para auditoria geral. Não conservar regra duplicada na casca/runtime.

0131: NewConfiguredMechanismExtractor valida catálogo/digest e retorna o extrator ou ErrIntegrity. Composição L4 não deve executar extração fictícia para validar configuração.

0131: o adaptador do executor de mercado genérico preserva Candidates/Gaps e Status entregues pelo investigador, validando schema/status. A projeção de observações para os documentos não rederiva as conclusões de negócio da porta.
