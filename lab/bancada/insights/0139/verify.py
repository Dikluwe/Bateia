"""Verify the 0139 gate, source lineage and preservation of prior evidence."""
import hashlib
import json
from pathlib import Path

ROOT = Path(__file__).resolve().parents[4]
LAB = Path(__file__).resolve().parent
CHANGED = {
    '01_core/domain/insight_sites.go',
    '01_core/usecase/insight_site_plan.go',
    '01_core/usecase/insight_sites.go',
    '01_core/usecase/insight_reference_collection.go',
    '04_wiring/insight_sites.go',
}
NEW = {
    '01_core/domain/insight_site_reading.go',
    '01_core/domain/insight_site_reading_test.go',
    '01_core/usecase/insight_site_reading.go',
    '01_core/usecase/insight_site_reading_test.go',
    '03_infra/forge/insight_reader.go',
    '03_infra/forge/insight_reader_test.go',
    '04_wiring/insight_site_reading_test.go',
}


def sha(path):
    return hashlib.sha256(path.read_bytes()).hexdigest()


def main():
    baseline = json.loads((LAB/'baseline.json').read_text())
    changed = set()
    for name, digest in baseline.items():
        path = ROOT/name
        assert path.is_file(), name
        if sha(path) != digest:
            changed.add(name)
    assert changed == CHANGED, changed
    current = {str(path.relative_to(ROOT)) for folder in ['01_core','02_shell','03_infra','04_wiring'] for path in (ROOT/folder).rglob('*.go')}
    assert current - baseline.keys() == NEW
    lineage = []
    for name in sorted(CHANGED | NEW):
        path = ROOT/name
        prompt = '00_nucleo/prompts/'+name.removesuffix('.go')+'.md'
        assert '@prompt '+prompt in path.read_text()
        assert '@updated 2026-10-28' in path.read_text()
        assert (ROOT/prompt).is_file()
        lineage.append(dict(source=name,source_sha256=sha(path),prompt=prompt,prompt_sha256=sha(ROOT/prompt)))
    gates = json.loads((LAB/'final/verification.json').read_text())
    assert all(gate['exit_code']==0 for gate in gates.values()), gates
    assert 'audited 378 productive Go files' in (LAB/'final/architecture.log').read_text()
    assert 'No violations found' in (LAB/'final/lint.log').read_text()
    replays = {}
    for step in ['0106','0134','0135','0136','0137','0138']:
        path = LAB/('replay-'+step+'.json')
        replays[step] = json.loads(path.read_text())
    assert replays['0138']['exploration_replay_equal']
    assert replays['0137']['exploration_replay_equal']
    assert replays['0136']['replay_byte_identical']
    assert replays['0135']['replay_byte_identical']
    assert replays['0134']['replay_byte_identical']
    assert replays['0106']['replay_byte_identical']
    (LAB/'lineage.json').write_text(json.dumps(lineage,indent=2)+'\n')
    print(json.dumps(dict(
        status='PASS',contract='0139',gate='make verify',gate_exit_codes={k:v['exit_code'] for k,v in gates.items()},
        productive_files_audited=378,changed_existing_go=sorted(CHANGED),new_go=sorted(NEW),
        preexisting_go_files=sum(name.endswith('.go') for name in baseline),
        preserved_prior_evidence_files=sum(name.startswith('lab/') for name in baseline),
        preexisting_tests_unchanged=True,historical_replays_passed=6,
        gitee_atomgit_http_fixture_passed=True,live_authenticated_availability_verified=False,
        candidate_code_executed=False,novelty_assessed=False,effect_verified=False,
        final_binary_sha256=sha(LAB/'final/bateia'),
    ),ensure_ascii=False,indent=2))


if __name__=='__main__':
    main()
