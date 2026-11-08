"""Replay the curated conc/target comparison; do not execute the candidate."""
import argparse
import copy
import hashlib
import json
from pathlib import Path
import subprocess
import tempfile


def sha(data):
    return hashlib.sha256(data).hexdigest()


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument('--binary', required=True)
    args = parser.parse_args()
    root = Path(__file__).resolve().parent
    binary = str(Path(args.binary).resolve())
    request = json.loads((root/'comparison-input.json').read_text())
    original = (root/'comparison-result.json').read_bytes()
    report = json.loads(original)
    assert request['investigation'] == json.loads((root.parent/'0134/source-insight-input.json').read_text())
    baseline = request['target_baseline']
    manifest = []
    for file in sorted(baseline['files'], key=lambda file: file['path']):
        assert file['content_hash'] == 'sha256:'+sha(file['content'].encode())
        manifest.append(dict(path=file['path'], content_hash=file['content_hash']))
    assert baseline['snapshot_sha256'] == 'sha256:'+sha(json.dumps(manifest, ensure_ascii=False, separators=(',', ':')).encode())
    assert report['schema_version'] == 'insight-baseline-comparison-v1'
    assert report['status'] == 'PARTIAL' and report['effect_state'] == 'UNVERIFIED'
    assert report['adoption_disposition'] == 'NOT_AN_ADOPTION_DECISION'
    assert 'hypotheses' not in report
    assert report['target_baseline']['scope'] == 'SUPPLIED_FILES_ONLY'
    assert 'commit' not in report['target_baseline']
    assert report['target_baseline']['snapshot_sha256'] == baseline['snapshot_sha256']
    source = json.loads((root.parent/'0134/source-insight-result.json').read_text())
    assert report['source_report_id'] == source['id']
    assert report['source_analyses'] == source['source_analyses']
    source_obs = {o['id']: o for a in report['source_analyses'] for o in a['observations']}
    target_obs = {o['id']: o for o in report['target_baseline']['observations']}
    assert len(source_obs) == 6 and len(target_obs) == 5
    files = {file['path']: file for file in baseline['files']}
    for observation in target_obs.values():
        assert observation['origin'] == 'DECLARED_NOT_REASSESSED'
        assert observation['id'].startswith('target_obs_')
        for ref in observation['references']:
            file = files[ref['path']]
            assert file['content_hash'] == ref['content_hash']
            assert ref['excerpt'] == '\n'.join(file['content'].removesuffix('\n').split('\n')[ref['start_line']-1:ref['end_line']])
    assert len(report['comparisons']) == 5
    for row in report['comparisons']:
        assert row['declaration_status'] == 'BOTH_DECLARED'
        assert 'TARGET_EFFECT_NOT_VERIFIED' in row['gaps']
        for ref in row['source_observation_refs']:
            assert source_obs[ref]['kind'] == row['kind']
        for ref in row['target_observation_refs']:
            assert target_obs[ref]['kind'] == row['kind']
            assert target_obs[ref]['mechanism_class'] == row['mechanism_class']
    command = [binary, 'idea', 'insights', 'compare', '--input']
    with tempfile.TemporaryDirectory(prefix='bateia-insights-0135-') as directory:
        input_path, output_path = Path(directory)/'input.json', Path(directory)/'protected.json'

        def run(value, output=False):
            input_path.write_text(json.dumps(value, ensure_ascii=False))
            return subprocess.run(command+[str(input_path)]+(['--output', str(output_path)] if output else []), capture_output=True, timeout=30)

        replay = run(request)
        assert replay.returncode == 0, replay.stderr.decode()
        assert replay.stdout == original
        reordered = copy.deepcopy(request)
        b = reordered['target_baseline']
        b['files'].reverse()
        b['observations'].reverse()
        b['snapshot_sha256'] = b['snapshot_sha256'].removeprefix('sha256:')
        for file in b['files']:
            file['content_hash'] = file['content_hash'].removeprefix('sha256:')
        for observation in b['observations']:
            observation['references'].reverse()
            for ref in observation['references']:
                ref['content_hash'] = ref['content_hash'].removeprefix('sha256:')
        stable = run(reordered)
        assert stable.returncode == 0 and stable.stdout == original
        for mode in ('content', 'snapshot', 'line'):
            altered = copy.deepcopy(request)
            if mode == 'content':
                altered['target_baseline']['files'][0]['content'] += 'altered'
            elif mode == 'snapshot':
                altered['target_baseline']['snapshot_sha256'] = '0'*64
            else:
                altered['target_baseline']['observations'][0]['references'][0]['end_line'] = 999999
            output_path.write_bytes(original)
            rejection = run(altered, output=True)
            assert rejection.returncode != 0 and output_path.read_bytes() == original
        without = copy.deepcopy(request)
        without['target_baseline'] = None
        unknown = run(without)
        assert unknown.returncode == 0
        unknown_report = json.loads(unknown.stdout)
        assert unknown_report['status'] == 'INCONCLUSIVE'
        assert all(row['declaration_status'] == 'SOURCE_ONLY' for row in unknown_report['comparisons'])
        changed = copy.deepcopy(request)
        changed['target_baseline']['observations'][0]['statement'] = 'Declaração oposta fornecida sem avaliação semântica.'
        mutation = run(changed)
        assert mutation.returncode == 0
        mutated = json.loads(mutation.stdout)
        assert mutated['id'] != report['id']
        assert mutated['effect_state'] == 'UNVERIFIED'
        assert [r['declaration_status'] for r in mutated['comparisons']] == [r['declaration_status'] for r in report['comparisons']]
    print(json.dumps(dict(
        schema_version='1', result_sha256=sha(original), snapshot_sha256=baseline['snapshot_sha256'],
        replay_byte_identical=True, source_v2_projection_unchanged=True,
        source_observations=6, target_observations=5, paired_dimensions=5,
        target_hashes_and_excerpts_verified=True, ordering_and_prefix_stable=True,
        tampered_content_snapshot_line_rejected=True, output_preserved=True,
        missing_baseline_control='INCONCLUSIVE', changed_statement_changes_identity=True,
        declaration_availability_is_not_semantic_verdict=True,
        effect_verified=False, candidate_code_executed=False,
    ), ensure_ascii=False, indent=2))


if __name__ == '__main__':
    main()
