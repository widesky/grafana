package migrations

import (
	. "github.com/grafana/grafana/pkg/services/sqlstore/migrator"
)

func addWideSkyProvisionerMigrations(mg *Migrator) {
	permissionsV1 := Table{
		Name: "ws_team_plugin_permissions",
		Columns: []*Column{
			{Name: "id", Type: DB_BigInt, IsPrimaryKey: true, IsAutoIncrement: true},
			{Name: "team_id", Type: DB_BigInt, Nullable: false},
			{Name: "plugin_id", Type: DB_NVarchar, Length: 190, Nullable: false},
			{Name: "last_modified", Type: DB_BigInt, Nullable: false},
			{Name: "page_access", Type: DB_Varchar, Length: 200, Nullable: false},
		},
		Indices: []*Index{
			{Cols: []string{"team_id"}, Type: IndexType},
		},
	}

	permissionsV2 := Table{
		Name: "widesky_permissions",
		Columns: []*Column{
			{Name: "id", Type: DB_BigInt, IsPrimaryKey: true, IsAutoIncrement: true},
			{Name: "team_id", Type: DB_BigInt, Nullable: false},
			{Name: "plugin_id", Type: DB_NVarchar, Length: 190, Nullable: false},
		},
	}

	// V1
	mg.AddMigration("create WideSky permissions v1", NewAddTableMigration(permissionsV1))

	// V2
	addTableReplaceMigrations(mg, permissionsV1, permissionsV2, 2, map[string]string{
		"id":        "id",
		"team_id":   "team_id",
		"plugin_id": "plugin_id",
	})
	mg.AddMigration("rename plugin_id name column to permission_id", NewRenameColumnMigration(
		permissionsV2, permissionsV2.Columns[2], "permission_id",
	))
	mg.AddMigration("add column type to permission rows", NewAddColumnMigration(permissionsV2, &Column{
		Name: "type", Type: DB_NVarchar, Nullable: false, Default: "plugin",
	}))
}
