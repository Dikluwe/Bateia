"""Run only an authored runtime model and compare normalized witnesses with frozen source observations."""
import hashlib
import json
from pathlib import Path
import subprocess
import sys
HERE = Path(__file__).resolve().parent
ROOT = HERE.parents[3]
OBSERVER = sys.argv[1] if len(sys.argv) == 2 else "/tmp/bateia-observer-0116"
command = ["go", "test", "-json", "./lab/bancada/concorrencia/0117", "-count=1"]
def model_run():
    raw = subprocess.check_output(command, cwd=ROOT)
    witnesses = []
    passed = []
    for line in raw.splitlines():
        event = json.loads(line)
        output = event.get("Output", "")
        if "WITNESS " in output:
            witnesses.append(json.loads(output.split("WITNESS ", 1)[1]))
        if event.get("Action") == "pass" and event.get("Test"):
            passed.append(event["Test"])
    assert len(passed) == 4 and len(witnesses) == 4
    return sorted(witnesses, key=lambda w: w["scenario"]), raw
first, raw = model_run();second, replay = model_run();assert first == second
(HERE / "results/model-test.jsonl").write_bytes(raw)
(HERE / "results/model-replay.jsonl").write_bytes(replay)
expected = {"CONFIGURE_THEN_INITIALIZE": (True, 3), "INITIALIZE_THEN_CONFIGURE": (False, 3), "CHECK_THEN_INITIALIZE_THEN_WRITE": (True, 7), "INITIALIZE_BEFORE_CONCURRENT_CHECK": (False, 0)}
for w in first:
    assert w["scope"] == "AUTHOR_MODEL_ONLY" and w["initialized"]
    assert (w["accepted"], w["limit"]) == expected[w["scenario"]]
args = [OBSERVER, "observe", "--input", str(HERE.parent / "0111/testdata/conc-input.json"), "--environment", str(HERE.parent / "0111/testdata/environment.json")]
source_raw = subprocess.check_output(args, cwd=ROOT)
assert source_raw == subprocess.check_output(args, cwd=ROOT)
source = json.loads(source_raw)
assert len(source["configuration_guards"]) == 1 and source["configuration_guards"][0]["atomicity"] == "UNRESOLVED"
(HERE / "results/source-observation.json").write_bytes(source_raw)
report = {"schema_version": "phase-scenario-evaluation-v1", "reference_origin": "AUTHOR_MODEL_NOT_INDEPENDENT", "model_code_executed": True, "candidate_code_executed": False, "independent_semantic_accuracy_measured": False, "race_free_model_implies_candidate_race_free": False, "same_source_observation_applies_to_sequential_and_concurrent_hypotheses": True, "candidate_supported_usage_contract_reviewed": False, "atomicity_in_candidate_verified": False, "normalized_witness_replay_identical": True, "raw_test_logs_expected_to_differ_in_timing": True, "source_analyzer": source["analyzer"], "source_input_sha256": source["input_sha256"], "environment_sha256": source["environment_sha256"], "witnesses": first, "assumptions": ["The initialization marker and configured limit are represented by separate atomic variables.", "A test-only rendezvous pauses the writer after its guard and before its store.", "The model omits candidate channels, sync.Once, worker execution, callbacks and receiver aliases.", "Only model witnesses are executed; candidate sources are observed statically."]}
(HERE / "results/scenario-evaluation.json").write_text(json.dumps(report, indent=2) + "\n")
print(json.dumps({"model_scenarios": len(first), "all_expected_witnesses_observed": True, "candidate_atomicity": "UNRESOLVED"}))
