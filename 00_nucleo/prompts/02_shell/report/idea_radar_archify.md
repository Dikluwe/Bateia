# Prompt L0 — `02_shell/report/idea_radar_archify.go`

Implemente em L2 o apresentador puro `ExportIdeaRadarArchify(IdeaRadarInput) ([]byte,error)`
conforme o passo/contrato 0069. Recalcule 0066 e 0067; UNKNOWN ou inconsistência falha sem
parcial. Emita somente Archify Dataflow v1, com IDs tipados em hex, ordenação canônica,
observações canônicas como flows e outcomes/probes em cards. Não invente relações, não
mute input, não faça I/O e não execute/importa Node ou Archify. Respeite os limites 0069.
Antes de emitir, imponha no máximo 5 origens únicas, 5 perspectivas e 5 targets; a sexta
falha `render_budget_exceeded`. Use somente rows 0..4 e labels ASCII de até 10 bytes pela
função `short` r1; preserve identidades completas em IDs hex, flow labels e cards.
