package mysql

import (
	"database/sql"
	"fmt"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/migrator"
	"gorm.io/gorm/schema"
)

type Migrator struct {
	migrator.Migrator
	Dialector
}

func (m Migrator) FullDataTypeOf(field *schema.Field) clause.Expr {
	expr := m.Migrator.FullDataTypeOf(field)

	if value, ok := field.TagSettings["COMMENT"]; ok {
		expr.SQL += " COMMENT " + m.Dialector.Explain("?", value)
	}

	return expr
}

func (m Migrator) AlterColumn(value interface{}, field string) error {
	return m.RunWithValue(value, func(stmt *gorm.Statement) error {
		if field := stmt.Schema.LookUpField(field); field != nil {
			return m.DB.Exec(
				"ALTER TABLE ? MODIFY COLUMN ? ?",
				clause.Table{Name: stmt.Table}, clause.Column{Name: field.DBName}, m.FullDataTypeOf(field),
			).Error
		}
		return fmt.Errorf("failed to look up field with name: %s", field)
	})
}

func (m Migrator) RenameColumn(value interface{}, oldName, newName string) error {
	return m.RunWithValue(value, func(stmt *gorm.Statement) error {
		if !m.Dialector.DontSupportRenameColumn {
			return m.Migrator.RenameColumn(value, oldName, newName)
		}

		var field *schema.Field
		if f := stmt.Schema.LookUpField(oldName); f != nil {
			oldName = f.DBName
			field = f
		}

		if f := stmt.Schema.LookUpField(newName); f != nil {
			newName = f.DBName
			field = f
		}

		if field != nil {
			return m.DB.Exec(
				"ALTER TABLE ? CHANGE ? ? ?",
				clause.Table{Name: stmt.Table}, clause.Column{Name: oldName},
				clause.Column{Name: newName}, m.FullDataTypeOf(field),
			).Error
		}

		return fmt.Errorf("failed to look up field with name: %s", newName)
	})
}

func (m Migrator) RenameIndex(value interface{}, oldName, newName string) error {
	if !m.Dialector.DontSupportRenameIndex {
		return m.RunWithValue(value, func(stmt *gorm.Statement) error {
			return m.DB.Exec(
				"ALTER TABLE ? RENAME INDEX ? TO ?",
				clause.Table{Name: stmt.Table}, clause.Column{Name: oldName}, clause.Column{Name: newName},
			).Error
		})
	}

	return m.RunWithValue(value, func(stmt *gorm.Statement) error {
		err := m.DropIndex(value, oldName)
		if err != nil {
			return err
		}

		if idx := stmt.Schema.LookIndex(newName); idx == nil {
			if idx = stmt.Schema.LookIndex(oldName); idx != nil {
				opts := m.BuildIndexOptions(idx.Fields, stmt)
				values := []interface{}{clause.Column{Name: newName}, clause.Table{Name: stmt.Table}, opts}

				createIndexSQL := "CREATE "
				if idx.Class != "" {
					createIndexSQL += idx.Class + " "
				}
				createIndexSQL += "INDEX ? ON ??"

				if idx.Type != "" {
					createIndexSQL += " USING " + idx.Type
				}

				return m.DB.Exec(createIndexSQL, values...).Error
			}
		}

		return m.CreateIndex(value, newName)
	})

}

func (m Migrator) DropTable(values ...interface{}) error {
	values = m.ReorderModels(values, false)
	return m.DB.Connection(func(tx *gorm.DB) error {
		tx.Exec("SET FOREIGN_KEY_CHECKS = 0;")
		for i := len(values) - 1; i >= 0; i-- {
			if err := m.RunWithValue(values[i], func(stmt *gorm.Statement) error {
				return tx.Exec("DROP TABLE IF EXISTS ? CASCADE", clause.Table{Name: stmt.Table}).Error
			}); err != nil {
				return err
			}
		}
		return tx.Exec("SET FOREIGN_KEY_CHECKS = 1;").Error
	})
}

