package whois

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/coopernurse/gorp"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"sync"
)

type Result struct {
	Emails []string          `json:"emails"`
	Data   map[string]string `json:"data"`
	Raw    string            `json:"raw"`
	Domain string
}

type SqlResult struct {
	Domain string
	Raw    string
	Emails string
}

func ParseCSV(f io.Reader) [][]string {
	csvReader := csv.NewReader(f)
	cont, _ := csvReader.ReadAll()
	fmt.Printf("Content: %#v", cont)
	return cont
}

func Retrieve(query string, db *gorp.DbMap) (*Result, error) {
	q := cleanDomain(query)
	data := []SqlResult{}
	_, err := db.Select(&data, "select * from WhoisResults where Domain=?", q)
	var result Result

	if len(data) == 0 || err != nil {
		fmt.Println(q, " not found in database")
		jwhois := fmt.Sprintf("./whois.sh")
		cmd := exec.Command(jwhois, q)

		record, err := cmd.Output()
		sRecord := strings.TrimSpace(string(record))

		if err != nil {
			result = Result{
				Emails: []string{},
				Data:   make(map[string]string),
				Raw:    "",
				Domain: q,
			}
		} else {
			result = newResult(sRecord, q)
		}
		emailString := strings.Join(result.Emails, ", ")
		dbResult := SqlResult{Domain: result.Domain, Raw: result.Raw, Emails: emailString}
		err = db.Insert(&dbResult)
	} else {
		result = newDBResult(data[0])
	}

	return &result, err
}

func newDBResult(record SqlResult) Result {
	lines := strings.Split(record.Emails, "\n")
	ourMap := toMap(lines)
	emailArray := strings.Split(record.Emails, ", ")

	result := Result{
		Emails: emailArray,
		Data:   ourMap,
		Raw:    record.Raw,
		Domain: record.Domain,
	}
	return result
}

func newResult(sRecord string, query string) Result {
	lines := strings.Split(sRecord, "\n")
	ourMap := toMap(lines)
	emailArray := emails(sRecord)

	result := Result{
		Emails: emailArray,
		Data:   ourMap,
		Raw:    sRecord,
		Domain: query,
	}
	return result
}

func AsyncRetrieve(domain string, db *gorp.DbMap, messages chan<- string, wg *sync.WaitGroup) {
	rec, err := Retrieve(domain, db)
	if err == nil {
		emails := strings.Join(rec.Emails, " ")
		output := strings.Join([]string{domain, emails}, ", ")
		messages <- output
	} else {
		fmt.Println(err)
	}
	wg.Done()
}

func RetrieveJSON(query string, db *gorp.DbMap) string {
	record, err := Retrieve(query, db)
	check(err)

	js, _ := json.Marshal(&record)
	return string(js)
}

func emails(rawRecord string) []string {
	lines := strings.Split(rawRecord, "\n")
	hash := toMap(lines)
	emails := getEmails(hash)

	if len(emails) == 0 {
		emails = grepEmails(rawRecord)
		fmt.Printf("slice - %#v\n", emails)
	}

	return emails
}

func grepEmails(c string) []string {
	emailRegex := "[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\\.[a-zA-Z0-9-.]+"
	r, _ := regexp.Compile(emailRegex)
	return r.FindAllString(c, -1)
}

func getEmails(h map[string]string) []string {

	fmt.Printf("emails - %#v\n", h)
	s := make([]string, 0)
	for k, v := range h {
		match, _ := regexp.MatchString("email", k)
		if match {
			s = append(s, v)
		}
	}
	return s

}

func dPrint(s ...interface{}) {
	if os.Getenv("WHOIS_DEBUG") != "" {
		fmt.Printf("%#v", s)
	}
}

func toMap(lines []string) map[string]string {

	hash := make(map[string]string)

	for _, line := range lines {
		if strings.Contains(line, ": ") {
			a := strings.SplitN(line, ":", 2)

			key := strings.ToLower(strings.TrimSpace(a[0]))
			val := strings.TrimSpace(a[1])

			hash[key] = val
		}
	}
	return hash
}

func getData(file string) string {
	dat, err := ioutil.ReadFile(file)
	check(err)

	return string(dat)
}

func check(e error) {
	if e != nil {
		fmt.Printf("ERROR: %#v", e)
		panic(1)
	}
}
