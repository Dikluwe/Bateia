#!/usr/bin/env python3
import gzip, json, pathlib, sys
p=pathlib.Path(sys.argv[1]);r=json.loads(gzip.decompress(p.read_bytes()) if p.suffix=='.gz' else p.read_bytes())
repo=r['results'][0];e=repo['query_evidence'];unit=next(u for u in repo['semantic']['units'] if u['path']==e['context'])
f=unit['flow']
assert unit['status']=='complete' and f['status']=='complete'
assert f['schema']=='bateia-flow-v1' and f['tool_sha256'] and not f.get('reason')
assert f['prune_trivially_false_edges'] is False
functions={x['name']:x for x in f['functions']}
for name in ['renderSingle','doSortedLayersBlendRendering','renderOneBlendLayer','renderSortedLayersFP','renderSortedLayersNV']:
    assert name in functions,name
for fn in f['functions']:
    ids={b['id'] for b in fn['blocks']}
    if fn['status']=='complete':
        assert fn['entry'] in ids and fn['exit'] in ids
        for b in fn['blocks']:
            assert all(n is None or n in ids for n in b['successors'])
blend=functions['doSortedLayersBlendRendering']
calls={s.get('callee') for b in blend['blocks'] for s in b['statements']}
assert {'initSortedLayersBlendRendering','renderOneBlendLayer','renderSortedLayersFP','renderSortedLayersNV'} <= calls
assert any(b.get('terminator',{}).get('kind')=='ForStmt' for b in blend['blocks'])
md=pathlib.Path(sys.argv[2]).read_text()
assert 'CFG observado' in md and 'não comprovam viabilidade' in md
print(json.dumps({'status':'passed','flow_functions':len(f['functions']),'skipped_ranges':len(f['skipped']),'context':e['context'],'markdown_bytes':len(md.encode())},indent=2))
