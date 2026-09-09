/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/03_infra/effects/analyzer_identity.md
 * @layer L3
 * @updated 2026-09-08
 */

package effects

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"runtime"
	"runtime/debug"
	"strconv"

	"bateia/01_core/domain"
)

func analyzerIdentity() domain.EffectAnalyzerIdentity {
	identity := domain.EffectAnalyzerIdentity{Status: "unknown", GoVersion: runtime.Version()}
	if info, ok := debug.ReadBuildInfo(); ok {
		identity.ModulePath = info.Main.Path
		identity.ModuleVersion = info.Main.Version
		for _, setting := range info.Settings {
			switch setting.Key {
			case "vcs.revision":
				identity.VCSRevision = setting.Value
			case "vcs.modified":
				if modified, err := strconv.ParseBool(setting.Value); err == nil {
					identity.VCSModified = &modified
				}
			}
		}
	}
	path := "/proc/self/exe"
	if runtime.GOOS != "linux" {
		var err error
		path, err = os.Executable()
		if err != nil {
			identity.Reason = "analyzer executable path unavailable: " + err.Error()
			return identity
		}
	}
	executable, err := os.Open(path)
	if err != nil {
		identity.Reason = "analyzer executable unavailable: " + err.Error()
		return identity
	}
	defer executable.Close()
	hash := sha256.New()
	if _, err := io.Copy(hash, executable); err != nil {
		identity.Reason = "analyzer executable hash unavailable: " + err.Error()
		return identity
	}
	identity.BinarySHA256 = hex.EncodeToString(hash.Sum(nil))
	identity.Status = "recorded"
	identity.Reason = "SHA-256 observed from the running analyzer executable"
	return identity
}
