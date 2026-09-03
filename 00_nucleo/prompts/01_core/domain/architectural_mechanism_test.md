# Architectural mechanism domain checks

**Camada:** L1  
**Arquivo proprietário:** `01_core/domain/architectural_mechanism_test.go`

Verifique que os enums públicos preservam as strings contratadas e que slices
vazias de um mecanismo são serializadas como arrays, nunca `null`.
