// Code generated by pkg/sql/lexbase/allkeywords. DO NOT EDIT.

package lexbase

var KeywordsCategories = map[string]string{
"abort": "U",
"absolute": "U",
"access": "U",
"action": "U",
"add": "U",
"admin": "U",
"after": "U",
"aggregate": "U",
"all": "R",
"alter": "U",
"always": "U",
"analyse": "R",
"analyze": "R",
"and": "R",
"annotate_type": "C",
"any": "R",
"array": "R",
"as": "R",
"asc": "R",
"asensitive": "U",
"asymmetric": "R",
"as_json": "U",
"at": "U",
"atomic": "U",
"attribute": "U",
"authorization": "T",
"automatic": "U",
"availability": "U",
"backup": "U",
"backups": "U",
"backward": "U",
"before": "U",
"begin": "U",
"between": "C",
"bigint": "C",
"binary": "U",
"bit": "C",
"boolean": "C",
"both": "R",
"box2d": "C",
"bucket_count": "U",
"bundle": "U",
"by": "U",
"cache": "U",
"called": "U",
"cancel": "U",
"cancelquery": "U",
"capabilities": "U",
"capability": "U",
"cascade": "U",
"case": "R",
"cast": "R",
"changefeed": "U",
"char": "C",
"character": "C",
"characteristics": "C",
"check": "R",
"check_files": "U",
"close": "U",
"cluster": "U",
"coalesce": "C",
"collate": "R",
"collation": "T",
"column": "R",
"columns": "U",
"comment": "U",
"comments": "U",
"commit": "U",
"committed": "U",
"compact": "U",
"complete": "U",
"completions": "U",
"concurrently": "R",
"configuration": "U",
"configurations": "U",
"configure": "U",
"conflict": "U",
"connection": "U",
"connections": "U",
"constraint": "R",
"constraints": "U",
"controlchangefeed": "U",
"controljob": "U",
"conversion": "U",
"convert": "U",
"copy": "U",
"cost": "U",
"covering": "U",
"create": "R",
"createdb": "U",
"createlogin": "U",
"createrole": "U",
"cross": "T",
"csv": "U",
"cube": "U",
"current": "U",
"current_catalog": "R",
"current_date": "R",
"current_role": "R",
"current_schema": "R",
"current_time": "R",
"current_timestamp": "R",
"current_user": "R",
"cursor": "U",
"cycle": "U",
"data": "U",
"database": "U",
"databases": "U",
"day": "U",
"deallocate": "U",
"debug_dump_metadata_sst": "U",
"debug_ids": "U",
"debug_pause_on": "U",
"dec": "C",
"decimal": "C",
"declare": "U",
"default": "R",
"defaults": "U",
"deferrable": "R",
"deferred": "U",
"definer": "U",
"delete": "U",
"delimiter": "U",
"depends": "U",
"desc": "R",
"destination": "U",
"detached": "U",
"details": "U",
"discard": "U",
"distinct": "R",
"do": "R",
"domain": "U",
"double": "U",
"drop": "U",
"else": "R",
"encoding": "U",
"encrypted": "U",
"encryption_info_dir": "U",
"encryption_passphrase": "U",
"end": "R",
"enum": "U",
"enums": "U",
"escape": "U",
"except": "R",
"exclude": "U",
"excluding": "U",
"execute": "U",
"execution": "U",
"exists": "C",
"experimental": "U",
"experimental_audit": "U",
"experimental_fingerprints": "U",
"experimental_relocate": "U",
"experimental_replica": "U",
"expiration": "U",
"explain": "U",
"export": "U",
"extension": "U",
"external": "U",
"extract": "C",
"extract_duration": "C",
"extremes": "U",
"failure": "U",
"false": "R",
"family": "T",
"fetch": "R",
"files": "U",
"filter": "U",
"first": "U",
"float": "C",
"following": "U",
"for": "R",
"force": "U",
"force_index": "U",
"force_not_null": "U",
"force_null": "U",
"force_quote": "U",
"force_zigzag": "U",
"foreign": "R",
"format": "U",
"forward": "U",
"freeze": "U",
"from": "R",
"full": "T",
"function": "U",
"functions": "U",
"generated": "U",
"geography": "C",
"geometry": "C",
"geometrycollection": "U",
"geometrycollectionm": "U",
"geometrycollectionz": "U",
"geometrycollectionzm": "U",
"geometrym": "U",
"geometryz": "U",
"geometryzm": "U",
"global": "U",
"goal": "U",
"grant": "R",
"grants": "U",
"greatest": "C",
"group": "R",
"grouping": "C",
"groups": "U",
"hash": "U",
"having": "R",
"header": "U",
"high": "U",
"histogram": "U",
"hold": "U",
"hour": "U",
"identity": "U",
"if": "C",
"iferror": "C",
"ifnull": "C",
"ignore_foreign_keys": "U",
"ilike": "T",
"immediate": "U",
"immutable": "U",
"import": "U",
"in": "R",
"include": "U",
"include_all_secondary_tenants": "U",
"including": "U",
"increment": "U",
"incremental": "U",
"incremental_location": "U",
"index": "U",
"indexes": "U",
"index_after_order_by_before_at": "R",
"index_before_name_then_paren": "R",
"index_before_paren": "R",
"inherits": "U",
"initially": "R",
"inject": "U",
"inner": "T",
"inout": "C",
"input": "U",
"insensitive": "U",
"insert": "U",
"int": "C",
"integer": "C",
"intersect": "R",
"interval": "C",
"into": "R",
"into_db": "U",
"inverted": "U",
"invisible": "U",
"invoker": "U",
"is": "T",
"iserror": "C",
"isnull": "T",
"isolation": "U",
"job": "U",
"jobs": "U",
"join": "T",
"json": "U",
"key": "U",
"keys": "U",
"kms": "U",
"kv": "U",
"label": "U",
"language": "U",
"last": "U",
"lateral": "R",
"latest": "U",
"lc_collate": "U",
"lc_ctype": "U",
"leading": "R",
"leakproof": "U",
"lease": "U",
"least": "C",
"left": "T",
"less": "U",
"level": "U",
"like": "T",
"limit": "R",
"linestring": "U",
"linestringm": "U",
"linestringz": "U",
"linestringzm": "U",
"list": "U",
"local": "U",
"locality": "U",
"localtime": "R",
"localtimestamp": "R",
"locked": "U",
"login": "U",
"lookup": "U",
"low": "U",
"match": "U",
"materialized": "U",
"maxvalue": "U",
"merge": "U",
"method": "U",
"minute": "U",
"minvalue": "U",
"modifyclustersetting": "U",
"modifysqlclustersetting": "U",
"month": "U",
"move": "U",
"multilinestring": "U",
"multilinestringm": "U",
"multilinestringz": "U",
"multilinestringzm": "U",
"multipoint": "U",
"multipointm": "U",
"multipointz": "U",
"multipointzm": "U",
"multipolygon": "U",
"multipolygonm": "U",
"multipolygonz": "U",
"multipolygonzm": "U",
"names": "U",
"nan": "U",
"natural": "T",
"never": "U",
"new_db_name": "U",
"new_kms": "U",
"next": "U",
"no": "U",
"nocancelquery": "U",
"nocontrolchangefeed": "U",
"nocontroljob": "U",
"nocreatedb": "U",
"nocreatelogin": "U",
"nocreaterole": "U",
"nologin": "U",
"nomodifyclustersetting": "U",
"none": "T",
"nonvoters": "U",
"normal": "U",
"nosqllogin": "U",
"not": "R",
"nothing": "U",
"nothing_after_returning": "R",
"notnull": "T",
"noviewactivity": "U",
"noviewactivityredacted": "U",
"noviewclustersetting": "U",
"nowait": "U",
"no_full_scan": "U",
"no_index_join": "U",
"no_zigzag_join": "U",
"null": "R",
"nullif": "C",
"nulls": "U",
"numeric": "C",
"of": "U",
"off": "U",
"offset": "R",
"oids": "U",
"old_kms": "U",
"on": "R",
"only": "R",
"operator": "U",
"opt": "U",
"option": "U",
"options": "U",
"or": "R",
"order": "R",
"ordinality": "U",
"others": "U",
"out": "C",
"outer": "T",
"over": "U",
"overlaps": "T",
"overlay": "C",
"owned": "U",
"owner": "U",
"parallel": "U",
"parent": "U",
"partial": "U",
"partition": "U",
"partitions": "U",
"password": "U",
"pause": "U",
"paused": "U",
"physical": "U",
"placement": "U",
"placing": "R",
"plan": "U",
"plans": "U",
"point": "C",
"pointm": "U",
"pointz": "U",
"pointzm": "U",
"polygon": "C",
"polygonm": "U",
"polygonz": "U",
"polygonzm": "U",
"position": "C",
"preceding": "U",
"precision": "C",
"prepare": "U",
"preserve": "U",
"primary": "R",
"prior": "U",
"priority": "U",
"privileges": "U",
"public": "U",
"publication": "U",
"queries": "U",
"query": "U",
"quote": "U",
"range": "U",
"ranges": "U",
"read": "U",
"real": "C",
"reason": "U",
"reassign": "U",
"recurring": "U",
"recursive": "U",
"redact": "U",
"ref": "U",
"references": "R",
"refresh": "U",
"region": "U",
"regional": "U",
"regions": "U",
"reindex": "U",
"relative": "U",
"release": "U",
"relocate": "U",
"rename": "U",
"repeatable": "U",
"replace": "U",
"replication": "U",
"reset": "U",
"restart": "U",
"restore": "U",
"restrict": "U",
"restricted": "U",
"resume": "U",
"retention": "U",
"retry": "U",
"return": "U",
"returning": "R",
"returns": "U",
"revision_history": "U",
"revoke": "U",
"right": "T",
"role": "U",
"roles": "U",
"rollback": "U",
"rollup": "U",
"routines": "U",
"row": "C",
"rows": "U",
"rule": "U",
"running": "U",
"savepoint": "U",
"scans": "U",
"scatter": "U",
"schedule": "U",
"schedules": "U",
"schema": "U",
"schemas": "U",
"schema_only": "U",
"scroll": "U",
"scrub": "U",
"search": "U",
"second": "U",
"secondary": "U",
"security": "U",
"select": "R",
"sequence": "U",
"sequences": "U",
"serializable": "U",
"server": "U",
"service": "U",
"session": "U",
"sessions": "U",
"session_user": "R",
"set": "U",
"setof": "C",
"sets": "U",
"setting": "U",
"settings": "U",
"share": "U",
"shared": "U",
"show": "U",
"similar": "T",
"simple": "U",
"skip": "U",
"skip_localities_check": "U",
"skip_missing_foreign_keys": "U",
"skip_missing_sequences": "U",
"skip_missing_sequence_owners": "U",
"skip_missing_udfs": "U",
"skip_missing_views": "U",
"smallint": "C",
"snapshot": "U",
"some": "R",
"split": "U",
"sql": "U",
"sqllogin": "U",
"stable": "U",
"start": "U",
"state": "U",
"statements": "U",
"statistics": "U",
"status": "U",
"stdin": "U",
"stdout": "U",
"stop": "U",
"storage": "U",
"store": "U",
"stored": "U",
"storing": "U",
"stream": "U",
"strict": "U",
"string": "C",
"subscription": "U",
"substring": "C",
"super": "U",
"support": "U",
"survival": "U",
"survive": "U",
"symmetric": "R",
"syntax": "U",
"system": "U",
"table": "R",
"tables": "U",
"tablespace": "U",
"temp": "U",
"template": "U",
"temporary": "U",
"tenant": "U",
"tenants": "U",
"tenant_name": "U",
"testing_relocate": "U",
"text": "U",
"then": "R",
"throttling": "U",
"ties": "U",
"time": "C",
"timestamp": "C",
"timestamptz": "C",
"timetz": "C",
"to": "R",
"trace": "U",
"tracing": "U",
"trailing": "R",
"transaction": "U",
"transactions": "U",
"transfer": "U",
"transform": "U",
"treat": "C",
"trigger": "U",
"trim": "C",
"true": "R",
"truncate": "U",
"trusted": "U",
"type": "U",
"types": "U",
"unbounded": "U",
"uncommitted": "U",
"union": "R",
"unique": "R",
"unknown": "U",
"unlisten": "U",
"unlogged": "U",
"unsafe_restore_incompatible_version": "U",
"unset": "U",
"unsplit": "U",
"until": "U",
"update": "U",
"upsert": "U",
"use": "U",
"user": "R",
"users": "U",
"using": "R",
"valid": "U",
"validate": "U",
"value": "U",
"values": "C",
"varbit": "C",
"varchar": "C",
"variadic": "R",
"varying": "U",
"verify_backup_table_data": "U",
"view": "U",
"viewactivity": "U",
"viewactivityredacted": "U",
"viewclustermetadata": "U",
"viewclustersetting": "U",
"viewdebug": "U",
"virtual": "C",
"visible": "U",
"volatile": "U",
"voters": "U",
"when": "R",
"where": "R",
"window": "R",
"with": "R",
"within": "U",
"without": "U",
"work": "C",
"write": "U",
"year": "U",
"zone": "U",
}

