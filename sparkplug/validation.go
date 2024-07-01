package sparkplug

import (
	"regexp"
)

// topicRegexp contains all the regular expression patterns for SparkplugB 3.0 MQTT commands
var topicRegexp = []*regexp.Regexp{}

func init() {

	patterns := []string{
		`^` + SPB_NS + `/` + string(STATE) + `/[^/]+$`,
		`^` + SPB_NS + `/[^/]+/` + string(DBIRTH) + `/[^/]+/[^/]+$`,
		`^` + SPB_NS + `/[^/]+/` + string(DDATA) + `/[^/]+/[^/]+$`,
		`^` + SPB_NS + `/[^/]+/` + string(DDEATH) + `/[^/]+/[^/]+$`,
		`^` + SPB_NS + `/[^/]+/` + string(NBIRTH) + `/[^/]+$`,
		`^` + SPB_NS + `/[^/]+/` + string(NDATA) + `/[^/]+$`,
		`^` + SPB_NS + `/[^/]+/` + string(NCMD) + `/[^/]+$`,
		`^` + SPB_NS + `/[^/]+/` + string(NDEATH) + `/[^/]+$`,
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

	if len(topic) < len(SPB_NS) {
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
