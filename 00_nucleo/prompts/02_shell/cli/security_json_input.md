# Common JSON file admission

readIdeaJSONInput reads via the nonblocking regular-file reader (16 MiB with growth guard), then delegates syntax/structure checks to validateIdeaJSONStructure, itself parser.ValidateBoundedJSON. All typed deserialization happens only after this guard: one JSON value, depth <=64, tokens <=100000, valid UTF8 and no duplicate object member names including escaped/casefold aliases. This is deliberately stricter than encoding/json's last-write-wins behavior.

Use this helper for advanced idea JSON families (replay, rules, discovery/evidence plans, screening, ontology, corpus, catalog). IdeaSpec and strict opt-in model readers may invoke the guard explicitly to preserve their existing error messages. Raw Git text and file/tree byte hashes keep readIdeaInput without JSON normalization. Valid input bytes and outputs/hashes remain unchanged.