// KeywordNames contains all keywords sorted, so that pg_get_keywords returns
// deterministic results.
var KeywordNames = []string{
"abort",
"absolute",
"access",
"action",
"add",
"admin",
"after",
"aggregate",
"all",
"alter",
"always",
"analyse",
"analyze",
"and",
"annotate_type",
"any",
"array",
"as",
"asc",
"asensitive",
"asymmetric",
"as_json",
"at",
"atomic",
"attribute",
"authorization",
"automatic",
"availability",
"backup",
"backups",
"backward",
"before",
"begin",
"between",
"bigint",
"binary",
"bit",
"boolean",
"both",
"box2d",
"bucket_count",
"bundle",
"by",
"cache",
"called",
"cancel",
"cancelquery",
"capabilities",
"capability",
"cascade",
"case",
"cast",
"changefeed",
"char",
"character",
"characteristics",
"check",
"check_files",
"close",
"cluster",
"coalesce",
"collate",
"collation",
"column",
"columns",
"comment",
"comments",
"commit",
"committed",
"compact",
"complete",
"completions",
"concurrently",
"configuration",
"configurations",
"configure",
"conflict",
"connection",
"connections",
"constraint",
"constraints",
"controlchangefeed",
"controljob",
"conversion",
"convert",
"copy",
"cost",
"covering",
"create",
"createdb",
"createlogin",
"createrole",
"cross",
"csv",
"cube",
"current",
"current_catalog",
"current_date",
"current_role",
"current_schema",
"current_time",
"current_timestamp",
"current_user",
"cursor",
"cycle",
"data",
"database",
"databases",
"day",
"deallocate",
"debug_dump_metadata_sst",
"debug_ids",
"debug_pause_on",
"dec",
"decimal",
"declare",
"default",
"defaults",
"deferrable",
"deferred",
"definer",
"delete",
"delimiter",
"depends",
"desc",
"destination",
"detached",
"details",
"discard",
"distinct",
"do",
"domain",
"double",
"drop",
"else",
"encoding",
"encrypted",
"encryption_info_dir",
"encryption_passphrase",
"end",
"enum",
"enums",
"escape",
"except",
"exclude",
"excluding",
"execute",
"execution",
"exists",
"experimental",
"experimental_audit",
"experimental_fingerprints",
"experimental_relocate",
"experimental_replica",
"expiration",
"explain",
"export",
"extension",
"external",
"extract",
"extract_duration",
"extremes",
"failure",
"false",
"family",
"fetch",
"files",
"filter",
"first",
"float",
"following",
"for",
"force",
"force_index",
"force_not_null",
"force_null",
"force_quote",
"force_zigzag",
"foreign",
"format",
"forward",
"freeze",
"from",
"full",
"function",
"functions",
"generated",
"geography",
"geometry",
"geometrycollection",
"geometrycollectionm",
"geometrycollectionz",
"geometrycollectionzm",
"geometrym",
"geometryz",
"geometryzm",
"global",
"goal",
"grant",
"grants",
"greatest",
"group",
"grouping",
"groups",
"hash",
"having",
"header",
"high",
"histogram",
"hold",
"hour",
"identity",
"if",
"iferror",
"ifnull",
"ignore_foreign_keys",
"ilike",
"immediate",
"immutable",
"import",
"in",
"include",
"include_all_secondary_tenants",
"including",
"increment",
"incremental",
"incremental_location",
"index",
"indexes",
"index_after_order_by_before_at",
"index_before_name_then_paren",
"index_before_paren",
"inherits",
"initially",
"inject",
"inner",
"inout",
"input",
"insensitive",
"insert",
"int",
"integer",
"intersect",
"interval",
"into",
"into_db",
"inverted",
"invisible",
"invoker",
"is",
"iserror",
"isnull",
"isolation",
"job",
"jobs",
"join",
"json",
"key",
"keys",
"kms",
"kv",
"label",
"language",
"last",
"lateral",
"latest",
"lc_collate",
"lc_ctype",
"leading",
"leakproof",
"lease",
"least",
"left",
"less",
"level",
"like",
"limit",
"linestring",
"linestringm",
"linestringz",
"linestringzm",
"list",
"local",
"locality",
"localtime",
"localtimestamp",
"locked",
"login",
"lookup",
"low",
"match",
"materialized",
"maxvalue",
"merge",
"method",
"minute",
"minvalue",
"modifyclustersetting",
"modifysqlclustersetting",
"month",
"move",
"multilinestring",
"multilinestringm",
"multilinestringz",
"multilinestringzm",
"multipoint",
"multipointm",
"multipointz",
"multipointzm",
"multipolygon",
"multipolygonm",
"multipolygonz",
"multipolygonzm",
"names",
"nan",
"natural",
"never",
"new_db_name",
"new_kms",
"next",
"no",
"nocancelquery",
"nocontrolchangefeed",
"nocontroljob",
"nocreatedb",
"nocreatelogin",
"nocreaterole",
"nologin",
"nomodifyclustersetting",
"none",
"nonvoters",
"normal",
"nosqllogin",
"not",
"nothing",
"nothing_after_returning",
"notnull",
"noviewactivity",
"noviewactivityredacted",
"noviewclustersetting",
"nowait",
"no_full_scan",
"no_index_join",
"no_zigzag_join",
"null",
"nullif",
"nulls",
"numeric",
"of",
"off",
"offset",
"oids",
"old_kms",
"on",
"only",
"operator",
"opt",
"option",
"options",
"or",
"order",
"ordinality",
"others",
"out",
"outer",
"over",
"overlaps",
"overlay",
"owned",
"owner",
"parallel",
"parent",
"partial",
"partition",
"partitions",
"password",
"pause",
"paused",
"physical",
"placement",
"placing",
"plan",
"plans",
"point",
"pointm",
"pointz",
"pointzm",
"polygon",
"polygonm",
"polygonz",
"polygonzm",
"position",
"preceding",
"precision",
"prepare",
"preserve",
"primary",
"prior",
"priority",
"privileges",
"public",
"publication",
"queries",
"query",
"quote",
"range",
"ranges",
"read",
"real",
"reason",
"reassign",
"recurring",
"recursive",
"redact",
"ref",
"references",
"refresh",
"region",
"regional",
"regions",
"reindex",
"relative",
"release",
"relocate",
"rename",
"repeatable",
"replace",
"replication",
"reset",
"restart",
"restore",
"restrict",
"restricted",
"resume",
"retention",
"retry",
"return",
"returning",
"returns",
"revision_history",
"revoke",
"right",
"role",
"roles",
"rollback",
"rollup",
"routines",
"row",
"rows",
"rule",
"running",
"savepoint",
"scans",
"scatter",
"schedule",
"schedules",
"schema",
"schemas",
"schema_only",
"scroll",
"scrub",
"search",
"second",
"secondary",
"security",
"select",
"sequence",
"sequences",
"serializable",
"server",
"service",
"session",
"sessions",
"session_user",
"set",
"setof",
"sets",
"setting",
"settings",
"share",
"shared",
"show",
"similar",
"simple",
"skip",
"skip_localities_check",
"skip_missing_foreign_keys",
"skip_missing_sequences",
"skip_missing_sequence_owners",
"skip_missing_udfs",
"skip_missing_views",
"smallint",
"snapshot",
"some",
"split",
"sql",
"sqllogin",
"stable",
"start",
"state",
"statements",
"statistics",
"status",
"stdin",
"stdout",
"stop",
"storage",
"store",
"stored",
"storing",
"stream",
"strict",
"string",
"subscription",
"substring",
"super",
"support",
"survival",
"survive",
"symmetric",
"syntax",
"system",
"table",
"tables",
"tablespace",
"temp",
"template",
"temporary",
"tenant",
"tenants",
"tenant_name",
"testing_relocate",
"text",
"then",
"throttling",
"ties",
"time",
"timestamp",
"timestamptz",
"timetz",
"to",
"trace",
"tracing",
"trailing",
"transaction",
"transactions",
"transfer",
"transform",
"treat",
"trigger",
"trim",
"true",
"truncate",
"trusted",
"type",
"types",
"unbounded",
"uncommitted",
"union",
"unique",
"unknown",
"unlisten",
"unlogged",
"unsafe_restore_incompatible_version",
"unset",
"unsplit",
"until",
"update",
"upsert",
"use",
"user",
"users",
"using",
"valid",
"validate",
"value",
"values",
"varbit",
"varchar",
"variadic",
"varying",
"verify_backup_table_data",
"view",
"viewactivity",
"viewactivityredacted",
"viewclustermetadata",
"viewclustersetting",
"viewdebug",
"virtual",
"visible",
"volatile",
"voters",
"when",
"where",
"window",
"with",
"within",
"without",
"work",
"write",
"year",
"zone",
}

