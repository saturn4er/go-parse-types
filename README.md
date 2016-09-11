## go-parse types ##
Library to parse types defined in golang source code 

### Installation ###

#### Install package ####
    
    go get github.com/saturn4er/go-parse-types
    
### Usage example###

    package main  
    
    import "github.com/saturn4er/go-parse-types"
    
    type SomeType struct {
        a, b int
        c bool
    }
    
    func main() {
        parser, err = tparser.New("./test_package")
        if err != nil {
            panic(err)
        }
        err = parser.Parse()
        if err != nil {
            panic(err)
        }
        type := tparser.GetTypeByName("SomeType")
        fmt.Println(type.Kind)                    // struct
        fmt.Println(type.Name)                    // SomeType
        fmt.Println(len(type.Fields))             // 3
    }
 
     