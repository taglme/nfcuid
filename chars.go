package main

type CharFlag int

const (
	CharFlagNone CharFlag = iota
	CharFlagSpace
	CharFlagTab
	CharFlagHyphen
	CharFlagEnter
	CharFlagSemiColon
	CharFlagColon
	CharFlagComma
)

type CharFlagDef struct {
	name   string
	output string
}

var charFlagStrings = map[CharFlag]CharFlagDef{
	CharFlagNone:      CharFlagDef{"none", ""},
	CharFlagSpace:     CharFlagDef{"space", " "},
	CharFlagTab:       CharFlagDef{"tab", "\\t"},
	CharFlagHyphen:    CharFlagDef{"hyphen", "-"},
	CharFlagEnter:     CharFlagDef{"enter", "\\n"},
	CharFlagSemiColon: CharFlagDef{"semicolon", ";"},
	CharFlagColon:     CharFlagDef{"colon", ":"},
	CharFlagComma:     CharFlagDef{"comma", ","},
}

func StringToCharFlag(s string) (CharFlag, bool) {
	for k, v := range charFlagStrings {
		if v.name == s {
			return k, true
		}
	}
	return 0, false
}

func (charFlag CharFlag) Name() string {
	return charFlagStrings[charFlag].name
}

func (charFlag CharFlag) Output() string {
	return charFlagStrings[charFlag].output
}

func CharFlagOptions() string {
	var s string
	for _, v := range charFlagStrings {
		s = s + "'" + v.name + "', "
	}
	return s
}
