package sqlserver

import (
	"database/sql"
	"fmt"
	"regexp"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/migrator"
)

type Migrator struct {
	migrator.Migrator
}

func (m Migrator) GetTables() (tableList []string, err error) {
	return tableList, m.DB.Raw("SELECT table_name FROM INFORMATION_SCHEMA.tables WHERE  table_catalog = ?", m.CurrentDatabase()).Scan(&tableList).Error
}

func (m Migrator) HasTable(value interface{}) bool {
	var count int
	m.RunWithValue(value, func(stmt *gorm.Statement) error {
		return m.DB.Raw(
			"SELECT count(*) FROM INFORMATION_SCHEMA.tables WHERE table_name = ? AND table_catalog = ?",
			stmt.Table, m.CurrentDatabase(),
		).Row().Scan(&count)
	})
	return count > 0
}

func (m Migrator) DropTable(values ...interface{}) error {
	values = m.ReorderModels(values, false)
	for i := len(values) - 1; i >= 0; i-- {
		tx := m.DB.Session(&gorm.Session{})
		if err := m.RunWithValue(values[i], func(stmt *gorm.Statement) error {
			type constraint struct {
				Name   string
				Parent string
			}
			var constraints []constraint
			err := tx.Raw("SELECT name, OBJECT_NAME(parent_object_id) as parent FROM sys.foreign_keys WHERE referenced_object_id = object_id(?)", stmt.Table).Scan(&constraints).Error

			for _, c := range constraints {
				if err == nil {
					err = tx.Exec("ALTER TABLE ? DROP CONSTRAINT ?;", gorm.Expr(c.Parent), gorm.Expr(c.Name)).Error
				}
			}

			if err == nil {
				err = tx.Exec("DROP TABLE IF EXISTS ?", clause.Table{Name: stmt.Table}).Error
			}

			return err
		}); err != nil {
			return err
		}
	}
	return nil
}

func (m Migrator) RenameTable(oldName, newName interface{}) error {
	var oldTable, newTable string
	if v, ok := oldName.(string); ok {
		oldTable = v
	} else {
		stmt := &gorm.Statement{DB: m.DB}
		if err := stmt.Parse(oldName); err == nil {
			oldTable = stmt.Table
		} else {
			return err
		}
	}

	if v, ok := newName.(string); ok {
		newTable = v
	} else {
		stmt := &gorm.Statement{DB: m.DB}
		if err := stmt.Parse(newName); err == nil {
			newTable = stmt.Table
		} else {
			return err
		}
	}

	return m.DB.Exec(
		"sp_rename @objname = ?, @newname = ?;",
		clause.Table{Name: oldTable}, clause.Table{Name: newTable},
	).Error
}

func (m Migrator) HasColumn(value interface{}, field string) bool {
	var count int64
	m.RunWithValue(value, func(stmt *gorm.Statement) error {
		currentDatabase := m.DB.Migrator().CurrentDatabase()
		name := field
		if field := stmt.Schema.LookUpField(field); field != nil {
			name = field.DBName
		}

		return m.DB.Raw(
			"SELECT count(*) FROM INFORMATION_SCHEMA.columns WHERE table_catalog = ? AND table_name = ? AND column_name = ?",
			currentDatabase, stmt.Table, name,
		).Row().Scan(&count)
	})

	return count > 0
}

func (m Migrator) AlterColumn(value interface{}, field string) error {
	return m.RunWithValue(value, func(stmt *gorm.Statement) error {
		if field := stmt.Schema.LookUpField(field); field != nil {
			fieldType := clause.Expr{SQL: m.DataTypeOf(field)}
			if field.NotNull {
				fieldType.SQL += " NOT NULL"
			} else {
				fieldType.SQL += " NULL"
			}

			return m.DB.Exec(
				"ALTER TABLE ? ALTER COLUMN ? ?",
				clause.Table{Name: stmt.Table}, clause.Column{Name: field.DBName}, fieldType,
			).Error
		}
		return fmt.Errorf("failed to look up field with name: %s", field)
	})
}

func (m Migrator) RenameColumn(value interface{}, oldName, newName string) error {
	return m.RunWithValue(value, func(stmt *gorm.Statement) error {
		if field := stmt.Schema.LookUpField(oldName); field != nil {
			oldName = field.DBName
		}

		if field := stmt.Schema.LookUpField(newName); field != nil {
			newName = field.DBName
		}

		return m.DB.Exec(
			"sp_rename @objname = ?, @newname = ?, @objtype = 'COLUMN';",
			fmt.Sprintf("%s.%s", stmt.Table, oldName), clause.Column{Name: newName},
		).Error
	})
}

var defaultValueTrimRegexp = regexp.MustCompile("^\\('?([^']*)'?\\)$")

