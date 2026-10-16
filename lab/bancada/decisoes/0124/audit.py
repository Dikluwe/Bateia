"""Execute bounded upstream mechanisms from collected immutable bytes."""
import ast,hashlib,json,random
from pathlib import Path
from typing import Sequence,Any,Tuple,List
from transformers import AutoTokenizer
from laya.common import build_head,build_sequence,_encode_question_text
P=Path(__file__).resolve().parent
bundle=json.loads((P/'results/bundles.json').read_text())[0]
items={i['path']:i for i in bundle['items']}
code=items['laya/train.py']['content']; tree=ast.parse(code)
node=next(n for n in tree.body if isinstance(n,ast.FunctionDef) and n.name=='split_calibration')
ns={'random':random,'Sequence':Sequence,'Any':Any,'Tuple':Tuple,'List':List}
exec(compile(ast.Module(body=[node],type_ignores=[]),'collected:laya/train.py','exec'),ns)
data=[{'family':'one_repository','example_id':i} for i in range(10)]
train,cal=ns['split_calibration'](data,500,0.2,42)
overlap=sorted({v['family'] for v in train}&{v['family'] for v in cal})
assert overlap==['one_repository']
ck=Path(json.loads((P.parent/'0123/results/checkpoint.json').read_text())['local_path'])
tok=AutoTokenizer.from_pretrained(ck/'tokenizer',local_files_only=True)
q0=json.loads((P/'controls.json').read_text())['question']['declaration_relation']
q={'t':'choice','ins':q0['instructions']+' '+('Use the supplied source only. '*150),'crit':q0['criteria']}
full=_encode_question_text(tok,'choice question: '+q['ins'],add_special_tokens=False)
ids,markers,stats=build_head(tok,q,192)
used=ids[1:ids.index(tok.sep_token_id)]
sequence=build_sequence(tok,'Short source.',q,max_len=512,head_max_len=192,return_truncation_stats=True)
state_stats=sequence[-1]
assert len(used)<len(full) and not state_stats['truncated']
report={'scope':'AUTHORED_MECHANISM_CONTROLS_NOT_TRAINING_OR_ACCURACY','split_calibration':{'source_path':'laya/train.py','source_hash':items['laya/train.py']['content_hash'],'first_line':node.lineno,'last_line':node.end_lineno,'train_count':len(train),'calibration_count':len(cal),'overlapping_families':overlap,'interpretation':'Row shuffle is not repository-family isolation; expected implementation behavior, not proven defect.'},'question_budget':{'source_path':'laya/common.py','source_hash':items['laya/common.py']['content_hash'],'full_instruction_tokens':len(full),'used_instruction_tokens':len(used),'head_stats':stats,'state_stats':state_stats,'interpretation':'State truncation flag does not attest full instructions; audit the head separately.'}}
(P/'results/mechanism-controls.json').write_text(json.dumps(report,indent=2)+'\n')
print(json.dumps(report))
