# Gates selados

Os gates `step0028_gate`, `step0071_gate`, `step0072_gate` e `step0073_gate` preservam literalmente
os recibos, hashes e caminhos do commit-base
`d02eafefddef0c8fadf6b6751dfed121105963ac`.

Este módulo Go aninhado impede que a suíte corrente reinterprete os caminhos
históricos como se fossem relativos à nova localização. Para reproduzir um gate
selado, use o commit-base indicado no recibo do próprio gate.
