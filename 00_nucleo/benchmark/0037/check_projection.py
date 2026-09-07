#!/usr/bin/env python3
"""Verify the archived Coin answer, exact excerpts, and compiler-context isolation."""
import gzip
import json
import pathlib
import sys
md = pathlib.Path(sys.argv[1]).read_text()
p = pathlib.Path(sys.argv[2])
r = json.loads(gzip.decompress(p.read_bytes()) if p.suffix == '.gz' else p.read_bytes())
repo = r['results'][0]
e = repo['query_evidence']
answer = md.split('## Resumo do Funil')[0]
for term in ('SORTED_LAYERS_BLEND', 'depthbits >= 24', 'SORTED_OBJECT_BLEND',
             'sortedlayersblendpasses', 'has_arb_fragment_program',
             'renderSortedLayersFP', 'renderSortedLayersNV'):
    assert term in answer, term
expected = {'handleTransparency','doSortedLayersBlendRendering','initSortedLayersBlendRendering',
            'renderOneBlendLayer','renderSortedLayersFP','renderSortedLayersNV'}
assert expected <= {d['symbol']['name'] for d in e['definitions']}
assert 'glglue_resolve_symbols' not in answer and 'LoadTransparency' not in answer
assert 'Corpo examinado integralmente' not in md
assert '**Commit-base local**: ``' not in md
assert r['language_distribution'][0]['top_similarity'] == -1
root = repo['candidate']['url'].removeprefix('file://').rstrip('/')
files = {root+'/'+f['path']: f['content'] for f in repo['files']}
units = [u for u in repo['semantic']['units'] if u['path'] == e['context']]
assert len(units) == 1 and units[0]['status'] == 'complete'
actual = {(c['caller_id'],c['callee_id'],c['file'],c['line']) for c in units[0]['calls']}
for c in e['calls']:
    assert (c['caller_id'],c['callee_id'],c['file'],c['line']) in actual
for d in e['definitions']:
    lines = files[d['symbol']['file']].split('\n')
    for x in d['excerpts']:
        assert x['content'] == '\n'.join(lines[x['first_line']-1:x['last_line']])
assert len(md.encode()) < 65000
print(json.dumps({'status':'passed','markdown_bytes':len(md.encode()),
                 'answer_bytes':len(answer.encode()),'context':e['context'],
                 'definitions':len(e['definitions']),'calls':len(e['calls'])},indent=2))
