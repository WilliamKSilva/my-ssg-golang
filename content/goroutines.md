# Introduction

- This repository contains examples of concurrency using golang *Goroutines*.
- *Goroutines* are lightweight threads managed by golang runtime that provide concurrent
processing of functions.

# When to use Goroutines
- As mentioned before, Goroutines provide concurrency for golang, the right question should be, "When should I use concurrency?"
Concurrency is a strong method of processing that allows code to be executed "at the same time".
Example: You have a batch of multiple users payment that needs to be processed through an external API service.
You can write a function that calls this payment API and use Goroutines so all the processing is executed
"at the same time".
- Under the hood Goroutines use execution scheduling, so it just appear that everything is being executed at the same time, but what it is really happening is a smart use of processing time. The scheduler can check if a thread is currently not being used (waiting for IO, for example, an HTTP response) and put another thread to work. That way, the process is never really stopped, there's always some piece of code bein run.
- A core principle of concurrency is that multiple tasks can be in progress at once, but not necessarily *physically* at the same time. What's really happening is sharing off processing time. One task is independent of another and can finish before or after the others.
- This is how concurrency differs from paralelism: both of the methods exist to execute multible things "at once", but concurrency creates the appearence of simultaneous execution through scheduling, parallelism actually performs processing at the same time, using multiple processors or cores at the physical level.
- All parallelism is concurrent, since the tasks are independent of each other and will be executed each one on their own time, but not all concurrency is parallelism, since in many implementations what is being used is processing time share and not multiple physical processing units.

# Examples

### Page Download
- To run this example you can use `go run page_download/main.go`.

- If you want to personalize the test you can pass the flag `-concurrency=true or -concurrency=false` to use Goroutines or not, the default is `true`. You can also pass how much pages should be downloaded using `-pages=100`, the default is `50`.

- On the first example *page_download* the idea is to make multiple GET HTTP request to a random website to download its content and check how much time will take
with and without Goroutines.
Since HTTP requests are IO bound and in this case the response of one download webpage not depend on the others we can spawn a Goroutine to make this operations
concurrently.

- The results *with* Goroutines: running 100 requests our process ran on *~0.29* seconds.
- The results *without* Goroutines: running 100 requests our process ran on *~3.4* seconds. 

### User Registration
- To run this example you can use `go run user_registration/main.go`.

- If you want to personalize the test you can pass the flag `-concurrency=true or -concurrency=false` to use Goroutines or not, the default is `true`. You can also pass how much threads the program will use with `-threads=5`, the default is `4`.

- On this example the idea was to read 200 users from two CSV files and populate their address data
based on their CEP (a identifier of address in Brazil) with the street name and the city name using a third party REST API.

- At first glance I thought that the huge difference would be when we read the data from the CSV files, 
but the difference was actually minimal on this step of the test. With Goroutines the time to read all users was around ~0.9 seconds and without was ~0.8/~0.9. A good explanation for this is the fact that the IO required to read data from the files is actually prety fast, since them are local we have *0* latency and are only dependent on our disk speed. Spawning Goroutines in this step can actually make the performance a little worse in some cases, since we have to spend processing time managing our threads.

- The real difference in performance using Goroutines where on the "enrich" step of the process. Since we need to make a
request to an external API to search for the details of our users addresses, this step can take some time, since it is a IO bound and is heavily dependent on the latency of the request. My first idea to make this work was to take all the users we read from the CSV files and make small *batchs* of all the users and call our function that does the enrichment of the address data using a Goroutine.

- The results *with* Goroutines: reading 200 users from a CSV file and searching for their address took *~8.7* seconds.
- The results *without* Goroutines: reading 200 users from a CSV file and searching for their address took *~124.0* seconds.

### Worker Pool
- This example demonstrates a simple service check, where we send a GET request to a website and determine whether the response status code indicates success.

- The main idea here is to implement a thread pool pattern, where we define a fixed number of workers (goroutines) to process tasks concurrently. Each worker is responsible for executing a specific "job" from the queue.

- Since we’re using Go, we also leverage the concept introduced in the previous example: channels. Each worker listens for incoming jobs on the jobs channel, executes the assigned task, and then sends the result to the results channel.

- In this case, a worker executes a job — checking the availability of a specific website. Once the check is complete, the worker sends a custom result to the results channel indicating whether the check was successful.