// Crystalline Lineage
// @prompt 00_nucleo/prompts/03_infra/clang/flowtool/main.md
// @layer L3
// @updated 2026-09-07
// Bateia step 0038. Source-level facts only; no execution or feasibility claims.
#include <clang/AST/RecursiveASTVisitor.h>
#include <clang/Analysis/CFG.h>
#include <clang/Frontend/CompilerInstance.h>
#include <clang/Frontend/FrontendActions.h>
#include <clang/Lex/PPCallbacks.h>
#include <clang/Tooling/Tooling.h>
#include <clang/Tooling/CompilationDatabase.h>
#include <clang/Basic/Version.h>
#include <llvm/Support/JSON.h>
#include <llvm/Support/raw_ostream.h>
using namespace clang;
using namespace clang::tooling;
namespace {
llvm::json::Array functions, skipped;
llvm::json::Object location(SourceRange range, SourceManager &sm) {
 auto a=sm.getExpansionLoc(range.getBegin()), b=sm.getExpansionLoc(range.getEnd());
 return llvm::json::Object{{"file",sm.getFilename(a).str()},{"line",(int64_t)sm.getExpansionLineNumber(a)},
 {"end_line",(int64_t)sm.getExpansionLineNumber(b)},
 {"spelling_file",sm.getFilename(sm.getSpellingLoc(range.getBegin())).str()},
 {"spelling_line",(int64_t)sm.getSpellingLineNumber(range.getBegin())}};
}
class PreprocessorTrace : public PPCallbacks {
 SourceManager &sm;
public:
 explicit PreprocessorTrace(SourceManager &s):sm(s){}
 void SourceRangeSkipped(SourceRange r, SourceLocation) override {
  if(sm.isWrittenInMainFile(sm.getExpansionLoc(r.getBegin()))) skipped.push_back(location(r,sm));
 }
};
class Visitor : public RecursiveASTVisitor<Visitor> {
 ASTContext &ctx;
public:
 explicit Visitor(ASTContext &c):ctx(c){}
 bool VisitFunctionDecl(FunctionDecl *f) {
  auto &sm=ctx.getSourceManager();
  if(!f->doesThisDeclarationHaveABody() || !sm.isWrittenInMainFile(sm.getExpansionLoc(f->getLocation()))) return true;
  CFG::BuildOptions options; options.PruneTriviallyFalseEdges=false;
  auto cfg=CFG::buildCFG(f,f->getBody(),&ctx,options);
  auto obj=location(SourceRange(f->getLocation(),f->getEndLoc()),sm);
  obj["name"]=f->getNameAsString();obj["status"]=cfg?"complete":"not_evaluated";
  llvm::json::Array blocks;
  if(cfg) for(auto *block:*cfg) {
   llvm::json::Array succs, statements;
   for(auto it=block->succ_begin();it!=block->succ_end();++it) {
    if(auto *next=it->getReachableBlock()) succs.push_back((int64_t)next->getBlockID());
    else succs.push_back(nullptr);
   }
   for(auto &element:*block) if(auto stmt=element.getAs<CFGStmt>()) {
    const Stmt *s=stmt->getStmt();auto item=location(s->getSourceRange(),sm);
    item["kind"]=s->getStmtClassName();
    if(auto *call=dyn_cast<CallExpr>(s)) {
     auto *callee=call->getDirectCallee();item["callee"]=callee?callee->getNameAsString():"";
     item["resolution"]=callee?"direct":"unresolved";
    }
    statements.push_back(std::move(item));
   }
   llvm::json::Object b{{"id",(int64_t)block->getBlockID()},{"successors",std::move(succs)},
     {"statements",std::move(statements)}};
   if(auto *term=block->getTerminatorStmt()) {
    auto t=location(term->getSourceRange(),sm);t["kind"]=term->getStmtClassName();
    b["terminator"]=std::move(t);
   }
   blocks.push_back(std::move(b));
  }
  if(cfg) {obj["entry"]=(int64_t)cfg->getEntry().getBlockID();obj["exit"]=(int64_t)cfg->getExit().getBlockID();}
  obj["blocks"]=std::move(blocks); functions.push_back(std::move(obj));return true;
 }
};
class Consumer : public ASTConsumer {
 void HandleTranslationUnit(ASTContext &ctx) override {Visitor(ctx).TraverseDecl(ctx.getTranslationUnitDecl());}
};
class Action : public ASTFrontendAction {
 std::unique_ptr<ASTConsumer> CreateASTConsumer(CompilerInstance &ci, StringRef) override {
  ci.getPreprocessor().addPPCallbacks(std::make_unique<PreprocessorTrace>(ci.getSourceManager()));
  return std::make_unique<Consumer>();
 }
};
}
int main(int argc,const char **argv) {
 if(argc<3) return 2;
 std::vector<std::string> args;for(int i=2;i<argc;i++) args.push_back(argv[i]);
 FixedCompilationDatabase db(".",args);ClangTool tool(db,{argv[1]});
 int status=tool.run(newFrontendActionFactory<Action>().get());
 if(status) return status;
 llvm::outs()<<llvm::json::Value(llvm::json::Object{{"schema","bateia-flow-v1"},
 {"version",CLANG_VERSION_STRING},{"status","complete"},{"scope","intraprocedural-main-file"},
 {"prune_trivially_false_edges",false},{"functions",std::move(functions)},{"skipped",std::move(skipped)}})<<"\n";
}
