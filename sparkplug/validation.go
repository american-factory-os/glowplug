package sparkplug

import (
	"regexp"
)

// topicRegexp contains all the regular expression patterns for SparkplugB 3.0 MQTT commands
var topicRegexp = []*regexp.Regexp{}

func init() {

	// Regex component for Group ID, Edge Node ID, and Device ID
	// Allows alphanumeric, hyphen, underscore, and period; excludes MQTT reserved chars (+, #, /)
	idComponent := `[a-zA-Z0-9][a-zA-Z0-9_\-\.]*`

	patterns := []string{
		// STATE: spBv1.0/STATE/<ClientID>
		`^` + SPB_NS + `/` + string(STATE) + `/` + idComponent + `$`,
		// NBIRTH: spBv1.0/<GroupID>/NBIRTH/<EdgeNodeID>
		`^` + SPB_NS + `/` + idComponent + `/` + string(NBIRTH) + `/` + idComponent + `$`,
		// NDEATH: spBv1.0/<GroupID>/NDEATH/<EdgeNodeID>
		`^` + SPB_NS + `/` + idComponent + `/` + string(NDEATH) + `/` + idComponent + `$`,
		// NDATA: spBv1.0/<GroupID>/NDATA/<EdgeNodeID>
		`^` + SPB_NS + `/` + idComponent + `/` + string(NDATA) + `/` + idComponent + `$`,
		// NCMD: spBv1.0/<GroupID>/NCMD/<EdgeNodeID>
		`^` + SPB_NS + `/` + idComponent + `/` + string(NCMD) + `/` + idComponent + `$`,
		// DBIRTH: spBv1.0/<GroupID>/DBIRTH/<EdgeNodeID>/<DeviceID>
		`^` + SPB_NS + `/` + idComponent + `/` + string(DBIRTH) + `/` + idComponent + `/` + idComponent + `$`,
		// DDEATH: spBv1.0/<GroupID>/DDEATH/<EdgeNodeID>/<DeviceID>
		`^` + SPB_NS + `/` + idComponent + `/` + string(DDEATH) + `/` + idComponent + `/` + idComponent + `$`,
		// DDATA: spBv1.0/<GroupID>/DDATA/<EdgeNodeID>/<DeviceID>
		`^` + SPB_NS + `/` + idComponent + `/` + string(DDATA) + `/` + idComponent + `/` + idComponent + `$`,
		// DCMD: spBv1.0/<GroupID>/DCMD/<EdgeNodeID>/<DeviceID>
		`^` + SPB_NS + `/` + idComponent + `/` + string(DCMD) + `/` + idComponent + `/` + idComponent + `$`,
	}

	for _, pattern := range patterns {
		topicRegexp = append(topicRegexp, regexp.MustCompile(pattern))
	}

}

// IsValidSparkplugBTopic checks if a given topic is a valid SparkplugB 3.0 topic
func IsValidSparkplugBTopic(topic string) bool {

	if len(topic) == 0 {
		return false
	}

	if len(topic) < len(SPB_NS)+1 {
		return false
	}

	for _, re := range topicRegexp {
		match := re.MatchString(topic)
		if match {
			return true
		}
	}

	return false
}
