package main

import (
        "fmt"
        "io/ioutil"
        "os"
        "regexp"
        // "bytes"
        "strings"
)

func main(){
}

func cleanDomain(query string) string{

    parts := strings.Split(query, ".")

    var base string
    l := len(parts)
    if l == 2 {
        base = query
    } else if l >= 3  {
        // fancy calulations
        if isMultiTLD(parts[l - 1]) {
            base = lastN(parts, 3)
        } else {
            base = lastN(parts, 2)
        }
    } else { panic(1) }

    return base
}

func lastN(a []string, i int) string {
    l := len(a)
    b := a[(l - i):]
    all := strings.Join(b, ".")
    return all
}

func isMultiTLD(last string) bool {

    countries := map[string]bool{
    "cy": true,
    "ro": true,
    "ke": true,
    "kh": true,
    "ki": true,
    "cr": true,
    "km": true,
    "kn": true,
    "kr": true,
    "ck": true,
    "cn": true,
    "kw": true,
    "rs": true,
    "ca": true,
    "kz": true,
    "rw": true,
    "ru": true,
    "za": true,
    "zm": true,
    "bz": true,
    "je": true,
    "uy": true,
    "bs": true,
    "br": true,
    "jo": true,
    "us": true,
    "bh": true,
    "bo": true,
    "bn": true,
    "bb": true,
    "ba": true,
    "ua": true,
    "eg": true,
    "ec": true,
    "et": true,
    "er": true,
    "es": true,
    "pl": true,
    "in": true,
    "ph": true,
    "il": true,
    "pe": true,
    "co": true,
    "pa": true,
    "id": true,
    "py": true,
    "ug": true,
    "ky": true,
    "ir": true,
    "pt": true,
    "pw": true,
    "iq": true,
    "it": true,
    "pr": true,
    "sh": true,
    "sl": true,
    "sn": true,
    "sa": true,
    "sb": true,
    "sc": true,
    "sd": true,
    "se": true,
    "hk": true,
    "sg": true,
    "sy": true,
    "sz": true,
    "st": true,
    "sv": true,
    "om": true,
    "th": true,
    "ve": true,
    "tz": true,
    "vn": true,
    "vi": true,
    "pk": true,
    "fk": true,
    "fj": true,
    "fr": true,
    "ni": true,
    "ng": true,
    "nf": true,
    "re": true,
    "na": true,
    "qa": true,
    "tw": true,
    "nr": true,
    "np": true,
    "ac": true,
    "af": true,
    "ae": true,
    "ao": true,
    "al": true,
    "yu": true,
    "ar": true,
    "tj": true,
    "at": true,
    "au": true,
    "ye": true,
    "mv": true,
    "mw": true,
    "mt": true,
    "mu": true,
    "tr": true,
    "mz": true,
    "tt": true,
    "mx": true,
    "my": true,
    "mg": true,
    "me": true,
    "mc": true,
    "ma": true,
    "mn": true,
    "mo": true,
    "ml": true,
    "mk": true,
    "do": true,
    "dz": true,
    "ps": true,
    "lr": true,
    "tn": true,
    "lv": true,
    "ly": true,
    "lb": true,
    "lk": true,
    "gg": true,
    "uk": true,
    "gn": true,
    "gh": true,
    "gt": true,
    "gu": true,
    "jp": true,
    "gr": true,
    "nz": true,
}

    return countries[last]

}

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
