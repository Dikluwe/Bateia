# Evidence coverage domain

Definir o diagnóstico versionado de perdas entre julgamentos humanos e evidências
emitidas. Uma unidade amostral é o par repositório–regra; aliases da mesma regra
permanecem explicáveis, mas não inflam contagens. Causas são fechadas, ordenáveis e
nunca promovem evidência automaticamente.

O relatório preserva versões e hashes das regras, métricas de precisão e recall,
fontes de apoio, alvos ausentes e explicações determinísticas. Ausência de base
suficiente para atribuir causa deve permanecer `GROUND_TRUTH_AMBIGUOUS`.
