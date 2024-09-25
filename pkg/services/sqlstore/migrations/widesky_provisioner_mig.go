package migrations

import (
	. "github.com/grafana/grafana/pkg/services/sqlstore/migrator"
)

func addWideSkyProvisionerMigrations(mg *Migrator) {
	teamPluginPermV1 := Table{
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

	mg.AddMigration("create WideSky team plugin perm table v1", NewAddTableMigration(teamPluginPermV1))
}
