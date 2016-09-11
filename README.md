## go-parse types ##
Library to parse types defined in golang source code 

### Installation ###

#### Install package ####
    
    go get github.com/saturn4er/go-parse-types
    
### Usage example###

    package main  
    func main() {
        parser, err = New("./test_package")
        if err != nil {
            panic(err)
        }
        err = parser.Parse()
        if err != nil {
            panic(err)
        }
        type := parser.GetTypeByName("SomeType")
        
    }
 
     