func (m Migrator) DropConstraint(value interface{}, name string) error {
	return m.RunWithValue(value, func(stmt *gorm.Statement) error {
		constraint, chk, table := m.GuessConstraintAndTable(stmt, name)
		if chk != nil {
			return m.DB.Exec("ALTER TABLE ? DROP CHECK ?", clause.Table{Name: stmt.Table}, clause.Column{Name: chk.Name}).Error
		}
		if constraint != nil {
			name = constraint.Name
		}

		return m.DB.Exec(
			"ALTER TABLE ? DROP FOREIGN KEY ?", clause.Table{Name: table}, clause.Column{Name: name},
		).Error
	})
}

// ColumnTypes column types return columnTypes,error
func (m Migrator) ColumnTypes(value interface{}) ([]gorm.ColumnType, error) {
	columnTypes := make([]gorm.ColumnType, 0)
	err := m.RunWithValue(value, func(stmt *gorm.Statement) error {
		var (
			currentDatabase = m.DB.Migrator().CurrentDatabase()
			columnTypeSQL   = "SELECT column_name, column_default, is_nullable = 'YES', data_type, character_maximum_length, column_type, column_key, extra, column_comment, numeric_precision, numeric_scale "
			rows, err       = m.DB.Session(&gorm.Session{}).Table(stmt.Table).Limit(1).Rows()
		)

		if err != nil {
			return err
		}

		rawColumnTypes, err := rows.ColumnTypes()

		if err := rows.Close(); err != nil {
			return err
		}

		if !m.DisableDatetimePrecision {
			columnTypeSQL += ", datetime_precision "
		}
		columnTypeSQL += "FROM information_schema.columns WHERE table_schema = ? AND table_name = ? ORDER BY ORDINAL_POSITION"

		columns, rowErr := m.DB.Raw(columnTypeSQL, currentDatabase, stmt.Table).Rows()
		if rowErr != nil {
			return rowErr
		}

		defer columns.Close()

		for columns.Next() {
			var (
				column            migrator.ColumnType
				datetimePrecision sql.NullInt64
				extraValue        sql.NullString
				columnKey         sql.NullString
				values            = []interface{}{
					&column.NameValue, &column.DefaultValueValue, &column.NullableValue, &column.DataTypeValue, &column.LengthValue, &column.ColumnTypeValue, &columnKey, &extraValue, &column.CommentValue, &column.DecimalSizeValue, &column.ScaleValue,
				}
			)

			if !m.DisableDatetimePrecision {
				values = append(values, &datetimePrecision)
			}

			if scanErr := columns.Scan(values...); scanErr != nil {
				return scanErr
			}

			column.PrimaryKeyValue = sql.NullBool{Bool: false, Valid: true}
			column.UniqueValue = sql.NullBool{Bool: false, Valid: true}
			switch columnKey.String {
			case "PRI":
				column.PrimaryKeyValue = sql.NullBool{Bool: true, Valid: true}
			case "UNI":
				column.UniqueValue = sql.NullBool{Bool: true, Valid: true}
			}

			if strings.Contains(extraValue.String, "auto_increment") {
				column.AutoIncrementValue = sql.NullBool{Bool: true, Valid: true}
			}

			column.DefaultValueValue.String = strings.Trim(column.DefaultValueValue.String, "'")

			if datetimePrecision.Valid {
				column.DecimalSizeValue = datetimePrecision
			}

			for _, c := range rawColumnTypes {
				if c.Name() == column.NameValue.String {
					column.SQLColumnType = c
					break
				}
			}

			columnTypes = append(columnTypes, column)
		}

		return nil
	})

	return columnTypes, err
}

func (m Migrator) CurrentDatabase() (name string) {
	baseName := m.Migrator.CurrentDatabase()
	m.DB.Raw(
		"SELECT SCHEMA_NAME from Information_schema.SCHEMATA where SCHEMA_NAME LIKE ? ORDER BY SCHEMA_NAME=? DESC limit 1",
		baseName+"%", baseName).Scan(&name)
	return
}

func (m Migrator) GetTables() (tableList []string, err error) {
	err = m.DB.Raw("SELECT TABLE_NAME FROM information_schema.tables where TABLE_SCHEMA=?", m.CurrentDatabase()).
		Scan(&tableList).Error
	return
}