// ColumnTypes return columnTypes []gorm.ColumnType and execErr error
func (m Migrator) ColumnTypes(value interface{}) ([]gorm.ColumnType, error) {
	columnTypes := make([]gorm.ColumnType, 0)
	execErr := m.RunWithValue(value, func(stmt *gorm.Statement) (err error) {
		rows, err := m.DB.Session(&gorm.Session{}).Table(stmt.Table).Limit(1).Rows()
		if err != nil {
			return err
		}

		rawColumnTypes, _ := rows.ColumnTypes()
		rows.Close()

		{
			var (
				columnTypeSQL   = "SELECT column_name, data_type, column_default, is_nullable, character_maximum_length, numeric_precision, numeric_precision_radix, numeric_scale, datetime_precision FROM INFORMATION_SCHEMA.COLUMNS WHERE table_catalog = ? AND table_name = ?"
				columns, rowErr = m.DB.Raw(columnTypeSQL, m.CurrentDatabase(), stmt.Table).Rows()
			)

			if rowErr != nil {
				return rowErr
			}

			for columns.Next() {
				var (
					column = migrator.ColumnType{
						PrimaryKeyValue: sql.NullBool{Valid: true},
						UniqueValue:     sql.NullBool{Valid: true},
					}
					datetimePrecision sql.NullInt64
					radixValue        sql.NullInt64
					nullableValue     sql.NullString
					values            = []interface{}{
						&column.NameValue, &column.ColumnTypeValue, &column.DefaultValueValue, &nullableValue, &column.LengthValue, &column.DecimalSizeValue, &radixValue, &column.ScaleValue, &datetimePrecision,
					}
				)

				if scanErr := columns.Scan(values...); scanErr != nil {
					return scanErr
				}

				if nullableValue.Valid {
					column.NullableValue = sql.NullBool{Bool: strings.EqualFold(nullableValue.String, "YES"), Valid: true}
				}

				if datetimePrecision.Valid {
					column.DecimalSizeValue = datetimePrecision
				}

				if column.DefaultValueValue.Valid {
					matches := defaultValueTrimRegexp.FindStringSubmatch(column.DefaultValueValue.String)
					for len(matches) > 1 {
						column.DefaultValueValue.String = matches[1]
						matches = defaultValueTrimRegexp.FindStringSubmatch(column.DefaultValueValue.String)
					}
				} else {
					column.DefaultValueValue.Valid = true
				}

				for _, c := range rawColumnTypes {
					if c.Name() == column.NameValue.String {
						column.SQLColumnType = c
						break
					}
				}

				columnTypes = append(columnTypes, column)
			}

			columns.Close()
		}

		{
			columnTypeRows, err := m.DB.Raw("SELECT c.column_name, t.constraint_type FROM information_schema.table_constraints t JOIN information_schema.constraint_column_usage c ON c.constraint_name=t.constraint_name WHERE t.constraint_type IN ('PRIMARY KEY', 'UNIQUE') AND c.table_catalog = ? AND c.table_name = ?", m.CurrentDatabase(), stmt.Table).Rows()
			if err != nil {
				return err
			}

			for columnTypeRows.Next() {
				var name, columnType string
				columnTypeRows.Scan(&name, &columnType)
				for idx, c := range columnTypes {
					mc := c.(migrator.ColumnType)
					if mc.NameValue.String == name {
						switch columnType {
						case "PRIMARY KEY":
							mc.PrimaryKeyValue = sql.NullBool{Bool: true, Valid: true}
						case "UNIQUE":
							mc.UniqueValue = sql.NullBool{Bool: true, Valid: true}
						}
						columnTypes[idx] = mc
						break
					}
				}
			}

			columnTypeRows.Close()
		}

		return
	})

	return columnTypes, execErr
}

func (m Migrator) HasIndex(value interface{}, name string) bool {
	var count int
	m.RunWithValue(value, func(stmt *gorm.Statement) error {
		if idx := stmt.Schema.LookIndex(name); idx != nil {
			name = idx.Name
		}

		return m.DB.Raw(
			"SELECT count(*) FROM sys.indexes WHERE name=? AND object_id=OBJECT_ID(?)",
			name, stmt.Table,
		).Row().Scan(&count)
	})
	return count > 0
}

func (m Migrator) RenameIndex(value interface{}, oldName, newName string) error {
	return m.RunWithValue(value, func(stmt *gorm.Statement) error {

		return m.DB.Exec(
			"sp_rename @objname = ?, @newname = ?, @objtype = 'INDEX';",
			fmt.Sprintf("%s.%s", stmt.Table, oldName), clause.Column{Name: newName},
		).Error
	})
}

func (m Migrator) HasConstraint(value interface{}, name string) bool {
	var count int64
	m.RunWithValue(value, func(stmt *gorm.Statement) error {
		constraint, chk, table := m.GuessConstraintAndTable(stmt, name)
		if constraint != nil {
			name = constraint.Name
		} else if chk != nil {
			name = chk.Name
		}

		return m.DB.Raw(
			`SELECT count(*) FROM sys.foreign_keys as F inner join sys.tables as T on F.parent_object_id=T.object_id inner join information_schema.tables as I on I.TABLE_NAME = T.name WHERE F.name = ?  AND T.Name = ? AND I.TABLE_CATALOG = ?;`,
			name, table, m.CurrentDatabase(),
		).Row().Scan(&count)
	})
	return count > 0
}

func (m Migrator) CurrentDatabase() (name string) {
	m.DB.Raw("SELECT DB_NAME() AS [Current Database]").Row().Scan(&name)
	return
}
