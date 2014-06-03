package main

import (
        "fmt"
        "io/ioutil"
        "os"
        "regexp"
        // "bytes"
        "strings"
)

func Parse() {
    out, _ := ioutil.ReadAll(os.Stdin)
    outS := string(out)

    lines := strings.Split(outS, "\n")

    hash := toMap(lines)

    emails := getEmails(hash)

    if len(emails) == 0 {
        emails = grepEmails(outS)
        fmt.Printf("slice - %#v\n", emails)
    }


    // ['joe@example.com', ...]
    // if len(sArray) == 0 {
    //   emails = grepEmails(outS)
    // }



    // fmt.Println("STDIN: ", string(out))
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
