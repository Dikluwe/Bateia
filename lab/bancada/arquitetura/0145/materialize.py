"""One-shot extraction from the owned source snapshot; never rewrite historical evidence."""
from pathlib import Path
ROOT=Path(__file__).resolve().parents[4]

def take(s,name):
 start=s.index('func '+name+'(');end=s.find('\nfunc ',start+1)
 if end<0:end=len(s)
 return s[start:end].strip(),s[:start]+s[end:]
def header(name,layer):
 return '/**\n * Crystalline Lineage\n * @prompt 00_nucleo/prompts/'+name+'.md\n * @layer '+layer+'\n * @updated 2026-10-29\n */\npackage '+('cli' if layer=='L2' else 'usecase')+'\n\n'
def write(name,imports,body,layer='L1'):
 p=ROOT/(name+'.go');assert not p.exists(),p
 p.write_text(header(name,layer)+'import (\n'+''.join('\t"'+v+'"\n' for v in imports)+')\n\n'+body+'\n')

p=ROOT/'01_core/usecase/insight_neighborhood.go';s=p.read_text()
validate,s=take(s,'ValidateInsightNeighborhood');options,s=take(s,'neighborhoodSiteRequest');finish,s=take(s,'finishInsightNeighborhood')
validate=validate.replace('func ValidateInsightNeighborhood(request domain.InsightNeighborhoodRequest) (domain.InsightNeighborhoodRequest, error)', 'func admitInsightNeighborhood(request domain.InsightNeighborhoodRequest) (insightNeighborhoodAdmission, error)')
validate=validate.replace('return request,', 'return insightNeighborhoodAdmission{Request: request},')
validate=validate.replace('if _, err := insightNeighborhoodOrigins(input); err != nil {','origins, err := insightNeighborhoodOrigins(input)\n\tif err != nil {')
validate=validate.replace('return insightNeighborhoodAdmission{Request: request}, nil','return insightNeighborhoodAdmission{Request: request, Origins: origins}, nil')
public='''type insightNeighborhoodAdmission struct {
 Request domain.InsightNeighborhoodRequest
 Origins []insightLeadOrigin
}

// ValidateInsightNeighborhood admits declared history, without authenticating it.
func ValidateInsightNeighborhood(request domain.InsightNeighborhoodRequest) (domain.InsightNeighborhoodRequest,error) {
 admission,err:=admitInsightNeighborhood(request)
 return admission.Request,err
}'''
write('01_core/usecase/insight_neighborhood_plan',['bateia/01_core/domain','encoding/json'],public+'\n\n'+validate+'\n\n'+options)
s=s.replace('// ValidateInsightNeighborhood admits declared history, without authenticating it.\n','')
s=s.replace('request, err := ValidateInsightNeighborhood(request)','admission, err := admitInsightNeighborhood(request)')
mark='\tplan, err := PlanInsightSites(neighborhoodSiteRequest(request), bindings)'
s=s.replace(mark,'\trequest = admission.Request\n'+mark)
start=s.index('\treport := domain.InsightSiteReport{');end=s.index('\tbeginInsightExplorationMemory(&report)',start)
initial=s[start:end].strip().replace('report :=','return',1).replace('request.Input.ID','inputReportID')
s=s[:start]+'\treport := newInsightNeighborhoodReport(plan, request.Input.ID)\n'+s[end:]
start=s.index('\treport.Exploration, err = ExploreInsights');end=s.index('\torigins, _ := insightNeighborhoodOrigins',start)
context_block=s[start:end]
s=s[:start]+'\tif err := prepareInsightNeighborhoodContext(&report); err != nil {\n\t\treturn domain.InsightSiteReport{},err\n\t}\n'+s[end:]
start=s.index('\torigins, _ := insightNeighborhoodOrigins');end=s.index('\tbyID := ',start)
leads=s[start:end].split('\n',1)[1].rstrip()
leads=leads.replace('deriveInsightReferenceLeads(&report,','deriveInsightReferenceLeads(report,').replace('order := rememberInsightReferenceLeads(&report)','return rememberInsightReferenceLeads(report)')
s=s[:start]+'\torder := prepareInsightNeighborhoodLeads(&report,admission.Origins)\n'+s[end:]
start=s.index('\tfinishInsightNeighborhood(&report,');end=s.index('\treturn report, err',start)
s=s[:start]+'\terr = completeInsightNeighborhood(&report,len(pending),failed,ctx.Err())\n'+s[end:]
s=s.replace('\n\t"encoding/json"','').replace('\n\t"errors"','').replace('@updated 2026-10-28','@updated 2026-10-29')
p.write_text(s)
result='''func newInsightNeighborhoodReport(plan domain.InsightSitePlan,inputReportID string) domain.InsightSiteReport {
'''+initial+'''
}
func prepareInsightNeighborhoodContext(report *domain.InsightSiteReport) error {
 var err error
 report.Exploration,err=ExploreInsights(context.Background(),domain.InsightExplorationRequest{SchemaVersion:"1",Question:report.Plan.Request.Discovery.Question,FamiliarApproaches:report.Plan.Request.Discovery.FamiliarApproaches,MaxPassages:report.Plan.Request.Discovery.MaxPassages})
 return err
}
'''+finish+'''
func completeInsightNeighborhood(report *domain.InsightSiteReport,pending int,failed bool,executionErr error) error {
 finishInsightNeighborhood(report,pending,failed)
 var err error
 if report.Status=="FAILED" {err=domain.ErrInsightDiscoveryUnavailable}
 if executionErr!=nil {
  err,report.Status=executionErr,"CANCELLED"
  if errors.Is(err,context.DeadlineExceeded) {report.Status="TIMED_OUT"}
 }
 report.ID="site_discovery_"+insightHash(*report)
 return err
}
'''
write('01_core/usecase/insight_neighborhood_result',['bateia/01_core/domain','context','errors'],result)
p=ROOT/'01_core/usecase/insight_neighborhood_frontier.go';s=p.read_text();origins,s=take(s,'insightNeighborhoodOrigins')
s=s.replace('@updated 2026-10-28','@updated 2026-10-29')+'\nfunc prepareInsightNeighborhoodLeads(report *domain.InsightSiteReport,origins []insightLeadOrigin) map[string]int {\n'+leads+'\n}\n'
p.write_text(s)
write('01_core/usecase/insight_neighborhood_sources',['bateia/01_core/domain'],origins)
p=ROOT/'01_core/usecase/insight_followups.go';s=p.read_text();memory=[]
for name in ['beginInsightExplorationMemory','insightMemoryReadState','rememberInsightReferenceLeads']:
 body,s=take(s,name);memory.append(body)
