# Habit
 Is a simple tool to keep track of your habits. It can be run from the command-line or as an HTTP server.

## Usage
### CLI mode
To start a new habit, for example a daily piano habit, simply type `habit piano`  on your terminal. As you repeat your habit daily, 
retype `habit piano`. You can have multiple habits at the same time, simply type `habit surfing` to start a new surfing 
habit. You can list all your streaks with `habit`. Also, you can create a weekly habit by passing the `weekly` option
like so `habit -f weekly piano`
```
Usage of habit:
  -f string
    	Set the frecuency of the habit: daily, weekly. (shorthand) (default "daily")
  -frequency string
    	Set the frecuency of the habit: daily, weekly. (default "daily")
  -s	Runs habit as a HTTP Server (shorthand)
  -server
    	Runs habit as a HTTP Server
```

### Server Mode
To start Habit as a server type `habit -s 127.0.0.1:8080`. If no address is provided `127.0.0.1:8080` is set as the
default. Use your browser to talk to the server as follows:
* To create a new habit or continue your streak type `http://127.0.0.1:8080/?habit=HabitName`.
* To list all habits go to `http://127.0.0.1:8080/` or `http://127.0.0.1:8080/all`.
* By default habits are created as daily habits. You can specify a weekly habit by passing the `interval=weekly`
  `http://127.0.0.1:8080/?habit=HabitName&interval=weekly`.

 ## Installation
Install by running `go install https://github.com/crmejia/habit/cmd/habit@latest`