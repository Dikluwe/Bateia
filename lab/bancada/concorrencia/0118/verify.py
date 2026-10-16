"""Verify review provenance and anti-promotion invariants; semantic judgments remain authored."""
from copy import deepcopy
import hashlib
import json
from pathlib import Path
HERE = Path(__file__).resolve().parent
ROOT = HERE.parents[3]
BUNDLE_PATH = ROOT / "lab/bancada/insights/0106/github-bundles.json"
MODEL_PATH = HERE.parent / "0117/results/scenario-evaluation.json"
OBSERVATION_PATH = HERE.parent / "0116/results/conc-guard.json"
bundle = json.loads(BUNDLE_PATH.read_text())[0]
items = {f["path"]: f for f in bundle["items"]}
model = json.loads(MODEL_PATH.read_text())
observation = json.loads(OBSERVATION_PATH.read_text())
expected = {
 "CONFIGURE_BEFORE_GO": "DOCUMENTED_USAGE_ORDER",
 "CONFIGURE_AFTER_GO": "DECLARED_CONTRACT_WITH_UNEXECUTED_TEST_EXPECTATION",
 "CONFIGURE_AFTER_WAIT": "CONTRADICTED_BY_UNEXECUTED_TEST_EXPECTATION",
 "REUSE_AFTER_WAIT": "UNEXECUTED_TEST_EXPECTATION",
 "OVERLAP_SUPPORTED": "NOT_ESTABLISHED_IN_REVIEWED_SOURCES",
 "GUARD_IMPLIES_ATOMICITY": "UNSUPPORTED_INFERENCE",
 "MODEL_IMPLIES_DEFECT": "UNSUPPORTED_MODEL_TO_CANDIDATE_TRANSFER",
 "SENTINEL_IS_SPAWN_COUNT": "COMMENT_AND_PREDICATE_HAVE_DIFFERENT_SCOPES"
}
def validate(review):
    assert review["schema_version"] == "contract-review-v1"
    assert review["reference_origin"] == "AUTHOR_REVIEW_NOT_INDEPENDENT"
    assert review["semantic_judgments_automatically_derived"] is False
    assert review["candidate_code_executed"] is False
    assert review["independent_semantic_accuracy_measured"] is False
    assert review["source_bundle_sha256"] == hashlib.sha256(BUNDLE_PATH.read_bytes()).hexdigest()
    assert len(review["claims"]) == len(expected) and {c["id"] for c in review["claims"]} == set(expected)
    references = 0
    for c in review["claims"]:
        assert c["judgment"] == expected[c["id"]]  # Frozen authored judgments, not automatic semantic verification.
        assert c["limit"] and c["claim"]
        for e in c["evidence"]:
            f = items[e["path"]]
            assert f["ref"]["commit"] == review["repository_commit"]
            digest = hashlib.sha256(f["content"].encode()).hexdigest()
            assert digest == f["content_hash"] == e["content_sha256"]
            lines = f["content"].splitlines()
            assert 1 <= e["first_line"] <= e["last_line"] <= len(lines)
            segment = "\n".join(lines[e["first_line"] - 1:e["last_line"]])
            assert hashlib.sha256(segment.encode()).hexdigest() == e["segment_sha256"]
            assert e["role"] in {"SOURCE_EXAMPLE_NOT_EXECUTED", "SOURCE_CONTRACT_DECLARATION", "SOURCE_TEST_ASSERTION_NOT_EXECUTED", "STATIC_SOURCE_RELATION", "BOUNDED_CONTRACT_REVIEW"}
            references += 1
        if "model_reference" in c:
            assert c["model_reference"] in {w["scenario"] for w in model["witnesses"]}
            assert model["candidate_code_executed"] is False and model["reference_origin"] == "AUTHOR_MODEL_NOT_INDEPENDENT"
    assert len(observation["configuration_guards"]) == 1
    assert observation["configuration_guards"][0]["atomicity"] == "UNRESOLVED"
    return references
review = json.loads((HERE / "claims.json").read_text())
references = validate(review)
controls = []
for name, mutate in [
 ("promote_model_to_candidate_defect", lambda r: r["claims"][6].update(judgment="CANDIDATE_DEFECT_VERIFIED")),
 ("assert_overlap_supported_from_silence", lambda r: r["claims"][4].update(judgment="DOCUMENTED_USAGE_ORDER")),
 ("mislabel_candidate_tests_executed", lambda r: r.update(candidate_code_executed=True)),
 ("claim_independent_accuracy", lambda r: r.update(independent_semantic_accuracy_measured=True)),
 ("tamper_source_locator", lambda r: r["claims"][0]["evidence"][0].update(segment_sha256="tampered"))]:
    altered = deepcopy(review);mutate(altered)
    try:
        validate(altered)
    except AssertionError:
        controls.append({"control": name, "rejected": True})
    else:
        raise AssertionError("invalid promotion accepted: " + name)
receipt = {"review_schema": review["schema_version"], "claims_reviewed": len(review["claims"]), "references_integrity_verified": references, "semantic_judgments_automatically_verified": False, "reference_origin": review["reference_origin"], "independent_semantic_accuracy_measured": False, "candidate_code_executed": False, "negative_controls": controls, "claims_sha256": hashlib.sha256((HERE / "claims.json").read_bytes()).hexdigest(), "source_bundle_sha256": review["source_bundle_sha256"], "model_report_sha256": hashlib.sha256(MODEL_PATH.read_bytes()).hexdigest(), "source_observation_sha256": hashlib.sha256(OBSERVATION_PATH.read_bytes()).hexdigest()}
(HERE / "results/review-receipt.json").write_text(json.dumps(receipt, indent=2) + "\n")
print(json.dumps({"claims": len(review["claims"]), "references_verified": references, "negative_controls": len(controls)}))
