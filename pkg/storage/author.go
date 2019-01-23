package storage

import (
	fmt "fmt"
	"regexp"
	"strconv"
	"time"
)

var (
	authorLine = regexp.MustCompile(`(.+)\s\<(.+)\>\s(\d+)\s?([\+-]\d+)?`)
	epoch      = time.Unix(0, 0).In(time.UTC)
)

// Signature holds the information given
type Signature struct {
	Name  string
	Email string
	Date  time.Time
}

func (s Signature) Validate() error {
	if len(s.Name) == 0 {
		return fmt.Errorf("no name")
	}
	if len(s.Email) == 0 {
		return fmt.Errorf("no email")
	}
	if s.Date.Before(epoch) {
		return fmt.Errorf("invalid date")
	}
	return nil
}

func offsetToString(offset int) string {
	offsetHour := int(offset / 60 / 60)
	offsetMinute := int(offset / 60 % 60)
	return fmt.Sprintf("%+03d%02d", offsetHour, offsetMinute)
}

func (s Signature) String() string {
	_, offset := s.Date.Zone()
	return fmt.Sprintf("%s <%s> %d %s", s.Name, s.Email, s.Date.Unix(), offsetToString(offset))
}

func parseSignature(line string) (Signature, error) {
	committer := authorLine.FindStringSubmatch(line)

	if len(committer) == 0 {
		return Signature{}, fmt.Errorf("could not parse signature")
	}
	if len(committer) < 5 {
		return Signature{}, fmt.Errorf("not a valid signature")
	}
	t, err := strconv.ParseInt(committer[3], 10, 64)
	if err != nil {
		return Signature{}, err
	}

	// Gracefully stolen from https://github.com/src-d/go-git/blob/434611b74cb54538088c6aeed4ed27d3044064fa/plumbing/object/object.go#L141-L149
	//
	// Include a dummy year in this time.Parse() call to avoid a bug in Go:
	// https://github.com/golang/go/issues/19750
	//
	// Parsing the timezone with no other details causes the tl.Location() call
	// below to return time.Local instead of the parsed zone in some cases
	tl, err := time.Parse("2006 -0700", "1970 "+committer[4])
	if err != nil {
		return Signature{}, err
	}
	tu := time.Unix(t, 0).In(tl.Location())

	return Signature{
		Name:  committer[1],
		Email: committer[2],
		Date:  tu,
	}, nil
}
