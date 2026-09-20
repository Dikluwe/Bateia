# Passo 0082 — Modelo operacional opt-in

## Objetivo

Permitir experimentar um limiar de avanço `0.42` sem mudar o comportamento padrão e sem
conceder autoridade de descarte a um modelo não calibrado.

## Contrato

- o modelo fica em `00_nucleo/modelos/screening-operacional-v1.json`;
- a CLI só o usa com `idea recommend --screening-model <arquivo>`;
- parsing é estrito e ocorre antes de qualquer efeito externo;
- os bytes do modelo participam da identidade da execução e da retomada;
- sem a flag, fluxo, identidade e versão histórica permanecem iguais;
- o modelo opt-in nunca emite probabilidade ou `DROP` enquanto permanecer não calibrado.

## Aceite

Os testes da CLI cobrem carregamento, rejeição antecipada, identidade, retomada e
contenção. O modelo não é procurado implicitamente no repositório nem no ambiente.