selection,s=take(s,'sortInsightFollowupQueue')
start=s.index('\t// Reuse the same per-origin');end=s.index('\tmemory := sealInsightMemory',start)
context=s[start:end].replace('return failed, err','return err')
s=s[:start]+'\tif err := prepareInsightFollowupContext(report); err != nil {\n\t\treturn failed,err\n\t}\n'+s[end:]
s=s.replace('\n\t"sort"','').replace('@updated 2026-10-28','@updated 2026-10-29');p.write_text(s)
write('01_core/usecase/insight_followup_selection',['bateia/01_core/domain','sort'],selection)
write('01_core/usecase/insight_followup_context',['bateia/01_core/domain'],'func prepareInsightFollowupContext(report *domain.InsightSiteReport) error {\n'+context+'\nreturn nil\n}')
p=ROOT/'01_core/usecase/insight_exploration_memory.go';s=p.read_text().replace('@updated 2026-10-28','@updated 2026-10-29');p.write_text(s+'\n\n'+'\n\n'.join(memory)+'\n')
# Relocate existing consumers without altering their bodies or names.
for source,name,target,imports in [
 ('insight_neighborhood_test','TestInsightNeighborhoodRejectsTamperedHistoryBeforeEffects','insight_neighborhood_sources_test',['bateia/01_core/domain','encoding/json','testing']),
 ('insight_neighborhood_test','TestInsightNeighborhoodFeedbackRetryAndDryRun','insight_neighborhood_result_test',['bateia/01_core/domain','errors','testing']),
 ('insight_followups_test','TestInsightFollowupsPreservePartialOnCancellationAndShareCardBudget','insight_followup_context_test',['bateia/01_core/domain','context','testing'])]:
 p=ROOT/'01_core/usecase'/(source+'.go');s=p.read_text();body,s=take(s,name);p.write_text(s.replace('@updated 2026-10-28','@updated 2026-10-29'))
 write('01_core/usecase/'+target,imports,body)
