package core

import (
	"bytes"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDisableLevelTruncation(t *testing.T) {
	entry := &Entry{
		Time:    time.Now(),
		Message: "testing",
	}
	keys := []string{}
	timestampFormat := "Mon Jan 2 15:04:05 -0700 MST 2006"
	checkDisableTruncation := func(disabled bool, level Level) {
		tf := &TextFormatter{DisableLevelTruncation: disabled}
		var b bytes.Buffer
		entry.Level = level
		tf.printColored(&b, entry, keys, timestampFormat)
		logLine := (&b).String()
		if disabled {
			expected := strings.ToUpper(level.String())
			if !strings.Contains(logLine, expected) {
				t.Errorf("level string expected to be %s when truncation disabled", expected)
			}
		} else {
			expected := strings.ToUpper(level.String())
			if len(level.String()) > 4 {
				if strings.Contains(logLine, expected) {
					t.Errorf("level string %s expected to be truncated to %s when truncation is enabled", expected, expected[0:4])
				}
			} else {
				if !strings.Contains(logLine, expected) {
					t.Errorf("level string expected to be %s when truncation is enabled and level string is below truncation threshold", expected)
				}
			}
		}
	}

	checkDisableTruncation(true, DebugLevel)
	checkDisableTruncation(true, InfoLevel)
	checkDisableTruncation(false, ErrorLevel)
	checkDisableTruncation(false, InfoLevel)
}

func TestTextFormatterFieldMap(t *testing.T) {
	formatter := &TextFormatter{
		DisableColors: true,
		FieldMap: FieldMap{
			FieldKeyMsg:   "message",
			FieldKeyLevel: "somelevel",
			FieldKeyTime:  "timeywimey",
		},
	}

	entry := &Entry{
		Message: "oh hi",
		Level:   WarnLevel,
		Time:    time.Date(1981, time.February, 24, 4, 28, 3, 100, time.UTC),
		Data: Fields{
			"field1":     "f1",
			"message":    "messagefield",
			"somelevel":  "levelfield",
			"timeywimey": "timeywimeyfield",
		},
	}

	b, err := formatter.Format(entry)
	if err != nil {
		t.Fatal("Unable to format entry: ", err)
	}

	assert.Equal(t,
		`timeywimey="1981-02-24T04:28:03Z" `+
			`somelevel=warning `+
			`message="oh hi" `+
			`field1=f1 `+
			`fields.message=messagefield `+
			`fields.somelevel=levelfield `+
			`fields.timeywimey=timeywimeyfield`+"\n",
		string(b),
		"Formatted output doesn't respect FieldMap")
}

func TestTextFormatterIsColored(t *testing.T) {
	params := []struct {
		name               string
		expectedResult     bool
		isTerminal         bool
		disableColor       bool
		forceColor         bool
		envColor           bool
		clicolorIsSet      bool
		clicolorForceIsSet bool
		clicolorVal        string
		clicolorForceVal   string
	}{
		// Default values
		{
			name:               "testcase1",
			expectedResult:     false,
			isTerminal:         false,
			disableColor:       false,
			forceColor:         false,
			envColor:           false,
			clicolorIsSet:      false,
			clicolorForceIsSet: false,
		},
		// Output on terminal
		{
			name:               "testcase2",
			expectedResult:     true,
			isTerminal:         true,
			disableColor:       false,
			forceColor:         false,
			envColor:           false,
			clicolorIsSet:      false,
			clicolorForceIsSet: false,
		},
		// Output on terminal with color disabled
		{
			name:               "testcase3",
			expectedResult:     false,
			isTerminal:         true,
			disableColor:       true,
			forceColor:         false,
			envColor:           false,
			clicolorIsSet:      false,
			clicolorForceIsSet: false,
		},
		// Output not on terminal with color disabled
		{
			name:               "testcase4",
			expectedResult:     false,
			isTerminal:         false,
			disableColor:       true,
			forceColor:         false,
			envColor:           false,
			clicolorIsSet:      false,
			clicolorForceIsSet: false,
		},
		// Output not on terminal with color forced
		{
			name:               "testcase5",
			expectedResult:     true,
			isTerminal:         false,
			disableColor:       false,
			forceColor:         true,
			envColor:           false,
			clicolorIsSet:      false,
			clicolorForceIsSet: false,
		},
		// Output on terminal with clicolor set to "0"
		{
			name:               "testcase6",
			expectedResult:     false,
			isTerminal:         true,
			disableColor:       false,
			forceColor:         false,
			envColor:           true,
			clicolorIsSet:      true,
			clicolorForceIsSet: false,
			clicolorVal:        "0",
		},
		// Output on terminal with clicolor set to "1"
		{
			name:               "testcase7",
			expectedResult:     true,
			isTerminal:         true,
			disableColor:       false,
			forceColor:         false,
			envColor:           true,
			clicolorIsSet:      true,
			clicolorForceIsSet: false,
			clicolorVal:        "1",
		},
		// Output not on terminal with clicolor set to "0"
		{
			name:               "testcase8",
			expectedResult:     false,
			isTerminal:         false,
			disableColor:       false,
			forceColor:         false,
			envColor:           true,
			clicolorIsSet:      true,
			clicolorForceIsSet: false,
			clicolorVal:        "0",
		},
		// Output not on terminal with clicolor set to "1"
		{
			name:               "testcase9",
			expectedResult:     false,
			isTerminal:         false,
			disableColor:       false,
			forceColor:         false,
			envColor:           true,
			clicolorIsSet:      true,
			clicolorForceIsSet: false,
			clicolorVal:        "1",
		},
		// Output not on terminal with clicolor set to "1" and force color
		{
			name:               "testcase10",
			expectedResult:     true,
			isTerminal:         false,
			disableColor:       false,
			forceColor:         true,
			envColor:           true,
			clicolorIsSet:      true,
			clicolorForceIsSet: false,
			clicolorVal:        "1",
		},
		// Output not on terminal with clicolor set to "0" and force color
		{
			name:               "testcase11",
			expectedResult:     false,
			isTerminal:         false,
			disableColor:       false,
			forceColor:         true,
			envColor:           true,
			clicolorIsSet:      true,
			clicolorForceIsSet: false,
			clicolorVal:        "0",
		},
		// Output not on terminal with clicolor_force set to "1"
		{
			name:               "testcase12",
			expectedResult:     true,
			isTerminal:         false,
			disableColor:       false,
			forceColor:         false,
			envColor:           true,
			clicolorIsSet:      false,
			clicolorForceIsSet: true,
			clicolorForceVal:   "1",
		},
		// Output not on terminal with clicolor_force set to "0"
		{
			name:               "testcase13",
			expectedResult:     false,
			isTerminal:         false,
			disableColor:       false,
			forceColor:         false,
			envColor:           true,
			clicolorIsSet:      false,
			clicolorForceIsSet: true,
			clicolorForceVal:   "0",
		},
		// Output on terminal with clicolor_force set to "0"
		{
			name:               "testcase14",
			expectedResult:     false,
			isTerminal:         true,
			disableColor:       false,
			forceColor:         false,
			envColor:           true,
			clicolorIsSet:      false,
			clicolorForceIsSet: true,
			clicolorForceVal:   "0",
		},
	}

	cleanenv := func() {
		os.Unsetenv("CLICOLOR")
		os.Unsetenv("CLICOLOR_FORCE")
	}

	defer cleanenv()

	for _, val := range params {
		t.Run("textformatter_"+val.name, func(subT *testing.T) {
			tf := TextFormatter{
				IsTerminal:                val.isTerminal,
				DisableColors:             val.disableColor,
				ForceColors:               val.forceColor,
				EnvironmentOverrideColors: val.envColor,
			}
			cleanenv()
			if val.clicolorIsSet {
				os.Setenv("CLICOLOR", val.clicolorVal)
			}
			if val.clicolorForceIsSet {
				os.Setenv("CLICOLOR_FORCE", val.clicolorForceVal)
			}
			res := tf.isColored()
			assert.Equal(subT, val.expectedResult, res)
		})
	}
}

// TODO add tests for sorting etc., this requires a parser for the text
// formatter output.
