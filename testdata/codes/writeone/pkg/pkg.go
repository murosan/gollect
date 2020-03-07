package pkg

import "fmt"

import f "fmt"

var NonFormatted1=       func()int{
	return 100
}()



//  comment
func Abc()    {
	fmt.Println( "abc" )
}


func       Def(){f.Println("def" )}
