"""Verify automatically prepared context from frozen bundles; never call a model."""
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
    workspace = root.parents[3]
    binary = str(Path(args.binary).resolve())
    case = json.loads((root/'case.json').read_text())
    bundles = json.loads((root/'exploration-bundles.json').read_text())
    original = (root/'exploration-result.json').read_bytes()
    report = json.loads(original)
    manifest = json.loads((root/'corpus-manifest.json').read_text())
    for path, digest in manifest['upstream_sha256'].items():
        assert sha((workspace/path).read_bytes()) == digest
    assert len(bundles) == 4
    items = {}
    content_bytes = 0
    for bundle in bundles:
        ref = bundle['repository']
        identity = (ref['owner']+'/'+ref['repo'], ref['commit'])
        for item in bundle['items']:
            assert item['ref'] == ref
            assert sha(item['content'].encode()) == item['content_hash'].removeprefix('sha256:')
            key = identity+(item['path'],)
            assert key not in items
            items[key] = item
            content_bytes += len(item['content'].encode())
    assert report['schema_version'] == 'insight-exploration-v1'
    assert report['status'] == 'CONTEXT_PREPARED'
    assert report['novelty_state'] == 'NOT_ASSESSED' and report['effect_state'] == 'UNVERIFIED'
    assert len(report['cards']) == 8
    refs, excerpts, repositories, kinds = {}, set(), set(), set()
    for card in report['cards']:
        assert card['origin'] == 'UNINTERPRETED_EXCERPT'
        source, ref = card['source'], card['reference']
        item = items[(source['repository'], source['commit'], ref['path'])]
        lines = item['content'].removesuffix('\n').split('\n')
        assert ref['content_hash'].removeprefix('sha256:') == item['content_hash'].removeprefix('sha256:')
        assert 1 <= ref['start_line'] <= ref['end_line'] <= min(len(lines), 1024)
        assert ref['end_line']-ref['start_line'] < 16
        assert ref['excerpt'] == '\n'.join(lines[ref['start_line']-1:ref['end_line']])
        assert 0 < len(ref['excerpt'].encode()) <= 2048
        assert ref['excerpt'] not in excerpts
        excerpts.add(ref['excerpt'])
        refs[card['id']] = card
        repositories.add(source['repository'].lower())
        kinds.add(card['kind'])
    assert len(repositories) >= 2 and len(kinds) >= 3
    assert report['diagnostics']['selected_sources'] == len(repositories)
    assert report['diagnostics']['incomplete_bundles'] == 3
    assert 'INCOMPLETE_SOURCE_BUNDLES' in report['gaps']
    assert {task['perspective'] for task in report['tasks']} == {'ALTERNATIVES', 'INVERSION', 'CROSS_DOMAIN', 'RECOMBINATION', 'COUNTEREXAMPLE'}
    for task in report['tasks']:
        assert all(ref in refs for ref in task['card_refs'])
        if task['perspective'] == 'RECOMBINATION':
            assert len({refs[ref]['source']['repository'].lower() for ref in task['card_refs']}) == 2
        if task['perspective'] == 'CROSS_DOMAIN':
            assert 'SEMANTIC_DOMAINS_NOT_ASSESSED' in task['gaps']
    assert all(direction['status'] == 'PLANNED_NOT_EXECUTED' for direction in report['research_directions'])
    command = [binary, 'idea', 'insights', 'explore', '--question', case['question']]
    with tempfile.TemporaryDirectory(prefix='bateia-insights-0136-') as directory:
        input_path, output_path = Path(directory)/'bundles.json', Path(directory)/'protected.json'

        def run(value, output=False, familiar=None):
            input_path.write_text(json.dumps(value, ensure_ascii=False))
            flags = []
            for approach in (familiar if familiar is not None else case['familiar_approaches']):
                flags += ['--familiar', approach]
            return subprocess.run(command+['--bundles', str(input_path)]+flags+(['--output', str(output_path)] if output else []), capture_output=True, timeout=30)

        replay = run(bundles)
        assert replay.returncode == 0, replay.stderr.decode()
        assert replay.stdout == original
        reordered = copy.deepcopy(bundles)
        reordered.reverse()
        for bundle in reordered:
            bundle['items'].reverse()
            for item in bundle['items']:
                item['content_hash'] = 'sha256:'+item['content_hash'].removeprefix('sha256:')
        stable = run(reordered, familiar=list(reversed(case['familiar_approaches'])))
        assert stable.returncode == 0 and stable.stdout == original
        altered = copy.deepcopy(bundles)
        altered[0]['items'][0]['content'] += 'altered'
        output_path.write_bytes(original)
        rejected = run(altered, output=True)
        assert rejected.returncode != 0 and output_path.read_bytes() == original
        empty = subprocess.run(command, capture_output=True, timeout=30)
        assert empty.returncode == 0
        empty_report = json.loads(empty.stdout)
        assert empty_report['status'] == 'SEARCH_NEEDED' and not empty_report['cards']
        assert all(not task['card_refs'] for task in empty_report['tasks'])
    portuguese = json.loads((root/'conc-portuguese-exploration.json').read_text())
    assert 'QUESTION_TERMS_NOT_MATCHED' in portuguese['gaps']
    for name in ('conc-exploration.json', 'conc-portuguese-exploration.json'):
        expected = (root/name).read_bytes()
        variant = json.loads(expected)
        variant_command = [binary, 'idea', 'insights', 'explore', '--question', variant['question'], '--bundles', str(root.parent/'0106/github-bundles.json')]
        for familiar in variant['familiar_approaches']:
            variant_command += ['--familiar', familiar]
        variant_replay = subprocess.run(variant_command, capture_output=True, timeout=30)
        assert variant_replay.returncode == 0 and variant_replay.stdout == expected
    print(json.dumps(dict(
        schema_version='1', result_sha256=sha(original), corpus_items=len(items),
        corpus_content_bytes=content_bytes, selected_excerpt_bytes=sum(len(e.encode()) for e in excerpts),
        reading_cards=8, selected_repositories=len(repositories), evidence_kinds=len(kinds),
        manual_assessments_required=0, manual_observations_required=0,
        replay_byte_identical=True, upstream_artifacts_unchanged=True,
        hashes_and_exact_locations_verified=True, card_refs_resolve=True,
        ordering_and_hash_prefix_stable=True, tampered_content_rejected=True, output_preserved=True,
        empty_corpus_control='SEARCH_NEEDED', cross_language_lexical_gap_explicit=True,
        model_called=False, search_executed=False, candidate_code_executed=False,
        semantic_novelty_assessed=False, insight_utility_measured=False,
    ), ensure_ascii=False, indent=2))


if __name__ == '__main__':
    main()
