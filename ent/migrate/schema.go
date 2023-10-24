// Code generated by ent, DO NOT EDIT.

package migrate

import (
	"entgo.io/ent/dialect/sql/schema"
	"entgo.io/ent/schema/field"
)

var (
	// UsersColumns holds the columns for the "users" table.
	UsersColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "username", Type: field.TypeString, Unique: true},
		{Name: "password_client_hashing_options", Type: field.TypeJSON},
		{Name: "password_server_hashing_options", Type: field.TypeJSON},
		{Name: "password_double_hashed", Type: field.TypeBytes},
		{Name: "certificate", Type: field.TypeString, Unique: true},
		{Name: "public_key", Type: field.TypeString, Unique: true},
		{Name: "encrypted_private_key", Type: field.TypeBytes},
		{Name: "totp_secret", Type: field.TypeString},
	}
	// UsersTable holds the schema information for the "users" table.
	UsersTable = &schema.Table{
		Name:       "users",
		Columns:    UsersColumns,
		PrimaryKey: []*schema.Column{UsersColumns[0]},
		Indexes: []*schema.Index{
			{
				Name:    "user_username",
				Unique:  false,
				Columns: []*schema.Column{UsersColumns[1]},
			},
			{
				Name:    "user_public_key",
				Unique:  false,
				Columns: []*schema.Column{UsersColumns[6]},
			},
		},
	}
	// Tables holds all the tables in the schema.
	Tables = []*schema.Table{
		UsersTable,
	}
)

func init() {
}
