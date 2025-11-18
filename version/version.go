package version

import (
	"strconv"
	"strings"
)

// Compare compare > compared return 1;compare > compared return 0;compare < compared return -1
func Compare(compare, compared string) int {
	if compare == compared {
		return 0
	}

	if compare == "" {
		return -1
	}

	if compared == "" {
		return 1
	}

	compareVersions := strings.Split(compare, ".")
	compareVersionCount := len(compareVersions)

	comparedVersions := strings.Split(compared, ".")
	comparedVersionCount := len(comparedVersions)

	minVersionCount := compareVersionCount
	if comparedVersionCount < compareVersionCount {
		minVersionCount = comparedVersionCount
	}

	for i := 0; i < minVersionCount; i++ {
		compareVersion, err := strconv.ParseInt(compareVersions[i], 10, 64)
		if err != nil {
			return -1
		}

		comparedVersion, err := strconv.ParseInt(comparedVersions[i], 10, 64)
		if err != nil {
			return 1
		}

		if compareVersion < comparedVersion {
			return -1
		} else if compareVersion > comparedVersion {
			return 1
		}
	}

	if compareVersionCount > comparedVersionCount {
		return 1
	} else if compareVersionCount == comparedVersionCount {
		return 0
	} else {
		return -1
	}
}
