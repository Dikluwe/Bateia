/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/parser/polyglot_matrix.md
 * @layer L1
 * @updated 2026-08-04
 */

package parser

import (
	"strings"

	"bateia/01_core/domain"
)

type protocolPattern struct {
	protocol   string
	technology string
	keywords   []string
}

var knownProtocols = []protocolPattern{
	{protocol: "gRPC / Protobuf", technology: "gRPC", keywords: []string{"grpc", "protobuf", "proto"}},
	{protocol: "GraphQL API", technology: "GraphQL", keywords: []string{"graphql", "apollo"}},
	{protocol: "Message Broker", technology: "Kafka / RabbitMQ / NATS", keywords: []string{"kafka", "rabbitmq", "amqp", "nats"}},
	{protocol: "In-Memory Cache / PubSub", technology: "Redis", keywords: []string{"redis", "valkey"}},
	{protocol: "REST / HTTP Engine", technology: "Express / Gin / Axum", keywords: []string{"gin-gonic", "express", "axum", "actix", "fastapi"}},
}

// EvaluatePolyglotMatrix identifica protocolos e contratos de comunicação entre os repositórios.
func EvaluatePolyglotMatrix(results []domain.ExtractedRepository) []domain.PolyglotContract {
	var matrix []domain.PolyglotContract

	for _, prot := range knownProtocols {
		var matchedRepos []string

		for _, res := range results {
			repoFullName := res.Candidate.Owner + "/" + res.Candidate.Repo
			foundInRepo := false

			if res.Manifest != nil {
				for _, dep := range res.Manifest.Dependencies {
					lowerDep := strings.ToLower(dep.Name)
					for _, kw := range prot.keywords {
						if strings.Contains(lowerDep, kw) {
							foundInRepo = true
							break
						}
					}
					if foundInRepo {
						break
					}
				}
			}

			if foundInRepo {
				matchedRepos = append(matchedRepos, repoFullName)
			}
		}

		if len(matchedRepos) > 0 {
			matrix = append(matrix, domain.PolyglotContract{
				Protocol:     prot.protocol,
				Technology:   prot.technology,
				Repositories: matchedRepos,
			})
		}
	}

	return matrix
}
