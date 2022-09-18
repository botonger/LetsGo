# Automate code writing with Template

This code was created to avoid the hassle of writing the same code repeatedly.

Go has a built-in flag parsing library. 

Separate the names of the methods you want to create in the method flag with a space and pass the component types (service, serviceImpl, controller) in the component flag.

> go run main.go -methods="testOne testTwo" -component="controller"

or

> go build -o template
> 
> ./template -methods="testOne testTwo" -component="controller"

or 

> go install
> template -methods="testOne testTwo" -component="controller"