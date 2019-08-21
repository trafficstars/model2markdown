package main

type Job struct {
	InputFile string
}

type Config struct {
	OutputDirectory string
	Jobs            []Job
}
