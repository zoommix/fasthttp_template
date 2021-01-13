package utils

const (
	// PrefixInfo – light blue
	PrefixInfo = "\t" + CInfo + "INFO:" + CEnd + "\t"
	// PrefixNotice – purple
	PrefixNotice = "\t" + CNotice + "NOTICE:" + CEnd + "\t"
	// PrefixWarning – yellow
	PrefixWarning = "\t" + CWarning + "WARNING:" + CEnd + "\t"
	// PrefixError – red
	PrefixError = "\t" + CError + "ERROR:" + CEnd + "\t"
	// PrefixDebug – cyan
	PrefixDebug = "\t" + CDebug + "DEBUG:" + CEnd + "\t"
	// PrefixEnd – of the color
	PrefixEnd = "\t" + "\033[0m"

	// CInfo light blue
	CInfo = "\033[1;34m"
	// CNotice purple
	CNotice = "\033[1;35m"
	// CWarning yellow
	CWarning = "\033[1;33m"
	// CError red
	CError = "\033[1;31m"
	// CDebug cyan
	CDebug = "\033[0;36m"
	// CEnd of the color
	CEnd = "\033[0m"
)
