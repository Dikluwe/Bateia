#!/usr/bin/env python3
"""0048 frozen maintenance criteria from spec; no expected labels altered."""
import json,pathlib
root=pathlib.Path(__file__).resolve().parent
rows={}
for corpus in ['real0046','holdout']:
 a=json.loads((root/corpus/'baseline.json').read_text())['aggregate'];b=json.loads((root/corpus/'final.json').read_text())['aggregate']
 rows[corpus]={'baseline':a,'candidate':b,'no_recall_loss':all(b['at_'+str(k)]['macro_recall_labeled']>=a['at_'+str(k)]['macro_recall_labeled'] for k in [3,12]),'no_negative_increase':b['negative_with_indications']<=a['negative_with_indications']}
a=json.loads((root/'control0045-baseline.json').read_text());b=json.loads((root/'control0045-final.json').read_text())
controls=all(next(r for r in b if r['id']==x['id'])['tp']>=x['tp'] and next(r for r in b if r['id']==x['id'])['returned']<=x['returned'] for x in a if x['id']!='aggregate')
result={'maintenance_passed':controls and all(r['no_recall_loss'] and r['no_negative_increase'] for r in rows.values()),'controls0045_nonregression':controls,'corpora':rows,'limitations':['Convenience corpus and tiny holdout, lexical indications only; A/B without physical isolation attestation.','At_k uses distinct paths within 24-hit budget, not extraction slots actually executed.']}
(root/'final-comparison.json').write_text(json.dumps(result,indent=2)+'\n');print(json.dumps(result,indent=2))
