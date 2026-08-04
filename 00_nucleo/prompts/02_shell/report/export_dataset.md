# Specification: Dataset Exporter (L2 Shell)

## Purpose
Serialize domain SearchReport into structured JSON dataset and SQL DDL/DML script for SQLite database insertion.

## Rules
- L2 package (`02_shell/report`).
- Pure presentation formatting.
- Functions: `ExportJSONReport(report *domain.SearchReport) ([]byte, error)` and `ExportSQLiteSchema(report *domain.SearchReport) string`.
