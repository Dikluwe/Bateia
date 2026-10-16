"""Freeze an authored contract review against previously collected bytes; no automatic semantic labeling."""
import hashlib
import json
from pathlib import Path
HERE = Path(__file__).resolve().parent
ROOT = HERE.parents[3]
BUNDLE = ROOT / "lab/bancada/insights/0106/github-bundles.json"
bundle = json.loads(BUNDLE.read_text())[0]
items = {f["path"]: f for f in bundle["items"]}
def source(path, first, last, role):
    f = items[path]
    digest = hashlib.sha256(f["content"].encode()).hexdigest()
    assert digest == f["content_hash"]
    segment = "\n".join(f["content"].splitlines()[first - 1:last])
    return {"path": path, "first_line": first, "last_line": last, "content_sha256": digest, "segment_sha256": hashlib.sha256(segment.encode()).hexdigest(), "role": role}
claims = [
 {"id": "CONFIGURE_BEFORE_GO", "claim": "Configurar antes de enviar tarefas é a ordem demonstrada de uso.", "judgment": "DOCUMENTED_USAGE_ORDER", "evidence": [source("pool/pool_test.go", 15, 22, "SOURCE_EXAMPLE_NOT_EXECUTED"), source("pool/pool.go", 24, 25, "SOURCE_CONTRACT_DECLARATION")], "limit": "Exemplo e declaração não provam todos os cenários de execução."},
 {"id": "CONFIGURE_AFTER_GO", "claim": "Configuração após a primeira chamada Go deve provocar panic.", "judgment": "DECLARED_CONTRACT_WITH_UNEXECUTED_TEST_EXPECTATION", "evidence": [source("pool/pool.go", 24, 25, "SOURCE_CONTRACT_DECLARATION"), source("pool/pool_test.go", 49, 55, "SOURCE_TEST_ASSERTION_NOT_EXECUTED")], "limit": "Teste da fonte foi lido, não executado."},
 {"id": "CONFIGURE_AFTER_WAIT", "claim": "Wait após Go permite reconfigurar o mesmo pool.", "judgment": "CONTRADICTED_BY_UNEXECUTED_TEST_EXPECTATION", "evidence": [source("pool/pool_test.go", 57, 63, "SOURCE_TEST_ASSERTION_NOT_EXECUTED"), source("pool/pool.go", 108, 114, "SOURCE_CONTRACT_DECLARATION")], "limit": "A expectativa explícita é panic; execução real não foi verificada."},
 {"id": "REUSE_AFTER_WAIT", "claim": "O mesmo pool pode receber outro lote depois de Wait.", "judgment": "UNEXECUTED_TEST_EXPECTATION", "evidence": [source("pool/pool_test.go", 129, 145, "SOURCE_TEST_ASSERTION_NOT_EXECUTED"), source("pool/pool.go", 72, 81, "STATIC_SOURCE_RELATION")], "limit": "Reúso de execução não estabelece direito de reconfiguração."},
 {"id": "OVERLAP_SUPPORTED", "claim": "Configuração simultânea com a primeira chamada Go é uso suportado.", "judgment": "NOT_ESTABLISHED_IN_REVIEWED_SOURCES", "evidence": [source("pool/pool.go", 15, 25, "BOUNDED_CONTRACT_REVIEW"), source("pool/pool_test.go", 49, 63, "BOUNDED_CONTRACT_REVIEW")], "limit": "Silêncio não prova suporte nem proibição; nenhuma hipótese sobre operações permitidas foi autenticada."},
 {"id": "GUARD_IMPLIES_ATOMICITY", "claim": "A guarda antes da escrita comprova proteção atômica concorrente.", "judgment": "UNSUPPORTED_INFERENCE", "evidence": [source("pool/pool.go", 91, 104, "STATIC_SOURCE_RELATION"), source("pool/pool.go", 108, 114, "STATIC_SOURCE_RELATION")], "model_reference": "CHECK_THEN_INITIALIZE_THEN_WRITE", "limit": "Check anterior é observado; atomicidade do candidato permanece desconhecida."},
 {"id": "MODEL_IMPLIES_DEFECT", "claim": "A intercalação no modelo demonstra um defeito em conc.", "judgment": "UNSUPPORTED_MODEL_TO_CANDIDATE_TRANSFER", "evidence": [], "model_reference": "CHECK_THEN_INITIALIZE_THEN_WRITE", "limit": "Modelo omite partes do candidato e sua operação permitida não foi estabelecida."},
 {"id": "SENTINEL_IS_SPAWN_COUNT", "claim": "A sentinela tasks indica necessariamente que goroutines já foram iniciadas.", "judgment": "COMMENT_AND_PREDICATE_HAVE_DIFFERENT_SCOPES", "evidence": [source("pool/pool.go", 108, 114, "SOURCE_CONTRACT_DECLARATION"), source("pool/pool.go", 72, 75, "STATIC_SOURCE_RELATION"), source("pool/pool.go", 100, 105, "STATIC_SOURCE_RELATION")], "limit": "Wait também chama init; não foi executado o caso Wait antes de Go. É uma obrigação de investigação, não bug demonstrado."}
]
review = {"schema_version": "contract-review-v1", "review_date": "2026-10-07", "reference_origin": "AUTHOR_REVIEW_NOT_INDEPENDENT", "semantic_judgments_automatically_derived": False, "candidate_code_executed": False, "independent_semantic_accuracy_measured": False, "repository_commit": "5f936abd7ae87036af1f75c95fb9d0daaf00116b", "source_bundle_sha256": hashlib.sha256(BUNDLE.read_bytes()).hexdigest(), "scope": "BOUNDED_FROZEN_SOURCE_REVIEW", "claims": claims}
with (HERE / "claims.json").open("x") as out:
    json.dump(review, out, ensure_ascii=False, indent=2)
    out.write("\n")
print(json.dumps({"claims": len(claims), "reference_origin": review["reference_origin"]}))
