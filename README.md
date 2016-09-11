## go-parse types ##
Library to parse types defined in golang source code 

### Installation ###
    
    go get github.com/saturn4er/go-parse-types
    
### Usage example###

    package main  
    
    import "github.com/saturn4er/go-parse-types"
    
    type SomeType struct {
        a, b int
        c bool
    }
    
    func main() {
        parser, err := tparser.New("./test_package")
        if err != nil {
            panic(err)
        }
        err = parser.Parse()
        if err != nil {
            panic(err)
        }
        type, err := tparser.GetType("SomeType")
        if err != nil {
            panic(err)
        }
        fmt.Println(type.Kind)                    // struct
        fmt.Println(type.Name)                    // SomeType
        fmt.Println(len(type.Fields))             // 3
        
        _, err := tparser.GetType("ABC")
        fmt.Println(err)                          // No such type
    }
 
### TODO ###

- parse interfaces