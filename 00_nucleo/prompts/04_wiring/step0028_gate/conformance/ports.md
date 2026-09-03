# Prompt espelhado — ports concretos do O3

Forneça um adapter público determinístico que despache cada operação congelada a um callback independente: extract, start, lifecycle, replay, runtime, tools, refinement, hypothesis, safety, artifact, pagination, schema e full refinement. O verificador deve ligar os subjects reais do candidato a esses callbacks; operação sem binding deve falhar explicitamente e jamais receber `Preserved` por ausência de execução.
