package cmd

// comments       string

type commentsSorter []hbaRule

func (slice commentsSorter) Len() int           { return len(slice) }
func (slice commentsSorter) Less(i, j int) bool { return slice[i].comments < slice[j].comments }
func (slice commentsSorter) Swap(i, j int)      { slice[i], slice[j] = slice[j], slice[i] }

type authTypeSorter []hbaRule

func (slice authTypeSorter) Len() int           { return len(slice) }
func (slice authTypeSorter) Less(i, j int) bool { return slice[i].authType < slice[j].authType }
func (slice authTypeSorter) Swap(i, j int)      { slice[i], slice[j] = slice[j], slice[i] }

type networkMaskSorter []hbaRule

func (slice networkMaskSorter) Len() int           { return len(slice) }
func (slice networkMaskSorter) Less(i, j int) bool { return slice[i].networkMask < slice[j].networkMask }
func (slice networkMaskSorter) Swap(i, j int)      { slice[i], slice[j] = slice[j], slice[i] }

type ipAddressSorter []hbaRule

func (slice ipAddressSorter) Len() int           { return len(slice) }
func (slice ipAddressSorter) Less(i, j int) bool { return slice[i].ipAddress < slice[j].ipAddress }
func (slice ipAddressSorter) Swap(i, j int)      { slice[i], slice[j] = slice[j], slice[i] }

type userNameSorter []hbaRule

func (slice userNameSorter) Len() int           { return len(slice) }
func (slice userNameSorter) Less(i, j int) bool { return slice[i].userName < slice[j].userName }
func (slice userNameSorter) Swap(i, j int)      { slice[i], slice[j] = slice[j], slice[i] }

type lineNumberSorter []hbaRule

func (slice lineNumberSorter) Len() int           { return len(slice) }
func (slice lineNumberSorter) Less(i, j int) bool { return slice[i].lineNumber < slice[j].lineNumber }
func (slice lineNumberSorter) Swap(i, j int)      { slice[i], slice[j] = slice[j], slice[i] }

type databaseNameSorter []hbaRule

func (slice databaseNameSorter) Len() int { return len(slice) }
func (slice databaseNameSorter) Less(i, j int) bool {
	return slice[i].databaseName < slice[j].databaseName
}
func (slice databaseNameSorter) Swap(i, j int) { slice[i], slice[j] = slice[j], slice[i] }

type connectionTypeSorter []hbaRule

func (slice connectionTypeSorter) Len() int { return len(slice) }
func (slice connectionTypeSorter) Less(i, j int) bool {
	return slice[i].connectionType < slice[j].connectionType
}
func (slice connectionTypeSorter) Swap(i, j int) { slice[i], slice[j] = slice[j], slice[i] }