// GetKeywordID returns the lex id of the SQL keyword k or IDENT if k is
// not a keyword.
func GetKeywordID(k string) int32 {
	// The previous implementation generated a map that did a string ->
	// id lookup. Various ideas were benchmarked and the implementation below
	// was the fastest of those, between 3% and 10% faster (at parsing, so the
	// scanning speedup is even more) than the map implementation.
	switch k {
	case "abort": return ABORT
	case "absolute": return ABSOLUTE
	case "access": return ACCESS
	case "action": return ACTION
	case "add": return ADD
	case "admin": return ADMIN
	case "after": return AFTER
	case "aggregate": return AGGREGATE
	case "all": return ALL
	case "alter": return ALTER
	case "always": return ALWAYS
	case "analyse": return ANALYSE
	case "analyze": return ANALYZE
	case "and": return AND
	case "annotate_type": return ANNOTATE_TYPE
	case "any": return ANY
	case "array": return ARRAY
	case "as": return AS
	case "asc": return ASC
	case "asensitive": return ASENSITIVE
	case "asymmetric": return ASYMMETRIC
	case "as_json": return AS_JSON
	case "at": return AT
	case "atomic": return ATOMIC
	case "attribute": return ATTRIBUTE
	case "authorization": return AUTHORIZATION
	case "automatic": return AUTOMATIC
	case "availability": return AVAILABILITY
	case "backup": return BACKUP
	case "backups": return BACKUPS
	case "backward": return BACKWARD
	case "before": return BEFORE
	case "begin": return BEGIN
	case "between": return BETWEEN
	case "bigint": return BIGINT
	case "binary": return BINARY
	case "bit": return BIT
	case "boolean": return BOOLEAN
	case "both": return BOTH
	case "box2d": return BOX2D
	case "bucket_count": return BUCKET_COUNT
	case "bundle": return BUNDLE
	case "by": return BY
	case "cache": return CACHE
	case "called": return CALLED
	case "cancel": return CANCEL
	case "cancelquery": return CANCELQUERY
	case "capabilities": return CAPABILITIES
	case "capability": return CAPABILITY
	case "cascade": return CASCADE
	case "case": return CASE
	case "cast": return CAST
	case "changefeed": return CHANGEFEED
	case "char": return CHAR
	case "character": return CHARACTER
	case "characteristics": return CHARACTERISTICS
	case "check": return CHECK
	case "check_files": return CHECK_FILES
	case "close": return CLOSE
	case "cluster": return CLUSTER
	case "coalesce": return COALESCE
	case "collate": return COLLATE
	case "collation": return COLLATION
	case "column": return COLUMN
	case "columns": return COLUMNS
	case "comment": return COMMENT
	case "comments": return COMMENTS
	case "commit": return COMMIT
	case "committed": return COMMITTED
	case "compact": return COMPACT
	case "complete": return COMPLETE
	case "completions": return COMPLETIONS
	case "concurrently": return CONCURRENTLY
	case "configuration": return CONFIGURATION
	case "configurations": return CONFIGURATIONS
	case "configure": return CONFIGURE
	case "conflict": return CONFLICT
	case "connection": return CONNECTION
	case "connections": return CONNECTIONS
	case "constraint": return CONSTRAINT
	case "constraints": return CONSTRAINTS
	case "controlchangefeed": return CONTROLCHANGEFEED
	case "controljob": return CONTROLJOB
	case "conversion": return CONVERSION
	case "convert": return CONVERT
	case "copy": return COPY
	case "cost": return COST
	case "covering": return COVERING
	case "create": return CREATE
	case "createdb": return CREATEDB
	case "createlogin": return CREATELOGIN
	case "createrole": return CREATEROLE
	case "cross": return CROSS
	case "csv": return CSV
	case "cube": return CUBE
	case "current": return CURRENT
	case "current_catalog": return CURRENT_CATALOG
	case "current_date": return CURRENT_DATE
	case "current_role": return CURRENT_ROLE
	case "current_schema": return CURRENT_SCHEMA
	case "current_time": return CURRENT_TIME
	case "current_timestamp": return CURRENT_TIMESTAMP
	case "current_user": return CURRENT_USER
	case "cursor": return CURSOR
	case "cycle": return CYCLE
	case "data": return DATA
	case "database": return DATABASE
	case "databases": return DATABASES
	case "day": return DAY
	case "deallocate": return DEALLOCATE
	case "debug_dump_metadata_sst": return DEBUG_DUMP_METADATA_SST
	case "debug_ids": return DEBUG_IDS
	case "debug_pause_on": return DEBUG_PAUSE_ON
	case "dec": return DEC
	case "decimal": return DECIMAL
	case "declare": return DECLARE
	case "default": return DEFAULT
	case "defaults": return DEFAULTS
	case "deferrable": return DEFERRABLE
	case "deferred": return DEFERRED
	case "definer": return DEFINER
	case "delete": return DELETE
	case "delimiter": return DELIMITER
	case "depends": return DEPENDS
	case "desc": return DESC
	case "destination": return DESTINATION
	case "detached": return DETACHED
	case "details": return DETAILS
	case "discard": return DISCARD
	case "distinct": return DISTINCT
	case "do": return DO
	case "domain": return DOMAIN
	case "double": return DOUBLE
	case "drop": return DROP
	case "else": return ELSE
	case "encoding": return ENCODING
	case "encrypted": return ENCRYPTED
	case "encryption_info_dir": return ENCRYPTION_INFO_DIR
	case "encryption_passphrase": return ENCRYPTION_PASSPHRASE
	case "end": return END
	case "enum": return ENUM
	case "enums": return ENUMS
	case "escape": return ESCAPE
	case "except": return EXCEPT
	case "exclude": return EXCLUDE
	case "excluding": return EXCLUDING
	case "execute": return EXECUTE
	case "execution": return EXECUTION
	case "exists": return EXISTS
	case "experimental": return EXPERIMENTAL
	case "experimental_audit": return EXPERIMENTAL_AUDIT
	case "experimental_fingerprints": return EXPERIMENTAL_FINGERPRINTS
	case "experimental_relocate": return EXPERIMENTAL_RELOCATE
	case "experimental_replica": return EXPERIMENTAL_REPLICA
	case "expiration": return EXPIRATION
	case "explain": return EXPLAIN
	case "export": return EXPORT
	case "extension": return EXTENSION
	case "external": return EXTERNAL
	case "extract": return EXTRACT
	case "extract_duration": return EXTRACT_DURATION
	case "extremes": return EXTREMES
	case "failure": return FAILURE
	case "false": return FALSE
	case "family": return FAMILY
	case "fetch": return FETCH
	case "files": return FILES
	case "filter": return FILTER
	case "first": return FIRST
	case "float": return FLOAT
	case "following": return FOLLOWING
	case "for": return FOR
	case "force": return FORCE
	case "force_index": return FORCE_INDEX
	case "force_not_null": return FORCE_NOT_NULL
	case "force_null": return FORCE_NULL
	case "force_quote": return FORCE_QUOTE
	case "force_zigzag": return FORCE_ZIGZAG
	case "foreign": return FOREIGN
	case "format": return FORMAT
	case "forward": return FORWARD
	case "freeze": return FREEZE
	case "from": return FROM
	case "full": return FULL
	case "function": return FUNCTION
	case "functions": return FUNCTIONS
	case "generated": return GENERATED
	case "geography": return GEOGRAPHY
	case "geometry": return GEOMETRY
	case "geometrycollection": return GEOMETRYCOLLECTION
	case "geometrycollectionm": return GEOMETRYCOLLECTIONM
	case "geometrycollectionz": return GEOMETRYCOLLECTIONZ
	case "geometrycollectionzm": return GEOMETRYCOLLECTIONZM
	case "geometrym": return GEOMETRYM
	case "geometryz": return GEOMETRYZ
	case "geometryzm": return GEOMETRYZM
	case "global": return GLOBAL
	case "goal": return GOAL
	case "grant": return GRANT
	case "grants": return GRANTS
	case "greatest": return GREATEST
	case "group": return GROUP
	case "grouping": return GROUPING
	case "groups": return GROUPS
	case "hash": return HASH
	case "having": return HAVING
	case "header": return HEADER
	case "high": return HIGH
	case "histogram": return HISTOGRAM
	case "hold": return HOLD
	case "hour": return HOUR
	case "identity": return IDENTITY
	case "if": return IF
	case "iferror": return IFERROR
	case "ifnull": return IFNULL
	case "ignore_foreign_keys": return IGNORE_FOREIGN_KEYS
	case "ilike": return ILIKE
	case "immediate": return IMMEDIATE
	case "immutable": return IMMUTABLE
	case "import": return IMPORT
	case "in": return IN
	case "include": return INCLUDE
	case "include_all_secondary_tenants": return INCLUDE_ALL_SECONDARY_TENANTS
	case "including": return INCLUDING
	case "increment": return INCREMENT
	case "incremental": return INCREMENTAL
	case "incremental_location": return INCREMENTAL_LOCATION
	case "index": return INDEX
	case "indexes": return INDEXES
	case "index_after_order_by_before_at": return INDEX_AFTER_ORDER_BY_BEFORE_AT
	case "index_before_name_then_paren": return INDEX_BEFORE_NAME_THEN_PAREN
	case "index_before_paren": return INDEX_BEFORE_PAREN
	case "inherits": return INHERITS
	case "initially": return INITIALLY
	case "inject": return INJECT
	case "inner": return INNER
	case "inout": return INOUT
	case "input": return INPUT
	case "insensitive": return INSENSITIVE
	case "insert": return INSERT
	case "int": return INT
	case "integer": return INTEGER
	case "intersect": return INTERSECT
	case "interval": return INTERVAL
	case "into": return INTO
	case "into_db": return INTO_DB
	case "inverted": return INVERTED
	case "invisible": return INVISIBLE
	case "invoker": return INVOKER
	case "is": return IS
	case "iserror": return ISERROR
	case "isnull": return ISNULL
	case "isolation": return ISOLATION
	case "job": return JOB
	case "jobs": return JOBS
	case "join": return JOIN
	case "json": return JSON
	case "key": return KEY
	case "keys": return KEYS
	case "kms": return KMS
	case "kv": return KV
	case "label": return LABEL
	case "language": return LANGUAGE
	case "last": return LAST
	case "lateral": return LATERAL
	case "latest": return LATEST
	case "lc_collate": return LC_COLLATE
	case "lc_ctype": return LC_CTYPE
	case "leading": return LEADING
	case "leakproof": return LEAKPROOF
	case "lease": return LEASE
	case "least": return LEAST
	case "left": return LEFT
	case "less": return LESS
	case "level": return LEVEL
	case "like": return LIKE
	case "limit": return LIMIT
	case "linestring": return LINESTRING
	case "linestringm": return LINESTRINGM
	case "linestringz": return LINESTRINGZ
	case "linestringzm": return LINESTRINGZM
	case "list": return LIST
	case "local": return LOCAL
	case "locality": return LOCALITY
	case "localtime": return LOCALTIME
	case "localtimestamp": return LOCALTIMESTAMP
	case "locked": return LOCKED
	case "login": return LOGIN
	case "lookup": return LOOKUP
	case "low": return LOW
	case "match": return MATCH
	case "materialized": return MATERIALIZED
	case "maxvalue": return MAXVALUE
	case "merge": return MERGE
	case "method": return METHOD
	case "minute": return MINUTE
	case "minvalue": return MINVALUE
	case "modifyclustersetting": return MODIFYCLUSTERSETTING
	case "modifysqlclustersetting": return MODIFYSQLCLUSTERSETTING
	case "month": return MONTH
	case "move": return MOVE
	case "multilinestring": return MULTILINESTRING
	case "multilinestringm": return MULTILINESTRINGM
	case "multilinestringz": return MULTILINESTRINGZ
	case "multilinestringzm": return MULTILINESTRINGZM
	case "multipoint": return MULTIPOINT
	case "multipointm": return MULTIPOINTM
	case "multipointz": return MULTIPOINTZ
	case "multipointzm": return MULTIPOINTZM
	case "multipolygon": return MULTIPOLYGON
	case "multipolygonm": return MULTIPOLYGONM
	case "multipolygonz": return MULTIPOLYGONZ
	case "multipolygonzm": return MULTIPOLYGONZM
	case "names": return NAMES
	case "nan": return NAN
	case "natural": return NATURAL
	case "never": return NEVER
	case "new_db_name": return NEW_DB_NAME
	case "new_kms": return NEW_KMS
	case "next": return NEXT
	case "no": return NO
	case "nocancelquery": return NOCANCELQUERY
	case "nocontrolchangefeed": return NOCONTROLCHANGEFEED
	case "nocontroljob": return NOCONTROLJOB
	case "nocreatedb": return NOCREATEDB
	case "nocreatelogin": return NOCREATELOGIN
	case "nocreaterole": return NOCREATEROLE
	case "nologin": return NOLOGIN
	case "nomodifyclustersetting": return NOMODIFYCLUSTERSETTING
	case "none": return NONE
	case "nonvoters": return NONVOTERS
	case "normal": return NORMAL
	case "nosqllogin": return NOSQLLOGIN
	case "not": return NOT
	case "nothing": return NOTHING
	case "nothing_after_returning": return NOTHING_AFTER_RETURNING
	case "notnull": return NOTNULL
	case "noviewactivity": return NOVIEWACTIVITY
	case "noviewactivityredacted": return NOVIEWACTIVITYREDACTED
	case "noviewclustersetting": return NOVIEWCLUSTERSETTING
	case "nowait": return NOWAIT
	case "no_full_scan": return NO_FULL_SCAN
	case "no_index_join": return NO_INDEX_JOIN
	case "no_zigzag_join": return NO_ZIGZAG_JOIN
	case "null": return NULL
	case "nullif": return NULLIF
	case "nulls": return NULLS
	case "numeric": return NUMERIC
	case "of": return OF
	case "off": return OFF
	case "offset": return OFFSET
	case "oids": return OIDS
	case "old_kms": return OLD_KMS
	case "on": return ON
	case "only": return ONLY
	case "operator": return OPERATOR
	case "opt": return OPT
	case "option": return OPTION
	case "options": return OPTIONS
	case "or": return OR
	case "order": return ORDER
	case "ordinality": return ORDINALITY
	case "others": return OTHERS
	case "out": return OUT
	case "outer": return OUTER
	case "over": return OVER
	case "overlaps": return OVERLAPS
	case "overlay": return OVERLAY
	case "owned": return OWNED
	case "owner": return OWNER
	case "parallel": return PARALLEL
	case "parent": return PARENT
	case "partial": return PARTIAL
	case "partition": return PARTITION
	case "partitions": return PARTITIONS
	case "password": return PASSWORD
	case "pause": return PAUSE
	case "paused": return PAUSED
	case "physical": return PHYSICAL
	case "placement": return PLACEMENT
	case "placing": return PLACING
	case "plan": return PLAN
	case "plans": return PLANS
	case "point": return POINT
	case "pointm": return POINTM
	case "pointz": return POINTZ
	case "pointzm": return POINTZM
	case "polygon": return POLYGON
	case "polygonm": return POLYGONM
	case "polygonz": return POLYGONZ
	case "polygonzm": return POLYGONZM
	case "position": return POSITION
	case "preceding": return PRECEDING
	case "precision": return PRECISION
	case "prepare": return PREPARE
	case "preserve": return PRESERVE
	case "primary": return PRIMARY
	case "prior": return PRIOR
	case "priority": return PRIORITY
	case "privileges": return PRIVILEGES
	case "public": return PUBLIC
	case "publication": return PUBLICATION
	case "queries": return QUERIES
	case "query": return QUERY
	case "quote": return QUOTE
	case "range": return RANGE
	case "ranges": return RANGES
	case "read": return READ
	case "real": return REAL
	case "reason": return REASON
	case "reassign": return REASSIGN
	case "recurring": return RECURRING
	case "recursive": return RECURSIVE
	case "redact": return REDACT
	case "ref": return REF
	case "references": return REFERENCES
	case "refresh": return REFRESH
	case "region": return REGION
	case "regional": return REGIONAL
	case "regions": return REGIONS
	case "reindex": return REINDEX
	case "relative": return RELATIVE
	case "release": return RELEASE
	case "relocate": return RELOCATE
	case "rename": return RENAME
	case "repeatable": return REPEATABLE
	case "replace": return REPLACE
	case "replication": return REPLICATION
	case "reset": return RESET
	case "restart": return RESTART
	case "restore": return RESTORE
	case "restrict": return RESTRICT
	case "restricted": return RESTRICTED
	case "resume": return RESUME
	case "retention": return RETENTION
	case "retry": return RETRY
	case "return": return RETURN
	case "returning": return RETURNING
	case "returns": return RETURNS
	case "revision_history": return REVISION_HISTORY
	case "revoke": return REVOKE
	case "right": return RIGHT
	case "role": return ROLE
	case "roles": return ROLES
	case "rollback": return ROLLBACK
	case "rollup": return ROLLUP
	case "routines": return ROUTINES
	case "row": return ROW
	case "rows": return ROWS
	case "rule": return RULE
	case "running": return RUNNING
	case "savepoint": return SAVEPOINT
	case "scans": return SCANS
	case "scatter": return SCATTER
	case "schedule": return SCHEDULE
	case "schedules": return SCHEDULES
	case "schema": return SCHEMA
	case "schemas": return SCHEMAS
	case "schema_only": return SCHEMA_ONLY
	case "scroll": return SCROLL
	case "scrub": return SCRUB
	case "search": return SEARCH
	case "second": return SECOND
	case "secondary": return SECONDARY
	case "security": return SECURITY
	case "select": return SELECT
	case "sequence": return SEQUENCE
	case "sequences": return SEQUENCES
	case "serializable": return SERIALIZABLE
	case "server": return SERVER
	case "service": return SERVICE
	case "session": return SESSION
	case "sessions": return SESSIONS
	case "session_user": return SESSION_USER
	case "set": return SET
	case "setof": return SETOF
	case "sets": return SETS
	case "setting": return SETTING
	case "settings": return SETTINGS
	case "share": return SHARE
	case "shared": return SHARED
	case "show": return SHOW
	case "similar": return SIMILAR
	case "simple": return SIMPLE
	case "skip": return SKIP
	case "skip_localities_check": return SKIP_LOCALITIES_CHECK
	case "skip_missing_foreign_keys": return SKIP_MISSING_FOREIGN_KEYS
	case "skip_missing_sequences": return SKIP_MISSING_SEQUENCES
	case "skip_missing_sequence_owners": return SKIP_MISSING_SEQUENCE_OWNERS
	case "skip_missing_udfs": return SKIP_MISSING_UDFS
	case "skip_missing_views": return SKIP_MISSING_VIEWS
	case "smallint": return SMALLINT
	case "snapshot": return SNAPSHOT
	case "some": return SOME
	case "split": return SPLIT
	case "sql": return SQL
	case "sqllogin": return SQLLOGIN
	case "stable": return STABLE
	case "start": return START
	case "state": return STATE
	case "statements": return STATEMENTS
	case "statistics": return STATISTICS
	case "status": return STATUS
	case "stdin": return STDIN
	case "stdout": return STDOUT
	case "stop": return STOP
	case "storage": return STORAGE
	case "store": return STORE
	case "stored": return STORED
	case "storing": return STORING
	case "stream": return STREAM
	case "strict": return STRICT
	case "string": return STRING
	case "subscription": return SUBSCRIPTION
	case "substring": return SUBSTRING
	case "super": return SUPER
	case "support": return SUPPORT
	case "survival": return SURVIVAL
	case "survive": return SURVIVE
	case "symmetric": return SYMMETRIC
	case "syntax": return SYNTAX
	case "system": return SYSTEM
	case "table": return TABLE
	case "tables": return TABLES
	case "tablespace": return TABLESPACE
	case "temp": return TEMP
	case "template": return TEMPLATE
	case "temporary": return TEMPORARY
	case "tenant": return TENANT
	case "tenants": return TENANTS
	case "tenant_name": return TENANT_NAME
	case "testing_relocate": return TESTING_RELOCATE
	case "text": return TEXT
	case "then": return THEN
	case "throttling": return THROTTLING
	case "ties": return TIES
	case "time": return TIME
	case "timestamp": return TIMESTAMP
	case "timestamptz": return TIMESTAMPTZ
	case "timetz": return TIMETZ
	case "to": return TO
	case "trace": return TRACE
	case "tracing": return TRACING
	case "trailing": return TRAILING
	case "transaction": return TRANSACTION
	case "transactions": return TRANSACTIONS
	case "transfer": return TRANSFER
	case "transform": return TRANSFORM
	case "treat": return TREAT
	case "trigger": return TRIGGER
	case "trim": return TRIM
	case "true": return TRUE
	case "truncate": return TRUNCATE
	case "trusted": return TRUSTED
	case "type": return TYPE
	case "types": return TYPES
	case "unbounded": return UNBOUNDED
	case "uncommitted": return UNCOMMITTED
	case "union": return UNION
	case "unique": return UNIQUE
	case "unknown": return UNKNOWN
	case "unlisten": return UNLISTEN
	case "unlogged": return UNLOGGED
	case "unsafe_restore_incompatible_version": return UNSAFE_RESTORE_INCOMPATIBLE_VERSION
	case "unset": return UNSET
	case "unsplit": return UNSPLIT
	case "until": return UNTIL
	case "update": return UPDATE
	case "upsert": return UPSERT
	case "use": return USE
	case "user": return USER
	case "users": return USERS
	case "using": return USING
	case "valid": return VALID
	case "validate": return VALIDATE
	case "value": return VALUE
	case "values": return VALUES
	case "varbit": return VARBIT
	case "varchar": return VARCHAR
	case "variadic": return VARIADIC
	case "varying": return VARYING
	case "verify_backup_table_data": return VERIFY_BACKUP_TABLE_DATA
	case "view": return VIEW
	case "viewactivity": return VIEWACTIVITY
	case "viewactivityredacted": return VIEWACTIVITYREDACTED
	case "viewclustermetadata": return VIEWCLUSTERMETADATA
	case "viewclustersetting": return VIEWCLUSTERSETTING
	case "viewdebug": return VIEWDEBUG
	case "virtual": return VIRTUAL
	case "visible": return VISIBLE
	case "volatile": return VOLATILE
	case "voters": return VOTERS
	case "when": return WHEN
	case "where": return WHERE
	case "window": return WINDOW
	case "with": return WITH
	case "within": return WITHIN
	case "without": return WITHOUT
	case "work": return WORK
	case "write": return WRITE
	case "year": return YEAR
	case "zone": return ZONE
	default: return IDENT
	}
}
