#!/usr/bin/env python3
"""Frozen promotion criteria. Does not alter expected labels or observations."""
import json,pathlib,sys
here=pathlib.Path(__file__).resolve().parent
label=sys.argv[1] if len(sys.argv)>1 else 'candidate'
rows={}
for kind in ['real0046','holdout']:
 a=json.loads((here/kind/'baseline.json').read_text())['aggregate'];b=json.loads((here/kind/(label+'.json')).read_text())['aggregate']
 deltas={str(k):b['at_'+str(k)]['macro_recall_labeled']-a['at_'+str(k)]['macro_recall_labeled'] for k in [3,12]}
 rows[kind]={'baseline':a,'candidate':b,'recall_delta':deltas,'no_recall_loss':all(d>=0 for d in deltas.values()),'recall_gain':any(d>0 for d in deltas.values()),'no_negative_increase':b['negative_with_indications']<=a['negative_with_indications']}
a=json.loads((here/'control0045-baseline.json').read_text());b=json.loads((here/('control0045-'+label+'.json')).read_text())
controls=all(next(r for r in b if r['id']==x['id'])['tp']>=x['tp'] and next(r for r in b if r['id']==x['id'])['returned']<=x['returned'] for x in a if x['id']!='aggregate')
passed=rows['real0046']['recall_gain'] and all(r['no_recall_loss'] and r['no_negative_increase'] for r in rows.values()) and controls
out={'label':label,'promotion_passed':passed,'controls0045_nonregression':controls,'corpora':rows,'limitations':['Recall/precision are against frozen manually labeled paths; at_k means distinct paths among first 24 hits, not extraction execution.','Convenience corpus and tiny holdout do not establish generalization. Negative hits are indications, not verified capability.','A/B role separation, no physical independence attestation.']}
path=here/(label+'-comparison.json')
if path.exists(): raise SystemExit('Refusing overwrite '+str(path))
path.write_text(json.dumps(out,indent=2)+'\n');print(json.dumps(out,indent=2))
