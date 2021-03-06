# Habit
Is a simple tool to keep track of your habits. It can be run from the command-line or as an HTTP server.

## CLI mode
### Installation
Install by running `go install github.com/crmejia/habit/cmd/habit@latest`

### Usage
To start a new habit, for example a daily piano habit, fire up your terminal and type: 
```
$habit piano
Good luck with your new habit 'piano'! Don't forget to do it again tomorrow.  
```
As you repeat your daily habit, log your activity with:
```
$habit piano
Nice work: you've done the habit 'piano' for 1 days in a row now. Keep it up!
```
You can have multiple habits at the same time, simply type `habit surfing` to start a new surfing 
habit. You can list all your streaks with `habit all`. Also, you can create a weekly habit by passing the `weekly` option
like so `habit -f weekly piano`
```
Usage: habit <Option Flags> <HABIT_NAME> -- to create/update a new habit
       habit all   --   to list all habits
Option Flags:
  -d string
    	Set the store directory. User's home directory is the default (default "/Users/crismar")
  -f string
    	Set the frequency of the habit: daily(default), weekly. (default "daily")
  -s string
    	Set the store backend for habit tracker: db(default), file (default "db")
```

## Server Mode
### Installation
Install by running `go install github.com/crmejia/habit/cmd/server@latest`

### Usage
To start Habit as a server type:
```
$ habit 127.0.0.1:8080
Starting HTTP server
```
Use your browser to talk to the server as follows:
* To create a new habit or continue your streak type `http://127.0.0.1:8080/?habit=HabitName`.
* To list all habits go to `http://127.0.0.1:8080/all`.
* By default, habits are created as daily habits. You can specify a weekly habit by passing the `interval=weekly`
  `http://127.0.0.1:8080/?habit=HabitName&interval=weekly`.

