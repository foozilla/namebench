package io

import (    
    "bufio"
    "os"
)

//GetSites returns list of popular websites
func GetSites()(sites []string, err error) {

    ss := make([]string, 100)
    inFile, _ := os.Open("popular_websites.txt")
    defer inFile.Close()
    scanner := bufio.NewScanner(inFile)
        scanner.Split(bufio.ScanLines) 
    i := 0
    for scanner.Scan() {        
        ss[i] = scanner.Text()
        i++
    }

    return ss, nil    
}