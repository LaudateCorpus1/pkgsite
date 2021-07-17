// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated with go generate -run gen_query.go. DO NOT EDIT.

package symbolsearch

// QuerySymbol is used when the search query is only one word, with no dots.
// In this case, the word must match a symbol name and ranking is completely
// determined by the path_tokens.
const QuerySymbol = `
WITH results AS (
	SELECT
			s.name AS symbol_name,
			sd.package_path,
			sd.module_path,
			sd.version,
			sd.name AS package_name,
			sd.synopsis,
			sd.license_types,
			sd.commit_time,
			sd.imported_by_count,
			ssd.package_symbol_id,
			ssd.goos,
			ssd.goarch,
			(ln(exp(1)+imported_by_count)
		* CASE WHEN sd.redistributable THEN 1 ELSE 0.500000 END
		* CASE WHEN COALESCE(has_go_mod, true) THEN 1 ELSE 0.800000 END) AS score
	FROM symbol_search_documents ssd
	INNER JOIN search_documents sd ON sd.unit_id = ssd.unit_id
	INNER JOIN symbol_names s ON s.id = ssd.symbol_name_id
	WHERE s.tsv_name_tokens @@ to_tsquery('symbols', replace($1, '_', '-'))
)
SELECT
	r.package_path,
	r.module_path,
	r.version,
	r.package_name,
	r.synopsis,
	r.license_types,
	r.commit_time,
	r.imported_by_count,
	r.symbol_name,
	r.goos,
	r.goarch,
	ps.type AS symbol_type,
	ps.synopsis AS symbol_synopsis,
	COUNT(*) OVER() AS total
FROM results r
INNER JOIN package_symbols ps ON r.package_symbol_id = ps.id
WHERE r.score > 0.1
ORDER BY
	score DESC,
	commit_time DESC,
	symbol_name,
	package_path
LIMIT $2
OFFSET $3;`
