# Launch Checklist

SimpleRAG-Go is a small Retrieval-Augmented Generation system built in Go.

The ingest path accepts plain text and markdown documents, splits them into chunks,
and stores the chunks in PostgreSQL.

The current query path uses PostgreSQL full-text search against chunk content and
document titles. The response returns a stitched answer plus citations for the
matched chunks.

Use this file for manual verification queries:

- Ask about startup order for the API and database.
- Ask what storage systems are required for the local demo.
- Ask what the current query path returns.
