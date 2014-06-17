package whois

import (
	// "github.com/codegangsta/cli"
    // "bytes"
    "fmt"
    "io/ioutil"
    "os"
    "os/exec"
    "regexp"
    "strings"
    "encoding/json"
	"encoding/csv"
	"io"
)

type Result struct {
    Emails []string          `json:"emails"`
    Data   map[string]string `json:"data"`
    Raw    string            `json:"raw"`
}

func ParseCSV (f io.Reader) [][]string {
	// f, _ := os.Open("sample.csv")
	// defer f.Close()

	csvReader := csv.NewReader(f)
	cont, _   := csvReader.ReadAll()
	fmt.Printf("Content: %#v", cont)
	return cont
}

func Retrieve(query string) (*Result, error) {
    // record := getData("fixtures/civet.com_verisign")
    jwhois := fmt.Sprintf("./whois.sh")
    q := cleanDomain(query)
    cmd := exec.Command(jwhois, q)

    record, e := cmd.Output()
    check(e)

    // dPrint(jwhois)
    // dPrint(record)
    // dPrint(e)

    sRecord := strings.TrimSpace(string(record))
    lines := strings.Split(sRecord, "\n")
    ourMap := toMap(lines)
    emailArray := emails(sRecord)

    result := &Result{
        Emails: emailArray,
        Data:   ourMap,
        Raw:    sRecord,
    }

    fmt.Printf("STRUCT: %#v", result)

    return result, nil
}


func AsyncRetrieve(domain string, messages chan<- string, done chan<- bool){
	rec, _ := Retrieve(domain)
	emails := strings.Join(rec.Emails, " ")
	output := strings.Join([]string{domain, emails}, ", ")
	messages <- output
	done <- true
}

func RetrieveJSON(query string) string {
    record, err := Retrieve(query)
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

func check(e error){
    if e != nil {
        fmt.Printf("ERROR: %#v", e)
        panic(1)
    }
}
