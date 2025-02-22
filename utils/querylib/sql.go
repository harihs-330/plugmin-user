package querylib

import (
	"fmt"
	"strings"
)

// generic function to add condition to query
func AddCondition(query string, column string, args []interface{}) string {
	if strings.Contains(query, "WHERE") {
		return fmt.Sprintf("%s AND %s = $%d", query, column, len(args))
	}
	return fmt.Sprintf("%s WHERE %s = $%d", query, column, len(args))
}

// adds groupBy statement
func GroupBy(sql string, columns ...string) string {
	if len(columns) > 0 {
		sql += " GROUP BY "
	}
	for i, column := range columns {
		sql = fmt.Sprintf("%s %s", sql, column)
		if i < len(columns)-1 {
			sql += ","
		}
	}

	return sql
}

// GeneratePlaceholders generates a string of placeholders for a multi-value insert query.
// `recordCount` specifies the number of records sets to insert.
// `columnCount` specifies the number of columns in each record set.
func GeneratePlaceholders(recordCount, columnCount int) string {
	// Create a slice to hold the placeholder strings for each value
	valueStrings := make([]string, recordCount)

	for count := 0; count < recordCount; count++ {
		// Generate placeholders for each value set
		placeholders := make([]string, columnCount)
		for j := 0; j < columnCount; j++ {
			placeholders[j] = fmt.Sprintf("$%d", count*columnCount+j+1)
		}
		// Join placeholders for the current value set and add to valueStrings
		valueStrings[count] = fmt.Sprintf("(%s)", strings.Join(placeholders, ", "))
	}

	// Join all value strings with commas
	return strings.Join(valueStrings, ", ")
}
