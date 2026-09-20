#!/usr/bin/env python3
"""Capture JSON metrics from go test output or compare captured arms; external to L1."""
import json,re,sys
from pathlib import Path

def capture(path):
    rows=[json.loads(line.split('METRIC0045 ',1)[1]) for line in Path(path).read_text().splitlines() if 'METRIC0045 ' in line]
    if len(rows)!=7: raise SystemExit('expected six cases and aggregate')
    for row in rows:
        row['recall']=row['tp']/row['expected'] if row['expected'] else None
        row['precision']=row['tp']/row['returned'] if row['returned'] else None
    return rows
if sys.argv[1]=='capture':
    Path(sys.argv[3]).write_text(json.dumps(capture(sys.argv[2]),indent=2)+'\n')
elif sys.argv[1]=='compare':
    a=json.loads(Path(sys.argv[2]).read_text()); b=json.loads(Path(sys.argv[3]).read_text())
    aa,bb=a[-1],b[-1]
    negatives=all(x['returned']==0 for x in b if x['id'] in ('absent','invalid'))
    exact=next(x for x in b if x['id']=='exact')['tp']==1
    gain=bb['recall']>=aa['recall'] and (bb['recall']>aa['recall'] or bb['precision']>aa['precision'])
    result={'baseline':aa,'candidate':bb,'no_recall_loss_and_gain':gain,'negative_paths_preserved':negatives,'exact_file_preserved':exact,'passed':gain and negatives and exact,'scope':'six frozen synthetic cases, three projects, top three distinct paths; exact priority also asserted by Go test'}
    print(json.dumps(result,indent=2));sys.exit(0 if result['passed'] else 1)
else: raise SystemExit('capture LOG OUTPUT or compare BASELINE CANDIDATE')
