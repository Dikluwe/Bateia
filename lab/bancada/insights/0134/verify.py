"""Replay the curated source-observation case; never execute candidate code."""
import argparse
import copy
import hashlib
import json
from pathlib import Path
import subprocess
import tempfile


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument('--binary', required=True)
    args = parser.parse_args()
    root = Path(__file__).resolve().parent
    binary = str(Path(args.binary).resolve())
    request = json.loads((root / 'source-insight-input.json').read_text())
    original = (root / 'source-insight-result.json').read_bytes()
    legacy_root = root.parent / '0106'
    legacy_request = json.loads((legacy_root / 'github-insight-input.json').read_text())
    legacy_bytes = (legacy_root / 'github-insight-result.json').read_bytes()
    legacy = json.loads(legacy_bytes)
    assert request['sources'][0]['bundle'] == legacy_request['sources'][0]['bundle']
    assert request['sources'][0]['assessments'] == legacy_request['sources'][0]['assessments']
    report = json.loads(original)
    assert report['schema_version'] == 'insight-source-analysis-v2'
    assert report['status'] == 'PARTIAL'
    assert len(report['source_analyses']) == len(report['hypotheses']) == 1
    analysis, hypothesis = report['source_analyses'][0], report['hypotheses'][0]
    assert analysis['mechanism_ref'] == legacy['mechanisms'][0]['id']
    assert analysis['mechanism_status'] == 'PARTIAL'
    assert analysis['evidence'] == legacy['mechanisms'][0]['evidence']
    assert analysis['counter_evidence'] == legacy['mechanisms'][0]['counter_evidence']
    assert analysis['catalog_definition']['origin'] == 'CATALOG_TEMPLATE'
    assert analysis['catalog_definition']['strategy'] == legacy['mechanisms'][0]['strategy']
    assert analysis['extraction_scope']['languages'] == legacy['mechanisms'][0]['preconditions']
    items = {item['path']: item for item in request['sources'][0]['bundle']['items']}
    for item in items.values():
        assert hashlib.sha256(item['content'].encode()).hexdigest() == item['content_hash'].removeprefix('sha256:')
    observations = analysis['observations']
    assert len(observations) == 6
    assert {o['statement'] for o in observations} == {o['statement'] for o in request['sources'][0]['observations']}
    references = 0
    for observation in observations:
        assert observation['origin'] == 'DECLARED_NOT_REASSESSED'
        for ref in observation['references']:
            item = items[ref['path']]
            assert ref['content_hash'].removeprefix('sha256:') == item['content_hash'].removeprefix('sha256:')
            lines = item['content'].removesuffix('\n').split('\n')
            assert ref['excerpt'] == '\n'.join(lines[ref['start_line']-1:ref['end_line']])
            assert 1 <= ref['start_line'] <= ref['end_line'] <= len(lines)
            assert ref['end_line']-ref['start_line'] < 20
            references += 1
    assert hypothesis['status'] == 'INCONCLUSIVE'
    assert hypothesis['target_state'] == 'UNKNOWN'
    assert hypothesis['adoption_disposition'] == 'NOT_AN_ADOPTION_DECISION'
    assert hypothesis['preconditions'] == [o['statement'] for o in observations if o['kind'] == 'PRECONDITION']
    assert hypothesis['tradeoffs'] == [o['statement'] for o in observations if o['kind'] == 'TRADEOFF']
    assert set(hypothesis['evidence_refs']) == {analysis['mechanism_ref']} | {o['id'] for o in observations}
    command = [binary, 'idea', 'insights', 'analyze', '--input']
    with tempfile.TemporaryDirectory(prefix='bateia-insights-0134-') as directory:
        scratch = Path(directory)
        input_path, output_path = scratch/'input.json', scratch/'protected.json'

        def run(value, output=False):
            input_path.write_text(json.dumps(value, ensure_ascii=False))
            return subprocess.run(command+[str(input_path)]+(['--output', str(output_path)] if output else []), capture_output=True, timeout=30)

        replay = run(request)
        assert replay.returncode == 0, replay.stderr.decode()
        assert replay.stdout == original
        reordered = copy.deepcopy(request)
        reordered['sources'][0]['observations'].reverse()
        for observation in reordered['sources'][0]['observations']:
            observation['references'].reverse()
            for ref in observation['references']:
                ref['content_hash'] = 'sha256:'+ref['content_hash'].removeprefix('sha256:')
        stable = run(reordered)
        assert stable.returncode == 0 and stable.stdout == original
        changed = copy.deepcopy(request)
        changed['sources'][0]['observations'][0]['statement'] += ' Interpretação adicional fornecida.'
        changed_result = run(changed)
        assert changed_result.returncode == 0
        changed_report = json.loads(changed_result.stdout)
        assert changed_report['id'] != report['id']
        assert changed_report['hypotheses'][0]['id'] != hypothesis['id']
        for mode in ('hash', 'line', 'content'):
            invalid = copy.deepcopy(request)
            if mode == 'hash':
                invalid['sources'][0]['observations'][0]['references'][0]['content_hash'] = '0'*64
            elif mode == 'line':
                invalid['sources'][0]['observations'][0]['references'][0]['end_line'] = 999999
            else:
                invalid['sources'][0]['bundle']['items'][0]['content'] += ' altered'
            output_path.write_bytes(original)
            rejected = run(invalid, output=True)
            assert rejected.returncode != 0
            assert output_path.read_bytes() == original
        no_assessments = copy.deepcopy(request)
        no_assessments['sources'][0]['assessments'] = []
        control = run(no_assessments)
        assert control.returncode == 0
        control_report = json.loads(control.stdout)
        assert control_report['status'] == 'INCONCLUSIVE' and not control_report['hypotheses']
        assert len(control_report['source_analyses'][0]['observations']) == 6
        assert control_report['source_analyses'][0]['mechanism_status'] == 'INCONCLUSIVE'
        old = run(legacy_request)
        assert old.returncode == 0 and old.stdout == legacy_bytes
    print(json.dumps(dict(
        schema_version='1', curated_case='sourcegraph/conc@5f936abd7ae87036af1f75c95fb9d0daaf00116b',
        items_verified=len(items), observations_verified=len(observations), references_verified=references,
        source_bytes_and_assessments_unchanged=True, replay_byte_identical=True,
        ordering_and_hash_prefix_stable=True, declaration_changes_causal_ids=True,
        invalid_hash_line_and_content_rejected=True, previous_output_preserved=True,
        no_assessments_control='INCONCLUSIVE_WITH_OBSERVATIONS_WITHOUT_HYPOTHESES',
        legacy_replay_byte_identical=True, report_status='PARTIAL', mechanism_status='PARTIAL',
        hypothesis_status='INCONCLUSIVE', target_state='UNKNOWN',
        candidate_tests_executed=False, semantic_interpretation_reassessed=False,
        result_sha256=hashlib.sha256(original).hexdigest(),
    ), ensure_ascii=False, indent=2))


if __name__ == '__main__':
    main()
