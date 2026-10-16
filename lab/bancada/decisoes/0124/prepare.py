"""Prepare scoped investigation using actual collected bytes, without fabricated assessments."""
import hashlib,json
from pathlib import Path
P=Path(__file__).resolve().parent
b=json.loads((P/'results/bundles.json').read_text())[0]
for item in b['items']:
    assert item['content_hash'].removeprefix('sha256:')==hashlib.sha256(item['content'].encode()).hexdigest()
request={'schema_version':'1','question':'Como investigar melhorias de contexto e calibração do Laya preservando evidências?','target':{'name':'Laya','context':'Classificação de alegações frente a declarações de APIs','property':'qualidade da decisão documental','desired_effect':'Reduzir erros e tornar limites de entrada e calibração verificáveis','constraints':['Não usar previsões como gabarito','Não promover hipótese a efeito verificado']},'mechanism_classes':['APPEND_ONLY_LOG_RECOVERY','BOUNDED_QUEUE_BACKPRESSURE','CONTENT_ADDRESSED_STORAGE','PORT_ADAPTER_ISOLATION'],'sources':[{'bundle':b,'assessments':[]}],'max_hypotheses':4}
(P/'insights-input.json').write_text(json.dumps(request,indent=2)+'\n')
pilot=json.loads((P.parent/'0122/pilot-requests.json').read_text())['cases'][1]
original=pilot['request']['state']; claim,source=original.split('\nSOURCE: ',1)
variants=[('original',original,'KNOWN_AUTHOR_REFERENCE'),('family_claim','CLAIM: Configuration methods (With*) may reconfigure this pool after its first Go call.\nSOURCE: '+source,'SYNTHETIC_CLAIM_SAME_SOURCE'),('named_source',claim+'\nSOURCE: WithMaxGoroutines will panic if it is used after calling Go() for the first time.','SYNTHETIC_SOURCE_NOT_UPSTREAM'),('family_binding',original+'\nCONTEXT: WithMaxGoroutines is one of the configuration methods named With*.','AUTHOR_EXPLICIT_FAMILY_BINDING_NOT_UPSTREAM_QUOTE'),('missing_source',claim+'\nSOURCE: This package provides a pool.','SYNTHETIC_INSUFFICIENT_SOURCE'),('original_repeat',original,'REPEAT')]
(P/'controls.json').write_text(json.dumps({'question':pilot['request']['questions'],'source':pilot['source'],'expected_author_only':{'original':'contradicts','family_claim':'contradicts','named_source':'contradicts','family_binding':'contradicts','missing_source':'undetermined','original_repeat':'contradicts'},'cases':[{'id':i,'state':s,'origin':o} for i,s,o in variants]},indent=2)+'\n')
print('Prepared actual bundle and six diagnostic controls; assessments deliberately empty.')
