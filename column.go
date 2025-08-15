package carta

import (
	"database/sql"
	"sort"
	"strings"
)

// column represents the ith struct field of this mapper where the column is to be mapped
type column struct {
	typ         *sql.ColumnType
	name        string
	columnIndex int
	i           fieldIndex
}

func allocateColumns(m *Mapper, columns map[string]column) error {
	var (
		candidates map[string]bool
	)
	presentColumns := map[string]column{}
	for cName, c := range columns {
		if m.IsBasic {
			candidates = getColumnNameCandidates("", m.AncestorNames, m.Delimiter)
			if _, ok := candidates[cName]; ok {
				presentColumns[cName] = column{
					typ:         c.typ,
					name:        cName,
					columnIndex: c.columnIndex,
				}
				delete(columns, cName) // dealocate claimed column
			}
		} else {
			for i, field := range m.Fields {
				subMap, isSubMap := m.SubMaps[i]
				delimiter := m.Delimiter
				if isSubMap {
					delimiter = subMap.Delimiter
				}
				candidates = getColumnNameCandidates(field.Name, m.AncestorNames, delimiter)
				// can only allocate columns to basic fields
				if isBasicType(field.Typ) {
					if _, ok := candidates[cName]; ok {
						presentColumns[cName] = column{
							typ:         c.typ,
							name:        cName,
							columnIndex: c.columnIndex,
							i:           i,
						}
						delete(columns, cName) // dealocate claimed column
					}
				}
			}
		}
	}
	m.PresentColumns = presentColumns

	columnIds := []int{}
	for _, column := range m.PresentColumns {
		if _, ok := m.SubMaps[column.i]; ok {
			continue
		}
		columnIds = append(columnIds, column.columnIndex)
	}
	sort.Ints(columnIds)
	m.SortedColumnIndexes = columnIds

	ancestorNames := []string{}
	if len(m.AncestorNames) != 0 {
		ancestorNames = m.AncestorNames
	}

	for i, subMap := range m.SubMaps {
		subMap.AncestorNames = append(ancestorNames, m.Fields[i].Name)
		if err := allocateColumns(subMap, columns); err != nil {
			return err
		}
	}
	return nil
}

func getColumnNameCandidates(fieldName string, ancestorNames []string, delimiter string) map[string]bool {
	// empty field name means that the mapper is basic, since there is no struct assiciated with this slice, there is no field name
	candidates := map[string]bool{}
	if fieldName != "" {
		candidates[fieldName] = true
		candidates[toSnakeCase(fieldName)] = true
		candidates[strings.ToLower(fieldName)] = true
	}
	if len(ancestorNames) == 0 {
		return candidates
	}
	nameConcat := fieldName
	snakeConcat := toSnakeCase(fieldName)
	for i := len(ancestorNames) - 1; i >= 0; i-- {
		ancestor := ancestorNames[i]
		snakeAncestor := toSnakeCase(ancestor)

		if nameConcat == "" {
			nameConcat = ancestor
			snakeConcat = snakeAncestor
		} else {
			nameConcat = ancestor + delimiter + nameConcat
			snakeConcat = snakeAncestor + "_" + snakeConcat
		}
		candidates[nameConcat] = true
		candidates[strings.ToLower(nameConcat)] = true
		candidates[snakeConcat] = true
		candidates[strings.ToLower(snakeConcat)] = true
	}
	return candidates
}

func toSnakeCase(s string) string {
	delimiter := "_"
	s = strings.Trim(s, " ")
	n := ""
	for i, v := range s {
		nextCaseIsChanged := false
		if i+1 < len(s) {
			next := s[i+1]
			vIsCap := v >= 'A' && v <= 'Z'
			vIsLow := v >= 'a' && v <= 'z'
			nextIsCap := next >= 'A' && next <= 'Z'
			nextIsLow := next >= 'a' && next <= 'z'
			if (vIsCap && nextIsLow) || (vIsLow && nextIsCap) {
				nextCaseIsChanged = true
			}
		}

		if i > 0 && n[len(n)-1] != uint8(delimiter[0]) && nextCaseIsChanged {
			if v >= 'A' && v <= 'Z' {
				n += string(delimiter) + string(v)
			} else if v >= 'a' && v <= 'z' {
				n += string(v) + string(delimiter)
			}
		} else if v == ' ' || v == '-' {
			n += string(delimiter)
		} else {
			n = n + string(v)
		}
	}
	return strings.ToLower(n)
}
